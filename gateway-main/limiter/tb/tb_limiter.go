package tb

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strings"
)

const redisPrefix = "token_bucket_rate:"

// 令牌桶的限流器
type TBLimiter struct {
	client *redis.Client
}

// 令牌桶限流配置
type TBLimit struct {
	Expire int
	Rate   int
	Burst  int
}

type Result struct {
	Allowed   int
	Remaining int
}

func NewTBLimiter(ctx context.Context, client *redis.Client) (*TBLimiter, error) {
	l := &TBLimiter{
		client: client,
	}
	_, err := client.ScriptLoad(ctx, AllowNScript).Result()
	if err != nil {
		return nil, err
	}

	return l, nil
}

func (tb *TBLimiter) Allow(ctx context.Context, key string, limit *TBLimit) (*Result, error) {
	return tb.AllowN(ctx, key, limit, 1)
}

func (tb *TBLimiter) AllowN(ctx context.Context, key string, limit *TBLimit, n int) (*Result, error) {
	values := []interface{}{limit.Rate, limit.Burst, n, limit.Expire}
	r := AllowN.EvalSha(ctx, tb.client, []string{redisPrefix + key}, values...)
	if err := r.Err(); err != nil && strings.HasPrefix(err.Error(), "NOSCRIPT ") {
		r = AllowN.Eval(ctx, tb.client, []string{redisPrefix + key}, values...)
	}
	v, err := r.Result()
	if err != nil {
		return nil, err
	}
	values = v.([]interface{})
	result := &Result{
		Allowed:   int(values[0].(int64)),
		Remaining: int(values[1].(int64)),
	}
	return result, nil
}
