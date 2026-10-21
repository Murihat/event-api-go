package controllers

import (
	"context"
	"net/http"
	"os"
	"time"

	"example.com/event-app/config"
	"example.com/event-app/models"
	"github.com/gin-gonic/gin"
	"github.com/imagekit-developer/imagekit-go/v2"
	"github.com/imagekit-developer/imagekit-go/v2/option"
)

func initImageKit() *imagekit.Client {
	client := imagekit.NewClient(
		option.WithPrivateKey(os.Getenv("IMAGEKIT_PRIVATE_KEY")),
	)
	return &client
}

func CreateEvent(c *gin.Context) {
	userID, _ := c.Get("userID")

	// menerima file form data
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Image is required!",
			"error":   err.Error(),
		})
		return
	}

	defer file.Close()

	// 1. upload file ke imagekit
	fileName := header.Filename
	ik := initImageKit()
	uploadRes, errUpload := ik.Files.Upload(context.Background(), imagekit.FileUploadParams{
		File:     file,
		FileName: fileName,
	})

	if errUpload != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to upload Image",
			"error":   errUpload.Error(),
		})
		return
	}

	parsedTime, _ := time.Parse(time.RFC3339, c.PostForm("datetime"))
	//simpan ke database
	event := models.Event{
		Name:        c.PostForm("name"),
		Description: c.PostForm("description"),
		Location:    c.PostForm("location"),
		DateTIME:    parsedTime,
		Image:       uploadRes.URL,
		ImageID:     uploadRes.FileID,
		UserId:      userID.(int),
	}

	config.DB.Create(&event)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Successfully created event",
		"event":   event,
	})
}

func GetAllEvents(context *gin.Context) {
	var events []models.Event

	config.DB.Find(&events)
	context.JSON(http.StatusCreated, gin.H{
		"message": "Successfully get all events",
		"event":   events,
	})
}

func GetEventById(context *gin.Context) {
	var event models.Event
	paramsId := context.Param("id")

	var eventData = config.DB.First(&event, paramsId).Error
	if eventData != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"message": "Failed to get Event by ID",
			"error":   eventData.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Successfully get event by id",
		"event":   event,
	})
}

func UpdateEvent(c *gin.Context) {
	userID, _ := c.Get("userID")

	var event models.Event
	paramsId := c.Param("id")

	var eventData = config.DB.First(&event, paramsId).Error
	if eventData != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Failed to get Event by ID",
			"error":   eventData.Error(),
		})
		return
	}

	if event.UserId != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "You dont have access other user data",
		})
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err == nil {
		defer file.Close()

		ik := initImageKit()

		fileName := header.Filename

		uploadRes, errUpload := ik.Files.Upload(context.Background(), imagekit.FileUploadParams{
			File:     file,
			FileName: fileName,
		})

		if errUpload == nil {
			// hapus gambar lama
			if event.ImageID != "" {
				ik.Files.Delete(context.Background(), event.ImageID)
			}

			// upload file gambar
			event.Image = uploadRes.URL
			event.ImageID = uploadRes.FileID
		}
	}

	if name := c.PostForm("name"); name != "" {
		event.Name = name
	}

	if description := c.PostForm("description"); description != "" {
		event.Description = description
	}

	if location := c.PostForm("location"); location != "" {
		event.Location = location
	}

	if dateTimeStr := c.PostForm("datetime"); dateTimeStr != "" {
		parseTime, errParse := time.Parse(time.RFC3339, dateTimeStr)
		if errParse == nil {
			event.DateTIME = parseTime
		}
	}

	config.DB.Save(&event)
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully update event",
		"event":   event,
	})
}

func DeleteEvent(c *gin.Context) {
	userID, _ := c.Get("userID")

	var event models.Event
	paramsId := c.Param("id")

	var eventData = config.DB.First(&event, paramsId).Error
	if eventData != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Failed to get Event by ID",
			"error":   eventData.Error(),
		})
		return
	}

	if event.UserId != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "You dont have access other user data",
		})
		return
	}

	if event.ImageID != "" {
		ik := initImageKit()
		ik.Files.Delete(context.Background(), event.ImageID)
	}

	config.DB.Unscoped().Delete(&event)
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully delete event",
	})
}
