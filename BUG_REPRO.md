# BUG_REPRO：不存在的植物品种被误判成 500

## Bug 是什么

`PlantSpeciesRepository.FindByID` / `Delete` 用 `%v` 包装 `ErrNotFound` 断链；`PlantSpeciesService.Get/Delete/Update` 用 `==` 直接比较，`errors.Is` 失效，不存在品种返回 500。

## 如何触发

请求不存在的植物品种，或跑 `go test ./internal/service -run '^TestPlantGet404$' -count=1`。

## 错误信息

不存在的品种返回 500 而不是 404。
