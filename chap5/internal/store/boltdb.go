package store

import (
	"encoding/json"
	"errors"
	"time"

	// 引入日志模块
	"kvstore/internal/logger"

	bolt "go.etcd.io/bbolt"
)

var bucketName = []byte("kv")

type Store struct {
	db *bolt.DB
}

// entry 用于存储值和过期时间
type entry struct {
	Value     string    `json:"value"`
	ExpiresAt time.Time `json:"expires_at"`
}

// NewStore 初始化 Store
func NewStore(path string) (*Store, error) {
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	// 创建 bucket（仅首次）
	err = db.Update(func(tx *bolt.Tx) error {
		_, e := tx.CreateBucketIfNotExists(bucketName)
		return e
	})
	if err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

// Set 保存 key-value（可选 TTL，单位：秒，ttlSec = 0 表示永久）
func (s *Store) Set(key, value string, ttlSec int64) error {
	var expiresAt time.Time
	if ttlSec > 0 {
		expiresAt = time.Now().Add(time.Duration(ttlSec) * time.Second)
	}

	// 将值和过期时间序列化
	data, err := json.Marshal(entry{
		Value:     value,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return err
	}

	// 更新数据库
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return errors.New("bucket not found")
		}
		return b.Put([]byte(key), data)
	})
}

// UpdateTTL 更新已有 key 的 TTL
func (s *Store) UpdateTTL(key string, ttlSec int64) error {
	var expiresAt time.Time
	if ttlSec > 0 {
		expiresAt = time.Now().Add(time.Duration(ttlSec) * time.Second)
	}

	// 将新的过期时间更新到数据库
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return errors.New("bucket not found")
		}

		// 获取当前键值对
		raw := b.Get([]byte(key))
		if raw == nil {
			return errors.New("key not found")
		}

		var e entry
		if err := json.Unmarshal(raw, &e); err != nil {
			return err
		}

		// 更新 TTL
		e.ExpiresAt = expiresAt

		// 将更新后的值和过期时间序列化
		data, err := json.Marshal(e)
		if err != nil {
			return err
		}

		// 保存更新后的数据
		return b.Put([]byte(key), data)
	})
}

// Get 查询 key，对 TTL 做过期判断
func (s *Store) Get(key string) (string, error) {
	var raw []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return errors.New("bucket not found")
		}
		raw = b.Get([]byte(key))
		return nil
	})
	if err != nil {
		return "", err
	}
	if raw == nil {
		return "", errors.New("key not found")
	}

	var e entry
	if err := json.Unmarshal(raw, &e); err != nil {
		return "", err
	}
	if !e.ExpiresAt.IsZero() && time.Now().After(e.ExpiresAt) {
		return "", errors.New("key expired")
	}

	return e.Value, nil
}

// GetWithTTL 查询 key 并返回值和剩余的 TTL
func (s *Store) GetWithTTL(key string) (string, int64, error) {
	var raw []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return errors.New("bucket not found")
		}
		raw = b.Get([]byte(key))
		return nil
	})
	if err != nil {
		return "", 0, err
	}
	if raw == nil {
		return "", 0, errors.New("key not found")
	}

	var e entry
	if err := json.Unmarshal(raw, &e); err != nil {
		return "", 0, err
	}

	// 如果有过期时间，检查是否过期并返回剩余TTL
	if !e.ExpiresAt.IsZero() {
		// 使用 time.Until 来计算剩余的 TTL
		remainingTTL := int64(time.Until(e.ExpiresAt).Seconds())
		if remainingTTL <= 0 {
			return "", 0, errors.New("key expired")
		}
		return e.Value, remainingTTL, nil
	}

	// 如果没有设置过期时间，则表示永不过期
	return e.Value, 0, nil
}

// Delete 删除 key
func (s *Store) Delete(key string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return errors.New("bucket not found")
		}
		return b.Delete([]byte(key))
	})
}

// ListKeys 返回所有未过期的 key
func (s *Store) ListKeys() ([]string, error) {
	var keys []string
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return errors.New("bucket not found")
		}
		return b.ForEach(func(k, v []byte) error {
			var e entry
			if err := json.Unmarshal(v, &e); err != nil {
				return nil // 跳过损坏数据
			}
			if e.ExpiresAt.IsZero() || time.Now().Before(e.ExpiresAt) {
				keys = append(keys, string(k))
			}
			return nil
		})
	})
	return keys, err
}

// PurgeExpiredKeys 删除过期的键值对
func (s *Store) PurgeExpiredKeys() error {
	// 假设我们有逻辑清理过期数据，您可以在此实现
	logger.Log.Info("Purging expired keys...")

	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return errors.New("bucket not found")
		}

		// 假设每个键值对的过期时间是存储在 value 中的一部分，我们遍历并删除过期的项
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			// 在这里实现过期检查的逻辑，例如：
			if isExpired(v) {
				// 记录删除的键
				logger.Log.Infof("Deleting expired key: %s", string(k))
				err := b.Delete(k)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	return err
}

// isExpired 判断键值是否过期
func isExpired(data []byte) bool {
	var e entry
	err := json.Unmarshal(data, &e)
	if err != nil {
		return false
	}

	// 如果没有设置过期时间，或者过期时间在未来，表示没有过期
	if e.ExpiresAt.IsZero() || e.ExpiresAt.After(time.Now()) {
		return false
	}

	// 否则表示过期
	return true
}

// StartTTLGC 启动后台自动清理协程（每 interval 检查一次）
func (s *Store) StartTTLGC(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		logger.Log.Info("TTL GC started")

		for range ticker.C {
			logger.Log.Info("Running TTL GC...")

			// 执行过期数据清理
			err := s.PurgeExpiredKeys()
			if err != nil {
				logger.Log.Errorf("Error during TTL GC: %v", err)
			} else {
				logger.Log.Info("TTL GC completed successfully")
			}
		}
	}()
}

// Close 关闭数据库
func (s *Store) Close() error {
	return s.db.Close()
}
