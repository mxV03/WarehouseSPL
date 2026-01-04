//go:build auth

package auth

import (
	"context"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/mxV03/wms/ent"
	"github.com/mxV03/wms/ent/user"
)

var (
	ErrInvalidUsername = fmt.Errorf("invalid username")
	ErrInvalidPassword = fmt.Errorf("invalid password")
	ErrInvalidRole     = fmt.Errorf("invalid role")
	ErrUserExists      = fmt.Errorf("user already exists")
	ErrUserNotFound    = fmt.Errorf("user not found")
	ErrUserDisabled    = fmt.Errorf("user disabled")
	ErrAuthFaild       = fmt.Errorf("authentication failed")
	ErrForbidden       = fmt.Errorf("forbidden")
)

const (
	RoleAdmin    = "Admin"
	RoleWorker   = "Worker"
	RoleReadOnly = "ReadOnly"
)

func NormalizeRole(r string) (string, error) {
	r = strings.TrimSpace(r)
	switch strings.ToLower(r) {
	case "admin":
		return RoleAdmin, nil
	case "worker":
		return RoleWorker, nil
	case "readonly", "read-only", "read_only":
		return RoleReadOnly, nil
	default:
		return "", ErrInvalidRole
	}
}

type AuthService struct {
	client *ent.Client
}

func NewAuthService(client *ent.Client) *AuthService {
	return &AuthService{
		client: client,
	}
}

type UserDTO struct {
	Username string
	Role     string
	Active   bool
}

func (s *AuthService) AddUser(ctx context.Context, username, role, password string) (*UserDTO, error) {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)

	if username == "" {
		return nil, ErrInvalidUsername
	}
	if len(password) < 4 {
		return nil, ErrInvalidPassword
	}

	nRole, err := NormalizeRole(role)
	if err != nil {
		return nil, err
	}

	exists, err := s.client.User.Query().
		Where(user.Username(username)).
		Exist(ctx)
	if err != nil {
		return nil, fmt.Errorf("check user existence: %w", err)
	}
	if exists {
		return nil, ErrUserExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	u, err := s.client.User.Create().
		SetUsername(username).
		SetPasswordHash(string(hash)).
		SetRole(nRole).
		SetActive(true).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, ErrUserExists
		}
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &UserDTO{
		Username: u.Username,
		Role:     u.Role,
		Active:   u.Active,
	}, nil
}

func (s *AuthService) DisableUser(ctx context.Context, username string) error {
	username = strings.TrimSpace(username)
	if username == "" {
		return ErrInvalidUsername
	}

	u, err := s.client.User.Query().
		Where(user.Username(username)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return ErrUserNotFound
		}
		return fmt.Errorf("fetch user: %w", err)
	}

	if err := s.client.User.UpdateOne(u).SetActive(false).Exec(ctx); err != nil {
		return fmt.Errorf("disable user: %w", err)
	}
	return nil
}

func (s *AuthService) ListUser(ctx context.Context, limit int) ([]*UserDTO, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}

	us, err := s.client.User.Query().
		Order(ent.Asc(user.FieldUsername)).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	out := make([]*UserDTO, 0, len(us))
	for _, u := range us {
		out = append(out, &UserDTO{
			Username: u.Username,
			Role:     u.Role,
			Active:   u.Active,
		})
	}
	return out, nil
}

type Principal struct {
	Username string
	Role     string
}

func CredentialsFromEnv() (string, string) {
	return strings.TrimSpace(os.Getenv("WMS_USER")), strings.TrimSpace(os.Getenv("WMS_PASS"))
}

func (s *AuthService) Authenticate(ctx context.Context, username, password string) (*Principal, error) {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)

	if username == "" || password == "" {
		return nil, ErrAuthFaild
	}

	u, err := s.client.User.Query().Where(user.Username(username)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrAuthFaild
		}
		return nil, fmt.Errorf("fetch user: %w", err)
	}
	if !u.Active {
		return nil, ErrUserDisabled
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, ErrAuthFaild
	}

	return &Principal{
		Username: u.Username,
		Role:     u.Role,
	}, nil
}

func (s *AuthService) RequireRole(ctx context.Context, allowedRoles ...string) (*Principal, error) {
	u, p := CredentialsFromEnv()
	pr, err := s.Authenticate(ctx, u, p)
	if err != nil {
		return nil, err
	}
	for _, r := range allowedRoles {
		if pr.Role == r {
			return pr, nil
		}
	}
	return nil, ErrForbidden
}
