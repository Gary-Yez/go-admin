/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"gitee.com/mxcker/go-admin/server/core"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "启动服务器",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// 启动服务
		core.Start()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
