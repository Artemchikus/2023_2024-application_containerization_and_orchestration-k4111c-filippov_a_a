package metrix

import "sync"

type (
	CurrentRequests struct {
		mx   sync.RWMutex
		data int64
	}
)

func mustCurrentRequests() *CurrentRequests {
	return &CurrentRequests{
		mx: sync.RWMutex{},
	}
}

func (c *CurrentRequests) Inc() {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.data++
}

func (c *CurrentRequests) Dec() {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.data--
}

func (c *CurrentRequests) Get() *CurrentReqData {
	return &CurrentReqData{
		CurrentReq: c.data,
	}
}
