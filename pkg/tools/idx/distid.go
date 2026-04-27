package idx

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

// 分布式数字 ID 生成器（不依赖 Redis/DB，仅需每台机器配置不同的节点种子）
//
// 用法：机器1 传入 NodeSeed(1)，机器2 传入 NodeSeed(2)，依此类推。
// 支持 6/8/12 位十进制数字，同一节点内单调递增，多节点间通过种子区分不冲突。
//
// 位数与保证期：
//   - 12 位：时间按分钟，约 15 年内不重复，推荐生产使用。
//   - 8 位：时间按秒且会周期回绕，约 4.5 小时一周期，适合短生命周期场景。
//   - 6 位：时间按秒且回绕更快，约 17 分钟一周期，且数值范围为 100000～624287。
var (
	epoch = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
)

// DistIdGen 分布式数字 ID 生成器
type DistIdGen struct {
	digits int   // 6, 8, 12
	node   int64 // 节点种子，如机器1=1，机器2=2
	mu     sync.Mutex
	// 上一时间片与序列，用于同时间片内自增
	lastTime int64
	lastSeq  int64
}

// DistIdConfig 配置
type DistIdConfig struct {
	Digits int   // 位数：6、8 或 12
	Node   int64 // 节点种子，不同机器不同值，如 1、2、3
}

// NewDistId 创建分布式 ID 生成器。node 为机器种子（如 1、2、3），digits 为 6、8 或 12。
func NewDistId(node int64, digits int) *DistIdGen {
	if digits != 6 && digits != 8 && digits != 12 {
		digits = 12
	}
	if node < 0 {
		node = 0
	}
	return &DistIdGen{digits: digits, node: node}
}

// NewDistIdFromEnv 从环境变量 NODE_ID 读取节点种子创建生成器；未设置或非法时使用 0。digits 为 6、8 或 12。
// 部署时可在每台机器上配置不同 NODE_ID（如 1、2、3）即可保证分布式不重复。
func NewDistIdFromEnv(digits int) *DistIdGen {
	node := int64(0)
	if s := os.Getenv("NODE_ID"); s != "" {
		if n, err := strconv.ParseInt(s, 10, 64); err == nil && n >= 0 {
			node = n
		}
	}
	return NewDistId(node, digits)
}

// Next 生成下一个 ID，返回十进制数字字符串，长度固定为 digits（前导零补齐）
func (g *DistIdGen) Next() string {
	g.mu.Lock()
	defer g.mu.Unlock()

	switch g.digits {
	case 12:
		return g.next12()
	case 8:
		return g.next8()
	default:
		return g.next6()
	}
}

// NextInt64 生成下一个 ID 并返回 int64（仅 12 位时能完整表示，8/6 位会落在对应范围内）
func (g *DistIdGen) NextInt64() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	switch g.digits {
	case 12:
		return g.next12Int64()
	case 8:
		return g.next8Int64()
	default:
		return g.next6Int64()
	}
}

// --- 12 位：23 位时间(分钟) + 8 位节点 + 8 位序列 = 39 位，约 15 年不回绕 ---
const (
	t12TimeBits = 23
	t12NodeBits = 8
	t12SeqBits  = 8
	t12TimeMax  = 1<<t12TimeBits - 1
	t12NodeMax  = 1<<t12NodeBits - 1
	t12SeqMax   = 1<<t12SeqBits - 1
)

func (g *DistIdGen) next12() string {
	return fmt.Sprintf("%012d", g.next12Int64())
}

func (g *DistIdGen) next12Int64() int64 {
	now := time.Now().UTC()
	minutes := int64(now.Sub(epoch) / time.Minute)
	if minutes > t12TimeMax {
		minutes = t12TimeMax
	}
	node := g.node & t12NodeMax

	g.lastSeq++
	if g.lastTime != minutes {
		g.lastTime = minutes
		g.lastSeq = 0
	}
	if g.lastSeq > t12SeqMax {
		// 同一分钟内序列用完，忙等到下一分钟（极少发生）
		for int64(time.Now().UTC().Sub(epoch)/time.Minute) == minutes {
			time.Sleep(10 * time.Millisecond)
		}
		return g.next12Int64()
	}
	seq := g.lastSeq & t12SeqMax

	id := (minutes << (t12NodeBits + t12SeqBits)) | (node << t12SeqBits) | seq
	// 39 位最大值 549755813887，满足 12 位
	const max12 = 999999999999
	if id > max12 {
		id = id % (max12 + 1)
	}
	return id
}

// --- 8 位：14 位时间(秒，约 4.5h 回绕) + 6 位节点 + 6 位序列 = 26 位 ---
const (
	t8TimeBits = 14
	t8NodeBits = 6
	t8SeqBits  = 6
	t8TimeMax  = 1<<t8TimeBits - 1
	t8NodeMax  = 1<<t8NodeBits - 1
	t8SeqMax   = 1<<t8SeqBits - 1
)

func (g *DistIdGen) next8() string {
	return fmt.Sprintf("%08d", g.next8Int64())
}

func (g *DistIdGen) next8Int64() int64 {
	now := time.Now().UTC()
	sec := int64(now.Sub(epoch) / time.Second)
	slot := sec & t8TimeMax
	node := g.node & t8NodeMax

	g.lastSeq++
	if g.lastTime != slot {
		g.lastTime = slot
		g.lastSeq = 0
	}
	if g.lastSeq > t8SeqMax {
		// 当前时间片序列用完，等下一秒再生成
		for int64(time.Now().UTC().Sub(epoch)/time.Second) <= sec {
			time.Sleep(10 * time.Millisecond)
		}
		return g.next8Int64()
	}
	seq := g.lastSeq & t8SeqMax

	id := (slot << (t8NodeBits + t8SeqBits)) | (node << t8SeqBits) | seq
	const max8 = 99999999
	if id > max8 {
		id = id % (max8 + 1)
	}
	return id
}

// --- 6 位：10 位时间(秒，约 17min 回绕) + 5 位节点 + 4 位序列 = 19 位，输出 100000～624287 ---
const (
	t6TimeBits = 10
	t6NodeBits = 5
	t6SeqBits  = 4
	t6TimeMax  = 1<<t6TimeBits - 1
	t6NodeMax  = 1<<t6NodeBits - 1
	t6SeqMax   = 1<<t6SeqBits - 1
	t6Base     = 100000
	t6MaxVal   = 1<<19 - 1
)

func (g *DistIdGen) next6() string {
	return fmt.Sprintf("%06d", g.next6Int64())
}

func (g *DistIdGen) next6Int64() int64 {
	now := time.Now().UTC()
	sec := int64(now.Sub(epoch) / time.Second)
	slot := sec & t6TimeMax
	node := g.node & t6NodeMax

	g.lastSeq++
	if g.lastTime != slot {
		g.lastTime = slot
		g.lastSeq = 0
	}
	if g.lastSeq > t6SeqMax {
		for int64(time.Now().UTC().Sub(epoch)/time.Second) <= sec {
			time.Sleep(10 * time.Millisecond)
		}
		return g.next6Int64()
	}
	seq := g.lastSeq & t6SeqMax

	raw := (slot << (t6NodeBits + t6SeqBits)) | (node << t6SeqBits) | seq
	if raw > t6MaxVal {
		raw = raw % (t6MaxVal + 1)
	}
	return t6Base + raw
}
