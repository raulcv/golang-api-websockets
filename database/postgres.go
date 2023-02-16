package database

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/raulcv/goapiws/models"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{db}, nil
}

func (repo *PostgresRepository) AddUser(ctx context.Context, user *models.User) error {
	_, err := repo.db.ExecContext(ctx, "insert into users (id, email, password) values ($1, $2, $3)", user.Id, user.Email, user.Password)
	return err
}

func (repo *PostgresRepository) GetUserById(ctx context.Context, id string) (*models.User, error) {
	userRows, err := repo.db.QueryContext(ctx, "select id, email from users where id = $1", id)
	defer func() {
		err = userRows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var user = models.User{}
	for userRows.Next() {
		if err = userRows.Scan(&user.Id, &user.Email); err == nil {
			return &user, nil
		}
	}
	if err = userRows.Err(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	userRows, err := repo.db.QueryContext(ctx, "select id, email, password from users where email = $1", email)
	defer func() {
		err = userRows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var user = models.User{}
	for userRows.Next() {
		if err = userRows.Scan(&user.Id, &user.Email, &user.Password); err == nil {
			return &user, nil
		}
	}
	if err = userRows.Err(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *PostgresRepository) Close() error {
	return repo.db.Close()
}
