/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"gitee.com/mxcker/go-admin/server/core/cmd"
)

func main() {
	// 使用标准化日志
	//slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	cmd.Execute()
}
