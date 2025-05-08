package main

import "github.com/gin-gonic/gin"

func InitializeServer() (*gin.Engine, error) {
	return initializeServer("../../config")
}
