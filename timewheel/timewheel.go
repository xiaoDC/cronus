package timewheel

import (
	"errors"
	"fmt"
	"github.com/robfig/cron/v3"
	"strings"
	"time"
)

const (
	SEPARATOR = "@"
	MAX       = 99999
	WHEELSIZE = 64
	NUMBER    = 63
	BITS      = 6
)

var globalCron = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

type TimeWheel struct {
	interval   time.Duration // 指针每隔多久往前移动一格
	ticker     *time.Ticker
	currentPos int
	startTime  time.Time
	buckets    []*Bucket
	// TaskQueue  chan string
	exitC chan struct{}
}

func New() *TimeWheel {
	return newTimeWheel(time.Millisecond * 20)
}

func newTimeWheel(interval time.Duration) *TimeWheel {
	if interval <= time.Millisecond {
		panic(errors.New("interval be greater than 0"))
	}

	tw := &TimeWheel{
		interval:   interval,
		currentPos: 0,
		exitC:      make(chan struct{}),
		// TaskQueue:  make(chan string, MAX),
	}

	buckets := make([]*Bucket, WHEELSIZE)
	var i int
	for i = 0; i < WHEELSIZE; i++ {
		buckets[i] = newBucket(tw)
	}
	tw.buckets = buckets

	return tw
}

func (tw *TimeWheel) Start() {
	tw.ticker = time.NewTicker(tw.interval)
	tw.startTime = time.Now()

	go tw.start()
}

func (tw *TimeWheel) start() {
	for {
		select {
		case <-tw.ticker.C:
			tw.tickNext()
		case <-tw.exitC:
			tw.ticker.Stop()
			return
		}
	}
}

/**
 * 核心
 */
func (tw *TimeWheel) tickNext() {
	index := tw.currentPos & NUMBER
	bucket := tw.buckets[index]
	tw.currentPos += 1

	// get current bucket to do tasks
	go bucket.Run()
}

func (tw *TimeWheel) AddCronTask(spec string, id string, priority int, repeat bool) (string, error) {
	schedule, err := globalCron.Parse(spec)
	if err != nil {
		return "", err
	}

	if schedule == nil {
		fmt.Printf("%s crontab 表达式不合法", spec)
		return "", errors.New("crontab 表达式不合法")
	}

	now := time.Now()
	next := schedule.Next(now)
	diff := next.Sub(tw.startTime)

	num := int64(diff / tw.interval)

	bucketNum := num & NUMBER
	cycle := num >> BITS

	task := &Task{
		BucketNum: int(bucketNum),
		NextTime:  next,
		LastTime:  now,
		Schedule:  schedule,
		Cycle:     cycle,
		priority:  priority,
		Id:        id,
		Repeat:    repeat,
	}

	return task.Id + SEPARATOR + spec, tw.addTask(task)
}

func (tw *TimeWheel) addTask(task *Task) error {
	bucket := tw.buckets[task.BucketNum]
	return bucket.AddTask(task)
}

func (tw *TimeWheel) RemoveTask(id string) error {
	for _, bucket := range tw.buckets {
		task := bucket.FindTask(id)
		if task != nil {
			fmt.Printf("timewheel/timewheel.go-----------> %#v\n", task)
			return nil
		}
	}

	return errors.New("没有找到 id 为 " + id + " 的任务")
}

func (tw *TimeWheel) RemoveTaskWithCron(key string) error {
	id, spec, err := parse(key)

	if err != nil {
		fmt.Printf("%s -> task key is not legal", key)
		return errors.New("task key is not legal")
	}

	schedule, err := globalCron.Parse(spec)
	if err != nil {
		return err
	}
	if schedule == nil {
		fmt.Printf("%s crontab 表达式不合法", spec)
		return errors.New("crontab 表达式不合法")
	}

	next := schedule.Next(time.Now())
	diff := next.Sub(tw.startTime)

	num := int64(diff / tw.interval)
	bucketNum := num & NUMBER
	cycle := num >> BITS

	task := &Task{
		BucketNum: int(bucketNum),
		LastTime:  time.Now(),
		Schedule:  schedule,
		Cycle:     cycle,
		Id:        id,
	}

	return tw.buckets[task.BucketNum].RemoveTask(task)
}

func parse(key string) (string, string, error) {
	a := strings.Split(key, SEPARATOR)

	if len(a) < 2 {
		err := errors.New("task key is not legal")
		return "", "", err
	}
	return a[0], a[1], nil
}
