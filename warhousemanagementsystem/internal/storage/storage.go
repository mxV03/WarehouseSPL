package storage

import "github.com/mxV03/warhousemanagementsystem/ent"

var Client *ent.Client

func SetClient(c *ent.Client) {
	Client = c
}
