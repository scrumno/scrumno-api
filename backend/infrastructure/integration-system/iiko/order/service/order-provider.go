package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	iikoConfig "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/config"
	iikoHTTP "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/http"
	"github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/order/model"
	"github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/utils"
	helpers "github.com/scrumno/scrumno-api/infrastructure/integration-system/shared/helpers"
)

type OrderProvider struct {
	http   *http.Client
	config *iikoConfig.Config
}

func NewOrderProvider(cfg *iikoConfig.Config) *OrderProvider {
	return &OrderProvider{
		http:   iikoHTTP.NewClient(cfg),
		config: cfg,
	}
}

func (p *OrderProvider) Create(order any) (any, error) {
	body, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}

	base := strings.TrimRight(p.config.BaseURL, "/")
	url := fmt.Sprintf("%s/api/1/deliveries/create", base)

	resp, err := helpers.SendRequest(p.http, url, body, p.config.AccessToken)
	if err != nil {
		return nil, err
	}

	var response any
	if err := json.Unmarshal(resp, &response); err != nil {
		return nil, err
	}

	return response, nil
}

// GetList возвращает *model.TerminalOrdersList как any (см. type assertion у вызывающего).
// Окно: последние windowSeconds секунд до текущего момента (локальное время процесса).
// В результат включаются заказы с creationStatus != Success (ожидание/синхронизация с сервером, в т.ч. после офлайн-режима).
// Остальные — только если order.terminalGroupId совпадает с IIKO_TERMINAL_GROUP_ID и есть активность в этом окне
// (whenCreated / whenBillPrinted / whenClosed).
func (p *OrderProvider) GetList(windowSeconds int) (any, error) {
	if windowSeconds <= 0 {
		return nil, fmt.Errorf("windowSeconds должен быть больше 0")
	}

	tg := strings.TrimSpace(p.config.TerminalGroupID)
	if tg == "" {
		return nil, fmt.Errorf("не задан IIKO_TERMINAL_GROUP_ID (терминальная группа для отбора заказов)")
	}

	window := time.Duration(windowSeconds) * time.Second
	now := time.Now()
	winStart := now.Add(-window)
	dateFrom := winStart.Format("2006-01-02 15:04:05.000")
	dateTo := now.Format("2006-01-02 15:04:05.000")

	tableIDs, err := p.tableIDsForTerminalGroup(tg)
	if err != nil {
		return nil, err
	}
	if len(tableIDs) == 0 {
		return &model.TerminalOrdersList{
			TerminalGroupID: tg,
			WindowSeconds:   windowSeconds,
			DateFrom:        dateFrom,
			DateTo:          dateTo,
			TableIDs:        []string{},
			Orders:          []model.TableOrderListItem{},
			Note:            "для терминальной группы не найдено столов в available_restaurant_sections (проверьте настройки зала / режим кассы)",
		}, nil
	}

	reqBody := model.GetTableOrdersByTableRequest{
		OrganizationIDs: []string{p.config.OrganizationID.String()},
		TableIDs:        tableIDs,
		DateFrom:        &dateFrom,
		DateTo:          &dateTo,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	base := strings.TrimRight(p.config.BaseURL, "/")
	url := fmt.Sprintf("%s/api/1/order/by_table", base)

	resp, err := helpers.SendRequest(p.http, url, body, p.config.AccessToken)
	if err != nil {
		return nil, err
	}

	var raw struct {
		Orders []model.TableOrderListItem `json:"orders"`
	}
	if err := json.Unmarshal(resp, &raw); err != nil {
		return nil, err
	}

	filtered := make([]model.TableOrderListItem, 0, len(raw.Orders))
	for i := range raw.Orders {
		if p.keepTableOrder(&raw.Orders[i], tg, now, winStart) {
			filtered = append(filtered, raw.Orders[i])
		}
	}

	return &model.TerminalOrdersList{
		TerminalGroupID: tg,
		WindowSeconds:   windowSeconds,
		DateFrom:        dateFrom,
		DateTo:          dateTo,
		TableIDs:        tableIDs,
		Orders:          filtered,
	}, nil
}

func (p *OrderProvider) tableIDsForTerminalGroup(terminalGroupID string) ([]string, error) {
	reqBody := model.GetRestaurantSectionsRequest{
		TerminalGroupIDs: []string{terminalGroupID},
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	base := strings.TrimRight(p.config.BaseURL, "/")
	url := fmt.Sprintf("%s/api/1/reserve/available_restaurant_sections", base)

	resp, err := helpers.SendRequest(p.http, url, body, p.config.AccessToken)
	if err != nil {
		return nil, err
	}

	var parsed struct {
		RestaurantSections []struct {
			Tables []struct {
				ID string `json:"id"`
			} `json:"tables"`
		} `json:"restaurantSections"`
	}
	if err := json.Unmarshal(resp, &parsed); err != nil {
		return nil, err
	}

	seen := make(map[string]struct{})
	var ids []string
	for _, sec := range parsed.RestaurantSections {
		for _, t := range sec.Tables {
			id := strings.TrimSpace(t.ID)
			if id == "" {
				continue
			}
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			ids = append(ids, id)
		}
	}
	return ids, nil
}

func (p *OrderProvider) keepTableOrder(o *model.TableOrderListItem, terminalGroupID string, now, winStart time.Time) bool {
	if o == nil {
		return false
	}
	if o.CreationStatus != "" && o.CreationStatus != "Success" {
		return true
	}
	if o.Order == nil {
		return false
	}
	if !strings.EqualFold(strings.TrimSpace(o.Order.TerminalGroupID), terminalGroupID) {
		return false
	}
	for _, s := range []string{o.Order.WhenCreated, o.Order.WhenBillPrinted, o.Order.WhenClosed} {
		if t, ok := utils.ParseIikoLocalTime(s); ok && !t.Before(winStart) && !t.After(now) {
			return true
		}
	}
	return false
}
