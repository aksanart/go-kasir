package config

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func Response(c *gin.Context, Message interface{}, Data interface{}) {
	if Data != nil {
		c.JSON(200, gin.H{
			"success": true,
			"message": fmt.Sprintf("%v", Message),
			"data":    Data,
		})
	} else {
		c.JSON(200, gin.H{
			"success": true,
			"message": fmt.Sprintf("%v", Message),
		})
	}
}
func ErrorResponse(c *gin.Context, Message interface{}) {
	log.Println(Message)
	c.AbortWithStatusJSON(400, gin.H{
		"success": false,
		"message": fmt.Sprintf("%v", Message),
	})

}
