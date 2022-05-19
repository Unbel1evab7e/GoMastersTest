package usecase

import (
	"GoMastersTest/models"
	"GoMastersTest/models/DTOs"
	"GoMastersTest/models/entity"
	"GoMastersTest/user"
	"context"
	"github.com/google/uuid"
	"time"
)

type userUseCases struct {
	userRepository user.Repository
	contextTimeout time.Duration
}

func NewUserUseCase(a user.Repository, timeout time.Duration) user.UseCase {
	return &userUseCases{
		userRepository: a,
		contextTimeout: timeout,
	}
}

func (a *userUseCases) GetAllUsers(c context.Context) ([]*entity.User, error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err := a.userRepository.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *userUseCases) GetByID(c context.Context, id uuid.UUID) (*entity.User, error) {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err := a.userRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *userUseCases) Update(c context.Context, id uuid.UUID, m *DTOs.User) (*entity.User, error) {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	return a.userRepository.Update(ctx, id, m)
}

func (a *userUseCases) Create(c context.Context, m *DTOs.User) (uuid.UUID, error) {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	id, err := a.userRepository.Create(ctx, m)
	if err != nil {
		return uuid.NullUUID{}.UUID, err
	}
	return id, nil
}

func (a *userUseCases) Delete(c context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedUser, err := a.userRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existedUser == nil {
		return models.ErrNotFound
	}
	return a.userRepository.Delete(ctx, id)
}
