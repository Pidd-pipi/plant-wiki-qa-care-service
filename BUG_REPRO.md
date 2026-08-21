# BUG_REPRO：限流器并发统计访问 data race

## Bug 是什么

`internal/middleware` 里的 `RateLimiter` 维护一个 `limits map[string]*bucket`，用 `mu sync.Mutex` 保护。
但新加的调试统计方法 `Snapshot / Reset / Count / Peek / ClearAll` 全部直接读写这个 map，没有加锁；同时 `Limit()` 在 `r.mu.Lock()` 之前先执行了一次 `b, ok := r.limits[ip]` 读操作。并发请求和调试统计一起发生时，就出现 map 并发读写 data race。

## 如何触发

1. 启动服务（backend 下 `go run ./cmd/server`，环境变量 `DB_DRIVER=sqlite`、`DB_NAME=file::memory:?cache=shared`）。
2. 并发请求任意被限流的接口（如 `GET /api/v1/plants`），同时反复请求调试统计接口 `GET /api/v1/debug/rate-limits`。
3. 用 `go test -race ./internal/middleware -run '^TestRateLimiterConcurrentAccess$' -count=1` 复现最稳定。

## 错误信息

```
WARNING: DATA RACE
Write at 0x00c000171a40 by goroutine 10:
  github.com/gbplantwiki/gbplantwiki/internal/middleware.(*RateLimiter).Reset()
      rate_limiter_observability.go:14
Previous read at 0x00c000171a40 by goroutine 108:
  github.com/gbplantwiki/gbplantwiki/internal/middleware.(*RateLimiter).ClearAll()
      rate_limiter_observability.go:33
```
