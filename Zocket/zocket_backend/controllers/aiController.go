package controllers

import (
	"context"
	"fmt"
	"net/http"
	"taskmanagement/services"
	"time"

	"github.com/gin-gonic/gin"
)

func SuggestTask(c *gin.Context) {
	var input struct {
		Prompt string `json:"prompt"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	resultChan := make(chan string)
	errorChan := make(chan error)

	go func() {
		suggestion, err := services.SuggestTask(input.Prompt)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- suggestion
	}()

	select {
	case suggestion := <-resultChan:
		fmt.Println(suggestion)
		c.JSON(http.StatusOK, gin.H{"suggestion": suggestion})
	case err := <-errorChan:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI task suggestion failed", "details": err.Error()})
	case <-ctx.Done():
		c.JSON(http.StatusRequestTimeout, gin.H{"error": "Request timed out"})
	}
}

func BreakDown(c *gin.Context) {
	var input struct {
		Prompt string `json:"prompt"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	breakDown, err := services.TaskBreakDown(input.Prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI task breakdown failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Task BreakDown": breakDown})

}

func ProtizeTask(c *gin.Context) {
	var input struct {
		Tasks []struct {
			ID          int    `json:"id"`
			Title       string `json:"title"`
			Description string `json:"description"`
			Deadline    string `json:"deadline,omitempty"`
			Importance  int    `json:"importance"`
		} `json:"tasks"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	taskDetails := "Here are some tasks:\n"
	for _, task := range input.Tasks {
		taskDetails += fmt.Sprintf("- %s: %s (Importance: %d, Deadline: %s)\n",
			task.Title, task.Description, task.Importance, task.Deadline, task.ID)
	}

	prompt := "Prioritize these tasks by urgency and importance. Assign a priority level (High, Medium, Low):\n" + taskDetails

	priority, err := services.PrioritizeTasks(prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI task breakdown failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task_prioritization": priority})

}
