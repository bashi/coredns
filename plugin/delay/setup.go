package delay

import (
	"fmt"
	"time"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
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
		switch len(args) {
		case 0:
			delay.duration = 50 * time.Millisecond
		case 1:
			dur, err := time.ParseDuration(args[0])
			if err != nil {
				return handler, plugin.Error("delay", fmt.Errorf("invalid duration: %q", args[0]))
			}
			delay.duration = dur
		default:
			return handler, plugin.Error("delay", c.ArgErr())
		}
		handler.Delays = append(handler.Delays, delay)
	}

	return
}