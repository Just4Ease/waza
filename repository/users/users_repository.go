package usersrepository

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"time"
	"waza/models"
	"waza/repository"
	"waza/utils"
)

const tableSetup = `
CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY UNIQUE,
		firstName TEXT,
		lastName TEXT,
		phone TEXT,
		email TEXT,
		timeCreated DATETIME,
		timeUpdated DATETIME
	)
`

type userRepo struct {
	dataStore *sql.DB
}

func (u *userRepo) CreateUser(ctx context.Context, payload *models.User) (*models.User, error) {

	// check duplicate user by phone
	if userByPhone, _ := u.GetUserByPhone(ctx, payload.Phone); userByPhone != nil {
		return nil, repository.ErrDuplicateFound
	}

	// check duplicate user by email if available
	if payload.Email != nil {
		userByEmail, _ := u.GetUserByEmail(ctx, *payload.Email)
		if userByEmail != nil {
			return nil, repository.ErrDuplicateFound
		}
	}

	now := time.Now()
	payload.TimeCreated = now
	payload.TimeUpdated = now
	payload.Id = utils.GenerateId()

	_txt := `
INSERT INTO 
    users (id, firstName, lastName, phone, email, timeCreated, timeUpdated)
    VALUES (?, ?, ?, ?, ?, ?, ?)
`

	statement, err := u.dataStore.PrepareContext(ctx, _txt)
	if err != nil {
		return nil, err
	}

	_, err = statement.ExecContext(ctx,
		payload.Id,
		payload.FirstName,
		payload.LastName,
		payload.Phone,
		payload.Email,
		payload.TimeCreated,
		payload.TimeUpdated,
	)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (u *userRepo) GetUserById(ctx context.Context, id string) (*models.User, error) {
	row := u.dataStore.QueryRowContext(ctx, "SELECT * from users WHERE id = ?", id)
	return scanner(row)
}

func (u *userRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	row := u.dataStore.QueryRowContext(ctx, "SELECT * from users WHERE email = ? ", email)
	return scanner(row)
}

func (u *userRepo) GetUserByPhone(ctx context.Context, phone string) (*models.User, error) {
	row := u.dataStore.QueryRowContext(ctx, "SELECT * from users WHERE phone = ?", phone)
	return scanner(row)
}

func NewUserRepository(db *sql.DB) (repository.UserRepository, error) {

	// Setup table if not exists.
	if _, err := db.Exec(tableSetup); err != nil {
		return nil, err
	}

	return &userRepo{
		dataStore: db,
	}, nil
}

func scanner(row *sql.Row) (*models.User, error) {
	var user models.User
	if err := row.Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.Email,
		&user.TimeCreated,
		&user.TimeUpdated,
	); err != nil {
		return nil, err
	}
	return &user, nil
}
