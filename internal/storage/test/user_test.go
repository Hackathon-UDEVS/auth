package postgres_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/Hackaton-UDEVS/auth/internal/genproto/auth"
	"github.com/Hackaton-UDEVS/auth/internal/storage/postgres"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocking the database
type MockDB struct {
	mock.Mock
}

func (m *MockDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	args = append([]interface{}{ctx, query}, args...)
	m.Called(args...)
	return nil // You can return a mock row here
}

// Mocking Redis client
type MockRedis struct {
	mock.Mock
}

func (m *MockRedis) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(key, value, expiration)
	return args.Get(0).(*redis.StatusCmd)
}

func (m *MockRedis) Get(key string) *redis.StringCmd {
	args := m.Called(key)
	return args.Get(0).(*redis.StringCmd)
}

// Test for the Login function
func TestLogin(t *testing.T) {
	mockDB := new(MockDB)
	mockRedis := new(MockRedis)

	// Mocking the database query result
	mockDB.On("QueryRowContext", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockRedis.On("Get", mock.Anything).Return(nil)

	repo := postgres.NewUserRepo(mockDB, mockRedis)

	req := &auth.LoginReq{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Test for successful login
	t.Run("Successful login", func(t *testing.T) {
		mockDB.On("QueryRowContext", mock.Anything, mock.Anything, req.Email).Return(nil)

		res, err := repo.Login(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})

	// Test for invalid email or password
	t.Run("Invalid email or password", func(t *testing.T) {
		mockDB.On("QueryRowContext", mock.Anything, mock.Anything, req.Email).Return(errors.New("user not found"))

		res, err := repo.Login(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

// Test for the RegisterUser function
func TestRegisterUser(t *testing.T) {
	mockDB := new(MockDB)
	mockRedis := new(MockRedis)

	// Create a mock for Redis Set operation
	mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(&redis.StatusCmd{})

	repo := postgres.NewUserRepo(mockDB, mockRedis)

	req := &auth.RegisterUserReq{
		Email:    "test@example.com",
		Password: "password123",
		Role:     "client",
	}

	// Test for registering a new user
	t.Run("Register a new user", func(t *testing.T) {
		// Mock checking if email already exists
		mockDB.On("QueryRowContext", mock.Anything, mock.Anything, req.Email).Return(nil)

		res, err := repo.RegisterUser(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, "Verification code sent to your email", res.Message)
	})

	// Test for already existing email
	t.Run("Email already registered", func(t *testing.T) {
		// Mock checking if email already exists
		mockDB.On("QueryRowContext", mock.Anything, mock.Anything, req.Email).Return(errors.New("email already registered"))

		res, err := repo.RegisterUser(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}
