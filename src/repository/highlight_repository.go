package repository

import (
	"context"
	"fmt"
	"github.com/gocql/gocql"
)

const (
	CreateHighlightTable = "CREATE TABLE IF NOT EXISTS story_keyspace.Highlights (name, profile_id, posts list<text> PRIMARY KEY (profile_id, name));"
	InsertIntoHighlightTable = "INSERT INTO story_keyspace.Highlights (name, profile_id, posts) VALUES (?, ?, ?);"
	GetAllHighlightsByUser = "SELECT name, posts from story_keyspace.Highlights WHERE profile_id = ?);"
	GetStoriesInsideOneHighlight = "SELECT posts FROM story_keyspace.Highlights WHERE profile_id = ? AND name = ?;"
	UpdatePostsInHighlight = "UPDATE story_keyspace.Highlighte SET posts = ? WHERE profile_id = ? AND name = ?;"
)

type HighlightRepo interface {
	AddToHighlight(context context.Context, userId string, storyId string, highlightName string) error
	RemoveFromHighlight(context context.Context, userId string, storyId string, highlightName string) error
}

type highlightRepository struct {
	cassandraSession *gocql.Session
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