package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pahulgogna/pgexec/customTypes"
	"github.com/pahulgogna/pgexec/executor"
	"github.com/pahulgogna/pgexec/utils"
)

func main() {

	var router = gin.Default()

	router.GET("/ping", ping)
	router.POST("/run", handleExecution)

	router.Run(fmt.Sprintf("0.0.0.0:%s", "8080"))
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func handleExecution(c *gin.Context) {
	var req customtypes.Snippet

	if err := c.BindJSON(&req); err != nil {
		fmt.Printf("Error: could not parse the request body: %v\n", err)
		c.AbortWithStatus(400)
		return
	}

	if req.Code == "" || req.Language == "" {
		c.JSON(400, gin.H{
			"data": "incomplete data provided.",
		})
		return
	}

	tag := utils.GetTagFromLanguage(req.Language)

	if tag == "" {
		c.JSON(400, gin.H{
			"data": "unsupported language.",
		})
		return
	}

	dockerContainerID, err := executor.StartDockerContainer(tag)

	

	if err != nil {
		c.JSON(400, gin.H{
			"data": "could not start the isolated environment.",
		})
		return
	}

	response := executor.InstallDependencies(dockerContainerID, &req)

	if response != "Done" && response != "None" {
		c.JSON(400, gin.H{
			"data": response,
		})
	} else {
		output, err := executor.RunCode(dockerContainerID, &req)
		if err != nil {
			utils.WriteToLogFile(dockerContainerID, "ERROR: could not run code")
			c.JSON(400, gin.H{
				"data":  output,
				"error": err,
			})
		} else {
			c.JSON(200, gin.H{
				"data":  output,
				"error": err,
			})
		}
	}

	go executor.StopDockerContainer(dockerContainerID)
}
