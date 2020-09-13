package timewheel

import (
	// "container/heap"
	"github.com/robfig/cron/v3"
	"time"
)

type Bucket struct {
	tree      *RBTree
	size      int
	cycle     int64
	timeWheel *TimeWheel
}

func newBucket(tw *TimeWheel) *Bucket {
	bucket := &Bucket{
		size:  0,
		cycle: 0,
		tree:  NewRBTree(),
		// bucket 持有 timeWheel 对象，方便后续 task 做移动
		timeWheel: tw,
	}
	return bucket
}

func (b *Bucket) Run() {
	tw := b.timeWheel
	node := b.tree.Remove(b.cycle)
	b.cycle += 1

	if node != nil {
		var now, next time.Time
		var schedule cron.Schedule
		var diff time.Duration

		queue := node.Queue
		size := 0
		now = time.Now()
		// fmt.Println()
		// fmt.Println("now --------------->", now)
		len := queue.len
		for i := 0; i < len; i++ {
			// item := heap.Pop(queue).(*Task)
			item := queue.Pop()
			item.LastTime = now
			size++
			schedule = item.Schedule
			next = schedule.Next(now)
			diff = next.Sub(tw.startTime)
			num := int64(diff / tw.interval)
			bucketNum := num % tw.wheelSize
			// 下次执行时， timeWheel 的轮数
			item.Cycle = num / tw.wheelSize
			// 下次执行时， bucket 的 index
			item.BucketNum = int(bucketNum)

			if item.Repeat {
				b.timeWheel.addTask(item)
			}

			tw.TaskQueue <- item.Id
		}

		b.size -= size
	}
}

func (b *Bucket) AddTask(task *Task) error {
	err := b.tree.InsertTask(task)
	if err == nil {
		b.size += 1
	}
	return err
}

func (b *Bucket) RemoveTask(task *Task) error {
	err := b.tree.RemoveTask(task)
	if err == nil {
		b.size -= 1
	}

	return err
}

func (b *Bucket) FindTask(id string) *Task {
	return b.tree.FindTaskBy(id)
}
