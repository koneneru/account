package repository

import (
	"account/internal/model"
	"account/internal/repository/mapper"
	repomodel "account/internal/repository/model"
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

func NewRepository(db *gorm.DB, logger *zerolog.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

func (r *Repository) CreateUser(ctx context.Context, user model.User) error {
	repoUser := mapper.UserToRepoUser(user)
	res := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(&repoUser)
	if res.Error != nil {
		r.logger.Err(res.Error).Msg("failed to save user")
		return fmt.Errorf("failed to save user: %w", res.Error)
	}
	return nil
}

func (r *Repository) GetUser(ctx context.Context, userID uint64) (model.User, error) {
	var user repomodel.User
	res := r.db.WithContext(ctx).
		Model(&repomodel.User{}).
		Where("id = ?", userID).
		First(&user)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return model.User{}, fmt.Errorf("user not found")
	} else if res.Error != nil {
		r.logger.Err(res.Error).Msg(("failed to get user id"))
	}
	return mapper.RepoUserToUser(user), nil
}

func (r *Repository) GetUsers(ctx context.Context, limit int, offset int) ([]model.User, error) {
	var users []repomodel.User
	res := r.db.WithContext(ctx).
		Model(&repomodel.User{}).
		Offset(offset).
		Limit(limit).
		Find(&users)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("user not fount")
	} else if res.Error != nil {
		r.logger.Err(res.Error).Msg("failed to load users")
	}
	return mapper.RepoUsersToUsers(users), nil
}

func (r *Repository) DeleteUser(ctx context.Context, userID uint64) error {
	res := r.db.WithContext(ctx).
		Where("id = ?", userID).
		Delete(&repomodel.User{})
	if res.Error != nil {
		r.logger.Err(res.Error).Msg("failed to delete user")
		return fmt.Errorf("failed to delete user")
	}
	return nil
}

func (r *Repository) UpdateUser(ctx context.Context, userID uint64, user model.UpdateUser) error {
	res := r.db.WithContext(ctx).Where("id = ?", userID).Updates(user)
	if res.Error != nil {
		r.logger.Err(res.Error).Msg("failed to update user")
		return fmt.Errorf("failed to update user")
	}
	return nil
}
