package cmd

import (
	"fmt"
	"gitee.com/mxcker/go-admin/server/configs"
	"gitee.com/mxcker/go-admin/server/core"
	"gitee.com/mxcker/go-admin/server/core/internal/initialization"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		if !configs.IsDev() {
			gin.SetMode("release")
		}
		//创建服务器
		Server := gin.Default()
		corsConfig := cors.DefaultConfig()
		corsConfig.AllowAllOrigins = true
		corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
		Server.Use(cors.New(corsConfig))
		Server.Use(static.Serve(configs.Config.Server.AdminPrefix, static.LocalFile("./dist", false)))
		//挂载核心组件
		core.Load(Server, true)
		// 监听并启动服务
		listenAddr := viper.GetString("server.host") + ":" + viper.GetString("server.port")
		fmt.Println("程序运行在：" + listenAddr)
		fmt.Println("管理员入口：" + configs.Config.Server.AdminPrefix)
		err = Server.Run(listenAddr)
		if err != nil {
			panic(err)
		}
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
