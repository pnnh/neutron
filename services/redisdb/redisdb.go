package redisdb

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisAddress struct {
	Host     string `json:"host"`
	Password string `json:"password"`
	Db       int    `json:"db"`
}

func ParseConfig(redisUrl string) (*RedisAddress, error) {

	urlArr := strings.Split(redisUrl, " ")
	if len(urlArr) < 2 {
		return nil, fmt.Errorf("REDIS_URL格式错误，应该是<host:port password>格式")
	}
	redisHost := urlArr[0]
	redisPassword := urlArr[1]
	redisDb := 0 // 默认使用 Redis 数据库 0
	if len(urlArr) > 2 {
		db, err := strconv.Atoi(urlArr[2])
		if err != nil {
			return nil, fmt.Errorf("Redis DB配置错误: %v", err)
		}
		redisDb = db
	}
	redisAddress := &RedisAddress{
		Host:     redisHost,
		Password: redisPassword,
		Db:       redisDb,
	}
	return redisAddress, nil
}

var redisClientMap = make(map[string]*redis.Client)
var locker = &sync.Mutex{}

func ConnectRedis(ctx context.Context, redisUrl string) (*redis.Client, error) {
	locker.Lock()
	defer locker.Unlock()
	if redisClient, ok := redisClientMap[redisUrl]; ok {
		return redisClient, nil
	}
	redisAddress, err := ParseConfig(redisUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %v", err)
	}

	// 初始化 Redis 客户端
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddress.Host,     // Redis 地址
		Password: redisAddress.Password, // Redis 密码（默认空）
		DB:       redisAddress.Db,       // Redis 数据库
	})

	// 检查 Redis 连接
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	// 将客户端存入全局 map
	redisClientMap[redisUrl] = client
	return client, nil

}

// Produce 生产者：向 Redis 队列推送消息
func Produce(ctx context.Context, client *redis.Client, queueName string, contentData []byte) error {

	// 使用 LPUSH 将消息推送到队列
	err := client.LPush(ctx, queueName, contentData).Err()
	if err != nil {
		return fmt.Errorf("failed to push task to queue: %v", err)
	}
	return nil
}

// Consume 消费者：从 Redis 队列获取并处理消息
func Consume(ctx context.Context, client *redis.Client, queueName string) ([]byte, error) {

	// 使用 BRPOP 阻塞获取消息，超时 0 表示无限等待
	result, err := client.BRPop(ctx, 0*time.Second, queueName).Result()
	if err != nil {
		return nil, fmt.Errorf("error popping from queue: %v", err)
	}

	return []byte(result[1]), nil
}
