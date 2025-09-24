package main

// https://medium.com/@ar3s./priority-queue-in-go-lang-6185ad69c40a
type Item struct {
	value    string // Task name
	priority int64  // Priority (lower value = higher priority)
}

// PriorityQueue is a slice of *Item
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int {
	return len(pq)
}
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority // Min-heap
	// return pq[i].priority > pq[j].priority // Max-heap
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x any) {
	*pq = append(*pq, x.(*Item))
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}
