package tradeClient

import (
	"context"
	"fmt"

	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"

	"github.com/NguyenHiu/lightning-exchange/tradeApp"
)

type TradeChannel struct {
	ch *client.Channel
}

// newTradeChannel creates a new verify app channel.
func newTradeChannel(ch *client.Channel) *TradeChannel {
	return &TradeChannel{ch: ch}
}

func (g *TradeChannel) SendNewTrades(trades []*tradeApp.Trade, bidOrder, askOrder *tradeApp.Order, isBid bool) {
	err := g.ch.UpdateBy(context.TODO(), func(state *channel.State) error {
		_app, ok := state.App.(*tradeApp.TradeApp)
		if !ok {
			return fmt.Errorf("invalid app type: %T", _app)
		}

		return _app.SendNewTrades(state, trades, bidOrder, askOrder, isBid)
	})
	if err != nil {
		panic(err) // We panic on error to keep the code simple.
	}
}

// Settle settles the app channel and withdraws the funds.
func (g *TradeChannel) Settle() {
	// Channel should be finalized through last ("winning") move.
	// No need to set `isFinal` here.
	err := g.ch.Settle(context.TODO(), false)
	if err != nil {
		panic(err)
	}

	// Cleanup.
	g.ch.Close()
}
