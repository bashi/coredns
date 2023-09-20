package delay

import (
	"fmt"
	"time"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"

	"github.com/miekg/dns"
)

func init() { plugin.Register("delay", setup) }

func setup(c *caddy.Controller) error {
	handler, err := delayParse(c)
	if err != nil {
		return plugin.Error("delay", err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		handler.Next = next
		return handler
	})

	return nil
}

func delayParse(c *caddy.Controller) (handler Handler, err error) {
	handler.Delays = make([]Delay, 0)

	for c.Next() {
		delay := Delay{}
		args := c.RemainingArgs()
		nargs := len(args)
		if (nargs == 0) {
			delay.duration = 50 * time.Millisecond
		} else {
			dur, err := time.ParseDuration(args[0])
			if err != nil {
				return handler, plugin.Error("delay", fmt.Errorf("invalid duration: %q", args[0]))
			}
			delay.duration = dur
		}

		if (nargs >= 2) {
			qtype, ok := dns.StringToType[args[1]]
			if !ok {
				return handler, c.Errf("invalid RR class %s", args[1])
			}
			delay.qtype = qtype
		} else {
			delay.qtype = dns.TypeANY
		}

		if (nargs >= 3) {
			return handler, plugin.Error("delay", c.ArgErr())
		}

		handler.Delays = append(handler.Delays, delay)
	}

	return
}