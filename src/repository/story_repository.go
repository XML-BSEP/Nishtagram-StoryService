package repository

import (
	"context"
	"fmt"
	"github.com/gocql/gocql"
	logger "github.com/jelena-vlajkov/logger/logger"
	"story-service/domain"
	"story-service/dto"
	"time"
)

const (
	CreateStoryTable     = "CREATE TABLE IF NOT EXISTS story_keyspace.Stories (id text, profile_id text, image text, timestamp timestamp, mentions list<text>, close_friends boolean, type text, location_name text, longitude double, latitude double, deleted boolean,, is_campaign boolean, campaign_id text, link text,  PRIMARY KEY (profile_id, id));"
	InsertIntoStoryTable = "INSERT INTO story_keyspace.Stories (id, profile_id, image, timestamp, mentions, close_friends, type, location_name, longitude, latitude, deleted, is_campaign, campaign_id, link) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"
	GetAllStoryByUser    = "SELECT id, profile_id, image, timestamp, close_friends, type, mentions, location_name, latitude, longitude FROM story_keyspace.Stories WHERE profile_id = ? AND id = ? AND deleted = false;"
	DeleteStory          = "UPDATE story_keyspace.Stories SET deleted = ? WHERE profile_id = ? AND id = ?;"
	GetStoryById         = "SELECT id, profile_id, image, timestamp, close_friends, type, mentions, location_name, latitude, longitude, deleted, is_campaign, campaign_id, link FROM story_keyspace.Stories WHERE profile_id = ? AND id = ?;"
	SeeIfExists          = "SELECT count(*) FROM story_keyspace.Stories WHERE profile_id = ? AND id = ?;"
	GetMediaFromId 		 = "SELECT image FROM story_keyspace.Stories WHERE profile_id = ? AND id = ?;"
	GetStoriesByUserId         = "SELECT id, profile_id, image, timestamp, close_friends, type, mentions, location_name, latitude, longitude, deleted, is_campaign, campaign_id, link FROM story_keyspace.Stories WHERE profile_id = ?;"

)

type StoryRepo interface {
	AddStory(ctx context.Context, story domain.Story) error
	RemoveStory(ctx context.Context, userId string, storyId string) error
	GetStoryById(ctx context.Context, userId string, postId string) (dto.StoryDTO, error)
	SeeIfExists(ctx context.Context, userId string, storyId string) bool
	GetAllStoriesById(ctx context.Context, userId string) ([]dto.StoryDTO, error)
	GetStoryByAdmin(ctx context.Context, userId string, postId string) (dto.StoryDTO, error)
}

type storyRepository struct {
	cassandraClient *gocql.Session
	logger *logger.Logger
}

func (s storyRepository) GetStoryByAdmin(ctx context.Context, userId string, postId string) (dto.StoryDTO, error) {
	var location domain.Location
	var id, profileId, image, storyType, locationName,campaignId, link string
	var mentions []string
	var latitude, longitude float64
	var timestamp time.Time
	var closeFriends, isCampaign bool
	var deleted bool
	iter := s.cassandraClient.Query(GetStoryById, userId, postId).Iter().Scanner()
	//id, profile_id, image, timestamp, close_friends, type, mentions, location_name, latitude, longitude
	for iter.Next() {
		iter.Scan(&id, &profileId, &image, &timestamp, &closeFriends, &storyType, &mentions, &locationName, &latitude, &longitude, &deleted, &isCampaign, &campaignId, &link)
		location = domain.NewLocation(locationName, latitude, longitude)
		return dto.NewStoryDTO(id, profileId, mentions, domain.Media{Path: image, Timestamp: timestamp}, storyType, location, timestamp, closeFriends, isCampaign, campaignId, link), nil

	}
	return dto.StoryDTO{}, fmt.Errorf("no such story")
}

func (s storyRepository) GetAllStoriesById(ctx context.Context, userId string) ([]dto.StoryDTO, error) {
	var location domain.Location
	var id, profileId, image, storyType, locationName, campaignId, link string
	var mentions []string
	var latitude, longitude float64
	var timestamp time.Time
	var closeFriends, isCampaign bool
	var retVal []dto.StoryDTO
	var deleted bool
	iter := s.cassandraClient.Query(GetStoriesByUserId, userId).Iter().Scanner()
	//id, profile_id, image, timestamp, close_friends, type, mentions, location_name, latitude, longitude
	for iter.Next() {
		iter.Scan(&id, &profileId, &image, &timestamp, &closeFriends, &storyType, &mentions, &locationName, &latitude, &longitude, &deleted, &isCampaign, &campaignId, &link)
		if !deleted {
			location = domain.NewLocation(locationName, latitude, longitude)
			retVal = append(retVal, dto.NewStoryDTO(id, profileId, mentions, domain.Media{Path: image, Timestamp: timestamp}, storyType, location, timestamp, closeFriends, isCampaign, campaignId, link))
		}

	}
	return retVal, nil
}

func (s storyRepository) SeeIfExists(ctx context.Context, userId string, storyId string) bool {
	count := 0
	s.cassandraClient.Query(SeeIfExists, userId, storyId).Iter().Scan(&count)
	return count > 0
}

func (s storyRepository) GetStoryById(ctx context.Context, userId string, postId string) (dto.StoryDTO, error) {
	var location domain.Location
	var id, profileId, image, storyType, locationName, campaignId, link string
	var mentions []string
	var latitude, longitude float64
	var timestamp time.Time
	var closeFriends, isCampaign bool
	var deleted bool
	iter := s.cassandraClient.Query(GetStoryById, userId, postId).Iter().Scanner()
	//id, profile_id, image, timestamp, close_friends, type, mentions, location_name, latitude, longitude
	for iter.Next() {
		iter.Scan(&id, &profileId, &image, &timestamp, &closeFriends, &storyType, &mentions, &locationName, &latitude, &longitude, &deleted, &isCampaign, &campaignId, &link)
		if !deleted {
			location = domain.NewLocation(locationName, latitude, longitude)
			return dto.NewStoryDTO(id, profileId, mentions, domain.Media{Path: image, Timestamp: timestamp}, storyType, location, timestamp, closeFriends, isCampaign, campaignId, link), nil
		}

	}
	return dto.StoryDTO{}, fmt.Errorf("no such story")

}

func (s storyRepository) AddStory(ctx context.Context, story domain.Story) error {
	var mentions []string
	for _, st := range story.Mentions {
		mentions = append(mentions, st.Id)
	}
	err := s.cassandraClient.Query(InsertIntoStoryTable, story.Id, story.Profile.Id, story.Media.Path, story.Timestamp, mentions, story.CloseFriends,
		story.StoryType.Type, story.Location.Location, story.Location.Longitude, story.Location.Latitude, false, story.IsCampaign, story.CampaignId, story.Link).Exec()
	if err != nil {
		s.logger.Logger.Errorf("error while adding story for user %v, error: %v\n", story.Profile.Id, err)
		return fmt.Errorf("server error")
	}
	return nil
}

func (s storyRepository) RemoveStory(ctx context.Context, userId string, storyId string) error {
	err := s.cassandraClient.Query(DeleteStory, true, userId, storyId).Exec()
	if err != nil {
		s.logger.Logger.Errorf("error while removing")
		fmt.Println(err)
		return err
	}
	return nil
}

func NewStoryRepo(cassandraClient *gocql.Session, logger *logger.Logger) StoryRepo {
	err := cassandraClient.Query(CreateStoryTable).Exec()
	if err != nil {

		fmt.Println(err)
	}
	return &storyRepository{cassandraClient: cassandraClient, logger: logger}
}
