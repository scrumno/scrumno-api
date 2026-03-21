package order

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Repository interface {
	Create(ctx context.Context, request CreateOrderRequest) (*OrderResponse, error)
	AddItems(ctx context.Context, request AddOrderItemsRequest) (*CorrelationIDResponse, error)
	AddCustomer(ctx context.Context, request AddCustomerToOrderRequest) (*CorrelationIDResponse, error)
	AddPayments(ctx context.Context, request AddOrderPaymentsRequest) (*CorrelationIDResponse, error)
	ChangePayments(ctx context.Context, request ChangeOrderPaymentsRequest) (*CorrelationIDResponse, error)
	Close(ctx context.Context, request CloseOrderRequest) (*CorrelationIDResponse, error)
	Cancel(ctx context.Context, request CancelOrderRequest) (*CorrelationIDResponse, error)
	InitByPosOrder(ctx context.Context, request InitByPosOrderRequest) (*CorrelationIDResponse, error)
	GetByIDs(ctx context.Context, request GetByIDRequest) (*TableOrdersResponse, error)
	GetByPosOrderIDs(ctx context.Context, request GetByIDRequest) (*TableOrdersResponse, error)
}

type repository struct {
	baseURL string
	client  *http.Client
}

func NewRepository(baseURL string, client *http.Client) Repository {
	if client == nil {
		client = http.DefaultClient
	}

	return &repository{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  client,
	}
}

func (r *repository) Create(ctx context.Context, request CreateOrderRequest) (*OrderResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать тело запроса создания заказа iiko: %w", err)
	}

	u := r.baseURL + "/api/1/order/create"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("не удалось создать запрос создания заказа iiko: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса создания заказа iiko: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать ответ создания заказа iiko: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ошибка создания заказа iiko: статус %d: %s", resp.StatusCode, strings.TrimSpace(string(raw)))
	}

	var out OrderResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("не удалось декодировать ответ создания заказа iiko: %w", err)
	}

	return &out, nil
}

func (r *repository) AddItems(ctx context.Context, request AddOrderItemsRequest) (*CorrelationIDResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать тело запроса добавления позиций заказа iiko: %w", err)
	}

	u := r.baseURL + "/api/1/order/add_items"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("не удалось создать запрос добавления позиций заказа iiko: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса добавления позиций заказа iiko: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать ответ добавления позиций заказа iiko: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ошибка добавления позиций заказа iiko: статус %d: %s", resp.StatusCode, strings.TrimSpace(string(raw)))
	}

	var out CorrelationIDResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("не удалось декодировать ответ добавления позиций заказа iiko: %w", err)
	}

	return &out, nil
}

func (r *repository) AddCustomer(ctx context.Context, request AddCustomerToOrderRequest) (*CorrelationIDResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать тело запроса добавления клиента в заказ iiko: %w", err)
	}

	u := r.baseURL + "/api/1/order/add_customer"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("не удалось создать запрос добавления клиента в заказ iiko: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса добавления клиента в заказ iiko: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать ответ добавления клиента в заказ iiko: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ошибка добавления клиента в заказ iiko: статус %d: %s", resp.StatusCode, strings.TrimSpace(string(raw)))
	}

	var out CorrelationIDResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("не удалось декодировать ответ добавления клиента в заказ iiko: %w", err)
	}

	return &out, nil
}

func (r *repository) AddPayments(ctx context.Context, request AddOrderPaymentsRequest) (*CorrelationIDResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать тело запроса добавления оплат в заказ iiko: %w", err)
	}

	u := r.baseURL + "/api/1/order/add_payments"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("не удалось создать запрос добавления оплат в заказ iiko: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса добавления оплат в заказ iiko: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать ответ добавления оплат в заказ iiko: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ошибка добавления оплат в заказ iiko: статус %d: %s", resp.StatusCode, strings.TrimSpace(string(raw)))
	}

	var out CorrelationIDResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("не удалось декодировать ответ добавления оплат в заказ iiko: %w", err)
	}

	return &out, nil
}

func (r *repository) ChangePayments(ctx context.Context, request ChangeOrderPaymentsRequest) (*CorrelationIDResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать тело запроса изменения оплат заказа iiko: %w", err)
	}

	u := r.baseURL + "/api/1/order/change_payments"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("не удалось создать запрос изменения оплат заказа iiko: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса изменения оплат заказа iiko: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать ответ изменения оплат заказа iiko: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ошибка изменения оплат заказа iiko: статус %d: %s", resp.StatusCode, strings.TrimSpace(string(raw)))
	}

	var out CorrelationIDResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("не удалось декодировать ответ изменения оплат заказа iiko: %w", err)
	}

	return &out, nil
}

func (r *repository) Close(ctx context.Context, request CloseOrderRequest) (*CorrelationIDResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать тело запроса закрытия заказа iiko: %w", err)
	}

	u := r.baseURL + "/api/1/order/close"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("не удалось создать запрос закрытия заказа iiko: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса закрытия заказа iiko: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать ответ закрытия заказа iiko: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ошибка закрытия заказа iiko: статус %d: %s", resp.StatusCode, strings.TrimSpace(string(raw)))
	}

	var out CorrelationIDResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("не удалось декодировать ответ закрытия заказа iiko: %w", err)
	}

	return &out, nil
}

func (r *repository) Cancel(ctx context.Context, request CancelOrderRequest) (*CorrelationIDResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать тело запроса отмены заказа iiko: %w", err)
	}

	u := r.baseURL + "/api/1/order/cancel"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("не удалось создать запрос отмены заказа iiko: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса отмены заказа iiko: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать ответ отмены заказа iiko: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ошибка отмены заказа iiko: статус %d: %s", resp.StatusCode, strings.TrimSpace(string(raw)))
	}

	var out CorrelationIDResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("не удалось декодировать ответ отмены заказа iiko: %w", err)
	}

	return &out, nil
}

func (r *repository) GetByPosOrderIDs(ctx context.Context, request GetByIDRequest) (*TableOrdersResponse, error) {
	request.OrderIDs = nil
	return r.GetByIDs(ctx, request)
}

func (r *repository) GetByIDs(ctx context.Context, request GetByIDRequest) (*TableOrdersResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать тело запроса заказов по ids iiko: %w", err)
	}

	u := r.baseURL + "/api/1/order/by_id"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("не удалось создать запрос заказов по ids iiko: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса заказов по ids iiko: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать ответ заказов по ids iiko: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ошибка получения заказов по ids iiko: статус %d: %s", resp.StatusCode, strings.TrimSpace(string(raw)))
	}

	var out TableOrdersResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("не удалось декодировать ответ заказов по ids iiko: %w", err)
	}

	return &out, nil
}

func (r *repository) InitByPosOrder(ctx context.Context, request InitByPosOrderRequest) (*CorrelationIDResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать тело запроса init_by_posOrder iiko: %w", err)
	}

	u := r.baseURL + "/api/1/order/init_by_posOrder"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("не удалось создать запрос init_by_posOrder iiko: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса init_by_posOrder iiko: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать ответ init_by_posOrder iiko: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ошибка init_by_posOrder iiko: статус %d: %s", resp.StatusCode, strings.TrimSpace(string(raw)))
	}

	var out CorrelationIDResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("не удалось декодировать ответ init_by_posOrder iiko: %w", err)
	}

	return &out, nil
}
