package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/layemut/faceit-case-go/service"
)

// SaveUser saves a user by request
func SaveUser(collection service.ICollection) func(c *gin.Context) {
	return func(c *gin.Context) {
		var user service.User
		if err := c.BindJSON(&user); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		if len(user.ID) > 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User ID should be empty when creating a new user, if you wish to update an existing user, please use PUT method"})
			return
		}
		err := service.SaveUser(collection, &user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusCreated, &user)
	}
}

// UpdateUser updates a user by request
func UpdateUser(collection service.ICollection) func(c *gin.Context) {
	return func(c *gin.Context) {
		var user service.User
		if err := c.BindJSON(&user); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		if len(user.ID) == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		}
		err := service.UpdateUser(collection, &user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, &user)
	}
}

// RemoveUser removes a user by ID
func RemoveUser(collection service.ICollection) func(c *gin.Context) {
	return func(c *gin.Context) {
		ID := c.Param("id")
		if len(ID) == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
			return
		}
		if err := service.RemoveUser(collection, ID); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"message": "User removed"})
	}

}

// ListUsers returns a list of users based on the page request
func ListUsers(collection service.ICollection) func(c *gin.Context) {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.Query("page"))
		size, _ := strconv.Atoi(c.Query("size"))
		country := c.Query("country")

		userList, err := service.ListUsers(collection, &service.PageRequest{
			Page:    page,
			Size:    size,
			Country: country,
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, userList)
	}
}
