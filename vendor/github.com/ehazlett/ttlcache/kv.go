package ttlcache

import "time"

type KV struct {
	Key   string
	Value interface{}
	TTL   time.Duration
}
