package settings

import (
	"server/internal/domain/settings"
	"server/internal/platform/miscallenous"
)

type useCase struct {
	repo settings.Repository
}

func NewUseCase(repo settings.Repository) settings.UseCase {
	return &useCase{repo: repo}
}

func (u *useCase) GetOrgSettings(orgID uint, userRole string) (*settings.OrgSettings, error) {
	if userRole != "boss" && userRole != "manager" {
		return nil, settings.ErrNotAuthorized
	}
	return u.repo.GetOrgSettings(orgID)
}

func (u *useCase) UpdateOrgSettings(orgID uint, userRole string, s *settings.OrgSettings) error {
	if userRole != "boss" {
		return settings.ErrNotAuthorized
	}
	return u.repo.UpdateOrgSettings(orgID, s)
}

func (u *useCase) DeleteOrganisation(orgID uint, userRole string, confirmationCode string) error {
	if userRole != "boss" {
		return settings.ErrNotAuthorized
	}
	org, err := u.repo.GetOrgSettings(orgID)
	if err != nil {
		return err
	}
	if confirmationCode != org.CompanyCode {
		return settings.ErrInvalidInput
	}
	return u.repo.DeleteOrganisation(orgID)
}

func (u *useCase) GetPlanInfo(orgID uint, userRole string) (*settings.PlanInfo, error) {
	if userRole != "boss" && userRole != "manager" {
		return nil, settings.ErrNotAuthorized
	}
	return u.repo.GetPlanInfo(orgID)
}

func (u *useCase) GetRolePermissions(orgID uint, userRole string) ([]settings.RolePermissionsSet, error) {
	if userRole != "boss" {
		return nil, settings.ErrNotAuthorized
	}
	return u.repo.GetRolePermissions(orgID)
}

func (u *useCase) UpdateRolePermissions(orgID uint, userRole string, perms []settings.RolePermissionsSet) error {
	if userRole != "boss" {
		return settings.ErrNotAuthorized
	}
	return u.repo.UpdateRolePermissions(orgID, perms)
}

func (u *useCase) GetAttendancePolicy(orgID uint, userRole string) (*settings.AttendancePolicy, error) {
	if userRole != "boss" && userRole != "manager" {
		return nil, settings.ErrNotAuthorized
	}
	return u.repo.GetAttendancePolicy(orgID)
}

func (u *useCase) UpdateAttendancePolicy(orgID uint, userRole string, policy *settings.AttendancePolicy) error {
	if userRole != "boss" {
		return settings.ErrNotAuthorized
	}
	return u.repo.UpdateAttendancePolicy(orgID, policy)
}

func (u *useCase) GetLeavePolicies(orgID uint, userRole string) ([]settings.LeavePolicyEntry, error) {
	if userRole != "boss" && userRole != "manager" {
		return nil, settings.ErrNotAuthorized
	}
	return u.repo.GetLeavePolicies(orgID)
}

func (u *useCase) SaveLeavePolicy(orgID uint, userRole string, policy *settings.LeavePolicyEntry) error {
	if userRole != "boss" {
		return settings.ErrNotAuthorized
	}
	return u.repo.SaveLeavePolicy(orgID, policy)
}

func (u *useCase) DeleteLeavePolicy(orgID uint, userRole string, policyID uint) error {
	if userRole != "boss" {
		return settings.ErrNotAuthorized
	}
	return u.repo.DeleteLeavePolicy(orgID, policyID)
}

func (u *useCase) UpdateLeavePolicies(orgID uint, userRole string, policies []settings.LeavePolicyEntry) error {
	if userRole != "boss" {
		return settings.ErrNotAuthorized
	}
	return u.repo.UpdateLeavePolicies(orgID, policies)
}

func (u *useCase) GetTaskTypes(orgID uint, userRole string) ([]settings.TaskTypeEntry, error) {
	return u.repo.GetTaskTypes(orgID)
}

func (u *useCase) SaveTaskType(orgID uint, userRole string, taskType *settings.TaskTypeEntry) error {
	if userRole != "boss" && userRole != "manager" {
		return settings.ErrNotAuthorized
	}
	return u.repo.SaveTaskType(orgID, taskType)
}

func (u *useCase) DeleteTaskType(orgID uint, userRole string, typeID uint) error {
	if userRole != "boss" && userRole != "manager" {
		return settings.ErrNotAuthorized
	}
	return u.repo.DeleteTaskType(orgID, typeID)
}

func (u *useCase) GetUserProfile(userID uint) (*settings.UserProfile, error) {
	return u.repo.GetUserProfile(userID)
}

func (u *useCase) UpdateUserProfile(userID uint, profile *settings.UserProfile) error {
	return u.repo.UpdateUserProfile(userID, profile)
}

func (u *useCase) UpdatePassword(userID uint, currentPwd, newPwd string) error {
	hashed, err := u.repo.GetHashedPassword(userID)
	if err != nil {
		return err
	}
	if !miscallenous.VerifyPassword(hashed, currentPwd) {
		return settings.ErrInvalidInput
	}
	newHash, err := miscallenous.HashPassword(newPwd)
	if err != nil {
		return err
	}
	return u.repo.UpdatePassword(userID, newHash)
}

func (u *useCase) GetNotificationPreferences(userID uint) ([]settings.NotificationPreferenceEntry, error) {
	return u.repo.GetNotificationPreferences(userID)
}

func (u *useCase) UpdateNotificationPreferences(userID uint, prefs []settings.NotificationPreferenceEntry) error {
	return u.repo.UpdateNotificationPreferences(userID, prefs)
}

func (u *useCase) ExportOrgData(orgID uint, userRole string) (string, error) {
	if userRole != "boss" {
		return "", settings.ErrNotAuthorized
	}
	return u.repo.ExportOrgData(orgID)
}