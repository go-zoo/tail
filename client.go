package tail

import (
	"time"
)

type Client struct {
	ID     string
	Valid  bool
	TTL    time.Duration
	Expire *time.Timer
}

func NewClient(id string, ttl time.Duration) *Client {
	c := &Client{ID: id, Valid: true, TTL: ttl, Expire: time.NewTimer(ttl)}
	go c.watch()
	return c
}

func (c *Client) watch() {
	<-c.Expire.C
	c.Valid = false
}

func (c *Client) renewClient() {
	c.Expire.Reset(c.TTL)
	c.Valid = true
}

func (c *Client) delClient() {
	c.Valid = false
}
