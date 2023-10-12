# Changelog

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
