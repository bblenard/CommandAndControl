package types

import (
	"github.com/google/uuid"
)

const ServerAddr = "http://127.0.0.1:8888"

type Client struct {
	ID      string
	Details SystemDetails
}

type SystemDetails struct {
	Hostname string `json: Hostname`
}

func (c *Client) New() {
	c.ID = uuid.New().String()
}
