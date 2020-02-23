package main

import (
	"github.com/sodaling/rlfilter"
	"net/http"
	"time"
)

func main() {
	// 设置qps为3
	var limiter = rlfilter.NewTokenBucket(3, time.Second)
	mux := http.DefaultServeMux

	mux.HandleFunc("/", limiter.Limit(test))
	http.ListenAndServe(":8000", mux)

}

func test(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("u r welcome"))
}
