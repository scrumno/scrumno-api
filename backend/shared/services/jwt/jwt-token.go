package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrUnknownSigningMethod = errors.New("Неизвестный метод подписи")
	ErrInvalidToken = errors.New("Некорректный токен JWT")
	ErrFailedToGenerateToken = errors.New("Не удалось сгенерировать токен")
	ErrFailedToParseToken = errors.New("Не удалось распарсить токен")
)

type Manager struct {
	accessSecret string
	refreshSecret string
	accessTokenTtl time.Duration
	refreshTokenTtl time.Duration
}

type Config struct {
	AccessSecret string
	RefreshSecret string
	AccessTokenTtl time.Duration
	RefreshTokenTtl time.Duration
}

type TokenPair struct {
    AccessToken  string
    RefreshToken string
    SessionID   string
    ExpiresIn    int64
}

type Claims struct {
    UserID string `json:"user_id"`
    Phone  string `json:"phone"`
    UUID   string `json:"uuid"`
	SessionID string `json:"session_id"`
    jwt.RegisteredClaims
}

func NewManager(cfg Config) *Manager {
	return &Manager{
		accessSecret: cfg.AccessSecret,
		refreshSecret: cfg.RefreshSecret,
		accessTokenTtl: cfg.AccessTokenTtl,
		refreshTokenTtl: cfg.RefreshTokenTtl,
	}
}

func (m *Manager) generateToken(claims Claims, secret string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := t.SignedString([]byte(secret))
	if err != nil {
		return "", ErrFailedToGenerateToken
	}
	return signed, nil
}

func (m *Manager) GenerateTokenPair(userID, phone, sessionID string) (*TokenPair, error) {

	accessClaims := Claims{
		UserID: userID,
		Phone: phone,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
            ExpiresAt:  jwt.NewNumericDate(time.Now().Add(m.accessTokenTtl)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

    refreshClaims := Claims{
        UserID: userID,
        Phone:  phone,
		SessionID: sessionID,
        RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.refreshTokenTtl)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

	accessToken, err := m.generateToken(accessClaims, m.accessSecret)
	if err != nil {
		return nil, ErrFailedToGenerateToken
	}
	refreshToken, err := m.generateToken(refreshClaims, m.refreshSecret)
	if err != nil {
		return nil, ErrFailedToGenerateToken
	}
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		SessionID: sessionID,
		ExpiresIn: int64(m.accessTokenTtl.Seconds()),
	}, nil
}


func (m *Manager) validateToken(tokenString string, secret string) (*Claims, error) {
	tok, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnknownSigningMethod
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, ErrFailedToParseToken
	}
	claims, ok := tok.Claims.(*Claims)
	if !ok || !tok.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

func (m *Manager) RefreshExpiresAtUnix() int64 {
	return time.Now().Unix() + int64(m.refreshTokenTtl.Seconds())
}

func (m *Manager) ValidateRefreshToken(refreshToken string) (*Claims, error) {
	return m.validateToken(refreshToken, m.refreshSecret)
}

func (m *Manager) ValidateAccessToken(accessToken string) (*Claims, error) {
	return m.validateToken(accessToken, m.accessSecret)
}


