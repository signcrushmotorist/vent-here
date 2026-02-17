package user

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"github.com/signcrushmotorist/vent-here/internal/infrastructure/persistence"
)

type UserService struct {
	Repo *persistence.UserRepository
}

func NewUserService(repo *persistence.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func GenerateAlias() (string, error) {
	var b [4]byte
	_, err := rand.Read(b[:])
	if err != nil {
		return "", err
	}
	num := binary.BigEndian.Uint32(b[:]) % 1000000
	return fmt.Sprintf("anon_%06d", num), nil
}

// ChangeAlias enforces 7-day cooldown
func (s *UserService) ChangeAlias(userID int, newAlias string) error {
	user, err := s.Repo.FindByID(userID)
	if err != nil {
		return err
	}

	if user.LastAliasChange != nil {
		next := user.LastAliasChange.Add(7 * 24 * time.Hour)
		if time.Now().Before(next) {
			return errors.New("alias can only be changed once every 7 days")
		}
	}

	existing, _ := s.Repo.FindByAlias(newAlias)
	if existing != nil {
		return errors.New("alias already taken")
	}

	now := time.Now()
	user.PublicAlias = newAlias
	user.AliasChanges++
	user.LastAliasChange = &now

	return s.Repo.Update(user)
}
