# BUG_REPRO：不存在的文章被误判成 500

## Bug 是什么

`CareArticleRepository.FindByID` / `Delete` 在未命中时用 `%v` 包装 `ErrNotFound`，把错误链断掉；`CareArticleService.Get/Delete/Update` 又用 `==` 直接比较，导致 `errors.Is(err, ErrNotFound)` 恒为 false，不存在文章最终被当成系统错误返回 500，而不是 404。

## 如何触发

1. 启动服务（backend 下 `go run ./cmd/server`，`DB_DRIVER=sqlite`、`DB_NAME=file::memory:?cache=shared`）。
2. 请求一个不存在的文章：`GET /api/v1/articles/999999`。
3. 或直接跑 `go test ./internal/service -run '^TestCareArticleGetMissing404$' -count=1`。

## 错误信息

不存在的文章返回 500，而不是 404（响应体为系统错误码）。
