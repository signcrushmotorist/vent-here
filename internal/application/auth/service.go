package auth

import (
	"errors"

	"github.com/google/uuid"
	"github.com/signcrushmotorist/vent-here/internal/application/user"
	"github.com/signcrushmotorist/vent-here/internal/domain"
	"github.com/signcrushmotorist/vent-here/internal/infrastructure/persistence"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Repo *persistence.UserRepository
}

func NewAuthServie(repo *persistence.UserRepository) *AuthService {
	return &AuthService{Repo: repo}
}

// Register user with auto-generated public alias
func (s *AuthService) Register(email, password, username string) (*domain.User, error) {
	if email == "" || password == "" || username == "" {
		return nil, errors.New("email, password, and username are required")
	}

	if existing, _ := s.Repo.FindByEmail(email); existing != nil {
		return nil, errors.New("email already in use")
	}

	if existing, _ := s.Repo.FindByUsername(username); existing != nil {
		return nil, errors.New("username already in use")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Generate unique public alias
	var alias string
	for {
		alias, err = user.GenerateAlias()
		if err != nil {
			return nil, err
		}
		existing, _ := s.Repo.FindByAlias(alias)
		if existing == nil {
			break
		}
	}

	newUser := &domain.User{
		PublicID:     uuid.New().String(),
		Email:        email,
		PasswordHash: string(hash),
		Username:     username,
		PublicAlias:  alias,
	}

	if err := s.Repo.Save(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

// Login returns user if credentials match
func (s *AuthService) Login(email, password string) (*domain.User, error) {
	user, err := s.Repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
