package services

import (
	"matching_system/internal/api/dto"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchService_AddSinglePersonAndMatch(t *testing.T) {
	ms := NewMatchService()

	req := dto.AddPersonRequest{
		Name:        "Test Person",
		Height:      170,
		Gender:      "female",
		WantedDates: 3,
	}

	// test add success
	person, _ := ms.AddSinglePersonAndMatch(req)

	// verify the person is added
	assert.NotEmpty(t, person.ID, "the ID should be generated")
	assert.Equal(t, "Test Person", person.Name, "the name should match")
	assert.Equal(t, 170, person.Height, "the height should match")
	assert.Equal(t, "female", person.Gender, "the gender should match")
	assert.Equal(t, 3, person.WantedDates, "the WantedDates should match")

	// verify the person is in activePeople
	result := ms.QuerySinglePeople(0)
	assert.Equal(t, 1, len(result), "should have 1 person")
	assert.Equal(t, person.ID, result[0].ID, "the ID should match")
}

func TestMatchService_RemoveSinglePerson(t *testing.T) {
	ms := NewMatchService()

	// add test person
	req := dto.AddPersonRequest{
		Name:        "Test Person",
		Height:      170,
		Gender:      "female",
		WantedDates: 3,
	}

	person, _ := ms.AddSinglePersonAndMatch(req)

	// verify the person is in activePeople
	result := ms.QuerySinglePeople(0)
	assert.Equal(t, 1, len(result), "should have 1 person")

	// test remove success
	success := ms.RemoveSinglePerson(person.ID)
	assert.True(t, success, "should remove the person successfully")

	// verify the person is removed
	result = ms.QuerySinglePeople(0)
	assert.Equal(t, 0, len(result), "should have 0 person")
}

func TestMatchService_RemoveSinglePerson_NotFound(t *testing.T) {
	ms := NewMatchService()

	// test remove non-existent person
	success := ms.RemoveSinglePerson("non-existent-id")
	assert.False(t, success, "should return false")
}

func TestMatchService_QuerySinglePeople_Sorting(t *testing.T) {
	ms := NewMatchService()

	testPeople := []dto.AddPersonRequest{
		{Name: "Alice", Height: 180, Gender: "female", WantedDates: 3},
		{Name: "Bob", Height: 170, Gender: "male", WantedDates: 3},
		{Name: "Carol", Height: 185, Gender: "female", WantedDates: 3},
		{Name: "David", Height: 165, Gender: "male", WantedDates: 3},
		{Name: "Eve", Height: 175, Gender: "female", WantedDates: 2},
		{Name: "Frank", Height: 175, Gender: "male", WantedDates: 2},
		{Name: "Grace", Height: 180, Gender: "female", WantedDates: 1},
		{Name: "Henry", Height: 160, Gender: "male", WantedDates: 1},
	}

	for _, req := range testPeople {
		ms.AddSinglePersonAndMatch(req)
	}

	result := ms.QuerySinglePeople(0)

	// verify the result
	assert.Equal(t, 8, len(result), "should return 8 people")

	// verify the WantedDates sorting (main sorting)
	assert.Equal(t, 3, result[0].WantedDates, "the first should be WantedDates=3")
	assert.Equal(t, 3, result[1].WantedDates, "the second should be WantedDates=3")
	assert.Equal(t, 3, result[2].WantedDates, "the third should be WantedDates=3")
	assert.Equal(t, 3, result[3].WantedDates, "the fourth should be WantedDates=3")
	assert.Equal(t, 2, result[4].WantedDates, "the fifth should be WantedDates=2")
	assert.Equal(t, 2, result[5].WantedDates, "the sixth should be WantedDates=2")
	assert.Equal(t, 1, result[6].WantedDates, "the seventh should be WantedDates=1")
	assert.Equal(t, 1, result[7].WantedDates, "the eighth should be WantedDates=1")

	// verify the sorting when WantedDates=3
	// should sort by gender: female first, male second
	assert.Equal(t, "female", result[0].Gender, "the first should be female")
	assert.Equal(t, "female", result[1].Gender, "the second should be female")
	assert.Equal(t, "male", result[2].Gender, "the third should be male")
	assert.Equal(t, "male", result[3].Gender, "the fourth should be male")

	// verify the female internal sorting (from low to high)
	assert.Equal(t, 180, result[0].Height, "the first female height should be 160")
	assert.Equal(t, 185, result[1].Height, "the second female height should be 165")

	// verify the male internal sorting (from high to low)
	assert.Equal(t, 170, result[2].Height, "the first male height should be 180")
	assert.Equal(t, 165, result[3].Height, "the second male height should be 175")
}

func TestMatchService_QuerySinglePeople_SameWantedDates(t *testing.T) {
	ms := NewMatchService()

	// create test data with same WantedDates
	testPeople := []dto.AddPersonRequest{
		{Name: "Alice", Height: 160, Gender: "female", WantedDates: 3},
		{Name: "Bob", Height: 180, Gender: "male", WantedDates: 3},
		{Name: "Carol", Height: 165, Gender: "female", WantedDates: 3},
		{Name: "David", Height: 175, Gender: "male", WantedDates: 3},
	}

	for _, req := range testPeople {
		ms.AddSinglePersonAndMatch(req)
	}

	result := ms.QuerySinglePeople(0)

	expectedOrder := []struct {
		gender string
		height int
	}{
		{"female", 160},
		{"female", 165},
		{"male", 180},
		{"male", 175},
	}

	for i, expected := range expectedOrder {
		assert.Equal(t, expected.gender, result[i].Gender,
			"the %dth should be %s", i+1, expected.gender)
		assert.Equal(t, expected.height, result[i].Height,
			"the %dth height should be %d", i+1, expected.height)
	}
}

func TestMatchService_QuerySinglePeople_Limit(t *testing.T) {
	ms := NewMatchService()

	// create test data
	testPeople := []dto.AddPersonRequest{
		{Name: "Alice", Height: 180, Gender: "female", WantedDates: 3},
		{Name: "Bob", Height: 160, Gender: "male", WantedDates: 3},
		{Name: "Carol", Height: 185, Gender: "female", WantedDates: 2},
		{Name: "David", Height: 155, Gender: "male", WantedDates: 2},
	}

	// directly add to activePeople map
	for _, req := range testPeople {
		ms.AddSinglePersonAndMatch(req)
	}

	// test limit=2
	result := ms.QuerySinglePeople(2)
	assert.Equal(t, 2, len(result), "should return 2 people")
	assert.Equal(t, 3, result[0].WantedDates, "the first should be WantedDates=3")
	assert.Equal(t, 3, result[1].WantedDates, "the second should be WantedDates=3")

	// test limit=0 (return all)
	result = ms.QuerySinglePeople(0)
	assert.Equal(t, 4, len(result), "should return all 4 people")
}
