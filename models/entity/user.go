package entity

import (
	"GoMastersTest/models/DTOs"
	_ "github.com/bxcodec/faker"
	"github.com/google/uuid"
	"time"
)

type User struct {
	DTOs.User
	Created time.Time
	ID      uuid.UUID `faker:"UUID"`
}
