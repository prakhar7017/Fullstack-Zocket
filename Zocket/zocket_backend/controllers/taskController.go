package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"taskmanagement/config"
	models "taskmanagement/model"
	socketUtils "taskmanagement/socketUtils"
	"time"

	"github.com/gorilla/websocket"
)

func GetAllTasks(conn *websocket.Conn) {
	var tasks []models.Task
	config.DB.Find(&tasks)

	response := map[string]interface{}{
		"event": "task_list",
		"tasks": tasks,
	}

	message, _ := json.Marshal(response)
	conn.WriteMessage(websocket.TextMessage, message)
}
func CreateTask(conn *websocket.Conn, request map[string]interface{}) {

	title, _ := request["title"].(string)
	desc, _ := request["desc"].(string)

	
	var assigneeID uint
	if val, ok := request["assignee_id"].(float64); ok {
		assigneeID = uint(val)
	} else if val, ok := request["assignee_id"].(string); ok {
		if parsedID, err := strconv.Atoi(val); err == nil {
			assigneeID = uint(parsedID)
		} else {
			socketUtils.SendError(conn, "Invalid assignee_id format")
			return
		}
	} else {
		socketUtils.SendError(conn, "Assignee ID must be a number")
		return
	}

	
	var importance int
	if val, ok := request["importance"].(float64); ok {
		importance = int(val)
	} else if val, ok := request["importance"].(string); ok {
		if parsedImportance, err := strconv.Atoi(val); err == nil {
			importance = parsedImportance
		} else {
			socketUtils.SendError(conn, "Invalid importance format")
			return
		}
	} else {
		socketUtils.SendError(conn, "Importance must be a number")
		return
	}

	statusStr, _ := request["status"].(string)
	statusStr = strings.TrimSpace(strings.ToLower(statusStr))

	var status models.TaskStatus
	if err := status.FromString(statusStr); err != nil {
		fmt.Println("Invalid status:", statusStr)
		socketUtils.SendError(conn, "Invalid task status")
		return
	}

	
	var deadlineTime time.Time
	if deadlineStr, ok := request["deadline"].(string); ok {
		parsedTime, err := time.Parse("2006-01-02", deadlineStr)
		if err != nil {
			socketUtils.SendError(conn, "Invalid deadline format. Use YYYY-MM-DD")
			return
		}
		deadlineTime = parsedTime
	}

	
	if assigneeID != 0 {
		var user models.User
		if err := config.DB.First(&user, assigneeID).Error; err != nil {
			socketUtils.SendError(conn, "Assigned user not found")
			return
		}
	}

	
	task := models.Task{
		Title:       title,
		Description: desc,
		Status:      status,
		AssigneeID:  assigneeID,
		Deadline:    deadlineTime,
		Importance:  importance,
	}

	
	if err := config.DB.Create(&task).Error; err != nil {
		socketUtils.SendError(conn, "Failed to create task")
		return
	}

	
	response := map[string]interface{}{
		"event": "task_created",
		"task":  task,
	}

	message, _ := json.Marshal(response)
	socketUtils.BroadcastMessage(message)
	SendSlackNotification(fmt.Sprintf("New task created: %s", task.Title))
}


func UpdateTask(conn *websocket.Conn, request map[string]interface{}) {
	taskID, err := strconv.Atoi(fmt.Sprintf("%v", request["id"]))
	if err != nil {
		socketUtils.SendError(conn, "Invalid task ID")
		return
	}

	var task models.Task
	if err := config.DB.First(&task, taskID).Error; err != nil {
		socketUtils.SendError(conn, "Task not found")
		return
	}

	
	if title, ok := request["title"].(string); ok {
		task.Title = title
	}
	if desc, ok := request["desc"].(string); ok {
		task.Description = desc
	}
	var importanceval int
	if val, ok := request["importance"].(float64); ok {
		importanceval = int(val)
		task.Importance = importanceval
	} else if val, ok := request["importance"].(string); ok {
		if parsedImportance, err := strconv.Atoi(val); err == nil {
			importanceval = parsedImportance
			task.Importance = importanceval
		} else {
			socketUtils.SendError(conn, "Invalid importance format")
			return
		}
	} else {
		socketUtils.SendError(conn, "Importance must be a number")
		return
	}
	if deadline := request["deadline"].(string); deadline != "" {
		parsedTime, err := time.Parse("2006-01-02", deadline)
		if err != nil {
			socketUtils.SendError(conn, "Invalid deadline format. Use YYYY-MM-DD")
			return
		}
		task.Deadline = parsedTime
	}
	if statusStr, ok := request["status"].(string); ok {
		var status models.TaskStatus
		if err := status.FromString(statusStr); err != nil {
			fmt.Println("Invalid status", statusStr)
			socketUtils.SendError(conn, "Invalid status")
			return
		}
		task.Status = status
	}

	
	if err := config.DB.Save(&task).Error; err != nil {
		socketUtils.SendError(conn, "Failed to update task")
		return
	}

	
	response := map[string]interface{}{
		"event": "task_updated",
		"task":  task,
	}

	message, _ := json.Marshal(response)
	socketUtils.BroadcastMessage(message)
	SendSlackNotification(fmt.Sprintf("Task updated: %s", task.Title))
}

func SendSlackNotification(message string) {
	webhookURL := os.Getenv("SLACK_WEBHOOK_URL")

	payload := map[string]string{"text": message,}
	jsonData, _ := json.Marshal(payload)

	_, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending Slack notification:", err)
	}
}
