package settings

type UseCase interface {
	GetOrgSettings(orgID uint, userRole string) (*OrgSettings, error)
	UpdateOrgSettings(orgID uint, userRole string, settings *OrgSettings) error
	DeleteOrganisation(orgID uint, userRole string, confirmationCode string) error

	GetPlanInfo(orgID uint, userRole string) (*PlanInfo, error)

	GetRolePermissions(orgID uint, userRole string) ([]RolePermissionsSet, error)
	UpdateRolePermissions(orgID uint, userRole string, perms []RolePermissionsSet) error

	GetAttendancePolicy(orgID uint, userRole string) (*AttendancePolicy, error)
	UpdateAttendancePolicy(orgID uint, userRole string, policy *AttendancePolicy) error

	GetLeavePolicies(orgID uint, userRole string) ([]LeavePolicyEntry, error)
	SaveLeavePolicy(orgID uint, userRole string, policy *LeavePolicyEntry) error
	DeleteLeavePolicy(orgID uint, userRole string, policyID uint) error
	UpdateLeavePolicies(orgID uint, userRole string, policies []LeavePolicyEntry) error

	GetTaskTypes(orgID uint, userRole string) ([]TaskTypeEntry, error)
	SaveTaskType(orgID uint, userRole string, taskType *TaskTypeEntry) error
	DeleteTaskType(orgID uint, userRole string, typeID uint) error

	GetUserProfile(userID uint) (*UserProfile, error)
	UpdateUserProfile(userID uint, profile *UserProfile) error
	UpdatePassword(userID uint, currentPwd, newPwd string) error

	GetNotificationPreferences(userID uint) ([]NotificationPreferenceEntry, error)
	UpdateNotificationPreferences(userID uint, prefs []NotificationPreferenceEntry) error
	ExportOrgData(orgID uint, userRole string) (string, error)
}