package main

import (
	_ "image/jpeg"
	"log"
	"reflect"
	"runtime"
	"syscall/js"
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
	// convert js.value to go integer
	n := args[0].Int()
	num := float64(n)

	result := fib(num)
	final := intToFloat64(result)
	return final
}

func fib(n float64) []float64 {
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
	return a
}

// intToFloat64 is converting []uint64 to a new javascript Int32Array.
func intToFloat64(s []float64) js.Value {
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
