package timewheel

import (
	"github.com/robfig/cron/v3"
	"time"
)

type Task struct {
	Id        string
	Schedule  cron.Schedule
	LastTime  time.Time
	BucketNum int
	Cycle     int64 // 记录 timeWheel 运行了多少轮
	priority  int   // 元素在队列中的优先级。
	index     int   // 元素在队列中的索引
	Repeat    bool
}
