package auth

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	domain "server/internal/domain/auth"
	empDomain "server/internal/domain/employees"
	"server/internal/platform/miscallenous"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	_ = os.Setenv("SECRET", "test-secret")
	os.Exit(m.Run())
}

type mockRepository struct {
	createFunc                         func(ctx context.Context, user *domain.User) error
	getByEmailFunc                     func(ctx context.Context, email string) (*domain.User, error)
	getByPhoneNumberFunc               func(ctx context.Context, phone string) (*domain.User, error)
	getByIDFunc                        func(ctx context.Context, id uint) (*domain.User, error)
	updateFunc                         func(ctx context.Context, user *domain.User) error
	saveRefreshTokenFunc               func(ctx context.Context, token *domain.RefreshToken) error
	getRefreshTokenFunc                func(ctx context.Context, tokenHash string) (*domain.RefreshToken, error)
	revokeRefreshTokenFunc             func(ctx context.Context, tokenHash string) error
	revokeAllUserRefreshTokensFunc     func(ctx context.Context, userID uint) error
	createEmailVerificationFunc        func(ctx context.Context, ev *domain.EmailVerification) error
	getEmailVerificationByTokenFunc    func(ctx context.Context, token string) (*domain.EmailVerification, error)
	getActiveEmailVerificationByUserID func(ctx context.Context, userID uint) (*domain.EmailVerification, error)
	updateEmailVerificationFunc        func(ctx context.Context, ev *domain.EmailVerification) error
}

func (m *mockRepository) Create(ctx context.Context, user *domain.User) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, user)
	}
	return nil
}

func (m *mockRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.getByEmailFunc != nil {
		return m.getByEmailFunc(ctx, email)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockRepository) GetByPhoneNumber(ctx context.Context, phone string) (*domain.User, error) {
	if m.getByPhoneNumberFunc != nil {
		return m.getByPhoneNumberFunc(ctx, phone)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockRepository) Update(ctx context.Context, user *domain.User) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, user)
	}
	return nil
}

func (m *mockRepository) SaveRefreshToken(ctx context.Context, token *domain.RefreshToken) error {
	if m.saveRefreshTokenFunc != nil {
		return m.saveRefreshTokenFunc(ctx, token)
	}
	return nil
}

func (m *mockRepository) GetRefreshToken(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	if m.getRefreshTokenFunc != nil {
		return m.getRefreshTokenFunc(ctx, tokenHash)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	if m.revokeRefreshTokenFunc != nil {
		return m.revokeRefreshTokenFunc(ctx, tokenHash)
	}
	return nil
}

func (m *mockRepository) RevokeAllUserRefreshTokens(ctx context.Context, userID uint) error {
	if m.revokeAllUserRefreshTokensFunc != nil {
		return m.revokeAllUserRefreshTokensFunc(ctx, userID)
	}
	return nil
}

func (m *mockRepository) CreateEmailVerification(ctx context.Context, ev *domain.EmailVerification) error {
	if m.createEmailVerificationFunc != nil {
		return m.createEmailVerificationFunc(ctx, ev)
	}
	return nil
}

func (m *mockRepository) GetEmailVerificationByToken(ctx context.Context, token string) (*domain.EmailVerification, error) {
	if m.getEmailVerificationByTokenFunc != nil {
		return m.getEmailVerificationByTokenFunc(ctx, token)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockRepository) GetActiveEmailVerificationByUserID(ctx context.Context, userID uint) (*domain.EmailVerification, error) {
	if m.getActiveEmailVerificationByUserID != nil {
		return m.getActiveEmailVerificationByUserID(ctx, userID)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockRepository) UpdateEmailVerification(ctx context.Context, ev *domain.EmailVerification) error {
	if m.updateEmailVerificationFunc != nil {
		return m.updateEmailVerificationFunc(ctx, ev)
	}
	return nil
}

type mockEmployeeRepository struct {
	getByUserIDFunc func(ctx context.Context, userID uint) (*empDomain.Employee, error)
}

func (m *mockEmployeeRepository) GetByUserID(ctx context.Context, userID uint) (*empDomain.Employee, error) {
	if m.getByUserIDFunc != nil {
		return m.getByUserIDFunc(ctx, userID)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockEmployeeRepository) Create(ctx context.Context, employee *empDomain.Employee) error {
	return nil
}

func (m *mockEmployeeRepository) GetEmployeesByOrg(ctx context.Context, orgID uint) ([]empDomain.Employee, error) {
	return nil, nil
}

func (m *mockEmployeeRepository) GetByID(ctx context.Context, id uint) (*empDomain.Employee, error) {
	return nil, gorm.ErrRecordNotFound
}

func (m *mockEmployeeRepository) Update(ctx context.Context, employee *empDomain.Employee) error {
	return nil
}

func (m *mockEmployeeRepository) Delete(ctx context.Context, id uint) error {
	return nil
}

func (m *mockEmployeeRepository) Search(ctx context.Context, query string) ([]empDomain.Employee, error) {
	return nil, nil
}

func (m *mockEmployeeRepository) CreateManager(ctx context.Context, manager *empDomain.Manager) error {
	return nil
}

func (m *mockEmployeeRepository) GetManagerByUserID(ctx context.Context, userID uint) (*empDomain.Manager, error) {
	return nil, gorm.ErrRecordNotFound
}

func (m *mockEmployeeRepository) CreateBoss(ctx context.Context, boss *empDomain.Boss) error {
	return nil
}

func (m *mockEmployeeRepository) GetBossByUserID(ctx context.Context, userID uint) (*empDomain.Boss, error) {
	return nil, gorm.ErrRecordNotFound
}

type mockMailer struct {
	sendMailFunc func(to, subject, htmlBody string) error
}

func (m *mockMailer) SendMail(to, subject, htmlBody string) error {
	if m.sendMailFunc != nil {
		return m.sendMailFunc(to, subject, htmlBody)
	}
	return nil
}

func TestSignUp_CreatesUserAndIssuesTokenPair(t *testing.T) {
	var createdUser *domain.User
	var savedRefreshToken *domain.RefreshToken

	repo := &mockRepository{
		getByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
			return nil, nil
		},
		createFunc: func(ctx context.Context, user *domain.User) error {
			createdUser = user
			user.ID = 42
			return nil
		},
		saveRefreshTokenFunc: func(ctx context.Context, token *domain.RefreshToken) error {
			savedRefreshToken = token
			return nil
		},
	}

	svc := NewService(repo, nil, nil)

	pair, err := svc.SignUp(context.Background(), &domain.User{
		FirstName: "Ada",
		LastName:  "Lovelace",
		Email:     "ada@example.com",
	}, "password123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if createdUser == nil {
		t.Fatal("expected create to be called")
	}
	if createdUser.HashedPassword == "" {
		t.Error("expected hashed password to be populated")
	}
	if createdUser.LastLoginAt.IsZero() {
		t.Error("expected last login to be set")
	}
	if pair.AccessToken == "" || pair.RefreshToken == "" {
		t.Fatal("expected token pair to be returned")
	}
	if savedRefreshToken == nil {
		t.Fatal("expected refresh token to be persisted")
	}
	if savedRefreshToken.UserID != 42 {
		t.Fatalf("expected refresh token for user 42, got %d", savedRefreshToken.UserID)
	}
	if savedRefreshToken.TokenHash != miscallenous.HashRefreshToken(pair.RefreshToken) {
		t.Error("expected stored refresh token hash to match the raw token")
	}
}

func TestSignUp_ExistingRegisteredUserReturnsErrRegistered(t *testing.T) {
	repo := &mockRepository{
		getByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
			return &domain.User{ID: 1, Email: email, HashedPassword: "already-set"}, nil
		},
	}

	svc := NewService(repo, nil, nil)

	_, err := svc.SignUp(context.Background(), &domain.User{Email: "ada@example.com"}, "password123")
	if err != gorm.ErrRegistered {
		t.Fatalf("expected gorm.ErrRegistered, got %v", err)
	}
}

func TestSignIn_ReturnsInvitePendingForUninitializedUser(t *testing.T) {
	repo := &mockRepository{
		getByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
			return &domain.User{ID: 1, Email: email, HashedPassword: ""}, nil
		},
	}

	svc := NewService(repo, nil, nil)

	_, err := svc.SignIn(context.Background(), "ada@example.com", "password123")
	if !errors.Is(err, ErrInvitePending) {
		t.Fatalf("expected ErrInvitePending, got %v", err)
	}
}

func TestSignIn_RejectsWrongPassword(t *testing.T) {
	repo := &mockRepository{
		getByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
			hashed, err := miscallenous.HashPassword("correct-password")
			if err != nil {
				t.Fatalf("failed to hash password: %v", err)
			}
			return &domain.User{ID: 1, Email: email, HashedPassword: hashed}, nil
		},
	}

	svc := NewService(repo, nil, nil)

	_, err := svc.SignIn(context.Background(), "ada@example.com", "wrong-password")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestSignIn_ReturnsTokenPairForValidCredentials(t *testing.T) {
	var savedRefreshToken *domain.RefreshToken

	hashedPassword, err := miscallenous.HashPassword("password123")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	repo := &mockRepository{
		getByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
			return &domain.User{
				ID:             7,
				FirstName:      "Ada",
				Email:          email,
				HashedPassword: hashedPassword,
			}, nil
		},
		saveRefreshTokenFunc: func(ctx context.Context, token *domain.RefreshToken) error {
			savedRefreshToken = token
			return nil
		},
	}

	svc := NewService(repo, nil, nil)

	pair, err := svc.SignIn(context.Background(), "ada@example.com", "password123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if pair.AccessToken == "" || pair.RefreshToken == "" {
		t.Fatal("expected token pair to be returned")
	}
	if savedRefreshToken == nil {
		t.Fatal("expected refresh token to be persisted")
	}
}

func TestCompleteSignIn_FallsBackToPhoneNumber(t *testing.T) {
	var updatedUser *domain.User

	repo := &mockRepository{
		getByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
		getByPhoneNumberFunc: func(ctx context.Context, phone string) (*domain.User, error) {
			return &domain.User{ID: 9, Email: "ada@example.com", PhoneNumber: phone}, nil
		},
		updateFunc: func(ctx context.Context, user *domain.User) error {
			updatedUser = user
			return nil
		},
		saveRefreshTokenFunc: func(ctx context.Context, token *domain.RefreshToken) error {
			return nil
		},
	}

	svc := NewService(repo, nil, nil)

	pair, err := svc.CompleteSignIn(context.Background(), "ada@example.com", "9999999999", "password123", "London", "UK", "profile.png")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if pair.AccessToken == "" || pair.RefreshToken == "" {
		t.Fatal("expected token pair to be returned")
	}
	if updatedUser == nil {
		t.Fatal("expected update to be called")
	}
	if !updatedUser.IsVerified {
		t.Error("expected user to be marked verified")
	}
	if updatedUser.City != "London" || updatedUser.Country != "UK" || updatedUser.ProfilePic != "profile.png" {
		t.Error("expected profile fields to be updated")
	}
	if updatedUser.HashedPassword == "" {
		t.Error("expected password to be hashed")
	}
}

func TestValidateEmployee_ReturnsErrorWhenUserIsNotEmployee(t *testing.T) {
	repo := &mockRepository{}
	empRepo := &mockEmployeeRepository{
		getByUserIDFunc: func(ctx context.Context, userID uint) (*empDomain.Employee, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	svc := NewService(repo, empRepo, nil)

	_, err := svc.ValidateEmployee(context.Background(), 1)
	if !errors.Is(err, ErrNotAnEmployee) {
		t.Fatalf("expected ErrNotAnEmployee, got %v", err)
	}
}

func TestValidateEmployee_ReturnsTokenForEmployee(t *testing.T) {
	empRepo := &mockEmployeeRepository{
		getByUserIDFunc: func(ctx context.Context, userID uint) (*empDomain.Employee, error) {
			return &empDomain.Employee{UserID: userID, Type: "employee"}, nil
		},
	}

	svc := NewService(&mockRepository{}, empRepo, nil)

	token, err := svc.ValidateEmployee(context.Background(), 11)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Fatal("expected a token to be returned")
	}

	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil || !parsed.Valid {
		t.Fatalf("expected token to be valid, parse error: %v", err)
	}
}

func TestRefreshToken_RejectsUnknownToken(t *testing.T) {
	repo := &mockRepository{
		getRefreshTokenFunc: func(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	svc := NewService(repo, nil, nil)

	_, err := svc.RefreshToken(context.Background(), "raw-refresh-token")
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken, got %v", err)
	}
}

func TestRefreshToken_RejectedWhenRevokedTokenIsReused(t *testing.T) {
	var revokedAllForUser uint

	repo := &mockRepository{
		getRefreshTokenFunc: func(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
			return &domain.RefreshToken{UserID: 88, TokenHash: tokenHash, Revoked: true}, nil
		},
		revokeAllUserRefreshTokensFunc: func(ctx context.Context, userID uint) error {
			revokedAllForUser = userID
			return nil
		},
	}

	svc := NewService(repo, nil, nil)

	_, err := svc.RefreshToken(context.Background(), "raw-refresh-token")
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken, got %v", err)
	}
	if revokedAllForUser != 88 {
		t.Fatalf("expected all tokens to be revoked for user 88, got %d", revokedAllForUser)
	}
}

func TestRefreshToken_RotatesValidToken(t *testing.T) {
	var revokedHash string
	var savedRefreshToken *domain.RefreshToken

	repo := &mockRepository{
		getRefreshTokenFunc: func(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
			return &domain.RefreshToken{
				UserID:    21,
				TokenHash: tokenHash,
				ExpiresAt: time.Now().Add(time.Hour).Unix(),
				Revoked:   false,
			}, nil
		},
		revokeRefreshTokenFunc: func(ctx context.Context, tokenHash string) error {
			revokedHash = tokenHash
			return nil
		},
		getByIDFunc: func(ctx context.Context, id uint) (*domain.User, error) {
			return &domain.User{ID: id, Email: "ada@example.com", FirstName: "Ada"}, nil
		},
		saveRefreshTokenFunc: func(ctx context.Context, token *domain.RefreshToken) error {
			savedRefreshToken = token
			return nil
		},
	}

	svc := NewService(repo, nil, nil)

	pair, err := svc.RefreshToken(context.Background(), "raw-refresh-token")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if pair.AccessToken == "" || pair.RefreshToken == "" {
		t.Fatal("expected token pair to be returned")
	}
	if revokedHash != miscallenous.HashRefreshToken("raw-refresh-token") {
		t.Fatal("expected old refresh token hash to be revoked")
	}
	if savedRefreshToken == nil {
		t.Fatal("expected a new refresh token to be saved")
	}
	if savedRefreshToken.UserID != 21 {
		t.Fatalf("expected new refresh token for user 21, got %d", savedRefreshToken.UserID)
	}
}

func TestLogout_RevokesRefreshToken(t *testing.T) {
	var revokedHash string

	repo := &mockRepository{
		revokeRefreshTokenFunc: func(ctx context.Context, tokenHash string) error {
			revokedHash = tokenHash
			return nil
		},
	}

	svc := NewService(repo, nil, nil)

	if err := svc.Logout(context.Background(), "raw-refresh-token"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if revokedHash != miscallenous.HashRefreshToken("raw-refresh-token") {
		t.Fatalf("expected revoked hash to match raw token hash, got %s", revokedHash)
	}
}

func TestSendVerificationEmail_CreatesTokenAndSendsMail(t *testing.T) {
	t.Setenv("FRONTEND_URL", "https://frontend.test")

	var createdEV *domain.EmailVerification
	var sentTo string
	var sentSubject string
	var sentBody string

	repo := &mockRepository{
		getByIDFunc: func(ctx context.Context, id uint) (*domain.User, error) {
			return &domain.User{ID: id, FirstName: "Ada", Email: "ada@example.com"}, nil
		},
		getActiveEmailVerificationByUserID: func(ctx context.Context, userID uint) (*domain.EmailVerification, error) {
			return nil, gorm.ErrRecordNotFound
		},
		createEmailVerificationFunc: func(ctx context.Context, ev *domain.EmailVerification) error {
			createdEV = ev
			return nil
		},
	}
	mailer := &mockMailer{
		sendMailFunc: func(to, subject, htmlBody string) error {
			sentTo = to
			sentSubject = subject
			sentBody = htmlBody
			return nil
		},
	}

	svc := NewService(repo, nil, mailer)

	if err := svc.SendVerificationEmail(context.Background(), 1); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if createdEV == nil {
		t.Fatal("expected verification record to be created")
	}
	if createdEV.Token == "" {
		t.Fatal("expected verification token to be generated")
	}
	if createdEV.ExpiresAt.Before(time.Now()) {
		t.Error("expected verification token to expire in the future")
	}
	if sentTo != "ada@example.com" || sentSubject != "Verify your email address" {
		t.Fatalf("unexpected mail metadata: to=%s subject=%s", sentTo, sentSubject)
	}
	if sentBody == "" || createdEV.Token == "" || !contains(sentBody, createdEV.Token) {
		t.Error("expected verification email body to contain the verification token")
	}
}

func TestSendVerificationEmail_UsesExistingActiveToken(t *testing.T) {
	t.Setenv("FRONTEND_URL", "https://frontend.test")

	var createdCalled bool
	var sentBody string

	repo := &mockRepository{
		getByIDFunc: func(ctx context.Context, id uint) (*domain.User, error) {
			return &domain.User{ID: id, FirstName: "Ada", Email: "ada@example.com"}, nil
		},
		getActiveEmailVerificationByUserID: func(ctx context.Context, userID uint) (*domain.EmailVerification, error) {
			return &domain.EmailVerification{UserID: userID, Token: "existing-token"}, nil
		},
		createEmailVerificationFunc: func(ctx context.Context, ev *domain.EmailVerification) error {
			createdCalled = true
			return nil
		},
	}
	mailer := &mockMailer{
		sendMailFunc: func(to, subject, htmlBody string) error {
			sentBody = htmlBody
			return nil
		},
	}

	svc := NewService(repo, nil, mailer)

	if err := svc.SendVerificationEmail(context.Background(), 1); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if createdCalled {
		t.Error("expected existing active token to be reused without creating a new record")
	}
	if sentBody == "" || !contains(sentBody, "existing-token") {
		t.Error("expected email body to contain the existing token")
	}
}

func TestVerifyEmail_Succeeds(t *testing.T) {
	now := time.Now()
	var updatedUser *domain.User
	var updatedEV *domain.EmailVerification

	repo := &mockRepository{
		getEmailVerificationByTokenFunc: func(ctx context.Context, token string) (*domain.EmailVerification, error) {
			return &domain.EmailVerification{UserID: 5, Token: token, ExpiresAt: now.Add(time.Hour)}, nil
		},
		getByIDFunc: func(ctx context.Context, id uint) (*domain.User, error) {
			return &domain.User{ID: id, Email: "ada@example.com"}, nil
		},
		updateFunc: func(ctx context.Context, user *domain.User) error {
			updatedUser = user
			return nil
		},
		updateEmailVerificationFunc: func(ctx context.Context, ev *domain.EmailVerification) error {
			updatedEV = ev
			return nil
		},
	}

	svc := NewService(repo, nil, nil)

	if err := svc.VerifyEmail(context.Background(), "verify-token"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updatedUser == nil || !updatedUser.IsVerified {
		t.Fatal("expected user to be marked verified")
	}
	if updatedEV == nil || updatedEV.UsedAt == nil {
		t.Fatal("expected verification token to be marked used")
	}
}

func TestGetMe_ReturnsTokenPair(t *testing.T) {
	var savedRefreshToken *domain.RefreshToken

	repo := &mockRepository{
		getByIDFunc: func(ctx context.Context, id uint) (*domain.User, error) {
			return &domain.User{ID: id, Email: "ada@example.com", FirstName: "Ada"}, nil
		},
		saveRefreshTokenFunc: func(ctx context.Context, token *domain.RefreshToken) error {
			savedRefreshToken = token
			return nil
		},
	}

	svc := NewService(repo, nil, nil)

	pair, err := svc.GetMe(context.Background(), 12)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if pair.AccessToken == "" || pair.RefreshToken == "" {
		t.Fatal("expected token pair to be returned")
	}
	if savedRefreshToken == nil || savedRefreshToken.UserID != 12 {
		t.Fatal("expected refresh token to be persisted for the requested user")
	}
}

func contains(haystack, needle string) bool {
	return len(needle) == 0 || (len(haystack) >= len(needle) && func() bool {
		for i := 0; i+len(needle) <= len(haystack); i++ {
			if haystack[i:i+len(needle)] == needle {
				return true
			}
		}
		return false
	}())
}
