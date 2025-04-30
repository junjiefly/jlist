package jlist

import (
	"reflect"
	"testing"
)

func TestNewList(t *testing.T) {
	capacity := 5
	list := NewList[string, int](capacity)
	if list == nil {
		t.Fatal("Failed to create new list")
	}
	if list.cap != capacity {
		t.Errorf("Expected data length %d, got %d", capacity, list.cap)
	}
	if len(list.freeIdx) != capacity {
		t.Errorf("Expected freeIdx size %d, got %d", capacity, len(list.freeIdx))
	}
}

func TestPushFront(t *testing.T) {
	list := NewList[string, int](3)
	testCases := []struct {
		key   string
		value int
	}{
		{"A", 1},
		{"B", 2},
		{"C", 3},
	}

	for i, tc := range testCases {
		entry, err := list.PushFront(tc.key, tc.value)
		if err != nil {
			t.Fatalf("Test %d: %v", i, err)
		}
		if entry.Key != tc.key || entry.Value != tc.value {
			t.Errorf("Test %d: Invalid entry values", i)
		}
	}

	if list.Len() != 3 {
		t.Errorf("Expected length 3, got %d", list.Len())
	}
}

func TestRemove(t *testing.T) {
	list := NewList[string, int](3)
	e1, _ := list.PushFront("A", 1)
	e2, _ := list.PushFront("B", 2)
	e3, _ := list.PushFront("C", 3)

	// Normal removal
	val, err := list.Remove(e2)
	if err != nil {
		t.Fatal(err)
	}
	if val != 2 {
		t.Errorf("Expected value 2, got %d", val)
	}
	if list.Len() != 2 {
		t.Errorf("Expected length 2, got %d", list.Len())
	}

	// Remove head
	val, err = list.Remove(e3)
	if val != 3 {
		t.Errorf("Expected value 3, got %d", val)
	}
	if list.Len() != 1 {
		t.Errorf("Expected length 1, got %d", list.Len())
	}

	// Remove last node
	val, err = list.Remove(e1)
	if val != 1 {
		t.Errorf("Expected value 1, got %d", val)
	}
	if list.Len() != 0 {
		t.Errorf("Expected length 0, got %d", list.Len())
	}
}

func TestMoveOperations(t *testing.T) {
	list := NewList[string, int](3)
	e1, _ := list.PushBack("A", 1)
	e2, _ := list.PushBack("B", 2)
	_, _ = list.PushBack("C", 3)

	list.MoveToFront(*e2)
	expected := []int{2, 1, 3}
	if !reflect.DeepEqual(list.Iterate(), expected) {
		t.Errorf("Unexpected order after MoveToFront")
	}

	list.MoveToBack(*e1)
	expected = []int{2, 3, 1}
	if !reflect.DeepEqual(list.Iterate(), expected) {
		t.Errorf("Unexpected order after MoveToBack")
	}
}

func TestInsertOperations(t *testing.T) {
	list := NewList[string, int](5)
	e1, _ := list.PushBack("A", 1)
	e2, _ := list.PushBack("B", 2)

	// Insert before
	_, err := list.InsertBefore("C", 0, *e1)
	if err != nil {
		t.Fatal(err)
	}
	expected := []int{0, 1, 2}
	if !reflect.DeepEqual(list.Iterate(), expected) {
		t.Errorf("Unexpected order after InsertBefore")
	}

	// Insert after
	_, err = list.InsertAfter("D", 3, *e2)
	if err != nil {
		t.Fatal(err)
	}
	expected = []int{0, 1, 2, 3}
	if !reflect.DeepEqual(list.Iterate(), expected) {
		t.Errorf("Unexpected order after InsertAfter")
	}
}

func TestErrorConditions(t *testing.T) {
	list := NewList[string, int](2)
	e1, _ := list.PushFront("A", 1)

	// Memory exhaustion
	_, err := list.PushFront("B", 2)
	if err != nil {
		t.Fatal(err)
	}
	_, err = list.PushFront("C", 3)
	if err == nil {
		t.Error("Expected memory pool exhausted error")
	}

	// Invalid node operations
	_, err = list.Remove(e1)
	if err != nil {
		t.Fatal(err)
	}
	_, err = list.Remove(e1)
	if err == nil {
		t.Error("Expected invalid node error")
	}
}

func TestIterate(t *testing.T) {
	list := NewList[string, int](4)
	e1, _ := list.PushBack("A", 1)
	e2, _ := list.PushBack("B", 2)
	list.MoveToFront(*e1)
	list.MoveToBack(*e2)

	expected := []int{1, 2}
	if !reflect.DeepEqual(list.Iterate(), expected) {
		t.Errorf("Unexpected iteration result")
	}
}

func Test0Iterate(t *testing.T) {
	list := NewList[string, int](1)
	e, _ := list.PushBack("A", 1)
	list.Remove(e)
	vs := list.Iterate()
	for range vs {
		t.Errorf("Expected len 0")
	}
}

func TestUpdateEntry(t *testing.T) {
	list := NewList[string, int](3)
	e1, _ := list.PushFront("A", 1)
	entry, _ := list.Entry(e1.idx)
	entry.Value = 100

	err := list.UpdateEntry(e1.idx, entry)
	if err != nil {
		t.Fatal(err)
	}

	updatedEntry, _ := list.Entry(e1.idx)
	if updatedEntry.Value != 100 {
		t.Errorf("Expected value 100, got %d", updatedEntry.Value)
	}
}

func TestClear(t *testing.T) {
	list := NewList[string, int](3)
	list.PushFront("A", 1)
	list.Clear()
	if list.Len() != 0 || list.size != 0 || len(list.freeIdx) != 0 {
		t.Error("Clear operation failed")
	}
}

func BenchmarkNewList(b *testing.B) {
	capacity := 1000
	for i := 0; i < b.N; i++ {
		NewList[int, int](capacity)
	}
}

func BenchmarkPushFront(b *testing.B) {
	capacity := 1000
	ll := NewList[int, int](capacity)
	for i := 0; i < b.N; i++ {
		ll.PushFront(i, i)
	}
}

func BenchmarkPushBack(b *testing.B) {
	capacity := 1000
	ll := NewList[int, int](capacity)
	for i := 0; i < b.N; i++ {
		ll.PushBack(i, i)
	}
}

func BenchmarkRemove(b *testing.B) {
	capacity := 1000
	ll := NewList[int, int](capacity)
	for i := 0; i < capacity; i++ {
		ll.PushBack(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		entry := ll.Front()
		ll.Remove(entry)
	}
}

func BenchmarkIterate(b *testing.B) {
	capacity := 1000
	ll := NewList[int, int](capacity)
	for i := 0; i < capacity; i++ {
		ll.PushBack(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ll.Iterate()
	}
}
