# BUG_REPRO：删除不存在的花园条目空指针崩溃

## Bug 是什么

`UserGardenRepository.Find` / `FindByID` 未命中时返回 typed-nil 指针配 nil error；`UserGardenService.Remove` / `BindReminder` 判空失效，直接解引用触发 nil 指针 panic。

## 如何触发

删除一个不存在的花园条目，或给不存在的花园条目绑定提醒；直接跑 `go test ./internal/service -run '^TestUserGardenRemoveMissingNoPanic$' -count=1` 复现。

## 错误信息

panic: runtime error: invalid memory address or nil pointer dereference
