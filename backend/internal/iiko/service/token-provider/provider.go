package token_provider

import (
	"context"
	"sync"
)

type FetchTokenFunc func(ctx context.Context) (string, error)

type Provider struct {
	fetchToken FetchTokenFunc

	mu    sync.RWMutex
	token string
}

func NewProvider(fetchToken FetchTokenFunc) *Provider {
	return &Provider{fetchToken: fetchToken}
}

func (p *Provider) GetToken(ctx context.Context) (string, error) {
	p.mu.RLock()
	cached := p.token
	p.mu.RUnlock()
	if cached != "" {
		return cached, nil
	}
	return p.RefreshToken(ctx)
}

func (p *Provider) RefreshToken(ctx context.Context) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Повторная проверка после захвата локa:
	// другой goroutine мог уже обновить токен.
	if p.token != "" {
		return p.token, nil
	}

	token, err := p.fetchToken(ctx)
	if err != nil {
		return "", err
	}
	p.token = token
	return token, nil
}

func (p *Provider) ForceRefreshToken(ctx context.Context) (string, error) {
	p.mu.Lock()
	p.token = ""
	p.mu.Unlock()

	return p.RefreshToken(ctx)
}
