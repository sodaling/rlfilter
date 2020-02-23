# rlfilter

利用限流算法构建两个简单的HttpHandleFilter。

## 令牌桶算法：

demo：_examples/tokenbucket/main.go

```go
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
```

> [基本思想受这启发](https://stackoverflow.com/questions/667508/whats-a-good-rate-limiting-algorithm/668327#668327)
>
> 算法部分[参考](https://github.com/bsm/ratelimit)，值得注意的是以下几点：
>
> 1. 并没有采用锁来控制并发，而是采用了原子操作。
> 2. 在计算allowance时候，并没有简单按照**allowance += time_passed * (rate / per);**,**而是改为用*current := atomic.AddUint64(&rl.allowance, passed*rate)**。计算用乘法替代。
> 3. 在算出超出max时候，回头设置allowance时候，并不是简单**atomic.StoreUint64(&t.allowance, max)**。因为：“ since `max` may change and you want to subtract the `current` (i.e. the already spent credits from it)”

