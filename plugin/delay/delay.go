package delay

import (
	"context"
	"time"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
)

type Delay struct {
	duration time.Duration
	qtype uint16
	name string
}

type Handler struct {
	Delays []Delay

	Next plugin.Handler
}

func (h Handler) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	for _, delay := range h.Delays {
		if delay.qtype != dns.TypeANY && delay.qtype != state.QType() {
			continue
		}

		if len(delay.name) > 0 && delay.name != state.Name() {
			continue
		}

		time.Sleep(delay.duration)
	}

	return plugin.NextOrFailure(h.Name(), h.Next, ctx, w, r)
}

func (h Handler) Name() string { return "delay" }
