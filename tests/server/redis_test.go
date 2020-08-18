package server

import (
	"github.com/seknox/trasa/server/api/redis"
	"testing"
	"time"
)

func TestRedisSet(t *testing.T) {
	type args struct {
		key    string
		expiry time.Duration
		val    []string
	}

	key := "1"
	val := []string{"k1", "v1", "k2", "v2", "k3", "v3"}

	err := redis.Store.Set(key, time.Minute, val...)
	if err != nil {
		t.Fatalf("redis.Store.Set() err = %v", err)
		return
	}

	got, err := redis.Store.Get(key, val[0])
	if err != nil {
		t.Fatalf("redis.Store.Set() err = %v", err)
		return
	}

	if got != val[1] {
		t.Errorf(`redis.Store.Get() got=%s want=%s`, got, val[1])
	}

	gotArr, err := redis.Store.MGet(key, val...)
	if err != nil {
		t.Fatalf("redis.Store.Mget() err = %v", err)
		return
	}

	if len(val) != len(gotArr) {
		t.Errorf(`incorrect number of values. got=%d want=%d`, len(val), len(gotArr))
	}

}
