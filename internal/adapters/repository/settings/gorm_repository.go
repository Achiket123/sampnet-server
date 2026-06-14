package settings

import (
	"encoding/json"
	"server/internal/domain/settings"
	"server/internal/platform/database/models"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) settings.Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) GetOrgSettings(orgID uint) (*settings.OrgSettings, error) {
	var org models.Organisation
	if err := r.db.First(&org, orgID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, settings.ErrSettingsNotFound
		}
		return nil, err
	}

	return &settings.OrgSettings{
		ID:                 org.ID,
		CompanyName:        org.CompanyName,
		CompanyCode:        org.CompanyCode,
		PrimaryContactName: org.PrimaryContactName,
		PrimaryEmail:       org.PrimaryEmail,
		PhoneNumber:        org.PhoneNumber,
		OfficeAddress:      org.OfficeAddress,
		City:               org.City,
		State:              org.State,
		PostalCode:         org.PostalCode,
		Country:            org.Country,
		CompanyLogo:        nil, // Not implemented or handled as integer ID
		Industry:           org.Industry,
		BillingAddress:     org.BillingAddress,
		CompanySize:        org.CompanySize,
	}, nil
}

func (r *gormRepository) UpdateOrgSettings(orgID uint, s *settings.OrgSettings) error {
	var org models.Organisation
	if err := r.db.First(&org, orgID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return settings.ErrSettingsNotFound
		}
		return err
	}

	org.CompanyName = s.CompanyName
	org.OfficeAddress = s.OfficeAddress
	org.City = s.City
	org.State = s.State
	org.PostalCode = s.PostalCode
	org.Country = s.Country
	org.Industry = s.Industry
	org.BillingAddress = s.BillingAddress
	org.CompanySize = s.CompanySize
	org.PrimaryContactName = s.PrimaryContactName
	org.PhoneNumber = s.PhoneNumber

	return r.db.Save(&org).Error
}

func (r *gormRepository) DeleteOrganisation(orgID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.Organisation{}, orgID).Error; err != nil {
			return err
		}
		if err := tx.Where("organisation_id = ?", orgID).Delete(&models.Employee{}).Error; err != nil {
			return err
		}
		if err := tx.Where("organisation_id = ?", orgID).Delete(&models.Manager{}).Error; err != nil {
			return err
		}
		if err := tx.Where("organisation_id = ?", orgID).Delete(&models.Boss{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *gormRepository) GetPlanInfo(orgID uint) (*settings.PlanInfo, error) {
	var org models.Organisation
	if err := r.db.First(&org, orgID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, settings.ErrSettingsNotFound
		}
		return nil, err
	}

	planName := "Standard Tier"
	if org.PlanID == 2 {
		planName = "Enterprise Pro"
	}

	return &settings.PlanInfo{
		PlanID:        org.PlanID,
		PlanName:      planName,
		PlanStatus:    org.PlanStatus,
		MaxEmployees:  org.MaxEmployees,
		PlanStartDate: org.PlanStartDate,
		PlanEndDate:   org.PlanEndDate,
	}, nil
}

func (r *gormRepository) GetRolePermissions(orgID uint) ([]settings.RolePermissionsSet, error) {
	var list []models.RolePermissions
	if err := r.db.Where("organisation_id = ?", orgID).Find(&list).Error; err != nil {
		return nil, err
	}

	roleMap := make(map[string][]string)
	// Initialize default roles
	roleMap["employee"] = []string{}
	roleMap["manager"] = []string{}
	roleMap["boss"] = []string{}

	for _, item := range list {
		roleMap[item.Role] = append(roleMap[item.Role], item.Permission)
	}

	var result []settings.RolePermissionsSet
	for _, role := range []string{"employee", "manager", "boss"} {
		result = append(result, settings.RolePermissionsSet{
			Role:        role,
			Permissions: roleMap[role],
		})
	}

	return result, nil
}

func (r *gormRepository) UpdateRolePermissions(orgID uint, perms []settings.RolePermissionsSet) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, rp := range perms {
			if err := tx.Where("organisation_id = ? AND role = ?", orgID, rp.Role).Delete(&models.RolePermissions{}).Error; err != nil {
				return err
			}
			for _, perm := range rp.Permissions {
				newPerm := models.RolePermissions{
					OrganisationID: orgID,
					Role:           rp.Role,
					Permission:     perm,
				}
				if err := tx.Create(&newPerm).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (r *gormRepository) GetAttendancePolicy(orgID uint) (*settings.AttendancePolicy, error) {
	var m models.AttendancePolicy
	if err := r.db.Where("organisation_id = ?", orgID).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &settings.AttendancePolicy{
				OrganisationID:  orgID,
				CheckInTime:     "09:00",
				CheckOutTime:    "18:00",
				GracePeriodMins: 15,
			}, nil
		}
		return nil, err
	}

	return &settings.AttendancePolicy{
		ID:              m.ID,
		OrganisationID:  m.OrganisationID,
		CheckInTime:     m.CheckInTime,
		CheckOutTime:    m.CheckOutTime,
		GracePeriodMins: m.GracePeriodMins,
	}, nil
}

func (r *gormRepository) UpdateAttendancePolicy(orgID uint, policy *settings.AttendancePolicy) error {
	var m models.AttendancePolicy
	err := r.db.Where("organisation_id = ?", orgID).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			m = models.AttendancePolicy{
				OrganisationID:  orgID,
				CheckInTime:     policy.CheckInTime,
				CheckOutTime:    policy.CheckOutTime,
				GracePeriodMins: policy.GracePeriodMins,
			}
			return r.db.Create(&m).Error
		}
		return err
	}

	m.CheckInTime = policy.CheckInTime
	m.CheckOutTime = policy.CheckOutTime
	m.GracePeriodMins = policy.GracePeriodMins
	return r.db.Save(&m).Error
}

func (r *gormRepository) GetLeavePolicies(orgID uint) ([]settings.LeavePolicyEntry, error) {
	var list []models.LeavePolicy
	if err := r.db.Where("organisation_id = ?", orgID).Find(&list).Error; err != nil {
		return nil, err
	}

	result := make([]settings.LeavePolicyEntry, len(list))
	for i, m := range list {
		result[i] = settings.LeavePolicyEntry{
			ID:             m.ID,
			OrganisationID: m.OrganisationID,
			LeaveType:      m.LeaveType,
			MaxDays:        m.MaxDays,
			Description:    m.Description,
		}
	}
	return result, nil
}

func (r *gormRepository) SaveLeavePolicy(orgID uint, policy *settings.LeavePolicyEntry) error {
	var m models.LeavePolicy
	if policy.ID != 0 {
		if err := r.db.Where("id = ? AND organisation_id = ?", policy.ID, orgID).First(&m).Error; err != nil {
			return err
		}
		m.LeaveType = policy.LeaveType
		m.MaxDays = policy.MaxDays
		m.Description = policy.Description
		if err := r.db.Save(&m).Error; err != nil {
			return err
		}
	} else {
		m = models.LeavePolicy{
			OrganisationID: orgID,
			LeaveType:      policy.LeaveType,
			MaxDays:        policy.MaxDays,
			Description:    policy.Description,
		}
		if err := r.db.Create(&m).Error; err != nil {
			return err
		}
		policy.ID = m.ID
	}
	return nil
}

func (r *gormRepository) DeleteLeavePolicy(orgID uint, policyID uint) error {
	return r.db.Where("id = ? AND organisation_id = ?", policyID, orgID).Delete(&models.LeavePolicy{}).Error
}

func (r *gormRepository) UpdateLeavePolicies(orgID uint, policies []settings.LeavePolicyEntry) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Get current IDs
		var current []models.LeavePolicy
		if err := tx.Where("organisation_id = ?", orgID).Find(&current).Error; err != nil {
			return err
		}

		// Create a set of new IDs to keep
		keepIDs := make(map[uint]bool)
		for _, p := range policies {
			if p.ID != 0 {
				keepIDs[p.ID] = true
			}
		}

		// Delete ones not in the list
		for _, item := range current {
			if !keepIDs[item.ID] {
				if err := tx.Delete(&item).Error; err != nil {
					return err
				}
			}
		}

		// Save/Create new list
		for _, p := range policies {
			var m models.LeavePolicy
			if p.ID != 0 {
				if err := tx.Where("id = ? AND organisation_id = ?", p.ID, orgID).First(&m).Error; err == nil {
					m.LeaveType = p.LeaveType
					m.MaxDays = p.MaxDays
					m.Description = p.Description
					if err := tx.Save(&m).Error; err != nil {
						return err
					}
				}
			} else {
				m = models.LeavePolicy{
					OrganisationID: orgID,
					LeaveType:      p.LeaveType,
					MaxDays:        p.MaxDays,
					Description:    p.Description,
				}
				if err := tx.Create(&m).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (r *gormRepository) GetTaskTypes(orgID uint) ([]settings.TaskTypeEntry, error) {
	var list []models.TaskType
	if err := r.db.Where("organisation_id = ?", orgID).Find(&list).Error; err != nil {
		return nil, err
	}

	result := make([]settings.TaskTypeEntry, len(list))
	for i, m := range list {
		result[i] = settings.TaskTypeEntry{
			ID:             m.ID,
			OrganisationID: m.OrganisationID,
			Name:           m.Name,
			Description:    m.Description,
		}
	}
	return result, nil
}

func (r *gormRepository) SaveTaskType(orgID uint, taskType *settings.TaskTypeEntry) error {
	var m models.TaskType
	if taskType.ID != 0 {
		if err := r.db.Where("id = ? AND organisation_id = ?", taskType.ID, orgID).First(&m).Error; err != nil {
			return err
		}
		m.Name = taskType.Name
		m.Description = taskType.Description
		if err := r.db.Save(&m).Error; err != nil {
			return err
		}
	} else {
		m = models.TaskType{
			OrganisationID: orgID,
			Name:           taskType.Name,
			Description:    taskType.Description,
		}
		if err := r.db.Create(&m).Error; err != nil {
			return err
		}
		taskType.ID = m.ID
	}
	return nil
}

func (r *gormRepository) DeleteTaskType(orgID uint, typeID uint) error {
	return r.db.Where("id = ? AND organisation_id = ?", typeID, orgID).Delete(&models.TaskType{}).Error
}

func (r *gormRepository) GetUserProfile(userID uint) (*settings.UserProfile, error) {
	var m models.UserModel
	if err := r.db.First(&m, userID).Error; err != nil {
		return nil, err
	}
	return &settings.UserProfile{
		FirstName:   m.FirstName,
		LastName:    m.LastName,
		Email:       m.Email,
		PhoneNumber: m.PhoneNumber,
		City:        m.City,
		Country:     m.Country,
		DateOfBirth: m.DateOfBirth,
		ProfilePic:  m.ProfilePic,
	}, nil
}

func (r *gormRepository) UpdateUserProfile(userID uint, profile *settings.UserProfile) error {
	var m models.UserModel
	if err := r.db.First(&m, userID).Error; err != nil {
		return err
	}
	m.FirstName = profile.FirstName
	m.LastName = profile.LastName
	m.PhoneNumber = profile.PhoneNumber
	m.City = profile.City
	m.Country = profile.Country
	m.DateOfBirth = profile.DateOfBirth
	m.ProfilePic = profile.ProfilePic
	return r.db.Save(&m).Error
}

func (r *gormRepository) UpdatePassword(userID uint, hashedPwd string) error {
	return r.db.Model(&models.UserModel{}).Where("id = ?", userID).Update("hashed_password", hashedPwd).Error
}

func (r *gormRepository) GetHashedPassword(userID uint) (string, error) {
	var m models.UserModel
	if err := r.db.Select("hashed_password").First(&m, userID).Error; err != nil {
		return "", err
	}
	return m.HashedPassword, nil
}

func (r *gormRepository) GetNotificationPreferences(userID uint) ([]settings.NotificationPreferenceEntry, error) {
	var list []models.NotificationPreference
	if err := r.db.Where("user_id = ?", userID).Find(&list).Error; err != nil {
		return nil, err
	}

	categories := []string{"announcements", "tasks", "leaves", "chats"}
	prefMap := make(map[string]settings.NotificationPreferenceEntry)
	for _, item := range list {
		prefMap[item.Category] = settings.NotificationPreferenceEntry{
			Category: item.Category,
			Email:    item.Email,
			Push:     item.Push,
			InApp:    item.InApp,
		}
	}

	var result []settings.NotificationPreferenceEntry
	for _, cat := range categories {
		if entry, exists := prefMap[cat]; exists {
			result = append(result, entry)
		} else {
			result = append(result, settings.NotificationPreferenceEntry{
				Category: cat,
				Email:    true,
				Push:     true,
				InApp:    true,
			})
		}
	}
	return result, nil
}

func (r *gormRepository) UpdateNotificationPreferences(userID uint, prefs []settings.NotificationPreferenceEntry) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, p := range prefs {
			var m models.NotificationPreference
			err := tx.Where("user_id = ? AND category = ?", userID, p.Category).First(&m).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					m = models.NotificationPreference{
						UserID:   userID,
						Category: p.Category,
						Email:    p.Email,
						Push:     p.Push,
						InApp:    p.InApp,
					}
					if err := tx.Create(&m).Error; err != nil {
						return err
					}
					continue
				}
				return err
			}
			m.Email = p.Email
			m.Push = p.Push
			m.InApp = p.InApp
			if err := tx.Save(&m).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *gormRepository) ExportOrgData(orgID uint) (string, error) {
	var org models.Organisation
	if err := r.db.First(&org, orgID).Error; err != nil {
		return "", err
	}
	data, err := json.Marshal(org)
	if err != nil {
		return "", err
	}
	return string(data), nil
}