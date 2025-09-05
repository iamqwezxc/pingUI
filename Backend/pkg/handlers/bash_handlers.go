// pkg/handlers/bash_handlers.go
package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iamqwezxc/pingUI/Backend/pkg/docker"
)

type BashRequest struct {
	Command string `json:"command"`
	Timeout int    `json:"timeout"`
}

func BashExecuteHandler(c *gin.Context) {
	var req BashRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Валидация команды
	if strings.TrimSpace(req.Command) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Command cannot be empty"})
		return
	}

	// Ограничение на длину команды
	if len(req.Command) > 1000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Command too long"})
		return
	}

	// Устанавливаем timeout по умолчанию
	if req.Timeout <= 0 || req.Timeout > 300 {
		req.Timeout = 30
	}

	executor, err := docker.NewDockerExecutor()
	if err != nil {
		log.Printf("Docker error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Docker initialization failed"})
		return
	}
	defer executor.Cleanup()

	// Билдим образ если нужно (можно вынести в отдельный endpoint)
	if err := executor.BuildImage(); err != nil {
		log.Printf("Build error: %v", err)
	}

	output, err := executor.ExecuteCommand(req.Command, req.Timeout)
	if err != nil {
		log.Printf("Execution error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Command execution failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"output":  output,
		"command": req.Command,
	})
}

func BashHealthHandler(c *gin.Context) {
	executor, err := docker.NewDockerExecutor()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "Docker unavailable"})
		return
	}
	defer executor.Cleanup()

	// Простая проверка
	output, err := executor.ExecuteCommand("echo 'Docker is working'", 10)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "Docker test failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"output": strings.TrimSpace(output),
	})
}
