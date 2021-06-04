package repository

import (
	"context"
	"fmt"
	"github.com/gocql/gocql"
	"story-service/domain"
	"story-service/dto"
)

const (
	CreateHighlightTable = "CREATE TABLE IF NOT EXISTS story_keyspace.Highlights (name, profile_id, posts list<text>, main_story text PRIMARY KEY (profile_id, name));"
	GetAllHighlightsByUser       = "SELECT name, posts, main_story from story_keyspace.Highlights WHERE profile_id = ?);"
	GetStoriesInsideOneHighlight = "SELECT posts, main_story FROM story_keyspace.Highlights WHERE profile_id = ? AND name = ?;"
	UpdatePostsInHighlight       = "UPDATE story_keyspace.Highlighte SET posts = ? WHERE profile_id = ? AND name = ?;"
	GetAllStoryHighlights        = "SELECT name, main_story FROM story_keyspace.Stories WHERE profile_id = ?;"
)

type HighlightRepo interface {
	AddToHighlight(context context.Context, userId string, storyId string, highlightName string) error
	RemoveFromHighlight(context context.Context, userId string, storyId string, highlightName string) error
	GetHighlightByName(context context.Context, userId string, highlightName string) ([]string, string, error)
	GetAllHighlightsByUser(context context.Context, userId string) ([]dto.HighlightsPreviewDTO, error)
}

type highlightRepository struct {
	cassandraSession *gocql.Session
}

func (h highlightRepository) GetAllHighlightsByUser(context context.Context, userId string) ([]dto.HighlightsPreviewDTO, error) {
	var name, mainStory string
	var retVal []dto.HighlightsPreviewDTO
	iter := h.cassandraSession.Query(GetAllStoryHighlights, userId).Iter().Scanner()
	
	for iter.Next() {
		err := iter.Scan(&name, &mainStory)
		if err != nil {
			return nil, err
		}

		retVal = append(retVal, dto.HighlightsPreviewDTO{UserId: userId, HighlightName: name, MainStory: domain.Media{Path: mainStory}})
	}
	return retVal, nil
}


func (h highlightRepository) GetHighlightByName(context context.Context, userId string, highlightName string) ([]string, string, error) {
	var stories []string
	var mainStory string

	h.cassandraSession.Query(GetStoriesInsideOneHighlight, userId, highlightName).Iter().Scan(&stories, &mainStory)

	return stories, mainStory, nil
}

func (h highlightRepository) AddToHighlight(context context.Context, userId string, storyId string, highlightName string) error {
	iter := h.cassandraSession.Query(GetStoriesInsideOneHighlight, userId, highlightName).Iter()
	var stories []string
	iter.Scan(&stories)
	stories = append(stories, storyId)
	err := h.cassandraSession.Query(UpdatePostsInHighlight, stories, userId, highlightName).Exec()
	if err != nil {
		return nil
	}
	return nil
}

func (h highlightRepository) RemoveFromHighlight(context context.Context, userId string, storyId string, highlightName string) error {
	iter := h.cassandraSession.Query(GetStoriesInsideOneHighlight, userId, highlightName).Iter()
	var stories []string
	iter.Scan(&stories)
	isInSlice, position := stringInSlice(storyId, stories)
	if isInSlice {
		remove(stories, position)
	}
	err := h.cassandraSession.Query(UpdatePostsInHighlight, stories, userId, highlightName).Exec()
	if err != nil {
		return nil
	}
	return nil
}

func stringInSlice(a string, list []string) (bool, int) {
	for i, b := range list {
		if b == a {
			return true, i
		}
	}
	return false, -1
}

func remove(slice []string, index int) []string {
	copy(slice[index:], slice[index+1:])
	new := slice[:len(slice)-1]
	return new
}

func NewHighlightRepo(cassandraSession *gocql.Session) HighlightRepo {
	err := cassandraSession.Query(CreateHighlightTable).Exec()
	if err != nil {
		fmt.Println(err)
	}
	return &highlightRepository{cassandraSession: cassandraSession}
}