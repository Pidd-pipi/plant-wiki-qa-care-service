# BUG_REPRO：取消不存在的收藏被误判成 500

## Bug 是什么

`FavoriteRepository.Find` 用 `%v` 包装 `ErrNotFound` 断链；`FavoriteService.Remove` 用 `==` 直接比较，取消不存在的收藏返回 500。

## 如何触发

取消一个不存在的收藏，或跑 `go test ./internal/service -run '^TestFavoriteRemoveNotFound404$' -count=1`。

## 错误信息

不存在的收藏返回 500 而不是 404。
