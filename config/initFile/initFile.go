package initFile

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io"
	"os"
)

// 初始化日志
func InitFile() {
	// Disable Console Color, you don't need console color when writing the logs to file.
	//gin.DisableConsoleColor() // 禁用控制台颜色，将日志写入文件时不需要控制台颜色。

	// Logging to a file.
	logFile := viper.GetString("app.logFile") // 获取配置文件中日志文件的路径
	f, _ := os.Create(logFile)                // 创建一个文件句柄，打开或创建指定路径的文件，忽略可能的错误
	gin.DefaultWriter = io.MultiWriter(f)     // 设置gin框架的默认写入器为创建的文件句柄
}
