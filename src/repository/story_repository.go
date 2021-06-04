package repository

import (
	"context"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"story-service/domain"
	"story-service/dto"
	"time"
)

const (
	CreateStoryTable = "CREATE TABLE IF NOT EXISTS story_keyspace.Stories (id text, profile_id text, image text, timestamp timestamp, mentions list<text>, close_friends bool, type text, location_name text, longitude double, latitude double, deleted boolean, PRIMARY KEY (profile_id, id));"
	InsertIntoStoryTable = "INSERT INTO story_keyspace.Stories (id, profile_id, image, timestamp, mentions, close_friends, type, location_name, longitude, latitude, deleted) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"
	GetAllStoryByUser = "SELECT id, profile_id, image, timestamp, close_friends, type, mentions, location_name, latitude, longitude FROM story_keyspace.Stories WHERE profile_id = ? AND id = ? AND deleted = false;"
	DeleteStory = "UPDATE story_keyspace.Stories SET deleted = ? WHERE profile_id = ? AND id = ?;"
	GetAllStoriesByUser = "SELECT id, profile_id, image, timestamp, close_friends, type, mentions, location_name, latitude, longitude FROM story_keyspace.Stories WHERE profile_id = ? AND id = ? AND deleted = false;"
)

type StoryRepo interface {
	AddStory(ctx context.Context, story domain.Story) error
	RemoveStory(ctx context.Context, userId string, storyId string) error
	GetStoriesByUserFollowing(ctx context.Context, userId string, postId string) ([]dto.StoryDTO, error)
}

type storyRepository struct {
	cassandraClient *gocql.Session
}

func (s storyRepository) GetStoriesByUserFollowing(ctx context.Context, userId string, postId string) ([]dto.StoryDTO, error) {
	var stories []dto.StoryDTO
	var location domain.Location
	var id, profileId, image, storyType, locationName string
	var mentions []string
	var latitude, longitude float64
	var timestamp time.Time
	var closeFriends bool
	iter := s.cassandraClient.Query(GetAllStoriesByUser, userId, postId).Iter().Scanner()
	//id, profile_id, image, timestamp, close_friends, type, mentions, location_name, latitude, longitude
	for iter.Next() {
		iter.Scan(&id, &profileId, &image, &timestamp, &closeFriends, &storyType, &mentions, &locationName, &latitude, &longitude)
		location = domain.NewLocation(locationName, latitude, longitude)
		stories = append(stories, dto.NewStoryDTO(id, profileId, mentions, image, storyType, location, timestamp, closeFriends))
	}
	return stories, nil

}

func (s storyRepository) AddStory(ctx context.Context, story domain.Story) error {
	var mentions []string
	for _, st := range story.Mentions {
		mentions = append(mentions, st.Id)
	}
	err := s.cassandraClient.Query(InsertIntoStoryTable, uuid.NewString(), story.Profile.Id, story.Media.Path, story.Timestamp, mentions, story.CloseFriends,
		story.StoryType.Type, story.Location.Location, story.Location.Longitude, story.Location.Latitude, false).Exec()
	if err != nil {

		return fmt.Errorf("server error")
	}
	return nil
}

func (s storyRepository) RemoveStory(ctx context.Context, userId string, storyId string) error {
	err := s.cassandraClient.Query(DeleteStory, true, userId, storyId).Exec()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func NewStoryRepo(cassandraClient *gocql.Session) StoryRepo {
	err := cassandraClient.Query(CreateStoryTable).Exec()
	if err != nil {
		fmt.Println(err)
	}
	return &storyRepository{cassandraClient: cassandraClient}
}
