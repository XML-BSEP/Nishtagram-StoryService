package usecase

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	logger "github.com/jelena-vlajkov/logger/logger"
	"io/ioutil"
	"os"
	"story-service/domain"
	"story-service/dto"
	"story-service/gateway"
	"story-service/mapper"
	"story-service/repository"
	"strings"
	"time"
)

type StoryUseCase interface {
	AddStory(ctx context.Context, dto dto.StoryDTO) error
	RemoveStory(ctx context.Context, dto dto.RemoveStoryDTO) error
	GetAllStoriesForOneUser(ctx context.Context, userId string) ([]dto.StoryDTO, error)
	EncodeBase64(media string, userId string, ctx context.Context) (string, error)
	DecodeBase64(media string, userId string, ctx context.Context) (string, error)
	GetAllStoriesByUser(userId string, userRequested string, ctx context.Context) ([]dto.StoryDTO, error)
	GetActiveUsersStories(userId string, ctx context.Context) ([]dto.StoryDTO, error)

}

type storyUseCase struct {
	storyRepository repository.StoryRepo
	redisUseCase RedisUseCase
	logger *logger.Logger
}

func (s storyUseCase) GetAllStoriesByUser(userId string, userRequested string, ctx context.Context) ([]dto.StoryDTO, error) {
	s.logger.Logger.Infof("getting all stories by user %v\n", userId)
	if userId != userRequested {
		userFollowing, _ := gateway.GetAllUserFollowing(context.Background(), userRequested, s.logger)
		isOkay := false
		for  _, u := range userFollowing {
			if u.Id == userId {
				isOkay = true
				break
			}
		}
		if !isOkay {
			s.logger.Logger.Errorf("error while getting all stories by user %v, error: no followings\n", userId)
			return nil, fmt.Errorf("oh no i hope i don't fall")
		}
	}
	stories, _ := s.storyRepository.GetAllStoriesById(context.Background(), userId)
	var retVal []dto.StoryDTO

	for _, st := range stories {

		encoded, err := s.DecodeBase64(st.MediaPath.Path, st.UserId, context.Background())

		if err != nil {
			continue
		}
		st.MediaPath.Path = encoded
		if st.Type == "VIDEO" {
			st.StoryContent = dto.StoryContent{IsVideo: true, Content: st.MediaPath.Path}
		} else {
			st.StoryContent = dto.StoryContent{IsVideo: false, Content: st.MediaPath.Path}
		}

		profile, _ := gateway.GetUser(context.Background(), st.UserId, s.logger)
		st.User = domain.Profile{Id: st.UserId, ProfilePhoto: profile.ProfilePhoto, Username: profile.Username}
		st.Story = encoded
		retVal = append(retVal, st)
	}
	return retVal, nil
}

func (s storyUseCase) GetActiveUsersStories(userId string, ctx context.Context) ([]dto.StoryDTO, error) {
	s.logger.Logger.Infof("getting all active stories for user %v\n", userId)
	var retVal []dto.StoryDTO


	searchPattern := userId + "*"
	keys, _ := s.redisUseCase.ScanKeyByPattern(context.Background(), searchPattern)


	for _, k := range keys {
		storyId, err := s.redisUseCase.GetValueByKey(context.Background(), k)
		if err != nil {
			continue
		}
		story, err := s.storyRepository.GetStoryById(context.Background(), userId, storyId)
		if err != nil {
			continue
		}
		encoded, err := s.DecodeBase64(story.MediaPath.Path, story.UserId, context.Background())

		profile, _ := gateway.GetUser(context.Background(), story.UserId, s.logger)
		story.User = domain.Profile{Id: story.UserId, ProfilePhoto: profile.ProfilePhoto, Username: profile.Username}

		if err != nil {
			continue
		}
		story.MediaPath.Path = encoded
		if story.Type == "VIDEO" {
			story.StoryContent = dto.StoryContent{IsVideo: true, Content: story.MediaPath.Path}
		} else {
			story.StoryContent = dto.StoryContent{IsVideo: false, Content: story.MediaPath.Path}
		}

		story.User = domain.Profile{Id: userId, ProfilePhoto: encoded, Username: "aaa"}


		retVal = append(retVal, story)
	}

	return retVal, nil
}

func (s storyUseCase) EncodeBase64(media string, userId string, ctx context.Context) (string, error) {
	s.logger.Logger.Infof("encoding image for user %v\n", userId)
	workingDirectory, _ := os.Getwd()
	if !strings.HasSuffix(workingDirectory, "src") {
		firstPart := strings.Split(workingDirectory, "src")
		value := firstPart[0] + "src"
		workingDirectory = value
		os.Chdir(workingDirectory)
	}
	path1 := "./assets/images/"
	err := os.Chdir(path1)
	if err != nil {
		s.logger.Logger.Errorf("error while encoding image for user %v, error: %v\n", userId, err)
	}
	err = os.Mkdir(userId, 0755)
	fmt.Println(err)

	err = os.Chdir(userId)
	fmt.Println(err)


	st := strings.Split(media, ",")
	a := strings.Split(st[0], "/")
	format := strings.Split(a[1], ";")

	dec, err := base64.StdEncoding.DecodeString(st[1])

	if err != nil {
		s.logger.Logger.Errorf("error while encoding image for user %v, error: %v\n", userId, err)
	}
	uuid := uuid.NewString()
	f, err := os.Create(uuid + "." + format[0])

	if err != nil {
		s.logger.Logger.Errorf("error while encoding image for user %v, error: %v\n", userId, err)
	}

	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		s.logger.Logger.Errorf("error while encoding image for user %v, error: %v\n", userId, err)
	}

	err = os.Chdir(workingDirectory)

	workingDirectory, _ = os.Getwd()
	return userId + "/" + uuid + "." + format[0], nil
}

func (s storyUseCase) DecodeBase64(media string, userId string, ctx context.Context) (string, error) {
	s.logger.Logger.Infof("decoding image %v for user %v\n", media, userId)
	workingDirectory, _ := os.Getwd()
	if !strings.HasSuffix(workingDirectory, "src") {
		firstPart := strings.Split(workingDirectory, "src")
		value := firstPart[0] + "src"
		workingDirectory = value
		os.Chdir(workingDirectory)
	}

	path1 := "./assets/images/"
	err := os.Chdir(path1)
	fmt.Println(err)
	spliced := strings.Split(media, "/")
	var f *os.File
	if len(spliced) > 1 {
		err = os.Chdir(userId)
		f, _ = os.Open(spliced[1])
	} else {
		f, _ = os.Open(spliced[0])
	}




	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)


	encoded := base64.StdEncoding.EncodeToString(content)
	os.Chdir(workingDirectory)

	return "data:image/jpg;base64," + encoded, nil
}

func (s storyUseCase) AddStory(ctx context.Context, dto dto.StoryDTO) error {
	s.logger.Logger.Infof("adding story for user %v\n", dto.UserId)
	var keyToStore, valueToStore string
	expiresAtTime := time.Now().Add(time.Hour*24)
	expiresAt := time.Unix(expiresAtTime.Unix(), 0)
	now := time.Now()

	decoded, err := s.EncodeBase64(dto.Story, dto.UserId, context.Background())
	if err != nil {
		return err
	}

	dto.MediaPath.Path = decoded


	keyToStore = dto.UserId + "/" + dto.StoryId
	valueToStore = dto.StoryId
	dto.Timestamp = now
	if dto.IsVideo {
		dto.Type = "VIDEO"
	} else {
		dto.Type = "IMAGE"
	}


	s.redisUseCase.AddKeyValueSet(context.Background(), keyToStore, valueToStore, now.Sub(expiresAt))
	return s.storyRepository.AddStory(context.Background(), mapper.MapDTOToStory(dto))
}

func (s storyUseCase) RemoveStory(ctx context.Context, dto dto.RemoveStoryDTO) error {
	s.logger.Logger.Infof("removing story %v for user %v\n", dto.StoryId, dto.UserId)
	if !s.storyRepository.SeeIfExists(context.Background(), dto.UserId, dto.StoryId) {
		s.logger.Logger.Errorf("no such story %v for user %v\n", dto.StoryId, dto.UserId)
		return fmt.Errorf("no cuch story")
	}

	return s.storyRepository.RemoveStory(context.Background(), dto.UserId, dto.StoryId)
}

func (s storyUseCase) GetAllStoriesForOneUser(ctx context.Context, userId string) ([]dto.StoryDTO, error) {
	userFollowing, err := gateway.GetAllUserFollowing(context.Background(), userId, s.logger)
	if err != nil {
		return nil, err
	}
	var retVal []dto.StoryDTO

	for _, u := range userFollowing {
		searchPattern := u.Id + "*"
		keys, err := s.redisUseCase.ScanKeyByPattern(context.Background(), searchPattern)
		if err != nil {
			continue
		}

		for _, k := range keys {
			storyId, err := s.redisUseCase.GetValueByKey(context.Background(), k)
			if err != nil {
				continue
			}
			story, err := s.storyRepository.GetStoryById(context.Background(), u.Id, storyId)
			if err != nil {
				continue
			}
			profile, _ := gateway.GetUser(context.Background(), u.Id, s.logger)
			story.User = domain.Profile{Id: u.Id, ProfilePhoto: profile.ProfilePhoto, Username: profile.Username}
			encoded, err := s.DecodeBase64(story.MediaPath.Path, story.UserId, context.Background())
			if err != nil {
				continue
			}
			story.MediaPath.Path = encoded
			if story.Type == "VIDEO" {
				story.StoryContent = dto.StoryContent{IsVideo: true, Content: story.MediaPath.Path}
			} else {
				story.StoryContent = dto.StoryContent{IsVideo: false, Content: story.MediaPath.Path}
			}


			retVal = append(retVal, story)
		}

	}

	return retVal, nil

}

func NewStoryUseCase(storyRepository repository.StoryRepo, useCase RedisUseCase, logger *logger.Logger) StoryUseCase {
	return &storyUseCase{
		storyRepository: storyRepository,
		redisUseCase: useCase,
		logger: logger,
	}
}
