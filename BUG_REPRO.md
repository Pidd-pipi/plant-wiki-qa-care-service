# BUG_REPRO：新增 snoozed 状态后提醒状态错乱

## Bug 是什么

新增 `snoozed` 中间态后，状态机转换表、状态文案、日志模板四处未同步；`UpdateStatus` 把 snoozed 错误落到 pending，文案返回未知。

## 如何触发

给提醒执行 snoozed 状态流转，或跑 `go test ./internal/service -run '^TestReminderSnoozeTransition$' -count=1`。

## 错误信息

snoozed 后状态错误、文案为未知。
