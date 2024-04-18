package token

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Token struct {
	rdb          redis.Cmdable
	opts         *option
	tokenListKey string
	tokenKey     string
}

func New(rdb redis.Cmdable, opts ...Option) *Token {
	o := new(option)
	o.apply(opts...)
	o.Default()

	return &Token{
		rdb:          rdb,
		opts:         o,
		tokenListKey: "users:token:list:%d",
		tokenKey:     "users:token",
	}
}

// listKey gets the key of the user's token list
func (tk *Token) listKey(userId int64) string {
	return fmt.Sprintf(tk.tokenListKey, userId)
}

// key gets the key of the token
func (tk *Token) key(token string) string {
	return fmt.Sprintf("%s:%s", tk.tokenKey, token)
}

// Set sets the token
func (tk *Token) Set(ctx context.Context, val *Value) (err error) {
	listKey := tk.listKey(val.UserId)
	key := tk.key(val.Token)
	if tk.opts.expire > 0 {
		val.ExpiredAt = val.CreatedAt.Add(tk.opts.expire)
	}

	_, err = tk.rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) (err error) {
		err = pipe.LPush(ctx, listKey, val.Token).Err()
		if err != nil {
			return
		}
		err = pipe.HSet(ctx, key, val.Token, val).Err()
		return
	})
	return
}

// Remove removes the token
func (tk *Token) Remove(ctx context.Context, userId int64, tokens ...string) (err error) {
	listKey := tk.listKey(userId)
	if len(tokens) <= 0 {
		tokens, err = tk.rdb.LRange(ctx, listKey, 0, -1).Result()
		if err != nil {
			return
		}
	}

	for _, token := range tokens {
		err = tk.rdb.LRem(ctx, listKey, 0, token).Err()
		if err != nil {
			return
		}
		err = tk.rdb.Del(ctx, tk.key(token)).Err()
		if err != nil {
			return
		}
	}
	return
}

// Get gets the token
func (tk *Token) Get(ctx context.Context, token string) (val *Value, err error) {
	key := tk.key(token)
	val = new(Value)
	err = tk.rdb.HGet(ctx, key, token).Scan(val)
	if err == redis.Nil {
		err = nil
	}
	return
}
