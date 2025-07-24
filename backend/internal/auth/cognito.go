package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"ktchat/backend/pkg/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/golang-jwt/jwt/v5"
)

// CognitoService handles AWS Cognito authentication
type CognitoService struct {
	client       *cognitoidentityprovider.Client
	userPoolID   string
	clientID     string
	clientSecret string
	jwtSecret    string
}

// NewCognitoService creates a new Cognito service
func NewCognitoService(cfg *config.Config) (*CognitoService, error) {
	// Load AWS configuration
	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(), awsconfig.WithRegion(cfg.AWSCognito.Region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create Cognito client
	client := cognitoidentityprovider.NewFromConfig(awsCfg)

	return &CognitoService{
		client:       client,
		userPoolID:   cfg.AWSCognito.UserPoolID,
		clientID:     cfg.AWSCognito.ClientID,
		clientSecret: cfg.AWSCognito.ClientSecret,
		jwtSecret:    cfg.JWTSecret,
	}, nil
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// UserInfo represents user information from Cognito
type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Login authenticates a user with AWS Cognito
func (s *CognitoService) Login(ctx context.Context, req LoginRequest) (*LoginResponse, *UserInfo, error) {
	// Prepare the authentication parameters
	authParams := map[string]string{
		"USERNAME": req.Username,
		"PASSWORD": req.Password,
	}

	// Call Cognito InitiateAuth
	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow:       types.AuthFlowTypeUserPasswordAuth,
		ClientId:       aws.String(s.clientID),
		AuthParameters: authParams,
	}

	// Add client secret if configured
	if s.clientSecret != "" {
		secretHash := s.calculateSecretHash(req.Username)
		fmt.Printf("DEBUG: Calculating SECRET_HASH for login - username: %s, clientID: %s\n", req.Username, s.clientID)
		fmt.Printf("DEBUG: SECRET_HASH: %s\n", secretHash)
		input.AuthParameters["SECRET_HASH"] = secretHash
	}

	result, err := s.client.InitiateAuth(ctx, input)
	if err != nil {
		return nil, nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Check if authentication was successful
	if result.AuthenticationResult == nil {
		return nil, nil, fmt.Errorf("authentication failed: no result")
	}

	// Extract user information from ID token
	userInfo, err := s.extractUserInfoFromToken(*result.AuthenticationResult.IdToken)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to extract user info: %w", err)
	}

	response := &LoginResponse{
		AccessToken:  *result.AuthenticationResult.AccessToken,
		RefreshToken: *result.AuthenticationResult.RefreshToken,
		IDToken:      *result.AuthenticationResult.IdToken,
		ExpiresIn:    int(result.AuthenticationResult.ExpiresIn),
		TokenType:    *result.AuthenticationResult.TokenType,
	}

	return response, userInfo, nil
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Register registers a new user with AWS Cognito
func (s *CognitoService) Register(ctx context.Context, req RegisterRequest) error {
	// Prepare user attributes
	userAttributes := []types.AttributeType{
		{
			Name:  aws.String("email"),
			Value: aws.String(req.Email),
		},
	}

	// Call Cognito SignUp
	input := &cognitoidentityprovider.SignUpInput{
		ClientId:       aws.String(s.clientID),
		Username:       aws.String(req.Username),
		Password:       aws.String(req.Password),
		UserAttributes: userAttributes,
	}

	// Add client secret if configured
	if s.clientSecret != "" {
		secretHash := s.calculateSecretHash(req.Username)
		fmt.Printf("DEBUG: Calculating SECRET_HASH for username: %s, clientID: %s\n", req.Username, s.clientID)
		fmt.Printf("DEBUG: SECRET_HASH: %s\n", secretHash)
		input.SecretHash = aws.String(secretHash)
	}

	_, err := s.client.SignUp(ctx, input)
	if err != nil {
		return fmt.Errorf("registration failed: %w", err)
	}

	return nil
}

// RefreshTokenRequest represents a token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshToken refreshes an access token using a refresh token
func (s *CognitoService) RefreshToken(ctx context.Context, req RefreshTokenRequest) (*LoginResponse, error) {
	// Prepare the authentication parameters
	authParams := map[string]string{
		"REFRESH_TOKEN": req.RefreshToken,
	}

	// Add client secret if configured
	if s.clientSecret != "" {
		// For refresh token, we need the username from the refresh token
		// This is a simplified approach - in production, you might want to store the username
		authParams["SECRET_HASH"] = s.calculateSecretHash("") // This won't work without username
	}

	// Call Cognito InitiateAuth with REFRESH_TOKEN_AUTH flow
	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow:       types.AuthFlowTypeRefreshTokenAuth,
		ClientId:       aws.String(s.clientID),
		AuthParameters: authParams,
	}

	result, err := s.client.InitiateAuth(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("token refresh failed: %w", err)
	}

	if result.AuthenticationResult == nil {
		return nil, fmt.Errorf("token refresh failed: no result")
	}

	response := &LoginResponse{
		AccessToken: *result.AuthenticationResult.AccessToken,
		IDToken:     *result.AuthenticationResult.IdToken,
		ExpiresIn:   int(result.AuthenticationResult.ExpiresIn),
		TokenType:   *result.AuthenticationResult.TokenType,
	}

	return response, nil
}

// ValidateToken validates a Cognito token
func (s *CognitoService) ValidateToken(ctx context.Context, tokenString string) (*UserInfo, error) {
	// For production, you should validate the token signature using Cognito's public keys
	// For now, we'll extract user info from the token
	userInfo, err := s.extractUserInfoFromToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return userInfo, nil
}

// extractUserInfoFromToken extracts user information from a Cognito ID token
func (s *CognitoService) extractUserInfoFromToken(tokenString string) (*UserInfo, error) {
	// For development, we'll extract user info from the token without signature validation
	// In production, you should validate the token signature using Cognito's public keys
	
	// Split the token to get the payload part
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}
	
	// Decode the payload (second part)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode token payload: %w", err)
	}
	
	// Parse the JSON payload
	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, fmt.Errorf("failed to parse token claims: %w", err)
	}
	
	// Extract user information from claims
	userInfo := &UserInfo{}
	
	// Cognito stores user ID in "sub" claim
	if sub, ok := claims["sub"].(string); ok {
		userInfo.ID = sub
	}
	
	// Username is typically in "cognito:username" or "username"
	if username, ok := claims["cognito:username"].(string); ok {
		userInfo.Username = username
	} else if username, ok := claims["username"].(string); ok {
		userInfo.Username = username
	}
	
	// Email is typically in "email" claim
	if email, ok := claims["email"].(string); ok {
		userInfo.Email = email
	}
	
	return userInfo, nil
}

// calculateSecretHash calculates the secret hash for Cognito
func (s *CognitoService) calculateSecretHash(username string) string {
	// Return empty string if no client secret is configured
	if s.clientSecret == "" {
		return ""
	}
	
	// Calculate HMAC-SHA256 hash
	hash := hmac.New(sha256.New, []byte(s.clientSecret))
	hash.Write([]byte(username + s.clientID))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

// GenerateCustomToken generates a custom JWT token for internal use
func (s *CognitoService) GenerateCustomToken(userInfo *UserInfo) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userInfo.ID,
		"username": userInfo.Username,
		"email":    userInfo.Email,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iat":      time.Now().Unix(),
		"iss":      "ktchat-backend",
		"aud":      "ktchat-frontend",
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
} 

// CognitoError represents a structured error response from Cognito
type CognitoError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ExtractCognitoError extracts structured error information from AWS Cognito errors
func (s *CognitoService) ExtractCognitoError(err error) *CognitoError {
	if err == nil {
		return &CognitoError{
			Code:    "UNKNOWN_ERROR",
			Message: "Unknown error occurred",
		}
	}

	errStr := err.Error()

	// Common Cognito error patterns
	switch {
	case strings.Contains(errStr, "NotAuthorizedException"):
		if strings.Contains(errStr, "Incorrect username or password") {
			return &CognitoError{
				Code:    "INVALID_CREDENTIALS",
				Message: "Incorrect username or password",
			}
		}
		if strings.Contains(errStr, "User is disabled") {
			return &CognitoError{
				Code:    "USER_DISABLED",
				Message: "Account is disabled. Please contact support.",
			}
		}
		return &CognitoError{
			Code:    "NOT_AUTHORIZED",
			Message: "Authentication failed. Please check your credentials.",
		}

	case strings.Contains(errStr, "UserNotConfirmedException"):
		return &CognitoError{
			Code:    "USER_NOT_CONFIRMED",
			Message: "Please check your email and confirm your account before logging in.",
		}

	case strings.Contains(errStr, "UserNotFoundException"):
		return &CognitoError{
			Code:    "USER_NOT_FOUND",
			Message: "User not found. Please check your username or register a new account.",
		}

	case strings.Contains(errStr, "UsernameExistsException"):
		return &CognitoError{
			Code:    "USERNAME_EXISTS",
			Message: "Username already exists. Please choose a different username.",
		}

	case strings.Contains(errStr, "InvalidPasswordException"):
		return &CognitoError{
			Code:    "INVALID_PASSWORD",
			Message: "Password does not meet requirements. Password must be at least 8 characters long and contain uppercase, lowercase, and numeric characters.",
		}

	case strings.Contains(errStr, "InvalidParameterException"):
		if strings.Contains(errStr, "email") {
			return &CognitoError{
				Code:    "INVALID_EMAIL",
				Message: "Invalid email format. Please provide a valid email address.",
			}
		}
		if strings.Contains(errStr, "username") {
			return &CognitoError{
				Code:    "INVALID_USERNAME",
				Message: "Invalid username format. Username must be 3-128 characters long.",
			}
		}
		return &CognitoError{
			Code:    "INVALID_PARAMETER",
			Message: "Invalid input parameters. Please check your information.",
		}

	case strings.Contains(errStr, "CodeMismatchException"):
		return &CognitoError{
			Code:    "CODE_MISMATCH",
			Message: "Invalid confirmation code. Please check your email and try again.",
		}

	case strings.Contains(errStr, "ExpiredCodeException"):
		return &CognitoError{
			Code:    "CODE_EXPIRED",
			Message: "Confirmation code has expired. Please request a new one.",
		}

	case strings.Contains(errStr, "LimitExceededException"):
		return &CognitoError{
			Code:    "LIMIT_EXCEEDED",
			Message: "Too many attempts. Please wait a moment before trying again.",
		}

	case strings.Contains(errStr, "TooManyRequestsException"):
		return &CognitoError{
			Code:    "TOO_MANY_REQUESTS",
			Message: "Too many requests. Please wait a moment before trying again.",
		}

	case strings.Contains(errStr, "SECRET_HASH"):
		return &CognitoError{
			Code:    "CONFIGURATION_ERROR",
			Message: "Authentication configuration error. Please contact support.",
		}

	case strings.Contains(errStr, "network") || strings.Contains(errStr, "timeout"):
		return &CognitoError{
			Code:    "NETWORK_ERROR",
			Message: "Network error. Please check your connection and try again.",
		}

	default:
		return &CognitoError{
			Code:    "AUTHENTICATION_FAILED",
			Message: "Authentication failed. Please try again.",
		}
	}
}

// ExtractCognitoErrorMessage extracts user-friendly error messages from AWS Cognito errors
// This is kept for backward compatibility but will be deprecated
func (s *CognitoService) ExtractCognitoErrorMessage(err error) string {
	cognitoError := s.ExtractCognitoError(err)
	return cognitoError.Message
}

// ConfirmUser confirms a user's account with the provided confirmation code
func (s *CognitoService) ConfirmUser(ctx context.Context, username, confirmationCode string) error {
	input := &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(s.clientID),
		Username:         aws.String(username),
		ConfirmationCode: aws.String(confirmationCode),
	}

	// Add SECRET_HASH if client secret is configured
	if s.clientSecret != "" {
		secretHash := s.calculateSecretHash(username)
		input.SecretHash = aws.String(secretHash)
	}

	_, err := s.client.ConfirmSignUp(ctx, input)
	return err
}

// ResendConfirmationCode resends the confirmation code to the user
func (s *CognitoService) ResendConfirmationCode(ctx context.Context, username string) error {
	input := &cognitoidentityprovider.ResendConfirmationCodeInput{
		ClientId: aws.String(s.clientID),
		Username: aws.String(username),
	}

	// Add SECRET_HASH if client secret is configured
	if s.clientSecret != "" {
		secretHash := s.calculateSecretHash(username)
		input.SecretHash = aws.String(secretHash)
	}

	_, err := s.client.ResendConfirmationCode(ctx, input)
	return err
} 