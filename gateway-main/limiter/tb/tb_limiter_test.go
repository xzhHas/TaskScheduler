package tb

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-redis/redis"
)

func TestTBLimiter_Allow(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "",
		Password: "",
	})
	tbLimiter, err := NewTBLimiter(context.TODO(), redisClient)
	if err != nil {
		panic(err)
	}

	tbLimit := &TBLimit{Rate: 100, Burst: 100, Expire: 10}
	result, err := tbLimiter.Allow(context.TODO(), "/foo", tbLimit)

	if err != nil {
		panic(err)
	}

	if result.Allowed > 0 {
		fmt.Printf("limit allowed: %+v, remining: %+v\n", result.Allowed, result.Remaining)
	} else {
		fmt.Printf("limit not allowed: %+v, remining: %+v\n", result.Allowed, result.Remaining)
	}

}
