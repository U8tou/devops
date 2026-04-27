package idx

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/jaevor/go-nanoid"
)

var (
	canonicGen   func() string
	canonicOnce  sync.Once
	numericGens  = make(map[int]func() string)
	numericMutex sync.RWMutex
)

// 初始化标准 nanoid 生成器（21位，懒加载）
func getCanonic() func() string {
	canonicOnce.Do(func() {
		gen, err := nanoid.Canonic()
		if err != nil {
			panic("failed to create nanoid generator: " + err.Error())
		}
		canonicGen = gen
	})
	return canonicGen
}

// 获取指定长度的数字 nanoid 生成器（懒加载 + 缓存）
func getNumericGen(length int) func() string {
	numericMutex.RLock()
	gen, ok := numericGens[length]
	numericMutex.RUnlock()
	if ok {
		return gen
	}

	numericMutex.Lock()
	defer numericMutex.Unlock()
	// 双重检查
	if gen, ok = numericGens[length]; ok {
		return gen
	}
	gen, err := nanoid.CustomASCII("0123456789", length)
	if err != nil {
		panic("failed to create numeric nanoid generator: " + err.Error())
	}
	numericGens[length] = gen
	return gen
}

// NanoId 生成21位标准 nanoid（替代 UUID）
func NanoId() string {
	return getCanonic()()
}

// NanoIdNum 生成指定长度的纯数字 ID
func NanoIdNum(length int) string {
	return getNumericGen(length)()
}

// OrderId 生成订单ID（18位）
// 规则：前缀(2位) + 年月日时分(10位) + 随机数(6位)
func OrderId(prefix int) int64 {
	if prefix < 10 || prefix > 99 {
		return 0
	}
	now := time.Now()
	str := fmt.Sprintf("%d%02d%02d%02d%02d%02d%s",
		prefix,
		now.Year()%100,
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		NanoIdNum(6),
	)
	id, _ := strconv.ParseInt(str, 10, 64)
	return id
}
