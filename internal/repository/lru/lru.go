package lru

import (
	"time"

	"github.com/karlseguin/ccache/v2"
)

type lru struct {
	lru *ccache.Cache
}

type LRU interface {
	Set(key string, value interface{})
	Get(key string) (interface{}, bool)
}

func New(opts ...Options) LRU {
	o := &options{}

	for _, option := range opts {
		option(o)
	}

	cache := ccache.New(o.configure())

	l := &lru{
		lru: cache,
	}

	return l
}

func (l *lru) Get(key string) (interface{}, bool) {
	item := l.lru.Get(key)

	if item == nil || item.Expired() {
		return nil, false
	}

	return item.Value(), true
}

func (l *lru) Set(key string, value interface{}) {
	l.lru.Set(key, value, time.Hour*24*10)
}
