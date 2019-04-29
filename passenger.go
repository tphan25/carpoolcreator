package main

import (
	"container/heap"
	"fmt"
)

//For parsing of distancematrix data
type Passenger struct {
	index    int
	distance int
}

type PassengerHeap []Passenger

func (heap PassengerHeap) Len() int           { return len(heap) }
func (heap PassengerHeap) Less(i, j int) bool { return heap[i].distance < heap[j].distance }
func (heap PassengerHeap) Swap(i, j int)      { heap[i], heap[j] = heap[j], heap[i] }

func (h *PassengerHeap) Push(x interface{}) {
	*h = append(*h, x.(Passenger))
}

func (h *PassengerHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func testHeap() {
	h := &PassengerHeap{}

	p1 := Passenger{100, -2}
	p2 := Passenger{1, 1}
	p3 := Passenger{500, 3}
	p4 := Passenger{-100, 2}
	heap.Init(h)
	heap.Push(h, p1)
	heap.Push(h, p2)
	heap.Push(h, p3)
	heap.Push(h, p4)
	fmt.Printf("minimum: %d\n", (*h)[0])
	for h.Len() > 0 {
		fmt.Printf("%d ", heap.Pop(h))
	}
}
