package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func executeTask(c *gin.Context) {
	data, _ := c.GetRawData()

	log.Println(string(data))
}
