package client

import "perun.network/go-perun/client"

type VerifyChannel struct {
	ch *client.Channel
}

// newVerifyChannel creates a new verify app channel.
func newVerifyChannel(ch *client.Channel) *VerifyChannel {
	return &VerifyChannel{ch: ch}
}
