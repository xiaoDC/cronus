package timewheel

import (
	"sync"
)

type color bool

type RBNode struct {
	Cycle               int64
	color               color
	Queue               *PriorityQueue
	Left, Right, Parent *RBNode
}

func Comparator(a, b int64) int64 {
	return a - b
}

type RBTree struct {
	Root      *RBNode
	runningMu sync.Mutex
}

const (
	RED   color = true
	BLACK color = false
)

func (n *RBNode) grandParent() *RBNode {
	return n.Parent.Parent
}

func (n *RBNode) uncle() *RBNode {
	if n == nil || n.Parent == nil || n.Parent.Parent == nil {
		return nil
	}

	return n.Parent.sibling()
}

func (n *RBNode) sibling() *RBNode {
	if n == nil || n.Parent == nil {
		return nil
	}
	if n == n.Parent.Left {
		return n.Parent.Right
	}
	return n.Parent.Left
}

func NewRBTree() *RBTree {
	return &RBTree{
		Root:      nil,
		runningMu: sync.Mutex{},
	}
}

func (t *RBTree) GetNode(cycle int64) *RBNode {
	node := t.lookup(cycle)
	return node
}

func (t *RBTree) lookup(cycle int64) *RBNode {
	node := t.Root

	for node != nil {
		compare := Comparator(cycle, node.Cycle)
		switch {
		case compare == 0:
			return node
		case compare < 0:
			node = node.Left
		case compare > 0:
			node = node.Right
		}
	}
	return nil
}

func (t *RBTree) InsertTask(task *Task) error {
	t.runningMu.Lock()
	defer t.runningMu.Unlock()
	var insertedNode *RBNode

	if t.Root == nil {
		queue := NewQueue()
		queue.Push(task)
		t.Root = &RBNode{
			Cycle: task.Cycle,
			color: BLACK,
			Queue: queue,
		}
		insertedNode = t.Root
	} else {
		node := t.GetNode(task.Cycle)

		if node != nil {
			node.Queue.Push(task)
			return nil
		} else {
			node = t.Root
			queue := NewQueue()
			queue.Push(task)

			loop := true
			for loop {
				compare := Comparator(task.Cycle, node.Cycle)
				switch {
				case compare == 0:
					loop = true
					return nil
				case compare < 0:
					if node.Left == nil {
						node.Left = &RBNode{
							Cycle:  task.Cycle,
							color:  RED,
							Queue:  queue,
							Parent: node,
						}
						insertedNode = node.Left
						loop = false
					} else {
						node = node.Left
					}

				case compare > 0:
					if node.Right == nil {
						node.Right = &RBNode{
							Cycle:  task.Cycle,
							color:  RED,
							Queue:  queue,
							Parent: node,
						}
						insertedNode = node.Right
						loop = false
					} else {
						node = node.Right
					}
				}
			}

			insertedNode.Parent = node
		}
	}

	return t.insertCase1(insertedNode)
}

func (t *RBTree) insertCase1(node *RBNode) error {
	if node.Parent == nil {
		node.color = BLACK
		return nil
	}
	return t.insertCase2(node)
}

func (t *RBTree) insertCase2(node *RBNode) error {
	if nodeColor(node.Parent) == BLACK {
		return nil
	}
	return t.insertCase3(node)
}

func (t *RBTree) insertCase3(node *RBNode) error {
	uncle := node.uncle()
	if nodeColor(uncle) == RED {
		node.Parent.color = BLACK
		uncle.color = BLACK
		node.grandParent().color = RED
		return t.insertCase1(node.grandParent())
	} else {
		return t.insertCase4(node)
	}
}

func (t *RBTree) insertCase4(node *RBNode) error {
	grandParent := node.grandParent()

	if node == node.Parent.Right && node.Parent == grandParent.Left {
		t.rotateLeft(node.Parent)
		node = node.Left
	} else if node == node.Parent.Left && node.Parent == grandParent.Right {
		t.rotateRight(node.Parent)
		node = node.Right
	}

	return t.insertCase5(node)
}

func (t *RBTree) insertCase5(node *RBNode) error {
	node.Parent.color = BLACK
	grandParent := node.grandParent()
	grandParent.color = RED

	if node == node.Parent.Left && node.Parent == grandParent.Left {
		t.rotateRight(grandParent)
		return nil
	}

	if node == node.Parent.Right && node.Parent == grandParent.Right {
		t.rotateLeft(grandParent)
		return nil
	}

	return nil
}

func (t *RBTree) rotateLeft(node *RBNode) {
	right := node.Right
	t.replaceNode(node, right)
	node.Right = right.Left
	if right.Left != nil {
		right.Left.Parent = node
	}
	right.Left = node
	node.Parent = right
}

func (t *RBTree) rotateRight(node *RBNode) {
	left := node.Left
	t.replaceNode(node, left)
	node.Left = left.Right
	if left.Right != nil {
		left.Right.Parent = node
	}
	left.Right = node
	node.Parent = left
}

func (t *RBTree) replaceNode(old *RBNode, newNode *RBNode) {
	if old.Parent == nil {
		t.Root = newNode
	} else {
		if old == old.Parent.Left {
			old.Parent.Left = newNode
		} else {
			old.Parent.Right = newNode
		}
	}
	if newNode != nil {
		newNode.Parent = old.Parent
	}
}

func nodeColor(node *RBNode) color {
	if node == nil {
		return BLACK
	}
	return node.color
}

func (t *RBTree) Remove(cycle int64) *RBNode {
	t.runningMu.Lock()
	defer t.runningMu.Unlock()
	node := t.lookup(cycle)
	if node == nil {
		return node
	}

	var child *RBNode
	if node.Left != nil && node.Right != nil {
		pred := node.Left.maximumNode()
		node.Cycle = pred.Cycle
		node.Queue = pred.Queue
		node = pred
	}

	if node.Left == nil || node.Right == nil {
		if node.Right == nil {
			child = node.Left
		} else {
			child = node.Right
		}
		if node.color == BLACK {
			node.color = nodeColor(child)
			t.deleteCase1(node)
		}
		t.replaceNode(node, child)
		if node.Parent == nil && child != nil {
			child.color = BLACK
		}
	}

	return node
}

func (t *RBTree) deleteCase1(node *RBNode) {
	if node.Parent == nil {
		return
	}
	t.deleteCase2(node)
}

func (t *RBTree) deleteCase2(node *RBNode) {
	sibling := node.sibling()
	if nodeColor(sibling) == RED {
		node.Parent.color = RED
		sibling.color = BLACK
		if node == node.Parent.Left {
			t.rotateLeft(node.Parent)
		} else {
			t.rotateRight(node.Parent)
		}
	}
	t.deleteCase3(node)
}

func (t *RBTree) deleteCase3(node *RBNode) {
	sibling := node.sibling()

	if nodeColor(node.Parent) == BLACK &&
		nodeColor(sibling) == BLACK &&
		nodeColor(sibling.Left) == BLACK &&
		nodeColor(sibling.Right) == BLACK {
		sibling.color = RED
		t.deleteCase1(node.Parent)
	} else {
		t.deleteCase4(node)
	}
}

func (t *RBTree) deleteCase4(node *RBNode) {
	sibling := node.sibling()

	if nodeColor(node.Parent) == RED &&
		nodeColor(sibling) == BLACK &&
		nodeColor(sibling.Left) == BLACK &&
		nodeColor(sibling.Right) == BLACK {
		sibling.color = RED
		node.Parent.color = BLACK
	} else {
		t.deleteCase5(node)
	}
}

func (t *RBTree) deleteCase5(node *RBNode) {
	sibling := node.sibling()

	if node == node.Parent.Left &&
		nodeColor(sibling) == BLACK &&
		nodeColor(sibling.Left) == RED &&
		nodeColor(sibling.Right) == BLACK {
		sibling.color = RED
		sibling.Left.color = BLACK
		t.rotateRight(sibling)
	} else if node == node.Parent.Right &&
		nodeColor(sibling) == BLACK &&
		nodeColor(sibling.Right) == RED &&
		nodeColor(sibling.Left) == BLACK {
		sibling.color = RED
		sibling.Right.color = BLACK
		t.rotateLeft(sibling)
	}

	t.deleteCase6(node)
}

func (t *RBTree) deleteCase6(node *RBNode) {
	sibling := node.sibling()
	sibling.color = nodeColor(node.Parent)
	node.Parent.color = BLACK
	if node == node.Parent.Left && nodeColor(sibling.Right) == RED {
		sibling.Right.color = BLACK
		t.rotateLeft(node.Parent)
	} else if nodeColor(sibling.Left) == RED {
		sibling.Left.color = BLACK
		t.rotateRight(node.Parent)
	}
}

func (node *RBNode) maximumNode() *RBNode {
	if node == nil {
		return nil
	}

	for node.Right != nil {
		node = node.Right
	}
	return node
}

func (t *RBTree) RemoveTask(task *Task) error {
	t.runningMu.Lock()
	defer t.runningMu.Unlock()

	node := t.lookup(task.Cycle)
	if node != nil {
		queue := node.Queue
		target := queue.Find(task, func(b *Task) bool {
			return task.Id == b.Id
		})

		queue.Remove(target.index)
	}
	return nil
}

func (t *RBTree) FindTaskBy(id string) *Task {
	t.runningMu.Lock()
	defer t.runningMu.Unlock()
	return FindTaskById(t.Root, id)
}

func FindTaskById(node *RBNode, id string) *Task {
	if node == nil {
		return nil
	}
	queue := node.Queue

	task := queue.Find(&Task{
		Id: id,
	}, func(b *Task) bool {
		return id == b.Id
	})

	if task != nil {
		queue.Remove(task.index)
		return task
	}

	task = FindTaskById(node.Left, id)
	if task != nil {
		return task
	}

	task = FindTaskById(node.Right, id)
	if task != nil {
		return task
	}

	return nil
}
