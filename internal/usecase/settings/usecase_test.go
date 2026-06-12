package settings

import (
	"server/internal/domain/settings"
	"testing"
)

type mockRepository struct {
	settings.Repository
	getOrgSettingsFunc func(orgID uint) (*settings.OrgSettings, error)
	updatePasswordFunc func(userID uint, hashed string) error
	getHashedPwdFunc   func(userID uint) (string, error)
}

func (m *mockRepository) GetOrgSettings(orgID uint) (*settings.OrgSettings, error) {
	return m.getOrgSettingsFunc(orgID)
}

func (m *mockRepository) UpdatePassword(userID uint, hashed string) error {
	return m.updatePasswordFunc(userID, hashed)
}

func (m *mockRepository) GetHashedPassword(userID uint) (string, error) {
	return m.getHashedPwdFunc(userID)
}

func TestGetOrgSettings_AccessControl(t *testing.T) {
	repo := &mockRepository{
		getOrgSettingsFunc: func(orgID uint) (*settings.OrgSettings, error) {
			return &settings.OrgSettings{CompanyName: "Test"}, nil
		},
	}
	uc := NewUseCase(repo)

	// Boss can view
	_, err := uc.GetOrgSettings(1, "boss")
	if err != nil {
		t.Errorf("boss should be authorized to view org settings: %v", err)
	}

	// Manager can view
	_, err = uc.GetOrgSettings(1, "manager")
	if err != nil {
		t.Errorf("manager should be authorized to view org settings: %v", err)
	}

	// Employee cannot view
	_, err = uc.GetOrgSettings(1, "employee")
	if err != settings.ErrNotAuthorized {
		t.Errorf("employee should not be authorized to view org settings: %v", err)
	}
}
