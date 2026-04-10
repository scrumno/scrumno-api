package order_queue

// import (
// 	"math"
// 	"time"

// 	"github.com/google/uuid"
// )

// // MockEstimateQueueTime пример с моковыми данными.
// func MockEstimateQueueTime() QueueEstimationResult {
// 	config := QueueEstimationConfig{
// 		ID:                    uuid.New(),
// 		KitchenParallelSlots:  2,
// 		QueueGrowthFactor:     0.18,
// 		OrderReserveMinutes:   2,
// 		RestaurantOpenAt:      "10:00",
// 		RestaurantCloseAt:     "22:00",
// 		EmptyQueueWaitMinMins: 10,
// 		EmptyQueueWaitMaxMins: 12,
// 		QueueTimeMinFactor:    0.90,
// 		QueueTimeMaxFactor:    1.25,
// 		CreatedAt:             time.Now(),
// 		UpdatedAt:             time.Now(),
// 	}

// 	currentOrderID := uuid.New()
// 	current := QueueEstimationOrderMock{
// 		ID:           currentOrderID,
// 		ExternalID:   "CURRENT-001",
// 		SetupMinutes: 3,
// 		Items: []QueueEstimationItemMock{
// 			{
// 				ID:               uuid.New(),
// 				OrderID:          currentOrderID,
// 				ProductID:        101,
// 				Quantity:         3,
// 				BaseCookMinutes:  6,
// 				GrowthFactor:     0.20,
// 				ComplexityFactor: 1.00,
// 				CreatedAt:        time.Now(),
// 				UpdatedAt:        time.Now(),
// 			},
// 			{
// 				ID:               uuid.New(),
// 				OrderID:          currentOrderID,
// 				ProductID:        205,
// 				Quantity:         2,
// 				BaseCookMinutes:  4,
// 				GrowthFactor:     0.08,
// 				ComplexityFactor: 1.10,
// 				CreatedAt:        time.Now(),
// 				UpdatedAt:        time.Now(),
// 			},
// 		},
// 		CreatedAt: time.Now(),
// 		UpdatedAt: time.Now(),
// 	}

// 	aheadOrder1ID := uuid.New()
// 	aheadOrder2ID := uuid.New()
// 	ordersAhead := []QueueEstimationOrderMock{
// 		{
// 			ID:           aheadOrder1ID,
// 			ExternalID:   "AHEAD-001",
// 			SetupMinutes: 2,
// 			Items: []QueueEstimationItemMock{
// 				{
// 					ID:               uuid.New(),
// 					OrderID:          aheadOrder1ID,
// 					ProductID:        103,
// 					Quantity:         4,
// 					BaseCookMinutes:  3,
// 					GrowthFactor:     0.10,
// 					ComplexityFactor: 1.00,
// 					CreatedAt:        time.Now(),
// 					UpdatedAt:        time.Now(),
// 				},
// 			},
// 			CreatedAt: time.Now(),
// 			UpdatedAt: time.Now(),
// 		},
// 		{
// 			ID:           aheadOrder2ID,
// 			ExternalID:   "AHEAD-002",
// 			SetupMinutes: 1,
// 			Items: []QueueEstimationItemMock{
// 				{
// 					ID:               uuid.New(),
// 					OrderID:          aheadOrder2ID,
// 					ProductID:        304,
// 					Quantity:         2,
// 					BaseCookMinutes:  7,
// 					GrowthFactor:     0.25,
// 					ComplexityFactor: 1.15,
// 					CreatedAt:        time.Now(),
// 					UpdatedAt:        time.Now(),
// 				},
// 			},
// 			CreatedAt: time.Now(),
// 			UpdatedAt: time.Now(),
// 		},
// 	}

// 	return EstimateQueueTime(current, ordersAhead, config)
// }
