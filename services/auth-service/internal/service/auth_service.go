package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

type AuthService struct {
	mu       sync.RWMutex
	db       *sql.DB
	sessions map[string]Session // Still in-memory for sessions
	config   struct {
		sessionTTL      time.Duration
		cleanupInterval time.Duration
	}
}

type Session struct {
	Token     string
	Login     string
	CreatedAt time.Time
	ExpiresAt time.Time
}

func NewAuthService(db *sql.DB, sessionTTL, cleanupInterval time.Duration) *AuthService {
	service := &AuthService{
		db:       db,
		sessions: make(map[string]Session),
	}
	service.config.sessionTTL = sessionTTL
	service.config.cleanupInterval = cleanupInterval

	go service.startCleanup()
	return service
}

func (s *AuthService) Register(ctx context.Context, login, password string) (string, error) {
	// Check if user exists
	var exists bool
	err := s.db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM obauth.registered_client WHERE login = $1)",
		login).Scan(&exists)

	if err != nil {
		return "", &AuthError{Message: "Database error"}
	}
	if exists {
		return "", &AuthError{Message: "User already exists"}
	}

	// Create user
	_, err = s.db.ExecContext(ctx,
		"INSERT INTO obauth.registered_client (login, password) VALUES ($1, $2)",
		login, password)

	if err != nil {
		return "", &AuthError{Message: "Failed to create user"}
	}

	// Generate token and store in memory
	token := s.generateToken()

	s.mu.Lock()
	s.sessions[token] = Session{
		Token:     token,
		Login:     login,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(s.config.sessionTTL),
	}
	s.mu.Unlock()

	return token, nil
}

func (s *AuthService) Login(ctx context.Context, login, password string) (string, error) {
	// Check credentials in database
	var dbPassword string
	err := s.db.QueryRowContext(ctx,
		"SELECT password FROM obauth.registered_client WHERE login = $1",
		login).Scan(&dbPassword)

	if err == sql.ErrNoRows {
		return "", &AuthError{Message: "Invalid credentials"}
	}
	if err != nil {
		return "", &AuthError{Message: "Database error"}
	}

	// Check password (plain text for simplicity)
	if dbPassword != password {
		return "", &AuthError{Message: "Invalid credentials"}
	}

	// Generate token and store in memory
	token := s.generateToken()

	s.mu.Lock()
	s.sessions[token] = Session{
		Token:     token,
		Login:     login,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(s.config.sessionTTL),
	}
	s.mu.Unlock()

	return token, nil
}

func (s *AuthService) Logout(ctx context.Context, login, token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Verify session exists and belongs to user
	if session, exists := s.sessions[token]; exists && session.Login == login {
		delete(s.sessions, token)
		return nil
	}

	return &AuthError{Message: "Session not found"}
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) bool {
	s.mu.RLock()
	session, exists := s.sessions[token]
	s.mu.RUnlock()

	if !exists {
		return false
	}

	return time.Now().Before(session.ExpiresAt)
}

func (s *AuthService) generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func (s *AuthService) startCleanup() {
	ticker := time.NewTicker(s.config.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		s.cleanupExpiredSessions()
	}
}

func (s *AuthService) cleanupExpiredSessions() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for token, session := range s.sessions {
		if now.After(session.ExpiresAt) {
			delete(s.sessions, token)
		}
	}
}

type AuthError struct {
	Message string
}

func (e *AuthError) Error() string {
	return e.Message
}
