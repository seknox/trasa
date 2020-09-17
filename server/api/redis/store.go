package redis

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/seknox/trasa/server/consts"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
)

// Set implements reis HSet
func (s redisStore) Set(key string, expiry time.Duration, val ...string) error {
	client := s.RedisClient
	ctx := context.Background()
	err := client.HSet(ctx, key, val).Err()
	if err != nil {
		return err
	}

	return client.Expire(ctx, key, expiry).Err()
}

// Get implements redis HMGet
func (s redisStore) Get(key string, field string) (string, error) {
	client := s.RedisClient
	vals, err := client.HMGet(context.Background(), key, field).Result()
	if err != nil || len(vals) != 1 {
		return "", errors.Errorf("get from redis: %v", err)
	}

	strVals, err := utils.ToStringArr(vals)
	if err != nil {
		return "", errors.Errorf("convert fom interface to string redis: %v", err)
	}
	return strVals[0], nil
}

// MGet implements redis HMGet with multiple fields
func (s redisStore) MGet(key string, field ...string) ([]string, error) {
	client := s.RedisClient
	vals, err := client.HMGet(context.Background(), key, field...).Result()
	if err != nil {
		return nil, err
	}
	return utils.ToStringArr(vals)
}

// GetSession can be used to store any key value. Main key should be unique value while value key name will be "data" and value should be json encoded byte.
func (s redisStore) GetSession(key string) (userID, orgID, deviceID, browserID, auth string, err error) {
	client := s.RedisClient
	ctx := context.Background()
	var vals []interface{}
	vals, err = client.HMGet(ctx, key, "userID", "orgID", "deviceID", "browserID", "auth").Result()
	if err != nil {
		return
	}
	var strArr []string
	strArr, err = utils.ToStringArr(vals)
	if err != nil {
		return
	}

	if len(strArr) != 5 {
		err = errors.Errorf("not enough values")
		return
	}

	//TODO
	client.Expire(ctx, key, time.Second*900)
	userID = strArr[0]
	orgID = strArr[1]
	deviceID = strArr[2]
	browserID = strArr[3]

	auth = strArr[4]

	return
}

func (s redisStore) SetVerifyIntent(key string, expiry time.Duration, intent, field, val string) error {
	client := s.RedisClient
	ctx := context.Background()
	err := client.HMSet(ctx, key, "intent", intent, field, val).Err()
	if err != nil {
		return err
	}

	return client.Expire(ctx, key, expiry).Err()
}

//TODO @sshah check this logic
func (s redisStore) VerifyIntent(key string, intent consts.VerifyTokenIntent) error {
	client := s.RedisClient
	res, err := client.HMGet(context.Background(), key, "intent").Result()
	if err != nil {
		return err
	}

	if len(res) != 1 {
		return errors.Errorf("failed to verify token")
	}
	if string(intent) != res[0] {
		return errors.Errorf("invalid intent")
	}

	return nil

	//return client.HGet(context.Background(), key, field...).String()
}

func (s redisStore) Delete(key string) error {
	client := s.RedisClient
	ctx := context.Background()
	return client.Del(ctx, key).Err()
}

//TODO use pubsub instead of polling or go channels if possible

// WaitForStatusAndGet polls redis to until timeout or status field is true
// then returns a field.
func (s redisStore) WaitForStatusAndGet(key, field string) (success bool, val string) {
	timeout := time.After(60 * time.Second)
	tick := time.Tick(1000 * time.Millisecond)
	//var ret string
	//var err error
	// Keep trying until we're timed out or got a result or got an error
	for {
		select {
		// Got a timeout! fail with a timeout error
		case <-timeout:
			return false, ""
		// Got a tick, we should check on doSomething()
		case <-tick:
			status, _ := s.Get(key, "status")
			if status == "true" {
				success = true
				val, err := Store.Get(key, field)
				if err != nil {
					logrus.Error(err)
					return false, ""
				}
				return true, val
			}
		}

	}

}

// SetHTTPGatewaySession
func (s redisStore) SetHTTPGatewaySession(key, orgusr, authDataVal string, sessionRecord string) error {
	client := s.RedisClient
	ctx := context.Background()
	err := client.HSet(ctx, key, "user", orgusr, "auth", authDataVal, "sessionRecord", sessionRecord).Err()
	if err != nil {
		return err
	}

	return client.Expire(ctx, key, time.Second*900).Err()

}

// GetHTTPGatewaySession
func (s redisStore) GetHTTPGatewaySession(key string) (user, auth, sessionRecord string, err error) {
	client := s.RedisClient
	ctx := context.Background()
	var vals []interface{}
	vals, err = client.HMGet(ctx, key, "user", "auth", "sessionRecord").Result()
	if err != nil {
		return
	}
	var strArr []string
	strArr, err = utils.ToStringArr(vals)
	if err != nil {
		return
	}

	if len(strArr) != 3 {
		return "", "", "", errors.Errorf("not enough values")
	}

	//TODO
	client.Expire(ctx, key, time.Second*900)
	user = strArr[0]
	auth = strArr[1]
	sessionRecord = strArr[2]

	return

}
