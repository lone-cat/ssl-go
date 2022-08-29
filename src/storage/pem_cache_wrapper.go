package storage

import (
	"encoding/pem"
	"errors"
)

type pemCache struct {
	cache []*pem.Block
	store Pem
}

func WrapInCache(store Pem) (wrapper *pemCache, err error) {
	if store == nil {
		err = errors.New(`nil pem storage passed`)
		return
	}

	wrapper = &pemCache{
		store: store,
	}

	return
}

func (c *pemCache) Load() (data []*pem.Block, err error) {
	if c.cache == nil {
		var blocks []*pem.Block
		blocks, err = c.store.Load()
		if err != nil {
			return
		}
		c.cache = blocks
	}

	if c.cache == nil {
		return
	}

	data = make([]*pem.Block, len(c.cache))
	copy(data, c.cache)

	return
}

func (c *pemCache) Save(data []*pem.Block) (err error) {
	err = c.store.Save(data)
	if err != nil {
		return
	}

	c.cache = make([]*pem.Block, len(data))
	copy(c.cache, data)

	return
}

func (c *pemCache) Delete() (err error) {
	err = c.store.Delete()
	if err != nil {
		return
	}

	c.cache = nil

	return
}
