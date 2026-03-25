package eventmanager

import "sync"

type Event struct {
	Name    string
	Payload any
}

type EventManager struct {
	events map[string][]func(payload any)
	mu     sync.RWMutex
	start  sync.Once

	queue chan Event
}

var (
	eventManagerInstance *EventManager
	eventManagerOnce     sync.Once
)

func New() *EventManager {
	eventManagerOnce.Do(func() {
		eventManagerInstance = &EventManager{
			events: make(map[string][]func(payload any)),
			queue:  make(chan Event, 2),
		}
	})

	return eventManagerInstance
}

func (em *EventManager) AddEventListener(eventName string, listener func(payload any)) {
	em.mu.Lock()
	defer em.mu.Unlock()

	em.events[eventName] = append(em.events[eventName], listener)
}

func (em *EventManager) AddEvent(event Event) {
	em.queue <- event
}

func (em *EventManager) EmitEvent(eventName string, payload any) {
	em.AddEvent(Event{
		Name:    eventName,
		Payload: payload,
	})
}

func (em *EventManager) Start() {
	em.start.Do(func() {
		go func() {
			for event := range em.queue {
				em.mu.RLock()
				listeners := append([]func(payload any){}, em.events[event.Name]...)
				em.mu.RUnlock()

				for _, listener := range listeners {
					if listener == nil {
						continue
					}
					listener(event.Payload)
				}
			}
		}()
	})
}
