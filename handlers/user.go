package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muplat/muplat-backend/services/user"
)

func (h *HttpHandler) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := user.CreateSession(input.Username, input.Password, h.Db, h.Jwt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}

func (h *HttpHandler) AddUser(c *gin.Context) {
	var input UserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username, err := h.Jwt.ExtractTokenUsername(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = user.AddUser(
		input.Username,
		input.Password,
		username,
		*input.Admin,
		h.Db,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("added user %s", input.Username)})
}

func (h *HttpHandler) DeleteUser(c *gin.Context) {
	var input UserDeleteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username, err := h.Jwt.ExtractTokenUsername(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = user.DeleteUser(input.Username, username, h.Db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("deleted user %s", input.Username)})
}

// GetUser godoc
// @Summary      Get user by name
// @Description  Admin/owner only. Successful response shape:
//
//	{"username":"...", "isAdmin": "...", "projects": [{...},{...}]}
//
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string][]user.UserResponse
// @Router       /api/auth/user/{username} [get]
func (h *HttpHandler) GetUser(c *gin.Context) {
	username := c.Param("username")
	callerUsername, err := h.Jwt.ExtractTokenUsername(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response, err := user.GetUser(username, callerUsername, h.Db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetUsers godoc
// @Summary      Get all registered users
// @Description  Admin‑only. Successful response shape:
//
//	{"users":[{...}, {...}]}
//
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string][]user.UserResponse
// @Router       /api/auth/user [get]
func (h *HttpHandler) GetUsers(c *gin.Context) {
	callerUsername, err := h.Jwt.ExtractTokenUsername(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response, err := user.GetUsers(callerUsername, h.Db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
