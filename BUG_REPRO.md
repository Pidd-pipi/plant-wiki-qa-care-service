# BUG_REPRO：不存在的病虫害条目被误判成 500

## Bug 是什么

`DiseasePestRepository.FindByID` / `Delete` 用 `%v` 包装 `ErrNotFound` 断链；`DiseasePestService.Get/Delete/Update` 用 `==` 直接比较，不存在条目返回 500。

## 如何触发

请求不存在的病虫害条目，或跑 `go test ./internal/service -run '^TestPestGet404$' -count=1`。

## 错误信息

不存在的条目返回 500 而不是 404。
