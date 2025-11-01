package websocket

import (
	"TeamTrackerBE/internal/config"
	"TeamTrackerBE/internal/domain/model"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
)

var (
	FlushInterval    = 30 * time.Second
	RedisTTL         = 60 * time.Second
	RedisListStart   = int64(0)
	RedisListEnd     = int64(-1)
	MaxRecordPerPush = 0
)

func (h *Hub) CacheLocation(loc LocationMessage) {
	if config.RedisClient == nil {
		return
	}
	if loc.UserID == uuid.Nil {
		return
	}

	key := h.RedisKey(loc.UserID)
	payload, err := json.Marshal(loc)
	if err != nil {
		return
	}

	if err := config.RedisClient.RPush(config.Ctx, key, string(payload)).Err(); err != nil {
		log.Printf("redis rpush failed for %s: %v", key, err)
		return
	}

	ttl, err := config.RedisClient.TTL(config.Ctx, key).Result()
	if err == nil && ttl < 0 {
		_ = config.RedisClient.Expire(config.Ctx, key, RedisTTL).Err()
	}

	if _, ok := h.flushTickers[loc.UserID]; !ok {
		h.StartFlushLoop(loc.UserID)
	}
}

func (h *Hub) StartFlushLoop(userID uuid.UUID) {
	tk := time.NewTicker(FlushInterval)
	h.flushTickers[userID] = tk

	go func(uid uuid.UUID, t *time.Ticker) {
		for range t.C {
			h.FlushUserLocations(uid)
		}
	}(userID, tk)
}

func (h *Hub) StopFlushLoop(userID uuid.UUID, flushNow bool) {
	if tk, ok := h.flushTickers[userID]; ok {
		tk.Stop()
		delete(h.flushTickers, userID)
		if flushNow {
			h.FlushUserLocations(userID)
		}
	}
}

func (h *Hub) FlushUserLocations(userID uuid.UUID) {
	if config.RedisClient == nil || h.locationRepo == nil || userID == uuid.Nil {
		return
	}

	key := h.RedisKey(userID)

	values, err := config.RedisClient.LRange(config.Ctx, key, RedisListStart, RedisListEnd).Result()
	if err != nil {
		log.Printf("redis lrange failed for %s: %v", key, err)
		return
	}
	if len(values) == 0 {
		return
	}

	records := make([]model.Location, 0, len(values))
	for _, raw := range values {
		var msg LocationMessage
		if err := json.Unmarshal([]byte(raw), &msg); err != nil {
			continue
		}

		msg.UserID = userID
		records = append(records, model.Location{
			UserID:    msg.UserID,
			Latitude:  msg.Latitude,
			Longitude: msg.Longitude,
		})
	}

	if len(records) == 0 {
		_ = config.RedisClient.Del(config.Ctx, key).Err()
		return
	}

	if err := h.locationRepo.BulkCreate(records); err != nil {
		log.Printf("bulk insert failed for user %s: %v", userID.String(), err)
		_ = config.RedisClient.Expire(config.Ctx, key, RedisTTL).Err()
		return
	}

	if err := config.RedisClient.Del(config.Ctx, key).Err(); err != nil {
		log.Printf("redis del failed for %s: %v", key, err)
	}
}

func (h *Hub) RedisKey(userID uuid.UUID) string {
	return "ws:locations:" + userID.String()
}
