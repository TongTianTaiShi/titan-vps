package orders

import (
	"context"
	"reflect"

	"github.com/filecoin-project/go-statemachine"
	"golang.org/x/xerrors"
)

// Plan prepares a plan for asset pulling
func (m *Manager) Plan(events []statemachine.Event, user interface{}) (interface{}, uint64, error) {
	next, processed, err := m.plan(events, user.(*OrderInfo))
	if err != nil || next == nil {
		return nil, processed, nil
	}

	return func(ctx statemachine.Context, si OrderInfo) error {
		err := next(ctx, si)
		if err != nil {
			log.Errorf("unhandled error (%s): %+v", si.OrderID, err)
			return nil
		}

		return nil
	}, processed, nil
}

// maps asset states to their corresponding planner functions
var planners = map[OrderState]func(events []statemachine.Event, state *OrderInfo) (uint64, error){
	// external import
	Created: planOne(
		on(WaitingPaymentSent{}, WaitingPayment),
	),
	WaitingPayment: planOne(
		on(OrderTimeOut{}, Done),
		on(OrderCancel{}, Done),
		on(PaymentSucceed{}, BuyGoods),
		apply(PaymentResult{}),
	),
	BuyGoods: planOne(
		on(BuySucceed{}, Done),
	),
	Done: planOne(),
}

// plan creates a plan for the next asset pulling action based on the given events and asset state
func (m *Manager) plan(events []statemachine.Event, state *OrderInfo) (func(statemachine.Context, OrderInfo) error, uint64, error) {
	log.Debugf("state:%s , events:%v", state.State, events)
	p := planners[state.State]
	if p == nil {
		if len(events) == 1 {
			if _, ok := events[0].User.(globalMutator); ok {
				p = planOne() // in case we're in a really weird state, allow restart / update state / remove
			}
		}

		if p == nil {
			return nil, 0, xerrors.Errorf("planner for state %s not found", state.State)
		}
	}

	processed, err := p(events, state)
	if err != nil {
		return nil, processed, xerrors.Errorf("running planner for state %s failed: %w", state.State, err)
	}

	log.Debugf("%s: %s", state.OrderID, state.State)

	switch state.State {
	// Happy path
	case Created:
		return m.handleCreated, processed, nil
	case WaitingPayment:
		return m.handleWaitingPayment, processed, nil
	case BuyGoods:
		return m.handleBuyGoods, processed, nil
	case Done:
		return m.handleDone, processed, nil
	// Fatal errors
	default:
		log.Errorf("unexpected asset update state: %s", state.State)
	}

	return nil, processed, nil
}

// prepares a single plan for a given asset state, allowing for one event at a time
func planOne(ts ...func() (mut mutator, next func(info *OrderInfo) (more bool, err error))) func(events []statemachine.Event, state *OrderInfo) (uint64, error) {
	return func(events []statemachine.Event, state *OrderInfo) (uint64, error) {
	eloop:
		for i, event := range events {
			if gm, ok := event.User.(globalMutator); ok {
				gm.applyGlobal(state)
				return uint64(i + 1), nil
			}

			for _, t := range ts {
				mut, next := t()

				if reflect.TypeOf(event.User) != reflect.TypeOf(mut) {
					continue
				}

				if err, isErr := event.User.(error); isErr {
					log.Warnf("asset %s got error event %T: %+v", state.OrderID, event.User, err)
				}

				event.User.(mutator).apply(state)
				more, err := next(state)
				if err != nil || !more {
					return uint64(i + 1), err
				}

				continue eloop
			}

			_, ok := event.User.(Ignorable)
			if ok {
				continue
			}

			return uint64(i + 1), xerrors.Errorf("planner for state %s received unexpected event %T (%+v)", state.State, event.User, event)
		}

		return uint64(len(events)), nil
	}
}

// on is a utility function to handle state transitions
func on(mut mutator, next OrderState) func() (mutator, func(*OrderInfo) (bool, error)) {
	return func() (mutator, func(*OrderInfo) (bool, error)) {
		return mut, func(state *OrderInfo) (bool, error) {
			state.State = next
			return false, nil
		}
	}
}

// apply like `on`, but doesn't change state
func apply(mut mutator) func() (mutator, func(*OrderInfo) (bool, error)) {
	return func() (mutator, func(*OrderInfo) (bool, error)) {
		return mut, func(state *OrderInfo) (bool, error) {
			return true, nil
		}
	}
}

// initStateMachines init all asset state machines
func (m *Manager) initStateMachines(ctx context.Context) error {
	// initialization
	defer m.stateMachineWait.Done()

	list, err := m.ListAssets()
	if err != nil {
		return err
	}

	for _, order := range list {
		if err := m.orderStateMachines.Send(order.OrderID, OrderRestart{}); err != nil {
			log.Errorf("initStateMachines asset send %s , err %s", order.OrderID, err.Error())
			continue
		}

		m.recoverOutstandingOrders(order)
	}

	return nil
}

// ListAssets load asset pull infos from state machine
func (m *Manager) ListAssets() ([]OrderInfo, error) {
	var list []OrderInfo
	if err := m.orderStateMachines.List(&list); err != nil {
		return nil, err
	}

	return list, nil
}
