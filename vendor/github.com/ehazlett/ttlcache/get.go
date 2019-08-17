package ttlcache

import "time"

func (t *TTLCache) Get(key string) *KV {
	if k, ok := t.data[key]; ok {
		return &KV{
			Key:   key,
			Value: k.Value,
			TTL:   t.ttl - (time.Since(k.updated)),
		}
	}

	return nil
}
