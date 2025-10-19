package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johnson7543/chatsheet-assessment/internal/database"
	"github.com/johnson7543/chatsheet-assessment/internal/models"
)

// GetAccounts retrieves all linked accounts for the authenticated user
func GetAccounts(c *gin.Context) {
	userID := c.GetUint("user_id")

	var accounts []models.LinkedAccount
	if err := database.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&accounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch accounts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accounts": accounts,
		"count":    len(accounts),
	})
}

// DeleteAccount removes a linked account
func DeleteAccount(c *gin.Context) {
	userID := c.GetUint("user_id")
	accountID := c.Param("id")

	var account models.LinkedAccount
	if err := database.DB.Where("id = ? AND user_id = ?", accountID, userID).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Account not found"})
		return
	}

	if err := database.DB.Delete(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to delete account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}
