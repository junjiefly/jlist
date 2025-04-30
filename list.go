package jlist

import (
	"errors"
)

type Entry[K comparable, V any] struct {
	HashId uint64
	prev   int
	next   int
	idx    int
	Key    K
	Value  V
}

func (e Entry[K, V]) Idx() int {
	return e.idx
}

type List[K comparable, V any] struct {
	cap     int
	data    []Entry[K, V]
	freeIdx []int
	head    int
	tail    int
	size    int
}

func NewList[K comparable, V any](capacity int) *List[K, V] {
	ll := &List[K, V]{
		freeIdx: make([]int, capacity),
		head:    -1,
		tail:    -1,
		size:    0,
		cap:     capacity,
	}
	ll.data = make([]Entry[K, V], capacity)
	for i := 0; i < capacity; i++ {
		ll.freeIdx[i] = i
	}
	return ll
}

func (l *List[K, V]) getNodeIdx() (int, bool) {
	if len(l.freeIdx) == 0 {
		return -1, false
	}
	idx := l.freeIdx[len(l.freeIdx)-1]       // get last
	l.freeIdx = l.freeIdx[:len(l.freeIdx)-1] // eject last
	return idx, true
}

func (l *List[K, V]) putNodeIdx(idx int) {
	l.freeIdx = append(l.freeIdx, idx) // append to end
}

// Len returns the number of elements of list l.
// The complexity is O(1).
func (l *List[K, V]) Len() int {
	return l.size
}

func (l *List[K, V]) Cap() int {
	return l.cap
}

// Front returns the first element of list l or nil if the list is empty.
func (l *List[K, V]) Front() *Entry[K, V] {
	if l.size == 0 {
		return nil
	}
	if l.head != -1 {
		return &l.data[l.head]
	}
	return nil
}

// Back returns the last element of list l or nil if the list is empty.
func (l *List[K, V]) Back() *Entry[K, V] {
	if l.size == 0 {
		return nil
	}
	if l.tail != -1 {
		return &l.data[l.tail]
	}
	return nil
}

func (l *List[K, V]) remove(e *Entry[K, V]) (V, error) {
	if e == nil {
		var empty V
		return empty, errors.New("node null")
	}
	if e.next == -1 || e.prev == -1 {
		return e.Value, errors.New("unknown node")
	}
	if l.cap <= e.idx || e.idx < 0 {
		return e.Value, errors.New("invalid node")
	}
	node := l.data[e.idx]
	if node.prev == -1 || node.next == -1 {
		return e.Value, errors.New("invalid node")
	}
	if e.next != node.next || e.prev != node.prev {
		return e.Value, errors.New("list changed")
	}
	value := e.Value
	l.data[node.prev].next = node.next
	l.data[node.next].prev = node.prev

	if l.head == e.idx {
		l.head = node.next
	}
	if l.tail == e.idx {
		l.tail = node.prev
	}
	l.data[e.idx].prev = -1
	l.data[e.idx].next = -1
	l.putNodeIdx(e.idx)
	l.size--
	if l.size == 0 {
		l.head = -1
		l.tail = -1
	}
	return value, nil
}

func (l *List[K, V]) Remove(e *Entry[K, V]) (V, error) {
	return l.remove(e)
}

func (l *List[K, V]) PushFront(key K, value V) (*Entry[K, V], error) {
	idx, ok := l.getNodeIdx()
	if !ok {
		return nil, errors.New("memory pool exhausted")
	}
	l.data[idx].Key = key
	l.data[idx].Value = value
	l.data[idx].idx = idx

	if l.head != -1 && l.tail != -1 {
		l.data[l.head].prev = idx
		l.data[l.tail].next = idx
		l.data[idx].next = l.head
		l.data[idx].prev = l.tail
		l.head = idx
	} else {
		l.data[idx].next = idx
		l.data[idx].prev = idx
		l.head = idx
		l.tail = idx
	}
	l.size++
	return &l.data[idx], nil
}

func (l *List[K, V]) PushBack(key K, value V) (*Entry[K, V], error) {
	idx, ok := l.getNodeIdx()
	if !ok {
		return nil, errors.New("memory pool exhausted")
	}
	l.data[idx].Key = key
	l.data[idx].Value = value
	l.data[idx].idx = idx

	if l.head != -1 && l.tail != -1 {
		l.data[l.tail].next = idx
		l.data[l.head].prev = idx
		l.data[idx].prev = l.tail
		l.data[idx].next = l.head
		l.tail = idx
	} else {
		l.data[idx].next = idx
		l.data[idx].prev = idx
		l.head = idx
		l.tail = idx
	}
	l.size++
	return &l.data[idx], nil
}

func (l *List[K, V]) InsertBefore(k K, v V, mark Entry[K, V]) (*Entry[K, V], error) {
	return l.insertBefore(k, v, mark)
}

func (l *List[K, V]) insertBefore(k K, v V, mark Entry[K, V]) (*Entry[K, V], error) {
	idx, ok := l.getNodeIdx()
	if !ok {
		return nil, errors.New("memory pool exhausted")
	}
	if l.cap <= mark.idx || mark.idx < 0 {
		return nil, errors.New("invalid node")
	}
	markNode := l.data[mark.idx]
	if markNode.prev == -1 || markNode.next == -1 {
		return nil, errors.New("invalid node")
	}
	l.data[idx].Key = k
	l.data[idx].Value = v
	l.data[idx].idx = idx

	l.data[idx].next = markNode.idx
	l.data[idx].prev = markNode.prev

	l.data[markNode.prev].next = idx
	l.data[markNode.idx].prev = idx

	if l.head == markNode.idx || l.head == -1 {
		l.head = idx
	}
	l.size++
	return &l.data[idx], nil
}

func (l *List[K, V]) InsertAfter(k K, v V, mark Entry[K, V]) (*Entry[K, V], error) {
	return l.insertAfter(k, v, mark)
}

func (l *List[K, V]) insertAfter(k K, v V, mark Entry[K, V]) (*Entry[K, V], error) {
	idx, ok := l.getNodeIdx()
	if !ok {
		return nil, errors.New("memory pool exhausted")
	}
	if l.cap <= mark.idx || mark.idx < 0 {
		return nil, errors.New("invalid node")
	}
	markNode := l.data[mark.idx]
	if markNode.prev == -1 || markNode.next == -1 {
		return nil, errors.New("invalid node")
	}
	l.data[idx].Key = k
	l.data[idx].Value = v
	l.data[idx].idx = idx

	l.data[idx].prev = markNode.idx
	l.data[idx].next = markNode.next

	l.data[markNode.next].prev = idx
	l.data[markNode.idx].next = idx

	if l.tail == markNode.idx || l.tail == -1 {
		l.tail = idx
	}
	l.size++
	return &l.data[idx], nil
}

// MoveToFront moves element e to the front of list l.
// If e is not an element of l, the list is not modified.
// The element must not be nil.
func (l *List[K, V]) MoveToFront(e Entry[K, V]) error {
	if l.cap <= e.idx || e.idx < 0 {
		return errors.New("invalid node")
	}
	markNode := l.data[e.idx]
	if markNode.prev == -1 || markNode.next == -1 {
		return errors.New("invalid node")
	}
	if e.idx == l.head {
		return nil
	}

	l.data[e.prev].next = e.next
	l.data[e.next].prev = e.prev

	l.data[l.tail].next = e.idx
	l.data[l.head].prev = e.idx

	l.data[e.idx].next = l.head
	l.data[e.idx].prev = l.tail

	l.head = e.idx
	return nil
}

// MoveToBack moves element e to the back of list l.
// If e is not an element of l, the list is not modified.
// The element must not be nil.
func (l *List[K, V]) MoveToBack(e Entry[K, V]) error {
	if l.cap <= e.idx || e.idx < 0 {
		return errors.New("invalid node")
	}
	markNode := l.data[e.idx]
	if markNode.prev == -1 || markNode.next == -1 {
		return errors.New("invalid node")
	}
	if e.idx == l.tail {
		return nil
	}
	l.data[e.prev].next = e.next
	l.data[e.next].prev = e.prev

	l.data[l.tail].next = e.idx
	l.data[l.head].prev = e.idx

	l.data[e.idx].next = l.head
	l.data[e.idx].prev = l.tail

	l.tail = e.idx
	return nil
}

func (l *List[K, V]) Iterate() []V {
	var result []V
	current := l.head
	for current != -1 {
		result = append(result, l.data[current].Value)
		current = l.data[current].next
		if current == l.head {
			break
		}
	}
	return result
}

func (l *List[K, V]) Entry(idx int) (Entry[K, V], error) {
	if l.cap <= idx || idx < 0 {
		return Entry[K, V]{}, errors.New("invalid node")
	}
	markNode := l.data[idx]
	if markNode.prev == -1 || markNode.next == -1 {
		return Entry[K, V]{}, errors.New("invalid node")
	}
	return l.data[idx], nil
}

func (l *List[K, V]) UpdateEntry(idx int, e Entry[K, V]) error {
	if l.cap <= idx || idx < 0 {
		return errors.New("invalid node")
	}
	markNode := l.data[idx]
	if markNode.prev == -1 || markNode.next == -1 {
		return errors.New("invalid node")
	}
	l.data[idx].Key = e.Key
	l.data[idx].HashId = e.HashId
	l.data[idx].Value = e.Value
	return nil
}

func (l *List[K, V]) Clear() {
	l.data = nil
	l.freeIdx = nil
	l.head = -1
	l.tail = -1
	l.size = 0
}
