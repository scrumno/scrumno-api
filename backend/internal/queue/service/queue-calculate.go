package service

import (
	"math"
	"time"

	"github.com/scrumno/scrumno-api/internal/queue/entity"
	"gorm.io/gorm"
)

// DurationRange диапазон времени от минимального до максимального.
type DurationRange struct {
	Min time.Duration `json:"min"`
	Max time.Duration `json:"max"`
}

type OrdersQueueResult struct {
	CurrentOrderNoQueue DurationRange `json:"current_order_no_queue"`
	QueueWait           DurationRange `json:"queue_wait"`
	Total               DurationRange `json:"total"`
}

type OrdersQueueService interface {
	EstimateQueueTime(current entity.OrdersQueueOrder, ordersAhead []entity.OrdersQueueOrder) OrdersQueueResult
}

type ordersQueueService struct {
	cfg *entity.OrdersQueueConfigTable
	db  *gorm.DB
}

func NewOrdersQueueService(cfg *entity.OrdersQueueConfigTable, db *gorm.DB) OrdersQueueService {
	return &ordersQueueService{cfg: cfg, db: db}
}

// EstimateQueueTime оценивает итоговое время текущего заказа:
// время в очереди + время приготовления самого заказа.
func (s *ordersQueueService) EstimateQueueTime(
	current entity.OrdersQueueOrder,
	ordersAhead []entity.OrdersQueueOrder,
) OrdersQueueResult {
	// Время приготовления текущего заказа без учета чужой очереди.
	currentNoQueueBase := estimateOrderNoQueue(current, s.cfg.OrderReserveMinutes)
	// Время самого заказа считаем фиксированным; диапазон даем только по очереди.
	currentNoQueueRange := DurationRange{
		Min: currentNoQueueBase,
		Max: currentNoQueueBase,
	}

	// Базовые минуты добавляем только если очереди нет.
	// Это минимальное ожидание старта кухни при пустой очереди.
	if len(ordersAhead) == 0 {
		queueWaitRange := minutesToRange(s.cfg.EmptyQueueWaitMinMins, s.cfg.EmptyQueueWaitMaxMins)
		totalRange := addRanges(currentNoQueueRange, queueWaitRange)
		totalRange = applyRestaurantHoursToRange(totalRange, *s.cfg, time.Now())
		queueWaitRange = subtractFixedFromRange(totalRange, currentNoQueueBase)

		return OrdersQueueResult{
			CurrentOrderNoQueue: currentNoQueueRange,
			QueueWait:           queueWaitRange,
			Total:               totalRange,
		}
	}

	// Для непустой очереди считаем только фактическую нагрузку:
	// суммируем время всех заказов перед текущим.
	queueLoad := time.Duration(0)
	for _, order := range ordersAhead {
		queueLoad += estimateOrderNoQueue(order, 0)
	}

	// Если параллельность не задана/ошибочна, считаем как один канал кухни.
	slots := s.cfg.KitchenParallelSlots
	if slots <= 0 {
		slots = 1
	}

	// Делим нагрузку на число параллельных "слотов" кухни.
	queueByCapacity := time.Duration(float64(queueLoad) / float64(slots))
	// Плавный рост ожидания от длины очереди, без резких скачков по порогу.
	queuePressure := 1.0 + s.cfg.QueueGrowthFactor*math.Log1p(float64(len(ordersAhead)))
	queueWaitBase := time.Duration(float64(queueByCapacity) * queuePressure)
	queueWaitRange := scaleDurationToRange(
		queueWaitBase,
		s.cfg.QueueTimeMinFactor,
		s.cfg.QueueTimeMaxFactor,
	)
	totalRange := addRanges(currentNoQueueRange, queueWaitRange)
	totalRange = applyRestaurantHoursToRange(totalRange, *s.cfg, time.Now())
	queueWaitRange = subtractFixedFromRange(totalRange, currentNoQueueBase)

	return OrdersQueueResult{
		CurrentOrderNoQueue: currentNoQueueRange,
		QueueWait:           queueWaitRange,
		Total:               totalRange,
	}
}

func estimateOrderNoQueue(order entity.OrdersQueueOrder, orderReserveMinutes int) time.Duration {
	// База заказа: стартовая подготовка + резерв на сам заказ.
	total := time.Duration(order.SetupMinutes+orderReserveMinutes) * time.Minute
	for _, item := range order.Items {
		total += estimateItemDuration(item)
	}
	return total
}

func estimateItemDuration(item entity.OrderItem) time.Duration {
	quantity := item.Quantity
	if quantity <= 0 {
		return 0
	}

	base := time.Duration(item.BaseCookMinutes) * time.Minute
	complexity := item.ComplexityFactor
	if complexity <= 0 {
		complexity = 1
	}

	total := time.Duration(0)
	for piece := 1; piece <= quantity; piece++ {
		// Каждая следующая штука считается отдельно:
		// это даёт неравномерный рост времени по количеству.
		factor := 1.0 + item.GrowthFactor*math.Log1p(float64(piece-1))
		// Нижняя граница защищает от слишком агрессивного ускорения.
		if factor < 0.60 {
			factor = 0.60
		}
		onePiece := time.Duration(float64(base) * factor * complexity)
		total += onePiece
	}
	return total
}

func scaleDurationToRange(
	base time.Duration,
	minFactor float64,
	maxFactor float64,
) DurationRange {
	const fallbackMin = 0.90
	const fallbackMax = 1.25

	if minFactor <= 0 {
		minFactor = fallbackMin
	}
	if maxFactor <= 0 {
		maxFactor = fallbackMax
	}
	if minFactor > maxFactor {
		minFactor, maxFactor = maxFactor, minFactor
	}

	min := time.Duration(float64(base) * minFactor)
	max := time.Duration(float64(base) * maxFactor)
	return DurationRange{Min: min, Max: max}
}

func minutesToRange(minMins int, maxMins int) DurationRange {
	if minMins < 0 {
		minMins = 0
	}
	if maxMins < 0 {
		maxMins = 0
	}
	if minMins > maxMins {
		minMins, maxMins = maxMins, minMins
	}

	return DurationRange{
		Min: time.Duration(minMins) * time.Minute,
		Max: time.Duration(maxMins) * time.Minute,
	}
}

func addRanges(a DurationRange, b DurationRange) DurationRange {
	return DurationRange{
		Min: a.Min + b.Min,
		Max: a.Max + b.Max,
	}
}

// subtractFixedFromRange нужен, чтобы получить диапазон ожидания очереди из total:
// queue_wait = total - fixed_time_of_current_order.
func subtractFixedFromRange(r DurationRange, fixed time.Duration) DurationRange {
	min := r.Min - fixed
	max := r.Max - fixed
	if min < 0 {
		min = 0
	}
	if max < 0 {
		max = 0
	}
	if min > max {
		min, max = max, min
	}
	return DurationRange{Min: min, Max: max}
}

// applyRestaurantHoursToRange переносит диапазон через график работы ресторана.
// Если ресторан закрыт сейчас или закроется в процессе приготовления, функция добавит паузу.
func applyRestaurantHoursToRange(
	totalRange DurationRange,
	config entity.OrdersQueueConfigTable,
	now time.Time,
) DurationRange {
	openMins, closeMins, ok := parseRestaurantHours(config.RestaurantOpenAt, config.RestaurantCloseAt)
	if !ok {
		return totalRange
	}

	return DurationRange{
		Min: projectDurationThroughWorkingHours(totalRange.Min, now, openMins, closeMins),
		Max: projectDurationThroughWorkingHours(totalRange.Max, now, openMins, closeMins),
	}
}

// parseRestaurantHours парсит формат HH:MM -> минуты от начала дня.
func parseRestaurantHours(openAt string, closeAt string) (int, int, bool) {
	openMins, ok := parseClockMinutes(openAt)
	if !ok {
		return 0, 0, false
	}
	closeMins, ok := parseClockMinutes(closeAt)
	if !ok {
		return 0, 0, false
	}
	return openMins, closeMins, true
}

// parseClockMinutes разбирает "15:04" и возвращает минуту суток.
func parseClockMinutes(value string) (int, bool) {
	t, err := time.Parse("15:04", value)
	if err != nil {
		return 0, false
	}
	return t.Hour()*60 + t.Minute(), true
}

// projectDurationThroughWorkingHours считает, сколько "реального" времени пройдет до готовности
// с учетом того, что кухня работает только в интервалах open/close.
func projectDurationThroughWorkingHours(
	workDuration time.Duration,
	now time.Time,
	openMins int,
	closeMins int,
) time.Duration {
	if workDuration <= 0 {
		return 0
	}

	// open == close означает круглосуточную работу.
	if openMins == closeMins {
		return workDuration
	}

	current := now
	remaining := workDuration
	elapsed := time.Duration(0)

	for remaining > 0 {
		if !isOpenAt(current, openMins, closeMins) {
			nextOpen := nextOpenTime(current, openMins, closeMins)
			wait := nextOpen.Sub(current)
			elapsed += wait
			current = nextOpen
		}

		closeAt := currentCloseTime(current, openMins, closeMins)
		available := closeAt.Sub(current)
		if available <= 0 {
			// Защита от потенциального зацикливания.
			current = current.Add(time.Minute)
			elapsed += time.Minute
			continue
		}

		step := remaining
		if step > available {
			step = available
		}

		remaining -= step
		elapsed += step
		current = current.Add(step)
	}

	return elapsed
}

// isOpenAt проверяет, открыт ли ресторан в конкретный момент.
// Поддерживает и обычный график (10:00-22:00), и ночной (20:00-05:00).
func isOpenAt(t time.Time, openMins int, closeMins int) bool {
	minutes := t.Hour()*60 + t.Minute()
	if openMins < closeMins {
		return minutes >= openMins && minutes < closeMins
	}
	return minutes >= openMins || minutes < closeMins
}

// nextOpenTime возвращает ближайший момент открытия после времени t.
func nextOpenTime(t time.Time, openMins int, closeMins int) time.Time {
	if isOpenAt(t, openMins, closeMins) {
		return t
	}

	location := t.Location()
	dayStart := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, location)
	openToday := dayStart.Add(time.Duration(openMins) * time.Minute)

	if openMins < closeMins {
		if t.Before(openToday) {
			return openToday
		}
		return openToday.Add(24 * time.Hour)
	}

	// Ночная смена: закрытый разрыв только между close и open текущего дня.
	return openToday
}

// currentCloseTime возвращает ближайший момент закрытия для текущего открытого окна.
func currentCloseTime(t time.Time, openMins int, closeMins int) time.Time {
	location := t.Location()
	dayStart := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, location)

	if openMins < closeMins {
		return dayStart.Add(time.Duration(closeMins) * time.Minute)
	}

	// Ночная смена.
	minutes := t.Hour()*60 + t.Minute()
	if minutes >= openMins {
		return dayStart.Add(24*time.Hour + time.Duration(closeMins)*time.Minute)
	}
	return dayStart.Add(time.Duration(closeMins) * time.Minute)
}
