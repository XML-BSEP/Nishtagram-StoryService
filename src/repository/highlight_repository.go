package repository

import (
	"context"
	"fmt"
	"github.com/gocql/gocql"
	logger "github.com/jelena-vlajkov/logger/logger"
	"story-service/dto"
)

const (
	CreateHighlightTable = "CREATE TABLE IF NOT EXISTS story_keyspace.Highlights (name text, profile_id text, posts list<text>, main_story text, PRIMARY KEY (profile_id, name));"
	GetAllHighlightsByUser       = "SELECT name, posts, main_story from story_keyspace.Highlights WHERE profile_id = ?);"
	GetStoriesInsideOneHighlight = "SELECT posts, main_story FROM story_keyspace.Highlights WHERE profile_id = ? AND name = ?;"
	UpdatePostsInHighlight       = "UPDATE story_keyspace.Highlights SET posts = ?, main_story = ? WHERE profile_id = ? AND name = ?;"
	GetAllStoryHighlights        = "SELECT name, main_story FROM story_keyspace.Highlights WHERE profile_id = ?;"
	InsertIntoHighlights = "INSERT INTO story_keyspace.Highlights (name, profile_id, posts, main_story) VALUES (?, ?, ?, ?);"
	SeeIfHighlightExists = "SELECT posts, main_story FROM story_keyspace.Highlights WHERE profile_id = ? AND name = ?"
	DeleteHighlight = "DELETE FROM story_keyspace.Highlights WHERE profile_id = ? AND name = ?;"
	UpdatePostsAndMainImageInHighlight       = "UPDATE story_keyspace.Highlights SET posts = ?, main_story = ? WHERE profile_id = ? AND name = ?;"
)

type HighlightRepo interface {
	AddToHighlight(context context.Context, userId string, storyId string, highlightName string) error
	RemoveFromHighlight(context context.Context, userId string, storyId string, highlightName string) error
	GetHighlightByName(context context.Context, userId string, highlightName string) ([]string, string, error)
	GetAllHighlightsByUser(context context.Context, userId string) ([]dto.HighlightsPreviewDTO, error)
	SeeIfHighlightExists(ctx context.Context, userId string, highlightName string) bool
	CreateHighlight(userId string, highlightName string, ctx context.Context) error
	UpdatePostsInHighlight(userId string, highlightName string, posts []string, ctx context.Context) error
	DeleteHighlight(userId string, highlightName string, ctx context.Context) error

}

type highlightRepository struct {
	cassandraSession *gocql.Session
	logger *logger.Logger
}

func (h highlightRepository) DeleteHighlight(userId string, highlightName string, ctx context.Context) error {
	err := h.cassandraSession.Query(DeleteHighlight, userId, highlightName).Exec()
	if err != nil {
		h.logger.Logger.Errorf("error while deleting highlight %v for user %userId, error: %v\n", highlightName, userId, err)
	}
	return err
}

func (h highlightRepository) UpdatePostsInHighlight(userId string, highlightName string, posts []string, ctx context.Context) error {
	if len(posts) == 0 {
		h.cassandraSession.Query(DeleteHighlight, userId, highlightName).Exec()
		return nil
	}

	err := h.cassandraSession.Query(UpdatePostsInHighlight, posts, posts[0], userId, highlightName).Exec()
	if err != nil {
		h.logger.Logger.Errorf("error while updationg post in highlight %v for user %v, error: %v\n", highlightName, userId, err)
	}
	return nil
}

func (h highlightRepository) CreateHighlight(userId string, highlightName string, ctx context.Context) error {
	var posts []string
	var mainStory string
	err := h.cassandraSession.Query(InsertIntoHighlights, highlightName, userId, posts, mainStory).Exec()
	if err != nil {
		h.logger.Logger.Errorf("error while creating highlight %v for user %v, error: %v\n", highlightName, userId, err)
		return err
	}
	return nil
}

func (h highlightRepository) SeeIfHighlightExists(ctx context.Context, userId string, highlightName string) bool {
	var posts []string
	var mainStory string
	h.cassandraSession.Query(SeeIfHighlightExists, userId, highlightName).Iter().Scan(&posts, mainStory)
	if len(posts) > 0 {
		return true
	}
	return false

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
		var mainMedia string
		h.cassandraSession.Query(GetMediaFromId, userId, mainStory).Iter().Scan(&mainMedia)

		retVal = append(retVal, dto.HighlightsPreviewDTO{UserId: userId, HighlightName: name, HighlightPhoto: mainMedia})
	}
	return retVal, nil
}


func (h highlightRepository) GetHighlightByName(context context.Context, userId string, highlightName string) ([]string, string, error) {
	var stories []string
	var mainStory string

	h.cassandraSession.Query(GetStoriesInsideOneHighlight, userId, highlightName).Iter().Scan(&stories, &mainStory)
	var mainMedia string
	h.cassandraSession.Query(GetMediaFromId, userId, mainStory).Iter().Scan(&mainMedia)
	return stories, mainMedia, nil
}

func (h highlightRepository) AddToHighlight(context context.Context, userId string, storyId string, highlightName string) error {
	iter := h.cassandraSession.Query(GetStoriesInsideOneHighlight, userId, highlightName).Iter()
	var stories []string
	iter.Scan(&stories)
	stories = append(stories, storyId)
	err := h.cassandraSession.Query(UpdatePostsInHighlight, stories, userId, highlightName).Exec()
	if err != nil {
		h.logger.Logger.Errorf("error while adding to highlight %v for user %v, error: %v\n", highlightName, userId, err)
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
		h.logger.Logger.Errorf("error while removing post from highlight %v for user %v, error: %v\n", highlightName, userId, err)
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

func NewHighlightRepo(cassandraSession *gocql.Session, logger *logger.Logger) HighlightRepo {
	err := cassandraSession.Query(CreateHighlightTable).Exec()
	if err != nil {
		fmt.Println(err)
	}
	return &highlightRepository{cassandraSession: cassandraSession, logger: logger}
}