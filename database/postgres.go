package database

import (
	"context"
	"database/sql"
	"fmt"
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

/*
* POSTS OPERATIONS
 */
func (repo *PostgresRepository) AddPost(ctx context.Context, post *models.Post) error {
	_, err := repo.db.ExecContext(ctx, "insert into posts (id, content, user_id) values ($1, $2, $3)",
		post.Id, post.Content, post.UserId)
	return err
}

func (repo *PostgresRepository) GetPostById(ctx context.Context, id string) (*models.Post, error) {
	sqlQuery := "select id, content, user_id, created_at," +
		"updated_at, deleted_at from posts where id = $1 and deleted_at is null"
		// sqlQuery := "select id, content, user_id, created_at from posts where id = $1"
	postRows, err := repo.db.QueryContext(ctx, sqlQuery, id)
	defer func() {
		err = postRows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var post = models.Post{}
	for postRows.Next() {
		if err = postRows.Scan(&post.Id, &post.Content, &post.UserId,
			&post.CreatedAt, &post.UpdatedAt, &post.DeletedAt); err == nil {
			return &post, nil
		}
		// fmt.Println("each -")
	}
	if err = postRows.Err(); err != nil {
		// log.Fatalln("LOGGG", postRows)
		return nil, err
	}
	return &post, nil
}

func (repo *PostgresRepository) UpdatePost(ctx context.Context, post *models.Post) error {
	_, err := repo.db.ExecContext(ctx,
		"update posts set content = $1, updated_at = $2 where id = $3 and user_id = $4", post.Content, post.UpdatedAt, post.Id, post.UserId)
	return err
}

func (repo *PostgresRepository) DeletePost(ctx context.Context, id string, userId string) error {
	sqlQuery := "update posts set deleted_at = now() where id = $1 and user_id = $2 and deleted_at is null"
	_, err := repo.db.ExecContext(ctx, sqlQuery, id, userId)
	return err
}

func (repo *PostgresRepository) ActivatePost(ctx context.Context, id string) error {
	sqlQuery := "update posts set updated_at = now(), deleted_at = null where id = $1 and deleted_at is not null"
	_, err := repo.db.ExecContext(ctx, sqlQuery, id)
	return err
}

func (repo *PostgresRepository) ListPost(ctx context.Context, page uint64) ([]*models.Post, error) {
	sqlQuery := "select id, content, user_id, created_at, " +
		"coalesce(updated_at,to_timestamp(0)) as updated_at, coalesce(deleted_at,to_timestamp(0)) as deleted_at from posts " +
		" limit $1 offset $2"
		// sqlQuery := "select id, content, user_id, created_at from posts where id = $1"
	postRows, err := repo.db.QueryContext(ctx, sqlQuery, 2, page*2)
	defer func() {
		err = postRows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	if err != nil {
		fmt.Print("err.Error(): ", err.Error())
		return nil, err
	}

	fmt.Println("postRows: ", postRows)
	var times = uint64(0)
	var posts = []*models.Post{}
	for postRows.Next() {
		var post = models.Post{}
		err = postRows.Scan(&post.Id, &post.Content, &post.UserId, &post.CreatedAt, &post.UpdatedAt, &post.DeletedAt)
		if err == nil {
			posts = append(posts, &post)
			times += 1
		} else {
			fmt.Println("error", times)
			fmt.Println(err.Error())
		}
	}
	fmt.Println("times: ", times)

	if err = postRows.Err(); err != nil {
		// log.Fatalln("LOGGG", postRows)
		return nil, err
	}
	return posts, nil
}

func (repo *PostgresRepository) Close() error {
	return repo.db.Close()
}
