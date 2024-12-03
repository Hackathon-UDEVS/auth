package postgres

import (
	"context"
	"database/sql"
	"errors"
	pb "github.com/Hackaton-UDEVS/auth/internal/genproto/auth"
	logger "github.com/Hackaton-UDEVS/auth/internal/logger"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	_ "github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
)

type UserRepo struct {
	db *sql.DB
	rd *redis.Client
}

func NewUserRepo(db *sql.DB, rd *redis.Client) *UserRepo {
	return &UserRepo{
		db: db,
		rd: rd,
	}
}
func (s *UserRepo) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginRes, error) {
	logs, err := logger.NewLogger()
	if err != nil {
		return nil, err
	}

	// Query to get user details based on the email
	query := `
		SELECT id, email, password, role, created_at, updated_at
		FROM users
		WHERE email = $1 AND deleted_at = 0
	`
	var user pb.UserModel
	var hashedPassword string

	// Execute query
	err = s.db.QueryRowContext(ctx, query, req.Email).Scan(
		&user.Id, &user.Email, &hashedPassword, &user.Role,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// User not found
			logs.Warn("User not found", zap.String("email", req.Email))
			return nil, errors.New("user not found")
		}
		// Error querying user
		logs.Error("Error querying user", zap.Error(err))
		return nil, err
	}

	// Check if password matches
	if !checkPasswordHash(req.Password, hashedPassword) {
		// Incorrect password
		logs.Warn("Invalid password", zap.String("email", req.Email))
		return nil, errors.New("invalid email or password")
	}

	// Return user data if login is successful123456789
	return &pb.LoginRes{UserRes: &user}, nil
}

type UserModel struct {
	Email    string
	Password string
	Role     string
}

func (s *UserRepo) RegisterUser(ctx context.Context, req *pb.RegisterUserReq) (*pb.RegisterUserRes, error) {
	logs, err := logger.NewLogger()
	if err != nil {
		return nil, err
	}

	// Check if the email already exists in the database1234
	var count int
	queryCheck := `
		SELECT COUNT(*) 
		FROM users 
		WHERE email = $1`
	err = s.db.QueryRowContext(ctx, queryCheck, req.Email).Scan(&count)
	if err != nil {
		logs.Error("Error checking email existence", zap.Error(err))
		return nil, err
	}

	// If email exists, return duplicate email error 12
	if count > 0 {
		logs.Warn("Duplicate email found", zap.String("email", req.Email))
		return &pb.RegisterUserRes{Message: "Duplicate email"}, nil
	}

	// Proceed with inserting the new user if no duplicate e
	id := uuid.NewString()
	queryInsert := `
		INSERT INTO users (id, email, password, role) 
		VALUES ($1, $2, $3, $4)`

	hashedPassword, _ := hashPassword(req.Password)

	_, err = s.db.ExecContext(ctx, queryInsert,
		id,
		req.Email,
		hashedPassword,
		req.Role,
	)
	if err != nil {
		logs.Error("Error registering user", zap.Error(err))
		return nil, err
	}

	return &pb.RegisterUserRes{Message: "Success register user"}, nil
}

func (s *UserRepo) GetUserByID(ctx context.Context, req *pb.GetUserByIDReq) (*pb.GetUserByIDRes, error) {
	logs, err := logger.NewLogger()
	if err != nil {
		return nil, err
	}

	query := `SELECT id, email, role, created_at, updated_at 
              FROM users 
              WHERE id = $1 AND deleted_at = 0`

	user := pb.UserModel{}
	err = s.db.QueryRowContext(ctx, query, req.Userid).Scan(
		&user.Id,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		logs.Error("Error getting user", zap.Error(err))
		return nil, err
	}

	return &pb.GetUserByIDRes{UserRes: &user}, nil
}

func (s *UserRepo) GetAllUsers(ctx context.Context, req *pb.GetAllUserReq) (*pb.GetAllUserRes, error) {
	query := `SELECT id, email, role, created_at, updated_at 
              FROM users 
              WHERE deleted_at = 0`

	logs, err := logger.NewLogger()
	if err != nil {
		return nil, err
	}

	args := []interface{}{}
	argCounter := 1

	if req.UserReq.Id != "" && req.UserReq.Id != "string" {
		query += " AND id = $" + strconv.Itoa(argCounter)
		args = append(args, req.UserReq.Id)
		argCounter++
	}

	if req.UserReq.Email != "" && req.UserReq.Email != "string" {
		query += " AND email = $" + strconv.Itoa(argCounter)
		args = append(args, req.UserReq.Email)
		argCounter++
	}

	if req.UserReq.Role != "" && req.UserReq.Role != "string" {
		query += " AND role = $" + strconv.Itoa(argCounter)
		args = append(args, req.UserReq.Role)
		argCounter++
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		logs.Error("Error with get all users query")
		return nil, err
	}
	defer rows.Close()

	var users []*pb.UserModel
	for rows.Next() {
		user := pb.UserModel{}
		err = rows.Scan(
			&user.Id,
			&user.Email,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			logs.Error("Error scanning user", zap.Error(err))
			continue
		}
		users = append(users, &user)
	}

	return &pb.GetAllUserRes{UserRes: users}, nil
}

func (s *UserRepo) UpdateUser(ctx context.Context, req *pb.UpdateUserReq) (*pb.UpdateUserRes, error) {
	query := "UPDATE users SET"
	var args []interface{}
	var updates []string
	argCounter := 1

	logs, err := logger.NewLogger()
	if err != nil {
		return nil, err
	}

	if req.UserReq.Email != "" && req.UserReq.Email != "string" {
		updates = append(updates, " email = $"+strconv.Itoa(argCounter))
		args = append(args, req.UserReq.Email)
		argCounter++
	}

	if req.UserReq.Role != "" && req.UserReq.Role != "string" {
		if req.UserReq.Role != "client" && req.UserReq.Role != "contractors" {
			return nil, errors.New("invalid role: must be either 'client' or 'contractors'")
		}
		updates = append(updates, " role = $"+strconv.Itoa(argCounter))
		args = append(args, req.UserReq.Role)
		argCounter++
	}

	if len(updates) == 0 {
		return nil, errors.New("no fields to update")
	}

	updates = append(updates, " updated_at = now()")
	query += strings.Join(updates, ", ") + " WHERE id = $" + strconv.Itoa(argCounter)
	args = append(args, req.UserReq.Id)

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		logs.Error("Error updating user", zap.Error(err))
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logs.Error("Error getting rows affected", zap.Error(err))
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, errors.New("user not found")
	}

	return &pb.UpdateUserRes{UserRes: req.UserReq}, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
