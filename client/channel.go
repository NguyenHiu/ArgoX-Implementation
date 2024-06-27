package client

import (
	"context"
	"fmt"

	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"

	"github.com/NguyenHiu/lightning-exchange/app"
	"github.com/google/uuid"
)

type VerifyChannel struct {
	ch *client.Channel
}

// newVerifyChannel creates a new verify app channel.
func newVerifyChannel(ch *client.Channel) *VerifyChannel {
	return &VerifyChannel{ch: ch}
}

func (g *VerifyChannel) SendNewOrder(order *app.Order) {
	err := g.ch.UpdateBy(context.TODO(), func(state *channel.State) error {
		app, ok := state.App.(*app.VerifyApp)
		if !ok {
			return fmt.Errorf("invalid app type: %T", app)
		}

		return app.SendNewOrder(state, order)
	})
	if err != nil {
		panic(err) // We panic on error to keep the code simple.
	}
}

func (g *VerifyChannel) UpdateExistedOrder(orderID uuid.UUID, updatedData app.OrderUpdatedInfo) {
	err := g.ch.UpdateBy(context.TODO(), func(state *channel.State) error {
		app, ok := state.App.(*app.VerifyApp)
		if !ok {
			return fmt.Errorf("invalid app type: %T", app)
		}

		return app.UpdateExistedOrder(state, orderID, updatedData)
	})
	if err != nil {
		panic(err) // We panic on error to keep the code simple.
	}
}

func (g *VerifyChannel) SendNewMessage(message *app.Message) {
	err := g.ch.UpdateBy(context.TODO(), func(state *channel.State) error {
		app, ok := state.App.(*app.VerifyApp)
		if !ok {
			return fmt.Errorf("invalid app type: %T", app)
		}

		return app.SendNewMessage(state, message)
	})
	if err != nil {
		panic(err) // We panic on error to keep the code simple.
	}
}

// Settle settles the app channel and withdraws the funds.
func (g *VerifyChannel) Settle() {
	// Channel should be finalized through last ("winning") move.
	// No need to set `isFinal` here.
	err := g.ch.Settle(context.TODO(), false)
	if err != nil {
		panic(err)
	}

	// Cleanup.
	g.ch.Close()
}
