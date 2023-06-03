package main

import (
	"github.com/gin-gonic/gin"
	"github.com/prasanth-pn/GO-HLS_VideoStreaming/pkg/streamer"
	"github.com/prasanth-pn/GO-HLS_VideoStreaming/pkg/uploader"
)

func main() {
	route := gin.Default()

	route.POST("/upload", uploader.Upload)
	route.GET("play/:video_id/:playlist", streamer.Stream)
	route.Run(":8000")
}
