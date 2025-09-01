package handlers

import (
	"fmt"
	"matching_system/internal/api/dto"
	"matching_system/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MatchHandler struct {
	matchService services.MatchService
}

func NewMatchHandler() *MatchHandler {
	return &MatchHandler{
		matchService: services.NewMatchService(),
	}
}

func (h *MatchHandler) AddSinglePersonAndMatch(c *gin.Context) {
	var req dto.AddPersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	person, matches := h.matchService.AddSinglePersonAndMatch(req)

	c.JSON(http.StatusCreated, dto.AddPersonResponse{
		Person:  *person,
		Matches: matches,
		Message: "person added successfully",
	})
}

func (h *MatchHandler) RemoveSinglePerson(c *gin.Context) {
	personID := c.Param("id")
	fmt.Println("personID:", personID)
	if personID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "person ID is required"})
		return
	}

	success := h.matchService.RemoveSinglePerson(personID)
	if !success {
		c.JSON(http.StatusNotFound, dto.RemovePersonResponse{
			Success: false,
			Message: "person not found",
		})
		return
	}

	c.JSON(http.StatusOK, dto.RemovePersonResponse{
		Success: true,
		Message: "person removed successfully",
	})
}

func (h *MatchHandler) QuerySinglePeople(c *gin.Context) {
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.QueryPeopleResponse{
			Message: "limit is required",
		})
		return
	}
	people := h.matchService.QuerySinglePeople(limit)
	c.JSON(http.StatusOK, dto.QueryPeopleResponse{
		People:  people,
		Message: "people queried successfully",
	})
}
