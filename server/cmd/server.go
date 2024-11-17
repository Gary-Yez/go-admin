/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"gitee.com/mxcker/go-admin/server/global"
	"gitee.com/mxcker/go-admin/server/initialization"
	"gitee.com/mxcker/go-admin/server/router"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "启动服务器",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// 初始化全局变量
		global.Init()
		initialization.Init()
		router.Start()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
