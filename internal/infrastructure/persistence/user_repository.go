package persistence

import (
	"database/sql"

	"github.com/signcrushmotorist/vent-here/internal/domain"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// Save inserts a new user
func (r *UserRepository) Save(user *domain.User) error {
	query := `
	INSERT INTO users (public_id, email, password_hash, username, public_alias)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at`
	return r.DB.QueryRow(query, user.PublicID, user.Email, user.PasswordHash, user.Username, user.PublicAlias).
		Scan(&user.ID, &user.CreatedAt)
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, public_id, email, password_hash, username, public_alias, alias_changes, last_alias_change, created_at FROM users WHERE email=$1`
	err := r.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.PublicID,
		&user.Email,
		&user.PasswordHash,
		&user.Username,
		&user.PublicAlias,
		&user.AliasChanges,
		&user.LastAliasChange,
		&user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByUsername(username string) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, public_id, email, password_hash, username, public_alias, alias_changes, last_alias_change, created_at FROM users WHERE username=$1`
	err := r.DB.QueryRow(query, username).Scan(
		&user.ID,
		&user.PublicID,
		&user.Email,
		&user.PasswordHash,
		&user.Username,
		&user.PublicAlias,
		&user.AliasChanges,
		&user.LastAliasChange,
		&user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByAlias(alias string) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id FROM users WHERE public_alias=$1`
	err := r.DB.QueryRow(query, alias).Scan(&user.ID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByID(id int) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, public_id, email, password_hash, username, public_alias, alias_changes, last_alias_change, created_at FROM users WHERE id=$1`
	err := r.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.PublicID,
		&user.Email,
		&user.PasswordHash,
		&user.Username,
		&user.PublicAlias,
		&user.AliasChanges,
		&user.LastAliasChange,
		&user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Update(user *domain.User) error {
	query := `UPDATE users SET public_alias=$1, alias_changes=$2, last_alias_change=$3 WHERE id=$4`
	_, err := r.DB.Exec(query, user.PublicAlias, user.AliasChanges, user.LastAliasChange, user.ID)
	return err
}
