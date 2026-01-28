package repository

import (
	"errors"
	"fmt"
	"user-management-system/internal/models"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}
func (r *UserRepo) CreateUser(req *models.User) (*models.User, error) {
	err := r.db.Omit("updated_at", "updated_by").Create(req).Error
	if err != nil {
		return nil, err
	}
	return req, nil
}
func (r *UserRepo) GetUserByID(id int) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *UserRepo) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
func (r *UserRepo) UpdateUser(id int, req *models.User) (*models.User, error) {
	err := r.db.Omit("created_at", "created_by").Save(req).Error
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (r *UserRepo) GetUsersByDivision(divisionID int) ([]models.User, error) {
	var users []models.User
	result := r.db.Where("division_id = ?", divisionID).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r *UserRepo) GetUsersByBusinessUnit(businessUnitID int) ([]models.User, error) {
	var users []models.User
	result := r.db.Where("business_unit_id = ?", businessUnitID).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r *UserRepo) GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	err := r.db.Where("email=?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) UpdateUserPassword(email, hashedPassword string) error {
	return r.db.Model(&models.User{}).Where("email = ?", email).Update("password", hashedPassword).Error
}

func (r *UserRepo) CreatePasswordResetToken(token *models.PasswordResetToken) error {
	return r.db.Create(token).Error
}

func (r *UserRepo) GetResetToken(tokenString string) (*models.PasswordResetToken, error) {
	var token models.PasswordResetToken
	err := r.db.Where("token = ? AND used = false AND expires_at > NOW()", tokenString).First(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid or expired token")
		}
		return nil, err
	}

	return &token, nil
}

func (r *UserRepo) MarkTokenUsed(tokenString string) error {
	return r.db.Model(&models.PasswordResetToken{}).Where("token = ?", tokenString).Update("used", true).Error
}

func (r *UserRepo) ApproveUser(userID int, approverEmail string) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{"is_active": true, "approved_by": approverEmail}).Error
}

func (r *UserRepo) GetPendingUsers() ([]models.User, error) {
	var users []models.User

	err := r.db.Where("approved_by IS NULL").Find(&users).Error

	return users, err
}

func (r *UserRepo) GetUsersByRoleID(roleID int) ([]int, error) {
	var userIDs []int
	err := r.db.Model(&models.User{}).Where("JSON_CONTAINS(role_ids, ?)", fmt.Sprintf("[%d]", roleID)).Pluck("id", &userIDs).Error

	return userIDs, err
}

func (r *UserRepo) GetUserWithRoles(userID int) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Roles").Preload("Roles.Permissions").First(&user, userID).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}
