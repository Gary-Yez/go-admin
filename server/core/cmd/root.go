package cmd

import (
	"gitee.com/mxcker/go-admin/server/core/internal"
	"gitee.com/mxcker/go-admin/server/core/internal/initialization"
	"gitee.com/mxcker/go-admin/server/global"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := initialization.InitConfig(cmd.PersistentFlags())
		if err != nil {
			panic("读取配置失败：" + err.Error())
		}
		err = initialization.InitGlobal()
		if err != nil {
			panic("全局变量初始化失败：" + err.Error())
		}
		if !global.Config.IsDev() {
			gin.SetMode("release")
		}
		//挂载核心组件
		internal.Run(true)
	},
}

func init() {
	rootCmd.PersistentFlags().String("server.host", "0.0.0.0", "Web服务运行的IP，默认0.0.0.0")
	rootCmd.PersistentFlags().String("server.port", "8080", "Web服务运行的端口，默认8080")
	rootCmd.PersistentFlags().String("config", "config.yaml", "配置文件路径（默认使用 ./config.yaml）")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
