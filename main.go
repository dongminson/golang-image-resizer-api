package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"strconv"
)

func main() {
	r := gin.Default()

	r.POST("/resize", func(c *gin.Context) {
		file, _, err := c.Request.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer file.Close()

		img, format, err := image.Decode(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode image"})
			return
		}

		width := c.Request.FormValue("width")
		height := c.Request.FormValue("height")

		resizedImg := resize.Resize(uint(MustAtoi(width)), uint(MustAtoi(height)), img, resize.Lanczos3)

		var buffer bytes.Buffer
		switch format {
		case "jpeg":
			if err := jpeg.Encode(&buffer, resizedImg, nil); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode resized image"})
				return
			}
			c.Data(http.StatusOK, "image/jpeg", buffer.Bytes())
		case "png":
			if err := png.Encode(&buffer, resizedImg); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode resized image"})
				return
			}
			c.Data(http.StatusOK, "image/png", buffer.Bytes())
		case "gif":
			if err := gif.Encode(&buffer, resizedImg, nil); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode resized image"})
				return
			}
			c.Data(http.StatusOK, "image/gif", buffer.Bytes())
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported image format"})
			return
		}
	})

	r.Run(":8080")
}

func MustAtoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
} 