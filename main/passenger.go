package main

import (
	"container/heap"
	"fmt"
)

/*Passenger represents a tuple of address and distance*/
type Passenger struct {
	address  string
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

	p1 := Passenger{"p1", -2}
	p2 := Passenger{"p2", 1}
	p3 := Passenger{"p3", 3}
	p4 := Passenger{"p4", 2}
	heap.Init(h)
	heap.Push(h, p1)
	heap.Push(h, p2)
	heap.Push(h, p3)
	heap.Push(h, p4)
	fmt.Printf("minimum: %d\n", (*h)[0].distance)
	for h.Len() > 0 {
		fmt.Printf("%d ", heap.Pop(h))
	}
}
