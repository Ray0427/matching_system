package services

import (
	"encoding/json"
	"matching_system/internal/api/dto"
	"matching_system/internal/models"
	"matching_system/pkg/logger"
	"sync"

	"github.com/google/uuid"
)

type MatchService struct {
	mu           sync.RWMutex
	activePeople map[string]*models.Person
	logger       *logger.Logger
}

func NewMatchService() *MatchService {
	return &MatchService{
		activePeople: make(map[string]*models.Person),
		logger:       logger.New(),
	}
}

func (ms *MatchService) AddSinglePersonAndMatch(req dto.AddPersonRequest) (*models.Person, []*models.Match) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	person := &models.Person{
		ID:            uuid.New().String(),
		Name:          req.Name,
		Height:        req.Height,
		Gender:        req.Gender,
		NumberOfDates: req.NumberOfDates,
	}

	ms.activePeople[person.ID] = person
	jsonData, _ := json.MarshalIndent(ms.activePeople, "", "  ")
	ms.logger.Info("Active people:\n" + string(jsonData))

	matchs := ms.findMatches(person)

	return person, matchs
}

func (ms *MatchService) RemoveSinglePerson(personID string) bool {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, ok := ms.activePeople[personID]; ok {
		delete(ms.activePeople, personID)
	} else {
		return false
	}

	jsonData, _ := json.MarshalIndent(ms.activePeople, "", "  ")
	ms.logger.Info("Active people:\n" + string(jsonData))

	return true
}

func (ms *MatchService) QuerySinglePeople(limit int) ([]*models.Person, int) {
	people := make([]*models.Person, 0)
	total := 0

	for _, person := range ms.activePeople {
		people = append(people, person)
		total++
	}

	return people, total
}

func (ms *MatchService) findMatches(person *models.Person) []*models.Match {
	matchs := make([]*models.Match, 0)

	return matchs
}
