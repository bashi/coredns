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
	duration := 50 * time.Millisecond

	for c.Next() {
		args := c.RemainingArgs()
		switch len(args) {
		case 0:
		case 1:
			dur, err := time.ParseDuration(args[0])
			if err != nil {
				return plugin.Error("delay", fmt.Errorf("invalid duration: %q", args[0]))
			}
			duration = dur
		default:
			return plugin.Error("delay", c.ArgErr())
		}
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return Delay{duration: duration, Next: next}
	})

	return nil
}
