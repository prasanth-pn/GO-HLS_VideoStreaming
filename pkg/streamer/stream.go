package streamer

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Stream(c *gin.Context) {
	videoId := c.Param("video_id")
	playlist := c.Param("playlist")
	playlistDataChan := make(chan []byte)
	errChan := make(chan error)

	go func() {
		playlistData, err := readPlaylistData(videoId, playlist)
		if err != nil {
			errChan <- err
			return
		}
		playlistDataChan <- playlistData

	}()
	select {
	case playlistData := <-playlistDataChan:
		c.Header("Content-Type", "application/vnd.apple.mpegurl")
		c.Header("Content-Disposition", "inline")
		c.Writer.Write(playlistData)
	case err := <-errChan:
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "failed to read file from server",
			"errror":  err.Error(),
		})
	case <-time.After(5 * time.Second):
		c.JSON(http.StatusGatewayTimeout, gin.H{
			"message ": "request timeout",
		})
	}
}

func readPlaylistData(videoId, playlist string) ([]byte, error) {
	playlistPath := fmt.Sprintf("cmd/pkg/storage/%s/%s", videoId, playlist)
	playlistData, err := ioutil.ReadFile(playlistPath)
	if err != nil {
		return nil, err
	}
	return playlistData, nil

}
