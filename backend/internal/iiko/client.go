package iiko

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/scrumno/scrumno-api/config"
)

// Client определяет интерфейс работы с iikoRMS.
type Client interface {
	CreatePickupOrder(ctx context.Context, order PickupOrder) (*CreateOrderResult, error)
	GetOrganizations(ctx context.Context) (*CreateOrderResult, error)
	GetNomenclature(ctx context.Context, organizationID string) (*CreateOrderResult, error)
}

type client struct {
	cfg    config.IikoConfig
	client *http.Client
}

func NewClient(cfg config.IikoConfig) Client {
	httpClient := &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).DialContext,
		},
	}

	return &client{
		cfg:    cfg,
		client: httpClient,
	}
}

const (
	authPath          = "api/auth"
	createOrderPath   = "api/deliveries/create"
	organizationsPath = "api/organization/list"
	nomenclaturePath  = "api/nomenclature"
)

func (c *client) CreatePickupOrder(ctx context.Context, order PickupOrder) (*CreateOrderResult, error) {
	if c.cfg.BaseURL == "" || c.cfg.Login == "" || c.cfg.Password == "" {
		return nil, fmt.Errorf("iiko config is not fully set")
	}

	token, err := c.auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("iiko auth error: %w", err)
	}

	reqPayload := map[string]any{
		"organizationId": c.cfg.OrganizationID,
		"terminalId":     c.cfg.TerminalID,
		"type":           "DeliveryByClient",
		"phone":          order.Customer.Phone,
		"customerName":   order.Customer.Name,
		"comment":        order.Comment,
		"items":          order.Items,
	}

	if order.PickupAt != nil {
		reqPayload["pickupAt"] = order.PickupAt.Format(time.RFC3339)
	}

	body, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, fmt.Errorf("marshal iiko request: %w", err)
	}

	u, err := url.Parse(c.cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse base url: %w", err)
	}
	u.Path = path.Join(u.Path, createOrderPath)

	q := u.Query()
	q.Set("key", token)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build create order request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send create order request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read create order response: %w", err)
	}

	if resp.StatusCode >= 400 {
		slog.Error("iiko create pickup order failed",
			"status", resp.StatusCode,
			"body", string(respBody),
		)
		return &CreateOrderResult{
			StatusCode: resp.StatusCode,
			Body:       respBody,
		}, fmt.Errorf("iiko create pickup order failed with status %d", resp.StatusCode)
	}

	return &CreateOrderResult{
		StatusCode: resp.StatusCode,
		Body:       respBody,
	}, nil
}

func (c *client) auth(ctx context.Context) (string, error) {
	u, err := url.Parse(c.cfg.BaseURL)
	if err != nil {
		return "", fmt.Errorf("parse base url: %w", err)
	}
	u.Path = path.Join(u.Path, authPath)

	q := u.Query()
	q.Set("login", c.cfg.Login)
	q.Set("pass", c.cfg.Password)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("build auth request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("send auth request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read auth response: %w", err)
	}

	if resp.StatusCode >= 400 {
		slog.Error("iiko auth failed",
			"status", resp.StatusCode,
			"body", string(body),
		)
		return "", fmt.Errorf("iiko auth failed with status %d", resp.StatusCode)
	}

	token := strings.TrimSpace(string(body))
	if token == "" {
		return "", fmt.Errorf("empty token from iiko auth")
	}

	return token, nil
}

// GetOrganizations запрашивает список организаций из iiko.
func (c *client) GetOrganizations(ctx context.Context) (*CreateOrderResult, error) {
	if c.cfg.BaseURL == "" || c.cfg.Login == "" || c.cfg.Password == "" {
		return nil, fmt.Errorf("iiko config is not fully set")
	}

	token, err := c.auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("iiko auth error: %w", err)
	}

	u, err := url.Parse(c.cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse base url: %w", err)
	}
	u.Path = path.Join(u.Path, organizationsPath)

	q := u.Query()
	q.Set("key", token)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("build get organizations request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send get organizations request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read get organizations response: %w", err)
	}

	if resp.StatusCode >= 400 {
		slog.Error("iiko get organizations failed",
			"status", resp.StatusCode,
			"body", string(body),
		)
		return &CreateOrderResult{
			StatusCode: resp.StatusCode,
			Body:       body,
		}, fmt.Errorf("iiko get organizations failed with status %d", resp.StatusCode)
	}

	return &CreateOrderResult{
		StatusCode: resp.StatusCode,
		Body:       body,
	}, nil
}

// GetNomenclature запрашивает номенклатуру (товары) по организации.
func (c *client) GetNomenclature(ctx context.Context, organizationID string) (*CreateOrderResult, error) {
	if c.cfg.BaseURL == "" || c.cfg.Login == "" || c.cfg.Password == "" {
		return nil, fmt.Errorf("iiko config is not fully set")
	}

	if organizationID == "" {
		return nil, fmt.Errorf("organizationID is required")
	}

	token, err := c.auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("iiko auth error: %w", err)
	}

	u, err := url.Parse(c.cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse base url: %w", err)
	}
	u.Path = path.Join(u.Path, nomenclaturePath)

	q := u.Query()
	q.Set("key", token)
	q.Set("organizationId", organizationID)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("build get nomenclature request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send get nomenclature request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read get nomenclature response: %w", err)
	}

	if resp.StatusCode >= 400 {
		slog.Error("iiko get nomenclature failed",
			"status", resp.StatusCode,
			"body", string(body),
		)
		return &CreateOrderResult{
			StatusCode: resp.StatusCode,
			Body:       body,
		}, fmt.Errorf("iiko get nomenclature failed with status %d", resp.StatusCode)
	}

	return &CreateOrderResult{
		StatusCode: resp.StatusCode,
		Body:       body,
	}, nil
}

