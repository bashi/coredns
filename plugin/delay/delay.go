package delay

import (
	"context"
	"time"

	"github.com/coredns/coredns/plugin"

	"github.com/miekg/dns"
)

type Delay struct {
	duration time.Duration
	Next plugin.Handler
}

func (d Delay) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	time.Sleep(d.duration)
	return plugin.NextOrFailure(d.Name(), d.Next, ctx, w, r)
}

func (d Delay) Name() string { return "delay" }
