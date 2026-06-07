package analytics

import (
	"context"
	authDomain "server/internal/domain/auth"
	"server/internal/domain/analytics"
	"server/internal/domain/employees"
	orgDomain "server/internal/domain/organisation"
	"server/internal/platform/database/models"
	"time"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) analytics.Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) GetEmployeeAnalytics(ctx context.Context, userID uint, orgID uint, period string) (*analytics.EmployeeAnalyticsSummary, error) {
	// Parse period to get the start date
	now := time.Now()
	var startDate time.Time
	switch period {
	case "week":
		startDate = now.AddDate(0, 0, -7)
	case "quarter":
		startDate = now.AddDate(0, -3, 0)
	case "month":
		fallthrough
	default:
		startDate = now.AddDate(0, -1, 0)
	}
	startDateStr := startDate.Format("2006-01-02")

	// 1. Get Employee Info
	var empModel models.Employee
	if err := r.db.WithContext(ctx).Where("user_id = ? AND organisation_id = ?", userID, orgID).Preload("User").First(&empModel).Error; err != nil {
		return nil, err
	}

	emp := toEmployeeDomain(&empModel)

	summary := &analytics.EmployeeAnalyticsSummary{
		EmployeeInfo: emp,
		Attendance:   analytics.AnalyticsAttendanceBlock{},
		Leave: analytics.AnalyticsLeaveBlock{
			LeavesByType: make(map[string]int),
		},
		Tasks: analytics.AnalyticsTaskBlock{},
	}

	// 2. Aggregate Attendance
	type AttendanceStats struct {
		PresentDays          int
		TotalDurationMinutes int
		LateCheckins         int
	}
	var attStats AttendanceStats
	r.db.WithContext(ctx).Table("attendances").
		Select("COUNT(id) as present_days").
		Where("user_id = ? AND date >= ?", userID, startDateStr).
		Scan(&attStats.PresentDays)

	// Since SQLite doesn't natively have easy date diff for duration in minutes without some function,
	// let's do it via rows and calculate in memory for duration and late checkins
	rows, err := r.db.WithContext(ctx).Table("attendances").
		Select("check_in_time, check_out_time").
		Where("user_id = ? AND date >= ?", userID, startDateStr).
		Rows()

	if err == nil {
		defer rows.Close()
		var checkInTime, checkOutTime *time.Time
		var totalCheckInMinutes int
		for rows.Next() {
			rows.Scan(&checkInTime, &checkOutTime)
			if checkInTime != nil && checkOutTime != nil {
				attStats.TotalDurationMinutes += int(checkOutTime.Sub(*checkInTime).Minutes())
			}
			if checkInTime != nil {
				totalCheckInMinutes += checkInTime.Hour()*60 + checkInTime.Minute()
				if checkInTime.Hour() > 9 || (checkInTime.Hour() == 9 && checkInTime.Minute() > 30) {
					attStats.LateCheckins++
				}
			}
		}

		summary.Attendance.PresentDays = attStats.PresentDays
		summary.Attendance.TotalWorkingDays = int(now.Sub(startDate).Hours() / 24)
		summary.Attendance.AbsentDays = summary.Attendance.TotalWorkingDays - attStats.PresentDays
		if summary.Attendance.AbsentDays < 0 {
			summary.Attendance.AbsentDays = 0
		}

		if summary.Attendance.TotalWorkingDays > 0 {
			summary.Attendance.AttendanceRate = float64(attStats.PresentDays) / float64(summary.Attendance.TotalWorkingDays) * 100
		}
		if attStats.PresentDays > 0 {
			summary.Attendance.AverageDurationMinutes = attStats.TotalDurationMinutes / attStats.PresentDays

			avgCheckInMins := totalCheckInMinutes / attStats.PresentDays
			avgTime := time.Date(0, 1, 1, avgCheckInMins/60, avgCheckInMins%60, 0, 0, time.UTC)
			summary.Attendance.AverageCheckInTime = avgTime.Format("03:04 PM")
		}
		summary.Attendance.LateCheckins = attStats.LateCheckins
	}

	// 3. Aggregate Leaves
	leaveRows, err := r.db.WithContext(ctx).Table("leaves").
		Select("status, total_days, leave_type").
		Where("employee_id = ? AND start_date >= ?", userID, startDateStr).
		Rows()

	if err == nil {
		defer leaveRows.Close()
		var status string
		var totalDays int
		var leaveType string
		for leaveRows.Next() {
			leaveRows.Scan(&status, &totalDays, &leaveType)
			summary.Leave.TotalLeavesInPeriod += totalDays
			if status == "approved" {
				summary.Leave.Approved++
			} else if status == "pending" {
				summary.Leave.Pending++
			} else if status == "rejected" {
				summary.Leave.Rejected++
			}
			summary.Leave.LeavesByType[leaveType]++
		}
	}

	// 4. Aggregate Tasks
	type TaskStats struct {
		Status string
		Count  int
	}
	// For Tasks, they might be in `tasks` table or similar, assuming standard structure.
	// We check `tasks` table. But actually the `tasks` table is usually linked via `task_assignees` or `assigned_to`
	taskRows, err := r.db.WithContext(ctx).Table("tasks").
		Select("status, count(id) as count").
		Where("assigned_to = ? AND created_at >= ?", userID, startDate).
		Group("status").
		Rows()
	if err == nil {
		defer taskRows.Close()
		var status string
		var count int
		for taskRows.Next() {
			taskRows.Scan(&status, &count)
			summary.Tasks.TotalAssigned += count
			if status == "completed" {
				summary.Tasks.Completed += count
			} else if status == "in_progress" || status == "todo" {
				summary.Tasks.InProgress += count
			}
		}

		// Get overdue tasks
		var overdueCount int
		r.db.WithContext(ctx).Table("tasks").
			Where("assigned_to = ? AND due_date < ? AND status != 'completed'", userID, now).
			Select("count(id)").
			Scan(&overdueCount)

		summary.Tasks.Overdue = overdueCount
		if summary.Tasks.TotalAssigned > 0 {
			summary.Tasks.CompletionRate = float64(summary.Tasks.Completed) / float64(summary.Tasks.TotalAssigned) * 100
		}
	}

	return summary, nil
}

func (r *gormRepository) GetOrgEmployeeMonitor(ctx context.Context, orgID uint) ([]analytics.EmployeeMonitorResponse, error) {
	var modelsList []models.Employee
	if err := r.db.WithContext(ctx).Preload("User").Where("organisation_id = ?", orgID).Find(&modelsList).Error; err != nil {
		return nil, err
	}

	var employeesList []employees.Employee
	for _, m := range modelsList {
		employeesList = append(employeesList, toEmployeeDomain(&m))
	}

	var responses []analytics.EmployeeMonitorResponse

	now := time.Now()
	startDateStr := now.AddDate(0, -1, 0).Format("2006-01-02")
	totalWorkingDays := int(now.Sub(now.AddDate(0, -1, 0)).Hours() / 24)

	for _, emp := range employeesList {
		snap := analytics.EmployeeMonitorSnapshot{
			LastSeen: &emp.LastLoginAt,
		}

		// Calculate Attendance Rate this month
		var presentDays int
		r.db.WithContext(ctx).Table("attendances").
			Where("user_id = ? AND date >= ?", emp.UserID, startDateStr).
			Select("COUNT(id)").
			Scan(&presentDays)
		if totalWorkingDays > 0 {
			snap.AttendanceRateThisMonth = float64(presentDays) / float64(totalWorkingDays) * 100
		}

		// Calculate Pending leaves
		r.db.WithContext(ctx).Table("leaves").
			Where("employee_id = ? AND status = 'pending'", emp.UserID).
			Select("COUNT(id)").
			Scan(&snap.LeavesPending)

		// Calculate Overdue tasks
		r.db.WithContext(ctx).Table("tasks").
			Where("assigned_to = ? AND due_date < ? AND status != 'completed'", emp.UserID, now).
			Select("COUNT(id)").
			Scan(&snap.TasksOverdue)

		// Determine status
		snap.Status = "inactive"
		if !emp.LastLoginAt.IsZero() {
			daysSinceLogin := now.Sub(emp.LastLoginAt).Hours() / 24
			if daysSinceLogin <= 7 {
				snap.Status = "active"
			} else if daysSinceLogin <= 30 {
				snap.Status = "idle"
			}
		}

		responses = append(responses, analytics.EmployeeMonitorResponse{
			EmployeeInfo: emp,
			Snapshot:     snap,
		})
	}

	return responses, nil
}

func toEmployeeDomain(m *models.Employee) employees.Employee {
	return employees.Employee{
		UserID: m.UserID,
		User: authDomain.User{
			ID:          m.User.ID,
			FirstName:   m.User.FirstName,
			LastName:    m.User.LastName,
			Email:       m.User.Email,
			PhoneNumber: m.User.PhoneNumber,
		},
		EmploymentID:   m.EmploymentID,
		OrganisationID: m.OrganisationID,
		Organisation: orgDomain.Entity{
			ID:          m.Organisation.ID,
			CompanyName: m.Organisation.CompanyName,
		},
		Email:       m.Email,
		Type:        m.Type,
		Salary:      m.Salary,
		LastLoginAt: m.LastLoginAt,
	}
}
