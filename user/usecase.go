package user

import (
	"GoMastersTest/models/DTOs"
	"GoMastersTest/models/entity"
	"context"
	"github.com/google/uuid"
)

type UseCase interface {
	GetAllUsers(ctx context.Context) ([]*entity.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	Create(ctx context.Context, user *DTOs.User) (id uuid.UUID, err error)
	Update(ctx context.Context, id uuid.UUID, user *DTOs.User) (*entity.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
