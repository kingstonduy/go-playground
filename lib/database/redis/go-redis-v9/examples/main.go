package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	REDIS_HOST = "10.96.20.152"
	REDIS_PORT = "6379"
)

func NewRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", REDIS_HOST, REDIS_PORT),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rdb
}

type redisRepository struct {
	db *redis.Client
}

func NewRedisRepository() *redisRepository {
	return &redisRepository{
		db: NewRedisClient(),
	}
}

// return a list of key with the given prefix
//
// example repo.GetKeysWithPrefix(context.Background(), "00371000")
func (r *redisRepository) GetKeysWithPrefix(ctx context.Context, prefix string) []string {
	var results []string
	iter := r.db.Scan(ctx, 0, prefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		results = append(results, key)
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}
	return results
}

func (r *redisRepository) Set(ctx context.Context, key string, val string, dur time.Duration) (string, error) {
	return r.db.Set(ctx, key, val, dur).Result()
}

func (r *redisRepository) Get(ctx context.Context, key string) (string, error) {
	return r.db.Get(ctx, key).Result()
}

func (r *redisRepository) Del(ctx context.Context, key string) (int64, error) {
	return r.db.Del(ctx, key).Result()
}

func (r *redisRepository) Pub(ctx context.Context, channel string, message interface{}) (int64, error) {
	return r.db.Publish(ctx, channel, message).Result()
}

func (r *redisRepository) Sub(ctx context.Context, channel string) (string, error) {
	subscriber := r.db.Subscribe(ctx, channel)

	msg, err := subscriber.ReceiveMessage(ctx)
	if err != nil {
		return "", err
	}
	log.Printf("receive from redis message = %s", msg)
	return msg.String(), nil
}

// Set if not exist
//
// key ko tồn tại -> true
//
// key tồn tại -> false
func (r *redisRepository) SetNx(ctx context.Context, key string, value string, dur time.Duration) (bool, error) {
	isNotExist, err := r.db.SetNX(ctx, key, value, 0).Result()
	return isNotExist, err
}

func main() {
	repo := NewRedisRepository()
	ctx := context.Background()
	key := "duydk3"
	val := "duydk3"
	// test set no duration
	repo.db.Set(ctx, key, val, 0)
	//--------------------------------------------------------------------------------------
	// test delete key
	repo.db.Del(ctx, key, val)
	//--------------------------------------------------------------------------------------
	// test set 10s
	repo.db.Set(ctx, key, val, time.Duration(time.Second*10))
	//--------------------------------------------------------------------------------------
	// test GetKeysWithPrefix
	for i := 1; i <= 10; i++ {
		newKey := key + fmt.Sprintf("%d", i)
		repo.db.Set(ctx, newKey, val, time.Duration(time.Second*10))
	}
	l := repo.GetKeysWithPrefix(ctx, key)
	log.Println(l)
	//--------------------------------------------------------------------------------------
	// test publish subcrire
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		_, err := repo.Sub(ctx, key)
		if err != nil {
			log.Println(err)
		}
		defer wg.Done()
	}()

	time.Sleep(time.Second * 3)

	go func() {
		_, err := repo.Pub(ctx, key, val)
		if err != nil {
			log.Println(err)
		}
		defer wg.Done()
	}()

	wg.Wait()
	//--------------------------------------------------------------------------------------

}
