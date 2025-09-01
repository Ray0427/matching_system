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

// AddSinglePersonAndMatch godoc
// @Summary Add a single person and match
// @Description Add a single person and match
// @Tags match
// @Accept json
// @Produce json
// @Param person body dto.AddPersonRequest true "Person"
// @Success 201 {object} dto.AddPersonResponse
// @Router /add-single-person-and-match [post]
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

// RemoveSinglePerson godoc
// @Summary Remove a single person
// @Description Remove a single person
// @Tags match
// @Accept json
// @Produce json
// @Param id path string true "Person ID"
// @Success 200 {object} dto.RemovePersonResponse
// @Router /remove-single-person/{id} [delete]
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

// QuerySinglePeople godoc
// @Summary Query single people
// @Description Query single people
// @Tags match
// @Accept json
// @Produce json
// @Param limit query int true "Limit"
// @Success 200 {object} dto.QueryPeopleResponse
// @Router /query-single-people [get]
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
