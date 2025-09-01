package services

import (
	"matching_system/internal/api/dto"
	"matching_system/internal/models"
	"matching_system/pkg/logger"
	"sort"
	"sync"

	"github.com/google/uuid"
)

type MatchService interface {
	AddSinglePersonAndMatch(req dto.AddPersonRequest) (*models.Person, []models.Match)
	RemoveSinglePerson(personID string) bool
	QuerySinglePeople(limit int) []models.Person
}

type matchService struct {
	mu           sync.RWMutex
	activePeople map[string]*models.Person
	logger       *logger.Logger
}

func NewMatchService() MatchService {
	return &matchService{
		activePeople: make(map[string]*models.Person),
		logger:       logger.New(),
	}
}

func (ms *matchService) AddSinglePersonAndMatch(req dto.AddPersonRequest) (*models.Person, []models.Match) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	person := &models.Person{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Height:      req.Height,
		Gender:      req.Gender,
		WantedDates: req.WantedDates,
	}

	ms.activePeople[person.ID] = person
	// jsonData, _ := json.MarshalIndent(ms.activePeople, "", "  ")
	// ms.logger.Info("Active people:\n" + string(jsonData))

	matches := ms.findMatches(person)

	return person, matches
}

func (ms *matchService) RemoveSinglePerson(personID string) bool {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, ok := ms.activePeople[personID]; ok {
		delete(ms.activePeople, personID)
	} else {
		return false
	}

	// jsonData, _ := json.MarshalIndent(ms.activePeople, "", "  ")
	// ms.logger.Info("Active people:\n" + string(jsonData))

	return true
}

func (ms *matchService) QuerySinglePeople(limit int) []models.Person {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	// Convert map to slice
	people := make([]models.Person, 0, len(ms.activePeople))
	for _, person := range ms.activePeople {
		people = append(people, *person)
	}

	sort.Slice(people, func(i, j int) bool {
		// sort by wanted dates
		if people[i].WantedDates != people[j].WantedDates {
			return people[i].WantedDates > people[j].WantedDates
		}

		// sort by gender
		if people[i].Gender != people[j].Gender {
			return people[i].Gender == "female"
		}
		// sort by height
		if people[i].Gender == "female" {
			// female height from low to high
			return people[i].Height < people[j].Height
		} else {
			// male height from high to low
			return people[i].Height > people[j].Height
		}
	})

	// fmt.Println("people:", people)
	if limit > 0 && limit < len(people) {
		people = people[:limit]
	}

	return people
}

func (ms *matchService) findMatches(newPerson *models.Person) []models.Match {
	var matches []models.Match
	var potentialMatches []*models.Person

	// Find potential matches based on gender and height rules
	for _, person := range ms.activePeople {
		if person.ID == newPerson.ID {
			continue
		}

		if ms.isCompatible(newPerson, person) {
			potentialMatches = append(potentialMatches, person)
		}
	}

	if newPerson.Gender == "male" {
		sort.Slice(potentialMatches, func(i, j int) bool {
			return potentialMatches[i].Height < potentialMatches[j].Height
		})
	} else {
		sort.Slice(potentialMatches, func(i, j int) bool {
			return potentialMatches[i].Height > potentialMatches[j].Height
		})
	}

	for _, potentialMatch := range potentialMatches {
		if newPerson.WantedDates <= 0 {
			break
		}
		matches = append(matches, models.Match{
			Person1: *newPerson,
			Person2: *potentialMatch,
		})
		newPerson.WantedDates--
		potentialMatch.WantedDates--

		if potentialMatch.WantedDates <= 0 {
			delete(ms.activePeople, potentialMatch.ID)
		}
	}
	if newPerson.WantedDates <= 0 {
		delete(ms.activePeople, newPerson.ID)
	}
	return matches
}

func (ms *matchService) isCompatible(person1, person2 *models.Person) bool {
	if person1.Gender == person2.Gender {
		return false
	}

	if person1.Gender == "male" && person2.Gender == "female" {
		return person1.Height > person2.Height
	}

	if person1.Gender == "female" && person2.Gender == "male" {
		return person2.Height > person1.Height
	}

	return false
}
