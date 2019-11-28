package main

import (
	"container/heap"
	"fmt"
)

func main() {
	var queue PriorityQueue

	item := &PriorityItem{Priority: 100}

	queue.Push(&PriorityItem{Priority: 1})
	queue.Push(&PriorityItem{Priority: 6})
	queue.Push(item)
	queue.Push(&PriorityItem{Priority: 9})
	queue.Push(&PriorityItem{Priority: 3})
	queue.Push(&PriorityItem{Priority: 4})
	queue.Push(&PriorityItem{Priority: 2})
	queue.Push(&PriorityItem{Priority: 8})
	queue.Push(&PriorityItem{Priority: 7})

	item.Priority = 5
	queue.Update(item)

	for !queue.Empty() {
		fmt.Println(queue.Pop().Priority)
	}
}

type PriorityItem struct {
	Priority int
	Index int
}

type PriorityStorage []*PriorityItem

func (s PriorityStorage) Len() int {
	return len(s)
}

func (s PriorityStorage) Less(i, j int) bool {
	// Highest priority is popped first.
	return s[i].Priority > s[j].Priority
}

func (s PriorityStorage) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
	s[i].Index, s[j].Index = i, j
}

func (s *PriorityStorage) Push(x interface{}) {
	item := x.(*PriorityItem)
	item.Index = len(*s)
	*s = append(*s, item)
}

func (s *PriorityStorage) Pop() interface{} {
	len := len(*s)
	item := (*s)[len-1]
	item.Index = -1
	*s = (*s)[:len-1]
	return item
}

type PriorityQueue struct {
	storage PriorityStorage
}

func (q *PriorityQueue) Len() int {
	return len(q.storage)
}

func (q *PriorityQueue) Empty() bool {
	return len(q.storage) == 0
}

func (q *PriorityQueue) Push(item *PriorityItem) {
	heap.Push(&q.storage, item)
}

func (q *PriorityQueue) Pop() *PriorityItem {
	return heap.Pop(&q.storage).(*PriorityItem)
}

func (q *PriorityQueue) Update(item *PriorityItem) {
	heap.Fix(&q.storage, item.Index)
}
