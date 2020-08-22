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

	key := "someKey"

	err := redis.Store.Set(key, time.Minute, "k1", "v1", "k2", "v2")
	if err != nil {
		t.Fatalf("redis.Store.Set() err = %v", err)
		return
	}

	got, err := redis.Store.Get(key, "k1")
	if err != nil {
		t.Fatalf("redis.Store.Set() err = %v", err)
		return
	}

	if got != "v1" {
		t.Errorf(`redis.Store.Get() got=%s want=%s`, got, "v1")
	}

	gotArr, err := redis.Store.MGet(key, "k1", "k2")
	if err != nil {
		t.Fatalf("redis.Store.Mget() err = %v", err)
		return
	}

	if 2 != len(gotArr) {
		t.Errorf(`incorrect number of values. got=%d want=%d`, len(gotArr), 2)
	}

	if gotArr[0] != "v1" && gotArr[1] != "v2" {
		t.Errorf(`redis.Store.MGet() got=%v want=%v`, gotArr, []string{"v1", "v2"})

	}

}

func TestRedisVerifyIntent(t *testing.T) {
	err := redis.Store.SetVerifyIntent("someKey", time.Minute, "someIntent", "someField", "someValue")
	if err != nil {
		t.Fatalf(`failed to set verify token: %v`, err)
	}

	err = redis.Store.VerifyIntent("someKey", "someIntent")
	if err != nil {
		t.Fatalf(`failed to verify intent: %v`, err)
	}

}
