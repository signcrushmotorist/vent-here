package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
)

type SessionStore struct {
	sessions map[string]int
	mu       sync.RWMutex
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]int),
	}
}

func (s *SessionStore) Create(userID int) string {
	b := make([]byte, 32)
	rand.Read(b)
	sessionID := hex.EncodeToString(b)

	s.mu.Lock()
	s.sessions[sessionID] = userID
	s.mu.Unlock()

	return sessionID
}

func (s *SessionStore) Get(sessionID string) (int, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	userID, ok := s.sessions[sessionID]
	return userID, ok
}
