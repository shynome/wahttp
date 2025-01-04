# Changelog

## [0.0.8] - 2025-01-04

- 修复: GoResponse 状态码为 204 但 body 不为 null 导致出错的情况

## [0.0.7] - 2025-01-03

- 添加了 JsWriter, 用于构建 js 和 golang 程序之间的 yamux stdio

## [0.0.6] - 2023-10-12

### Fix

- GoRequest 的 url 不全, 预期为`http://127.0.0.1/`却只有`/`导致的 panic
- GoRequest ContentLength 为 0 时不必设置 body

## [0.0.5] - 2023-10-12

### Add

- 添加了许可说明

## [0.0.4] - 2023-10-12

### Add

- 添加 GoRequest 和 JsResponse, 这样 Go 和 JS 互操作更方便了

## [0.0.3] - 2023-02-28

### Fix

- 支持 js abort signal
