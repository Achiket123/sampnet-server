package settings

type Repository interface {
	GetOrgSettings(orgID uint) (*OrgSettings, error)
	UpdateOrgSettings(orgID uint, settings *OrgSettings) error
	DeleteOrganisation(orgID uint) error

	GetPlanInfo(orgID uint) (*PlanInfo, error)

	GetRolePermissions(orgID uint) ([]RolePermissionsSet, error)
	UpdateRolePermissions(orgID uint, perms []RolePermissionsSet) error

	GetAttendancePolicy(orgID uint) (*AttendancePolicy, error)
	UpdateAttendancePolicy(orgID uint, policy *AttendancePolicy) error

	GetLeavePolicies(orgID uint) ([]LeavePolicyEntry, error)
	SaveLeavePolicy(orgID uint, policy *LeavePolicyEntry) error
	DeleteLeavePolicy(orgID uint, policyID uint) error
	UpdateLeavePolicies(orgID uint, policies []LeavePolicyEntry) error

	GetTaskTypes(orgID uint) ([]TaskTypeEntry, error)
	SaveTaskType(orgID uint, taskType *TaskTypeEntry) error
	DeleteTaskType(orgID uint, typeID uint) error

	GetUserProfile(userID uint) (*UserProfile, error)
	GetHashedPassword(userID uint) (string, error)
	UpdateUserProfile(userID uint, profile *UserProfile) error
	UpdatePassword(userID uint, hashedPwd string) error

	GetNotificationPreferences(userID uint) ([]NotificationPreferenceEntry, error)
	UpdateNotificationPreferences(userID uint, prefs []NotificationPreferenceEntry) error
	ExportOrgData(orgID uint) (string, error)
}