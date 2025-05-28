package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ireuven89/auctions/auth-service/key"
	"github.com/ireuven89/auctions/auth-service/user"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type UserDB struct {
	id   string `db:"id"`
	name string `db:"name"`

	//only for persistence and verify password on login
	password string `db:"password"` // <-- do NOT include in public struct or JSON output
	// <-- do NOT include in public struct or JSON output
	email string `db:"email"`
}

func toUser(userDB UserDB) *user.User {

	return &user.User{
		ID:       userDB.id,
		Name:     userDB.name,
		Email:    userDB.email,
		Password: userDB.password,
	}
}

var refresh = "refresh:%s"
var refreshRate = "refresh:rate:%s"
var MaxRefreshRate = 3
var refreshRateTtl = 15 * time.Minute

type Repository interface {
	CreateUser(ctx context.Context, user user.User) error
	FindUser(ctx context.Context, id string) (*user.User, error)
	FindUserByCredentials(ctx context.Context, identifier string) (*user.User, error)
	SaveRefreshToken(ctx context.Context, token string, userInfo string, ttl time.Duration) error
	GetToken(ctx context.Context, token string) (string, error)
	GetRefreshRate(ctx context.Context, token string) (int, error)
	DeleteUser(ctx context.Context, id string) error
}

type UserRepo struct {
	db     *sql.DB
	logger *zap.Logger
	redis  *redis.Client
}

func New(logger *zap.Logger, db *sql.DB, redisDB *redis.Client) Repository {

	return &UserRepo{
		db:     db,
		logger: logger,
		redis:  redisDB,
	}
}

func (r *UserRepo) CreateUser(ctx context.Context, user user.User) error {
	_, err := r.db.ExecContext(ctx, "insert into users (id, name, password, email) values(?, ?, ?, ?)", user.ID, user.Name, user.Password, user.Email)

	if err != nil {
		r.logger.Error("UserRepo.CreateUser", zap.Error(err))
		return fmt.Errorf("UserRepo.CreateUser %w", err)
	}

	return nil
}

func (r *UserRepo) FindUser(ctx context.Context, id string) (*user.User, error) {
	var userDB UserDB
	row := r.db.QueryRowContext(ctx, "select id, name, email from users where id = ?", id)

	if row.Err() != nil {
		r.logger.Error("UserRepo.FindUser", zap.Error(row.Err()))

		return nil, fmt.Errorf("UserRepo.FindUser failed fetching user %w", row.Err())
	}

	if err := row.Scan(&userDB.id, &userDB.name, &userDB.email); err != nil {
		return nil, fmt.Errorf("UserRepo.FindUser failed fetching user %w", err)
	}

	return toUser(userDB), nil
}

func (r *UserRepo) SaveRefreshToken(ctx context.Context, token string, userInfo string, ttl time.Duration) error {

	//set refresh
	statusCmd := r.redis.HSet(ctx,
		fmt.Sprintf(refresh, token),
		map[string]interface{}{
			"user_info":    userInfo,
			"refresh_rate": MaxRefreshRate,
		})

	if statusCmd.Err() != nil {
		return fmt.Errorf("failed inserting to redis %w", statusCmd.Err())
	}

	r.redis.Expire(ctx, fmt.Sprintf(refresh, token), ttl)

	_, err := statusCmd.Result()
	if err != nil {
		return fmt.Errorf("SaveRefreshToken failed saving %w", err)
	}

	return nil
}

func (r *UserRepo) GetToken(ctx context.Context, token string) (string, error) {
	val, err := r.redis.Get(ctx, fmt.Sprintf(refresh, token)).Result()

	if err != nil {
		if err == redis.Nil {
			return "", key.ErrExpiredToken
		}

		return "", fmt.Errorf("UserRepo.GetToken")
	}

	refreshRateVal, err := r.redis.Get(ctx, fmt.Sprintf(refreshRate, token)).Int()

	if err != nil {
		return "", fmt.Errorf("UserRepo.GetToken failed setting refresh rate %w", err)
	}

	if err = r.redis.Set(ctx, fmt.Sprintf(refreshRate, token), refreshRateVal-1, refreshRateTtl).Err(); err != nil {
		return "", fmt.Errorf("UserRepo.GetToken failed setting refresh rate %w", err)
	}

	return val, nil
}

func (r *UserRepo) GetRefreshRate(ctx context.Context, token string) (int, error) {
	val, err := r.redis.Get(ctx, fmt.Sprintf(refreshRate, token)).Int()

	if err != nil {
		if err == redis.Nil {
			return -1, key.ErrExpiredToken
		}

		return -1, fmt.Errorf("UserRepo.GetRefreshRate %w", err)
	}

	return val, nil
}

func (r *UserRepo) SaveAccessToken(ctx context.Context, userId, token string, ttl time.Duration) error {

	statusCmd := r.redis.Set(ctx, userId, token, ttl)

	if statusCmd.Err() != nil {
		return fmt.Errorf("UserRepo.SaveAccessToken failed inserting token %w", statusCmd.Err())
	}

	return nil
}

func (r *UserRepo) FindUserByCredentials(ctx context.Context, identifier string) (*user.User, error) {
	var userDB UserDB
	row := r.db.QueryRowContext(ctx, "SELECT id, name, email, password FROM users WHERE name = ? OR email = ?", identifier, identifier)

	if row.Err() != nil {
		return nil, fmt.Errorf("failed fetching user %w", row.Err())
	}

	if err := row.Scan(&userDB.id, &userDB.name, &userDB.email, &userDB.password); err != nil {
		return nil, fmt.Errorf("failed scan user result %w", err)
	}

	userResult := toUser(userDB)

	return userResult, nil
}

func (r *UserRepo) DeleteUser(ctx context.Context, id string) error {
	q := "delete from users where id = ?"

	if _, err := r.db.ExecContext(ctx, q, id); err != nil {
		return fmt.Errorf("UserRepo.DeleteUser failed deleting user %w", err)
	}

	return nil
}
