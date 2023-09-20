package delay

import (
	"context"
	"time"

	"github.com/coredns/coredns/plugin"

	"github.com/miekg/dns"
)

type Delay struct {
	duration time.Duration
}

type Handler struct {
	Delays []Delay

	Next plugin.Handler
}

func (h Handler) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	for _, delay := range h.Delays {
		time.Sleep(delay.duration)
	}
	return plugin.NextOrFailure(h.Name(), h.Next, ctx, w, r)
}

func (h Handler) Name() string { return "delay" }
