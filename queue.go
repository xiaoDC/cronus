package timewheel

type PriorityQueue struct {
	data []*Task
	len  int
}

func NewQueue() *PriorityQueue {
	return &PriorityQueue{
		data: []*Task{},
		len:  0,
	}
}

func (pq *PriorityQueue) Push(task *Task) {
	task.index = pq.len
	pq.data = append(pq.data, task)
	pq.len += 1
	pq.up(pq.len - 1)
}

func (pq *PriorityQueue) Less(j, i int) bool {
	return pq.data[j].priority > pq.data[i].priority
}

func (pq *PriorityQueue) Swap(j, i int) {
	pq.data[i], pq.data[j] = pq.data[j], pq.data[i]
	pq.data[i].index = i
	pq.data[j].index = j
}

func (pq *PriorityQueue) Pop() *Task {
	pq.Swap(0, pq.len-1)
	pq.down(0, pq.len-1)

	old := pq.data
	item := old[pq.len-1]
	item.index = -1
	pq.len -= 1
	pq.data = old[0:pq.len]
	return item
}

func (pq *PriorityQueue) Remove(i int) *Task {
	n := pq.len - 1
	if n != i {
		pq.Swap(i, n)
		if !pq.down(i, n) {
			pq.up(i)
		}
	}

	old := pq.data
	item := old[pq.len-1]
	item.index = -1
	pq.len -= 1
	pq.data = old[0:pq.len]
	return item
}

func (pq *PriorityQueue) Find(task *Task, compartor func(a *Task) bool) *Task {
	for i := pq.len - 1; i >= 0; i-- {
		if compartor(pq.data[i]) {
			return pq.data[i]
		}
	}

	return nil
}

func (pq *PriorityQueue) up(j int) {
	for {
		i := (j - 1) / 2

		if i == j || !pq.Less(j, i) {
			break
		}
		pq.Swap(i, j)
		j = i
	}
}

func (pq *PriorityQueue) down(i0, n int) bool {
	i := i0

	for {
		j1 := 2*i + 1

		if j1 >= n || j1 < 0 {
			break
		}
		j := j1
		if j2 := j1 + 1; j2 < n && pq.Less(j2, j1) {
			j = j2
		}

		if !pq.Less(j, i) {
			break
		}
		pq.Swap(i, j)
		i = j
	}

	return i > i0
}
