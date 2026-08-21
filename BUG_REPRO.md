# BUG_REPRO：更新不存在的用户资料空指针崩溃

## Bug 是什么

`UserRepository.FindByID` / `FindByUsername` 未命中时返回 typed-nil 指针配 nil error；`UserService.UpdateProfile/GetByID/Login` 判空失效，直接解引用触发 panic。

## 如何触发

更新一个不存在的用户资料，或跑 `go test ./internal/service -run '^TestUserUpdateMissingNoCrash$' -count=1`。

## 错误信息

panic: runtime error: invalid memory address or nil pointer dereference
