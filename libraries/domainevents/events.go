package domainevents

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/reecerussell/distro-blog/libraries/database"
	"github.com/reecerussell/distro-blog/libraries/result"
)

var (
	mu       = sync.RWMutex{}
	handlers = make(map[reflect.Type]EventHandler)
)

type Event interface{}

// EventHandler is an interface used to implement specific event handlers.
type EventHandler interface {
	Invoke(ctx context.Context, tx *database.Transaction, e interface{}) result.Result
}

// RegisterEventHandler registers a mapping between an event and its handler.
func RegisterEventHandler(e Event, h EventHandler) {
	mu.RLock()
	defer mu.RUnlock()

	handlers[reflect.TypeOf(e)] = h
}

type Aggregate struct {
	raisedEvents []interface{}
}

// GetRaisedEvents returns a non-nil []interface{} of the aggregate's raise events.
func (a *Aggregate) GetRaisedEvents() []interface{} {
	if a.raisedEvents == nil {
		return []interface{}{}
	}

	return a.raisedEvents
}

// RaiseEvent appends an event to the Aggregate, assuming there is a handler mapped.
func (a *Aggregate) RaiseEvent(e Event) {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := handlers[reflect.TypeOf(e)]; !ok {
		panic(fmt.Errorf("no handler registered for event type '%s'", reflect.TypeOf(e)))
	}

	a.raisedEvents = append(a.raisedEvents, e)
}

// DispatchEvents synchronously executes each raise event from the
// aggregate, with the given transaction.
func (a *Aggregate) DispatchEvents(ctx context.Context, tx *database.Transaction) result.Result {
	mu.Lock()
	defer mu.Unlock()

	for _, e := range a.raisedEvents {
		h := handlers[reflect.TypeOf(e)]

		res := h.Invoke(ctx, tx, e)
		if !res.IsOk(){
			return res
		}
	}

	return result.Ok()
}