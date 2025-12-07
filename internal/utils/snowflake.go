package utils

import (
	"sync"
	"time"
)

// 雪花算法配置
const (
	epoch          = int64(1735660800000) // 起始时间戳：2025-01-01 00:00:00 Asia/Shanghai
	timestampBits  = 41                   // 时间戳41位
	datacenterBits = 5                    // 数据中心5位
	workerBits     = 5                    // 工作节点5位
	sequenceBits   = 12                   // 序列号12位

	maxSequence   = int64(-1) ^ (int64(-1) << sequenceBits)   // 4095
	maxWorker     = int64(-1) ^ (int64(-1) << workerBits)     // 31
	maxDatacenter = int64(-1) ^ (int64(-1) << datacenterBits) // 31

	timestampShift  = sequenceBits + workerBits + datacenterBits // 22
	datacenterShift = sequenceBits + workerBits                  // 17
	workerShift     = sequenceBits                               // 12
)

// Snowflake 雪花算法生成器
type Snowflake struct {
	mu            sync.Mutex
	lastTimestamp int64
	datacenterID  int64
	workerID      int64
	sequence      int64
}

var (
	// 全局唯一的ID生成器实例
	globalIDGenerator *Snowflake
	initOnce          sync.Once
)

// InitIDGenerator 初始化全局ID生成器
func InitIDGenerator(datacenterID, workerID int64) {
	if datacenterID < 0 || datacenterID > maxDatacenter {
		panic("datacenter ID 超出范围")
	}
	if workerID < 0 || workerID > maxWorker {
		panic("worker ID 超出范围")
	}

	globalIDGenerator = &Snowflake{
		datacenterID: datacenterID,
		workerID:     workerID,
	}
}

// Generate 生成全局唯一ID
func (s *Snowflake) Generate() uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	currentTime := time.Now().UnixNano() / 1e6
	timestamp := currentTime - epoch

	if timestamp < 0 {
		// 时钟回拨处理
		for timestamp < 0 {
			time.Sleep(100 * time.Microsecond)
			currentTime = time.Now().UnixNano() / 1e6
			timestamp = currentTime - epoch
		}
	}

	if timestamp > (1<<timestampBits)-1 {
		panic("时间戳溢出，请调整起始时间")
	}

	if timestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 {
			// 序列号用尽，等待下一毫秒
			for timestamp <= s.lastTimestamp {
				time.Sleep(100 * time.Microsecond)
				currentTime = time.Now().UnixNano() / 1e6
				timestamp = currentTime - epoch
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = timestamp

	return uint64(timestamp)<<timestampShift |
		uint64(s.datacenterID)<<datacenterShift |
		uint64(s.workerID)<<workerShift |
		uint64(s.sequence)
}

// GenerateID 全局唯一ID生成函数（简洁调用）
func GenerateID() uint64 {
	if globalIDGenerator == nil {
		// 默认初始化（数据中心1，工作节点1）
		InitIDGenerator(1, 1)
	}
	return globalIDGenerator.Generate()
}
