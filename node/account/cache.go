package account

import (
	"fmt"
	"sync"
	"time"
)

const ExpireTime = 60 * 5

type Cache struct {
	m          map[string]VerificationInfo
	mutex      sync.Mutex
	expireTime time.Duration
}

type VerificationInfo struct {
	code      string
	startTime int64
	isUse     bool
}

func NewCache() (*Cache, error) {
	c := &Cache{
		m:     make(map[string]VerificationInfo),
		mutex: sync.Mutex{},
	}

	return c, nil
}

func (c *Cache) Set(key, code string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if len(code) != 6 {
		return fmt.Errorf("验证码格式不正确")
	}

	now := time.Now().Unix()

	if info, ok := c.m[key]; ok {
		if info.startTime+ExpireTime > now {
			return fmt.Errorf("获取验证码间隔太快")
		}
	}

	c.m[key] = VerificationInfo{
		code:      code,
		startTime: now,
		isUse:     false,
	}

	if len(c.m) > 100 {
		var keysToDelete []string
		for key, value := range c.m {
			if value.startTime+ExpireTime*2 < now {
				keysToDelete = append(keysToDelete, key)
			}
		}

		// 遍历切片，从map中删除相应的键
		for i := range keysToDelete {
			delete(c.m, keysToDelete[i])
		}
	}

	return nil
}

func (c *Cache) Check(key, code string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now().Unix()
	if info, ok := c.m[key]; ok {
		if info.startTime+ExpireTime < now {
			return fmt.Errorf("验证码已经过期")
		} else if info.code != code {
			return fmt.Errorf("验证码不正确")
		} else if info.isUse {
			return fmt.Errorf("验证码已经使用")
		} else {
			info.isUse = true
			c.m[key] = info
			return nil
		}
	} else {
		return fmt.Errorf("请先获取验证码")
	}
}
