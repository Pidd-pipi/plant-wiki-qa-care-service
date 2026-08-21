# BUG_REPRO：采纳或点赞不存在的回答被误判成 500

## Bug 是什么

`AnswerRepository.FindByID` 用 `%v` 包装 `ErrNotFound` 断链；`AnswerService.Adopt/Like` 用 `==` 直接比较，不存在回答返回 500。

## 如何触发

采纳或点赞一个不存在的回答，或跑 `go test ./internal/service -run '^TestAnswerAdoptMissingAnswer404$' -count=1`。

## 错误信息

不存在的回答返回 500 而不是 404。
