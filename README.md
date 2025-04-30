# jList
Pointer-less Doubly Linked List with GO


# Benchmark
```go
  go test -bench=BenchmarkTest  -benchmem
  goos: windows
  goarch: amd64
  pkg: lruTest/jlist
  cpu: Intel(R) Core(TM) i5-6200U CPU @ 2.30GHz
  BenchmarkTest-4         26526543                49.23 ns/op           40 B/op          0 allocs/op
  PASS
  ok      lruTest/jlist   3.247s
  PS E:\lruTest\jlist>
```

# example
```go
    package main

    import (
        "flag"
        "github.com/junjiefly/jlist"
    )
    
    func main() {
      list := NewList[string, int](3)
    	e1, _ := list.PushBack("A", 1)
    	e2, _ := list.PushBack("B", 2)
    	_, _ = list.PushBack("C", 3)
    	list.MoveToFront(*e2)
    	list.MoveToBack(*e1)
    }
```
