package calendar

import (
	"context"
	"fmt"
	"time"

	domain "server/internal/domain/calendar"
	"server/internal/platform/database/models"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) domain.Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) GetCalendarEvents(ctx context.Context, orgID uint, userID *uint, startDate time.Time, endDate time.Time) ([]domain.CalendarEvent, error) {
	var events []domain.CalendarEvent

	// 1. Fetch Tasks
	var tasks []models.Task
	taskTx := r.db.WithContext(ctx).Preload("AssignedUser").
		Where("organisation_id = ? AND due_date BETWEEN ? AND ?", orgID, startDate, endDate)
	if userID != nil {
		taskTx = taskTx.Where("assigned_to = ?", *userID)
	}
	if err := taskTx.Find(&tasks).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch tasks: %w", err)
	}
	for _, t := range tasks {
		assignedName := ""
		if t.AssignedUser != nil {
			assignedName = t.AssignedUser.FirstName + " " + t.AssignedUser.LastName
		}
		eventType := "task"
		if t.Type != "" {
			eventType = t.Type
		}
		events = append(events, domain.CalendarEvent{
			ID:             t.ID,
			Title:          t.Title,
			Description:    t.Description,
			StartTime:      t.DueDate,
			EndTime:        t.DueDate,
			EventType:      eventType,
			EntityID:       t.ID,
			Status:         t.Status,
			AssignedToName: assignedName,
		})
	}

	// 2. Fetch Milestones
	if userID == nil {
		var milestones []models.Milestone
		milestoneTx := r.db.WithContext(ctx).
			Where("organisation_id = ? AND due_date BETWEEN ? AND ?", orgID, startDate, endDate)
		if err := milestoneTx.Find(&milestones).Error; err != nil {
			return nil, fmt.Errorf("failed to fetch milestones: %w", err)
		}
		for _, m := range milestones {
			events = append(events, domain.CalendarEvent{
				ID:             m.ID,
				Title:          m.Title,
				Description:    m.Description,
				StartTime:      m.DueDate,
				EndTime:        m.DueDate,
				EventType:      "milestone",
				EntityID:       m.ID,
				Status:         m.Status,
				AssignedToName: "",
			})
		}
	}

	// 3. Fetch Leaves
	var leaves []models.Leave
	leaveTx := r.db.WithContext(ctx).Preload("Employee").
		Where("organisation_id = ? AND ((start_date BETWEEN ? AND ?) OR (end_date BETWEEN ? AND ?) OR (start_date <= ? AND end_date >= ?))", orgID, startDate, endDate, startDate, endDate, startDate, endDate)
	if userID != nil {
		leaveTx = leaveTx.Where("employee_id = ?", *userID)
	}
	if err := leaveTx.Find(&leaves).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch leaves: %w", err)
	}
	for _, l := range leaves {
		assignedName := l.Employee.FirstName + " " + l.Employee.LastName
		events = append(events, domain.CalendarEvent{
			ID:             l.ID,
			Title:          fmt.Sprintf("Leave: %s", l.LeaveType),
			Description:    l.Reason,
			StartTime:      l.StartDate,
			EndTime:        l.EndDate,
			EventType:      "leave",
			EntityID:       l.ID,
			Status:         l.Status,
			AssignedToName: assignedName,
		})
	}

	// 4. Fetch Attendance
	var attendances []models.Attendance
	attendanceTx := r.db.WithContext(ctx).Preload("User").
		Where("organisation_id = ? AND date BETWEEN ? AND ?", orgID, startDate, endDate)
	if userID != nil {
		attendanceTx = attendanceTx.Where("user_id = ?", *userID)
	}
	if err := attendanceTx.Find(&attendances).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch attendance: %w", err)
	}
	for _, a := range attendances {
		assignedName := a.User.FirstName + " " + a.User.LastName
		title := "Attendance: Present"
		status := "present"
		if a.CheckInTime == nil {
			title = "Attendance: Absent/Not Checked In"
			status = "absent"
		}
		
		startTime := a.Date
		if a.CheckInTime != nil {
			startTime = *a.CheckInTime
		}
		endTime := a.Date
		if a.CheckOutTime != nil {
			endTime = *a.CheckOutTime
		}

		events = append(events, domain.CalendarEvent{
			ID:             a.ID,
			Title:          title,
			Description:    fmt.Sprintf("Date: %s", a.Date.Format("2006-01-02")),
			StartTime:      startTime,
			EndTime:        endTime,
			EventType:      "attendance",
			EntityID:       a.ID,
			Status:         status,
			AssignedToName: assignedName,
		})
	}

	// 5. Fetch Project Deadlines
	var projects []models.Project
	projectTx := r.db.WithContext(ctx).
		Where("organisation_id = ? AND end_date BETWEEN ? AND ?", orgID, startDate, endDate)
	if err := projectTx.Find(&projects).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch projects: %w", err)
	}
	for _, p := range projects {
		events = append(events, domain.CalendarEvent{
			ID:             p.ID,
			Title:          fmt.Sprintf("Project Deadline: %s", p.Name),
			Description:    p.Description,
			StartTime:      p.EndDate,
			EndTime:        p.EndDate,
			EventType:      "project_deadline",
			EntityID:       p.ID,
			Status:         p.Status,
			AssignedToName: "",
		})
	}

	return events, nil
}
