package repository

import (
	"GoMastersTest/models"
	"GoMastersTest/models/DTOs"
	"GoMastersTest/models/entity"
	"GoMastersTest/user"
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"time"
)

type postgreUserRepository struct {
	Conn *sql.DB
}

func NewPostgreUserRepository(Conn *sql.DB) user.Repository {
	return &postgreUserRepository{Conn}
}

func (m *postgreUserRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*entity.User, error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	result := make([]*entity.User, 0)
	for rows.Next() {
		t := new(entity.User)
		err = rows.Scan(
			&t.ID,
			&t.Firstname,
			&t.Lastname,
			&t.Email,
			&t.Age,
			&t.Created,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *postgreUserRepository) GetAllUsers(ctx context.Context) (res []*entity.User, err error) {
	query := `select id, first_name, last_name, email, age, created
						from users`

	list, err := m.fetch(ctx, query)
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		res = list
	} else {
		return nil, models.ErrNotFound
	}

	return
}

func (m *postgreUserRepository) GetByID(ctx context.Context, id uuid.UUID) (res *entity.User, err error) {
	query := `select id, first_name, last_name, email, age, created
						from users where id = $1`

	list, err := m.fetch(ctx, query, id.String())
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return nil, models.ErrNotFound
	}

	return
}
func (m *postgreUserRepository) Create(ctx context.Context, user *DTOs.User) (id uuid.UUID, err error) {
	query := `INSERT INTO users (id, first_name, last_name, email, age, created) VALUES ($1, $2, $3, $4, $5, $6)`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return uuid.NullUUID{}.UUID, err
	}

	var usEn entity.User

	usEn.User = *user
	usEn.ID, err = uuid.NewUUID()
	usEn.Created = time.Now()
	if err != nil {
		return uuid.NullUUID{}.UUID, err
	}
	res, err := stmt.ExecContext(ctx, usEn.ID.String(), usEn.Firstname, usEn.Lastname, usEn.Email, usEn.Age, usEn.Created)
	if err != nil {
		return uuid.NullUUID{}.UUID, err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return uuid.NullUUID{}.UUID, err
	}

	if ra != 0 {
		return usEn.ID, nil
	} else {
		return uuid.NullUUID{}.UUID, models.ErrInternalServerError
	}

}
func (m *postgreUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, id.String())
	if err != nil {
		return err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 0 {
		return nil
	} else {
		return models.ErrInternalServerError
	}
}
func (m *postgreUserRepository) Update(ctx context.Context, id uuid.UUID, user *DTOs.User) (*entity.User, error) {
	query := `UPDATE users SET first_name=$1, last_name=$2, email=$3, age=$4 WHERE id = $5`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	res, err := stmt.ExecContext(ctx, user.Firstname, user.Lastname, user.Email, user.Age, id.String())
	if err != nil {
		return nil, err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affect != 1 {
		err = fmt.Errorf("Weird  Behaviour. Total Affected: %d", affect)

		return nil, err
	}
	var usEn entity.User

	usEn.User = *user
	return &usEn, nil
}
