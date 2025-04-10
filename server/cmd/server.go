/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"gitee.com/mxcker/go-admin/server/core"
	"gitee.com/mxcker/go-admin/server/core/global"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"strconv"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "启动服务器",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if !global.IsDevelopment {
			gin.SetMode("release")
		}
		//创建服务器
		Server := gin.Default()
		corsConfig := cors.DefaultConfig()
		corsConfig.AllowAllOrigins = true
		corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
		Server.Use(cors.New(corsConfig))
		Server.Use(static.Serve(global.Config.Deploy.AdminPrefix, static.LocalFile("./dist", false)))
		//挂载核心组件
		core.Load(Server, true)
		// 监听并在 0.0.0.0:8080 上启动服务
		err := Server.Run(global.Config.Deploy.Listen.Host + ":" + strconv.Itoa(global.Config.Deploy.Listen.Port))
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
