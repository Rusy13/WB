package storage

import (
	"WB/internal/order/storage/database/dto"
	"encoding/json"
	"github.com/gomodule/redigo/redis"
)

func (s *OrderStorageDB) GetOrderFromCache(orderID string) (*dto.OrderFromCache, error) {
	key := constructRedisKey(orderID)
	orderFromRedis, err := redis.String(s.redisConn.Do("GET", key))
	if err != nil {
		s.logger.Errorf("redis error: %v", err)
		return nil, err
	}

	var orderCache dto.OrderFromCache
	err = json.Unmarshal([]byte(orderFromRedis), &orderCache)
	if err != nil {
		s.logger.Errorf("unmarshal error: %v", err)
		return nil, err
	}

	return &orderCache, nil
}
