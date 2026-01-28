package services

import (
	"errors"
	"fmt"
	"time"
	"user-management-system/internal/models"
	"user-management-system/internal/repository"
	"user-management-system/internal/shared"
)

type UserService struct {
	repo                     *repository.UserRepo
	tokenInvalidationService *TokenInvalidationService
	auditService             *AuditService
}

func NewUserService(repo *repository.UserRepo, tokenInvalidationService *TokenInvalidationService, auditService *AuditService) *UserService {
	return &UserService{repo: repo, tokenInvalidationService: tokenInvalidationService, auditService: auditService}
}

func (s *UserService) CreateUser(req *models.User) (*models.User, error) {
	if req == nil {
		return nil, fmt.Errorf("User request cannot be nil")
	}

	// Validate password strength
	if err := shared.ValidatePasswordStrength(req.Password); err != nil {
		return nil, fmt.Errorf("password validation failed: %w", err)
	}

	// Validate email format
	if err := shared.ValidateEmail(req.Email); err != nil {
		return nil, fmt.Errorf("email validation failed: %w", err)
	}

	// Validate username format
	if err := shared.ValidateUsername(req.UserName); err != nil {
		return nil, fmt.Errorf("username validation failed: %w", err)
	}

	hashedPassword, err := shared.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	req.Password = hashedPassword

	user, err := s.repo.CreateUser(req)
	if err != nil {
		return nil, err
	}
	return user, nil

}
func (s *UserService) GetUserByID(id int) (*models.User, error) {
	if id == 0 {
		return nil, fmt.Errorf("Invalid User ID")
	}
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	users, err := s.repo.GetAllUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) GetAllUsersPaginated(page, pageSize int) ([]models.User, int64, error) {
	return s.repo.GetAllUsersPaginated(page, pageSize)
}

func (s *UserService) UpdateUser(id int, req *models.User) (*models.User, error) {
	if id == 0 || req == nil {
		return nil, fmt.Errorf("Invalid input for updating User")
	}

	// Only hash password if it's being changed (non-empty)
	if req.Password != "" {
		// Validate password strength
		if err := shared.ValidatePasswordStrength(req.Password); err != nil {
			return nil, fmt.Errorf("password validation failed: %w", err)
		}

		hashedPassword, err := shared.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		req.Password = hashedPassword
	}

	// Validate email if provided
	if req.Email != "" {
		if err := shared.ValidateEmail(req.Email); err != nil {
			return nil, fmt.Errorf("email validation failed: %w", err)
		}
	}

	// Validate username if provided
	if req.UserName != "" {
		if err := shared.ValidateUsername(req.UserName); err != nil {
			return nil, fmt.Errorf("username validation failed: %w", err)
		}
	}

	user, err := s.repo.UpdateUser(id, req)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUsersByDivision(divisionID int) ([]models.User, error) {
	if divisionID == 0 || divisionID < 0 {
		return nil, fmt.Errorf("Invalid Division ID")
	}

	return s.repo.GetUsersByDivision(divisionID)
}

func (s *UserService) GetUsersByBusinessUnit(businessUnitID int) ([]models.User, error) {
	if businessUnitID == 0 || businessUnitID < 0 {
		return nil, fmt.Errorf("Invalid Business Unit ID")
	}
	return s.repo.GetUsersByBusinessUnit(businessUnitID)
}

func (s *UserService) ValidateLogin(email, password string) (*models.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if !shared.CheckPassword(password, user.Password) {
		return nil, errors.New("Invalid credentials")
	}

	if user.IsActive == nil || !*user.IsActive {
		return nil, errors.New("Account is not active")
	}

	return user, nil
}

func (s *UserService) ForgotPassword(email string) error {

	_, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	token, err := shared.GenerateResetToken()
	if err != nil {
		return err
	}

	resetToken := &models.PasswordResetToken{
		Email:     email,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Used:      false,
	}

	err = s.repo.CreatePasswordResetToken(resetToken)
	if err != nil {
		return err
	}

	if s.auditService != nil {
		s.auditService.LogChange("Password Reset", email, fmt.Sprintf("Password updated for email: %s", email), nil, nil)
	}

	emailService := NewEmailService()

	return emailService.SendPasswordResetLink(email, token)

}

func (s *UserService) ResetPassword(token, newPassword string) (int, error) {
	resetToken, err := s.repo.GetResetToken(token)

	if err != nil {
		return 0, errors.New("invalid or expired token")
	}

	if resetToken == nil {
		return 0, errors.New("invalid or expired token")
	}

	// Check token expiration
	if time.Now().After(resetToken.ExpiresAt) {
		return 0, errors.New("token has expired")
	}

	// Check if token already used
	if resetToken.Used {
		return 0, errors.New("token has already been used")
	}

	user, err := s.repo.GetUserByEmail(resetToken.Email)
	if err != nil {
		return 0, errors.New("user not found")
	}

	// Validate password strength
	if err := shared.ValidatePasswordStrength(newPassword); err != nil {
		return 0, fmt.Errorf("password validation failed: %w", err)
	}

	hashedPassword, err := shared.HashPassword(newPassword)

	if err != nil {
		return 0, err
	}

	err = s.repo.UpdateUserPassword(resetToken.Email, hashedPassword)
	if err != nil {
		return 0, err
	}

	err = s.repo.MarkTokenUsed(token)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}
