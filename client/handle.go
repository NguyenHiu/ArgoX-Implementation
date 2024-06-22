package client

import (
	"context"
	"fmt"

	"github.com/NguyenHiu/lightning-exchange/app"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
)

// HandleProposal is the callback for incoming channel proposals.
func (c *AppClient) HandleProposal(p client.ChannelProposal, r *client.ProposalResponder) {
	lcp, err := func() (*client.LedgerChannelProposal, error) {
		// Ensure that we got a ledger channel proposal.
		lcp, ok := p.(*client.LedgerChannelProposal)
		if !ok {
			return nil, fmt.Errorf("invalid proposal type: %T", p)
		}

		// Ensure the ledger channel proposal includes the expected app.
		if !lcp.App.Def().Equals(c.app.Def()) {
			return nil, fmt.Errorf("invalid app type ")
		}

		// Check that we have the correct number of participants.
		if lcp.NumPeers() != 2 {
			return nil, fmt.Errorf("invalid number of participants: %d", lcp.NumPeers())
		}

		// Check that the channel has the expected assets.
		err := channel.AssetsAssertEqual(lcp.InitBals.Assets, c.currencies)
		if err != nil {
			return nil, fmt.Errorf("invalid assets: %v", err)
		}

		// Check that the channel has the expected assets and funding balances.
		if err := channel.AssetsAssertEqual(lcp.InitBals.Assets, c.currencies); err != nil {
			return nil, fmt.Errorf("invalid assets: %v", err)
		}

		// const assetIdx, peerIdx = 0, 1
		for i := range lcp.NumPeers() {
			for idx := range c.currencies {
				if lcp.FundingAgreement[idx][i].Cmp(c.stakes[idx]) != 0 {
					return nil, fmt.Errorf("invalid funding balance")
				}
			}
		}
		return lcp, nil
	}()
	if err != nil {
		r.Reject(context.TODO(), err.Error()) //nolint:errcheck // It's OK if rejection fails.
	}

	// Create a channel accept message and send it.
	accept := lcp.Accept(
		c.account,                // The account we use in the channel.
		client.WithRandomNonce(), // Our share of the channel nonce.
	)
	ch, err := r.Accept(context.TODO(), accept)
	if err != nil {
		_logger.Error("Error accepting channel proposal: %v\n", err)
		return
	}

	// Start the on-chain event watcher. It automatically handles disputes.
	c.startWatching(ch)

	c.channels <- newVerifyChannel(ch)
}

// HandleUpdate is the callback for incoming channel updates.
func (c *AppClient) HandleUpdate(cur *channel.State, next client.ChannelUpdate, r *client.UpdateResponder) {
	// Perun automatically checks that the transition is valid.
	// We always accept.

	fromData := cur.Data.(*app.VerifyAppData)
	toData := next.State.Data.(*app.VerifyAppData)
	if c.UseTrigger && len(fromData.Orders) < len(toData.Orders) {
		for i := len(fromData.Orders); i < len(toData.Orders); i++ {
			c.TriggerChannel <- toData.Orders[i]
		}
	}

	err := r.Accept(context.TODO())
	if err != nil {
		panic(err)
	}
}

// HandleAdjudicatorEvent is the callback for smart contract events.
func (c *AppClient) HandleAdjudicatorEvent(e channel.AdjudicatorEvent) {
	// _logger.Info("Adjudicator event: type = %T, client = %v\n", e, c.account)
}
