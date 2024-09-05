package client

import (
	"context"
	"fmt"

	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"

	"github.com/NguyenHiu/lightning-exchange/app"
)

type VerifyChannel struct {
	ch *client.Channel
}

// newVerifyChannel creates a new verify app channel.
func newVerifyChannel(ch *client.Channel) *VerifyChannel {
	return &VerifyChannel{ch: ch}
}

func (g *VerifyChannel) SendNewOrders(orders []*app.Order) {
	err := g.ch.UpdateBy(context.TODO(), func(state *channel.State) error {
		_app, ok := state.App.(*app.VerifyApp)
		if !ok {
			return fmt.Errorf("invalid app type: %T", _app)
		}

		return _app.SendNewOrders(state, orders)
	})
	if err != nil {
		panic(err) // We panic on error to keep the code simple.
	}
}

func (g *VerifyChannel) SendNewTrades(trades []*app.Trade) {
	err := g.ch.UpdateBy(context.TODO(), func(state *channel.State) error {
		_app, ok := state.App.(*app.VerifyApp)
		if !ok {
			return fmt.Errorf("invalid app type: %T", _app)
		}

		return _app.SendNewTrades(state, trades)
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
