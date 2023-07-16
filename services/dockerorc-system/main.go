package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrhid6/dockerorc/services/dockerorc-system/nodes"
)

func main() {

	nodes.InitNodes()

	router := gin.Default()

	apiGroup := router.Group("/api")

	apiGroup.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	nodeGroup := apiGroup.Group("/nodes")
	nodeGroup.POST("/register", nodes.API_RegisterNode)
	nodeGroup.POST("/status", nodes.API_SetStatus)
	router.Run(":6443")
}
