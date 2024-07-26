package http

import (
	"fmt"
	"net/http"
	"server/config"

	"github.com/gin-gonic/gin"
)

const (
	httpPort = 8888
)

func init() {
	// gin.SetMode(gin.ReleaseMode)
}

func Start() {
	r := gin.Default()
	fmt.Println("Starting http server on port: ", httpPort)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello,World!",
		})
	})

	r.GET("/world-stat", func(c *gin.Context) {
		file, err := getWorldFile()

		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		fileInfo, err := file.Stat()
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		file.Close()

		fileSize := fileInfo.Size()

		c.JSON(200, gin.H{
			"size":  fileSize,
			"world": config.WorldName,
		})
	})

	r.GET("/download-world", func(c *gin.Context) {
		c.Header("World", "island")
		c.File(config.WorldFilePath)
	})

	r.GET("/Assets/*filepath", func(ctx *gin.Context) {
		localDir := "/Users/ice/MMO/ServerData"
		filePath := ctx.Param("filepath")
		http.ServeFile(ctx.Writer, ctx.Request, localDir+filePath)
	})

	r.Run(fmt.Sprintf(":%d", httpPort))
}
