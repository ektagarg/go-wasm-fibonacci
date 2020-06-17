package main

import (
	_ "image/jpeg"
	"log"
	"reflect"
	"runtime"
	"strconv"
	"syscall/js"
	"time"
	"unsafe"
)

func main() {
	done := make(chan struct{}, 0)

	log.Println("WASM Initialized")
	fibonacciFunc := js.FuncOf(fibonacci)
	js.Global().Set("fibonacci", fibonacciFunc)

	defer fibonacciFunc.Release()
	<-done
}

func fibonacci(this js.Value, args []js.Value) interface{} {
	doc := js.Global().Get("document")
	value1 := doc.Call("getElementById", args[0].String()).Get("value").String()
	n, _ := strconv.Atoi(value1)
	num := float64(n)

	result, timeTook := fib(num)
	// convert time to flaot
	finalResult := float64ToJs(result)

	doc.Call("getElementById", args[1].String()).Set("value", finalResult)
	doc.Call("getElementById", args[2].String()).Set("innerHTML", timeTook.String())

	return finalResult
}

func fib(n float64) ([]float64, time.Duration) {
	start := time.Now()
	var t1, t2, nextTerm, i float64
	t1 = 0
	t2 = 1
	nextTerm = 0
	a := []float64{}

	for i = 1; i <= n; i++ {
		if i == 1 {
			a = append(a, t1)
			continue
		}
		if i == 2 {
			a = append(a, t2)
			continue
		}
		nextTerm = t1 + t2
		t1 = t2
		t2 = nextTerm
		a = append(a, nextTerm)
	}
	return a, time.Since(start)
}

// float64ToJs is converting []uint64 to a new javascript js.Value.
func float64ToJs(s []float64) js.Value {
	a := js.Global().Get("Uint8Array").New(len(s) * 8)
	js.CopyBytesToJS(a, sliceToByteSlice(s))
	runtime.KeepAlive(s)
	buf := a.Get("buffer")
	return js.Global().Get("Float64Array").New(buf, a.Get("byteOffset"), a.Get("byteLength").Int()/8)
}

func sliceToByteSlice(s []float64) []byte {
	h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	h.Len *= 8
	h.Cap *= 8
	return *(*[]byte)(unsafe.Pointer(h))
}
