package events

import (
	"context"
	"waza/events/topics"
	"waza/models"
	"waza/setup"
	"waza/store"
	"waza/utils"
)

type EventHandlers struct {
	opts          *setup.ServiceDependencies
	subscriptions []func() error
}

func NewEventHandler(opts *setup.ServiceDependencies) EventHandlers {
	return EventHandlers{opts: opts}
}

func (e EventHandlers) Listen() {
	e.subscriptions = append(
		e.subscriptions,
		e.handleUserCreated,
	)

	for _, sub := range e.subscriptions {
		if err := sub(); err != nil {
			e.opts.Logger.WithError(err).Fatal("failed to mount subscription")
		}
	}

	<-make(chan bool)
}

func (e EventHandlers) handleUserCreated() error {
	return e.opts.EventStore.Subscribe(topics.UserCreated, func(event store.Event) error {
		e.opts.Logger.Infof("received event on %s", topics.UserCreated)
		ctx := context.Background()
		var user models.User
		if err := utils.UnPack(event.Data(), &user); err != nil {
			e.opts.Logger.WithError(err).Errorf("failed to decode event payload for %s", event.Topic())
		}

		if _, err := e.opts.AccountService.CreateAccount(ctx, user); err != nil {
			e.opts.Logger.WithError(err).Errorf("failed to create account from payload for %s", event.Topic())
		}

		return nil
	})
}
