package main

import (
	"fmt"
	"sync/atomic"
)

var cacheNumber int64

func setNumber(n int) {
	atomic.StoreInt64(&cacheNumber, int64(n))
}

func getNumber() int {
	return int(atomic.LoadInt64(&cacheNumber))
}

func main() {
	setNumber(1325)
	fmt.Println("Cached Number:", getNumber()) // Output: Cached Number: 123
}
