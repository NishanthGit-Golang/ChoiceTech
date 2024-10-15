package handler

import (
	"choice/config"
	"choice/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

func UploadExcel(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is uploaded"})
		return
	}
	savePath := filepath.Join(".", file.Filename)
	err = c.SaveUploadedFile(file, savePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
		return
	}
	excelFile, err := excelize.OpenFile(savePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open Excel file", "details": err.Error()})
		return
	}
	sheetMap := excelFile.GetSheetMap()
	if len(sheetMap) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No sheets found in Excel file"})
		return
	}
	sheetName := sheetMap[1]
	rows, err := excelFile.GetRows(sheetName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read Excel file", "details": err.Error()})
		return
	}
	var records []models.Record
	for row, columns := range rows[1:] {
		if len(columns) < 10 {
			continue
		}
		employee := models.Record{
			ID:          row + 1,
			FirstName:   columns[0],
			LastName:    columns[1],
			Company:     columns[2],
			Address:     columns[3],
			City:        columns[4],
			Country:     columns[5],
			Postal:      columns[6],
			PhoneNumber: columns[7],
			EmailID:     columns[8],
			WebLink:     columns[9],
		}

		records = append(records, employee)
	}
	go func() {
		if len(records) > 0 {
			if err := InsertRecordsBatch(records); err != nil {
				log.Println("Failed to insert records:", err)
				c.JSON(400, gin.H{"error": "Error in inserting the Records"})
				return
			}
		}
	}()
	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully and data stored"})
}
func InsertRecordsBatch(records []models.Record) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		for _, record := range records {
			if err := tx.Create(&record).Error; err != nil {
				return err
			}
			recordJSON, err := json.Marshal(record)
			if err != nil {
				log.Println("Failed to marshal record:", err)

				continue
			}
			cacheKey := fmt.Sprintf("record:%d", record.ID)
			if err := config.SetCache(cacheKey, string(recordJSON)); err != nil {
				log.Println("Failed to cache record in Redis:", err)
			}
		}
		return nil
	})
}

func GetRecords(c *gin.Context) {
	var records []models.Record
	cachedData, err := config.GetCache("record")
	if err == nil {
		json.Unmarshal([]byte(cachedData), &records)
		c.JSON(200, records)
		return
	}
	if err := config.DB.Find(&records).Error; err != nil {
		c.JSON(500, gin.H{"error": "Error fetching records from MySQL", "details": err.Error()})
		return
	}
	if len(records) == 0 {
		c.JSON(404, gin.H{"error": "No records found"})
		return
	}
	recordJSON, _ := json.Marshal(records)
	config.SetCache("record", string(recordJSON))
	c.JSON(200, records)
}
func GetRecordByID(c *gin.Context) {
	var record models.Record
	recordID := c.Param("id")
	cacheKey := fmt.Sprintf("record:%s", recordID)
	cachedData, err := config.GetCache(cacheKey)
	if err == nil {
		json.Unmarshal([]byte(cachedData), &record)
		c.JSON(200, record)
		return
	}
	if err := config.DB.Where("id = ?", recordID).First(&record).Error; err != nil {
		c.JSON(404, gin.H{"error": "Record not found"})
		return
	}
	recordJSON, _ := json.Marshal(record)
	config.SetCache(cacheKey, string(recordJSON))
	c.JSON(200, record)
}

func EditRecord(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var record models.Record
	if err := config.DB.First(&record, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Record not found"})
		return
	}
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	config.DB.Save(&record)
	recordJSON, _ := json.Marshal(record)
	cacheKey := fmt.Sprintf("record:%d", record.ID)
	config.SetCache(cacheKey, string(recordJSON))
	c.JSON(200, gin.H{"message": "Record updated successfully"})
}

func DeleteRecordByID(c *gin.Context) {
	recordID := c.Param("id")
	result := config.DB.Where("id = ?", recordID).Delete(&models.Record{})
	if result.Error != nil {
		c.JSON(400, gin.H{"error": "Invalid record ID", "details": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Record not found", "details": "No record with the given ID"})
		return
	}
	cacheKey := fmt.Sprintf("record:%s", recordID)
	err := config.RDB.Del(context.Background(), cacheKey).Err()
	if err != nil {
		log.Println("Warning: Error deleting record from Redis:", err)
	}
	c.JSON(200, gin.H{"message": "Record deleted successfully"})
}

func DeleteAllRecords(c *gin.Context) {
	rowsAffected := config.DB.Exec("DELETE FROM records").RowsAffected
	if rowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error deleting records"})
		return
	}
	err := config.RDB.Del(context.Background(), "record").Err()
	if err != nil {
		log.Println("Warning: Error deleting record from Redis:", err)
	}
	c.JSON(http.StatusOK, gin.H{"message": "All records deleted successfully", "rowsAffected": rowsAffected})
}
