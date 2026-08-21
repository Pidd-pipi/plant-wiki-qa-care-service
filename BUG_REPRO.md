# BUG_REPRO：不存在的提问被误判成 500

## Bug 是什么

`QuestionRepository.FindByID` 用 `%v` 包装 `ErrNotFound` 断链；`QuestionService.Get` 用 `==` 直接比较，不存在提问返回 500。

## 如何触发

打开一个不存在的提问详情，或跑 `go test ./internal/service -run '^TestQuestionGet404$' -count=1`。

## 错误信息

不存在的提问返回 500 而不是 404。
