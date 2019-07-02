package discovery

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/bilibili/discovery/conf"
	"github.com/bilibili/discovery/registry"
	http "github.com/bilibili/kratos/pkg/net/http/blademaster"
)

// Discovery discovery.
type Discovery struct {
	c         *conf.Config
	protected bool
	client    *http.Client
	registry  *registry.Registry
	nodes     atomic.Value
}

// New get a discovery.
func New(c *conf.Config) (d *Discovery, cancel context.CancelFunc) {
	d = &Discovery{
		c:        c,
		client:   http.NewClient(c.HTTPClient),
		registry: registry.NewRegistry(c),
	}
	d.nodes.Store(registry.NewNodes(c))
	d.syncUp()
	cancel = d.regSelf()
	go d.nodesproc()
	go d.exitProtect()
	return
}

func (d *Discovery) exitProtect() {
	time.Sleep(time.Second * 90)
	d.protected = false
}
