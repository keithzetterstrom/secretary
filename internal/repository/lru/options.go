package lru

import (
	"github.com/AlekSi/pointer"
	"github.com/karlseguin/ccache/v2"
)

var (
	defaultMaxSize = pointer.ToInt64(1_000_000)
	defaultBuckets = pointer.ToUint32(16)
)

type options struct {
	maxSize *int64
	buckets *uint32
}

type Options func(o *options)

func (o *options) configure() *ccache.Configuration {
	cfg := ccache.Configure()

	if o == nil {
		return cfg
	}

	if o.maxSize == nil {
		o.maxSize = defaultMaxSize
	}

	if o.buckets == nil {
		o.buckets = defaultBuckets
	}

	return cfg.
		MaxSize(*o.maxSize).
		Buckets(*o.buckets)
}

func WithMaxSize(maxSize int64) Options {
	return func(o *options) {
		o.maxSize = &maxSize
	}
}

func WithBucketsCount(buckets uint32) Options {
	return func(o *options) {
		o.buckets = &buckets
	}
}
