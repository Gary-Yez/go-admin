/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"gitee.com/mxcker/go-admin/server/core"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "启动服务器",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		//创建服务器
		Server := gin.Default()
		corsConfig := cors.DefaultConfig()
		corsConfig.AllowAllOrigins = true
		corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
		Server.Use(cors.New(corsConfig))
		// 启动服务
		core.Start(Server)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
