package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
)

var rClient *redis.Client

func init() {
	rClient = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
}

func main() {
	ctx := context.Background()

	// clean up
	if err := rClient.FlushAll(ctx).Err(); err != nil {
		log.Fatalf("FlushAll error: %v", err)
	}

	// sorted sets
	// set members and scores in redis
	key1 := "rank1"
	for i := 1; i <= 1000; i++ {
		score := rand.Intn(1000) * 100
		if err := rClient.ZAdd(ctx, key1, &redis.Z{
			Score:  float64(score),
			Member: i,
		}).Err(); err != nil {
			log.Fatalf("ZAdd error: %v", err)
		}
	}

	// get member count
	count, err := rClient.ZCard(ctx, key1).Result()
	if err != nil {
		log.Fatalf("ZScore error: %v", err)
	}
	log.Printf("all members count %d\n", int(count))

	// get a score of 500th member
	mem500score, err := rClient.ZScore(ctx, key1, "500").Result()
	if err != nil {
		log.Fatalf("ZScore error: %v", err)
	}
	// get 500th member ranking
	mem500rank, err := rClient.ZRevRank(ctx, key1, "500").Result()
	if err != nil {
		log.Fatalf("ZScore error: %v", err)
	}
	log.Printf("rank %d, user 500, score %d\n", mem500rank+1, int(mem500score))

	// increment a score of 500th member
	if err := rClient.ZIncrBy(ctx, key1, 100, "500").Err(); err != nil {
		log.Fatalf("ZIncrBy error: %v", err)
	}
	// get a score of 500th member
	mem500score2, err := rClient.ZScore(ctx, key1, "500").Result()
	if err != nil {
		log.Fatalf("ZScore error: %v", err)
	}
	log.Printf("rank -, user 500, score %d\n", int(mem500score2))

	// get members with the top 10 scores
	top10, err := rClient.ZRevRangeWithScores(ctx, key1, int64(0), int64(9)).Result()
	if err != nil {
		log.Fatalf("ZRevRangeWithScores error: %v", err)
	}
	for i, m := range top10 {
		log.Printf("rank %d, user %s, score %d\n", i+1, m.Member, int(m.Score))
	}

	// get members with scores between 45th and 54th
	around50, err := rClient.ZRevRangeWithScores(ctx, key1, int64(45), int64(54)).Result()
	if err != nil {
		log.Fatalf("ZRevRangeWithScores error: %v", err)
	}
	for i, m := range around50 {
		log.Printf("rank %d, user %s, score %d\n", i+45, m.Member, int(m.Score))
	}

	// string
	// set total in redis
	if err := rClient.Set(ctx, "total", 0, 0).Err(); err != nil {
		log.Fatalf("Set error: %v", err)
	}

	// get total
	total, err := rClient.Get(ctx, "total").Result()
	if err != nil {
		log.Fatalf("Get error: %v", err)
	}
	log.Printf("total %s\n", total)

	// increase total
	if err := rClient.IncrBy(ctx, "total", 100).Err(); err != nil {
		log.Fatalf("IncrBy error: %v", err)
	}
	total2, err := rClient.Get(ctx, "total").Result()
	if err != nil {
		log.Fatalf("Get error: %v", err)
	}
	log.Printf("total %s\n", total2)

	// string expire
	// set string value with limit
	if err := rClient.SetEX(ctx, "limit", 0, time.Minute).Err(); err != nil {
		log.Fatalf("SetEX error: %v", err)
	}

	// get limit value
	limit, err := rClient.Get(ctx, "limit").Result()
	if err != nil {
		log.Fatalf("Get error: %v", err)
	}
	log.Printf("limit %s\n", limit)

	// get ttl of limit
	d, err := rClient.TTL(ctx, "limit").Result()
	if err != nil {
		log.Fatalf("TTL error: %v", err)
	}
	log.Printf("limit duration %v\n", d)

	time.Sleep(3 * time.Second)

	// increase limit value
	if err := rClient.IncrBy(ctx, "limit", 1).Err(); err != nil {
		log.Fatalf("IncrBy error: %v", err)
	}
	limit2, err := rClient.Get(ctx, "limit").Result()
	if err != nil {
		log.Fatalf("Get error: %v", err)
	}
	log.Printf("limit %s\n", limit2)

	// get ttl of limit
	d2, err := rClient.TTL(ctx, "limit").Result()
	if err != nil {
		log.Fatalf("TTL error: %v", err)
	}
	log.Printf("limit duration %v\n", d2)

	// list
	// set list in redis
	lkey := "list1"
	for i := 1; i <= 500; i++ {
		r := time.Now()
		if err := rClient.LPush(ctx, lkey, r).Err(); err != nil {
			log.Fatalf("LPush error: %v", err)
		}
	}

	// get length of list
	llen, err := rClient.LLen(ctx, lkey).Result()
	if err != nil {
		log.Fatalf("LLen error: %v", err)
	}
	log.Printf("list count %d\n", int(llen))

	// get 10 values from the head
	ltop10, err := rClient.LRange(ctx, lkey, int64(0), int64(9)).Result()
	if err != nil {
		log.Fatalf("LRange error: %v", err)
	}
	for i, v := range ltop10 {
		log.Printf("index %d, value %v\n", i, v)
	}

	// get around 50 values
	laround50, err := rClient.LRange(ctx, lkey, int64(45), int64(54)).Result()
	if err != nil {
		log.Fatalf("LRange error: %v", err)
	}
	for i, v := range laround50 {
		log.Printf("index %d, value %v\n", i+44, v)
	}
}
