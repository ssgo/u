package u

import (
	"sync"
	"time"
)

type IdMaker struct {
	secCurrent   uint64
	secIndexNext uint64
	secIndexLock sync.Mutex
	Incr         func(sec uint64) uint64
}

func NewIdMaker(incr func(sec uint64) uint64) *IdMaker {
	return &IdMaker{Incr: incr, secIndexNext: 1}
}

var DefaultIdMaker = NewIdMaker(nil)

func MakeId(size int) string {
	return DefaultIdMaker.Get(size)
}

func (im *IdMaker) defaultIncr(sec uint64) uint64 {
	im.secIndexLock.Lock()
	defer im.secIndexLock.Unlock()
	if im.secCurrent == sec {
		im.secIndexNext++
	} else {
		if im.secCurrent == 0 {
			im.secIndexNext = GlobalRand1.Uint64N(1000) + 1000 // 防止服务在1秒内重启导致碰撞，我们假设单机模式1秒内并发数不超过2000个
		} else {
			im.secIndexNext = 1
		}
		im.secCurrent = sec
	}
	return im.secIndexNext
}

// ID 格式：secTag（1位） + sec（5位/6位） + secIndex（1~5位） + padding（随机填充）
//  1. sec：当前绝对秒数（从2000-01-01起算），5位Base62纯净容量为 901,356,495 秒（约28.56年为一个周期 n）。
//  2. secIndex：同1秒内的并发自增序号。长度记为 m（1~5位，最大支持单秒 9.16 亿并发）。
//  3. secTag：融合了时间跨度周期 n 和 索引长度 m。算法为 secTag = n*5 + (m-1)。
//     当 n 在 0~10 之间时（即前 314 年，直到公元 2314 年），sec 固定输出 5 位。
//     当 n = 11 时（魔数触发），sec 升级为 6 位原始 Base62，可额外支撑 1800 年不碰撞。
//  4. 数据库友好主键(MakeIdForPK)：截取最后一位随机字符并整体右旋至首位，完美实现 B+树的 62 个分片散列，避免写入热点，同时保证同一分片内具备时间局部单调性。
//  5. 输出长度并发支持情况：
//     8位(8-5-1=2)：可以支持每秒 3844个ID
//     9位(9-5-1=3)：可以支持每秒 9万个ID
//     10位(10-5-1=4)：可以支持每秒 1477万个ID
//     11位(11-5-1=6)：可以支持每秒 9亿个ID
func (im *IdMaker) get(size int, ordered bool, hashToHead bool) string {
	tm := time.Now()

	// 计算当前秒数（相对于 2000 年纪元：946656000）
	nowSec := uint64(tm.Unix() - 946656000)
	var n, sec uint64
	secCapacity := uint64(901356495) // 5位 Base62 纯净容量
	if nowSec < 11*secCapacity {
		// 314 年内的 5位 sec 逻辑 (n = 0 到 10)
		n = nowSec / secCapacity
		sec = (nowSec % secCapacity) + 14776336 // 加上 5位62进制最小值
	} else {
		// 314 年后（约公元2314年），激活 6位 sec 隐藏剧情
		n = 11       // 魔数 11 代表进入 6位扩展时代
		sec = nowSec // 此时的 nowSec 天然在 6位 Base62 范围内，直接使用
	}

	var secIndex uint64
	if im.Incr != nil {
		secIndex = im.Incr(sec)
	}
	if secIndex == 0 {
		secIndex = im.defaultIncr(sec) // 约定必须从1开始，返回0表示失败，使用 defaultIncr 兜底
	}

	// 计算 secTag
	intEncoder := DefaultIntEncoder
	if ordered {
		intEncoder = OrderedIntEncoder
	}
	secBytes := intEncoder.EncodeInt(sec)
	secLen := len(secBytes)
	inSecIndexBytes := intEncoder.EncodeInt(secIndex)

	// 计算 m (secIndex 的长度，最大支持到 5，即 9.16亿并发/秒), 极端兜底，超过 9.16亿/秒 时强制截断，防止 secTag 溢出
	m := min(uint64(len(inSecIndexBytes)), 5)

	secTagVal := n*5 + (m - 1)
	var uid = make([]byte, 0, size)
	uid = intEncoder.AppendInt(uid, secTagVal)

	// 追加 secBytes 和 inSecIndexBytes
	uid = append(uid, secBytes...)
	uid = append(uid, inSecIndexBytes...)
	uid = intEncoder.FillInt(uid, size) // 用随机数填充
	if !ordered {
		intEncoder.HashInt(ExchangeInt(uid)) // 整体交叉然后散列乱序
	} else {
		intEncoder.HashInt(ExchangeInt(uid[secLen+1:])) // 仅对 secIndex 之后的字符交叉然后散列乱序
		if hashToHead {                                 // 整体右旋至首位，实现 B+树的 62 个分片散列，避免写入热点（针对Mysql系数据库优化）
			size = len(uid)
			lastByte := uid[size-1]
			copy(uid[1:], uid[:size-1])
			uid[0] = lastByte
		}
	}
	return string(uid)
}

func (im *IdMaker) Get(size int) string {
	return im.get(size, false, false)
}

func (im *IdMaker) GetForMysql(size int) string {
	return im.get(size, true, true)
}

func (im *IdMaker) GetForPostgreSQL(size int) string {
	return im.get(size, true, false)
}
