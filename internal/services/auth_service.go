package services

import (
	"database/sql"
	"fmt"
	"time"

	"auth-backend/internal/models"
	"auth-backend/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db           *sql.DB
	emailService *EmailService
	jwtSecret    string
}

func NewAuthService(db *sql.DB, emailService *EmailService, jwtSecret string) *AuthService {
	return &AuthService{
		db:           db,
		emailService: emailService,
		jwtSecret:    jwtSecret,
	}
}

func (as *AuthService) RequestVerification(email string) error {
	// Generate and send verification code
	code, err := as.emailService.SendVerificationEmail(email)
	if err != nil {
		return fmt.Errorf("failed to send verification email: %v", err)
	}

	// Store verification code in database (using UTC)
	expiresAt := time.Now().UTC().Add(10 * time.Minute)
	_, err = as.db.Exec(`
		INSERT INTO verification_codes (email, code, expires_at) 
		VALUES ($1, $2, $3)`,
		email, code, expiresAt)

	if err != nil {
		return fmt.Errorf("failed to store verification code: %v", err)
	}

	return nil
}

func (as *AuthService) VerifyCode(email, code string) error {
	// Check if code is valid and not expired
	var codeRecord models.VerificationCode
	err := as.db.QueryRow(`
		SELECT id, email, code, expires_at, used 
		FROM verification_codes 
		WHERE email = $1 AND code = $2 AND used = FALSE
		ORDER BY created_at DESC LIMIT 1`,
		email, code).Scan(
		&codeRecord.ID, &codeRecord.Email, &codeRecord.Code,
		&codeRecord.ExpiresAt, &codeRecord.Used)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("invalid verification code")
		}
		return fmt.Errorf("database error: %v", err)
	}

	// Check if code is expired (using UTC)
	if time.Now().UTC().After(codeRecord.ExpiresAt) {
		return fmt.Errorf("verification code has expired")
	}

	// Mark code as used and set verified_at timestamp (using UTC)
	_, err = as.db.Exec(`UPDATE verification_codes SET used = TRUE, verified_at = $1 WHERE id = $2`, time.Now().UTC(), codeRecord.ID)
	if err != nil {
		return fmt.Errorf("failed to update verification code: %v", err)
	}

	return nil
}

func (as *AuthService) SetPassword(email, password string) (*models.AuthResponse, error) {
	// Check if user has a verified code (within last 30 minutes from when it was verified, using UTC)
	var count int
	err := as.db.QueryRow(`
		SELECT COUNT(*) FROM verification_codes 
		WHERE email = $1 AND used = TRUE AND verified_at > $2`,
		email, time.Now().UTC().Add(-30*time.Minute)).Scan(&count)

	// Fallback: if no verified_at timestamp, check for recent used codes (for backward compatibility)
	if err != nil || count == 0 {
		err = as.db.QueryRow(`
			SELECT COUNT(*) FROM verification_codes 
			WHERE email = $1 AND used = TRUE AND created_at > $2`,
			email, time.Now().UTC().Add(-2*time.Hour)).Scan(&count)

		if err != nil || count == 0 {
			return nil, fmt.Errorf("email not verified or verification expired")
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	// Create or update user
	var userID int
	err = as.db.QueryRow(`
		INSERT INTO users (email, password_hash, is_verified) 
		VALUES ($1, $2, TRUE) 
		ON CONFLICT (email) 
		DO UPDATE SET password_hash = $2, is_verified = TRUE, updated_at = CURRENT_TIMESTAMP
		RETURNING id`,
		email, string(hashedPassword)).Scan(&userID)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	// Generate JWT token
	token, err := jwt.GenerateToken(userID, email, as.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	// Get user details
	user, err := as.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user details: %v", err)
	}

	return &models.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (as *AuthService) Login(email, password string) (*models.AuthResponse, error) {
	// Get user from database
	var user models.User
	err := as.db.QueryRow(`
		SELECT id, email, password_hash, is_verified, created_at, updated_at 
		FROM users WHERE email = $1`, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.IsVerified,
		&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	// Check if user is verified
	if !user.IsVerified {
		return nil, fmt.Errorf("email not verified")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := jwt.GenerateToken(user.ID, user.Email, as.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	return &models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (as *AuthService) GetUserByID(userID int) (*models.User, error) {
	var user models.User
	err := as.db.QueryRow(`
		SELECT id, email, is_verified, created_at, updated_at 
		FROM users WHERE id = $1`, userID).Scan(
		&user.ID, &user.Email, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
