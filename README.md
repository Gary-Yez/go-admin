# go-admin

可嵌入 Go 项目的管理后台核心。

```go
package main

import "github.com/Gary-Yez/go-admin"

func main() {
	if err := admin.Run(false); err != nil {
		panic(err)
	}
}
```

开发环境使用 `admin.Run(true)`，默认读取 `config.dev.yaml`；生产环境使用
`admin.Run(false)`，默认读取 `config.yaml`。

业务模块通过 `admin.MustRegister` 注册：

```go
admin.MustRegister("order", new(order.Mounter))
```

开发脚手架位于独立的 `go-admin-template` 仓库。
