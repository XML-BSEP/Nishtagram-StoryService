package usecase

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
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
}

type storyUseCase struct {
	storyRepository repository.StoryRepo
	redisUseCase RedisUseCase
}

func (s storyUseCase) EncodeBase64(media string, userId string, ctx context.Context) (string, error) {

	workingDirectory, _ := os.Getwd()
	path1 := "./assets/images/"
	err := os.Chdir(path1)
	if err != nil {
		fmt.Println(err)
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
		panic(err)
	}
	uuid := uuid.NewString()
	f, err := os.Create(uuid + "." + format[0])

	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		panic(err)
	}
	if err := f.Sync(); err != nil {
		panic(err)
	}

	os.Chdir(workingDirectory)
	return userId + "/" + uuid + "." + format[0], nil
}

func (s storyUseCase) DecodeBase64(media string, userId string, ctx context.Context) (string, error) {
	workingDirectory, _ := os.Getwd()

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


	fmt.Println("ENCODED: " + encoded)
	os.Chdir(workingDirectory)

	return "data:image/jpg;base64," + encoded, nil
}

func (s storyUseCase) AddStory(ctx context.Context, dto dto.StoryDTO) error {
	var keyToStore, valueToStore string
	expiresAtTime := time.Now().Add(time.Hour*24)
	expiresAt := time.Unix(expiresAtTime.Unix(), 0)
	now := time.Now()

	decoded, err := s.EncodeBase64(dto.MediaPath.Path, dto.UserId, context.Background())
	if err != nil {
		return err
	}

	dto.MediaPath.Path = decoded

	keyToStore = dto.UserId + "/" + dto.StoryId
	valueToStore = dto.StoryId

	s.redisUseCase.AddKeyValueSet(context.Background(), keyToStore, valueToStore, now.Sub(expiresAt))
	return s.storyRepository.AddStory(context.Background(), mapper.MapDTOToStory(dto))
}

func (s storyUseCase) RemoveStory(ctx context.Context, dto dto.RemoveStoryDTO) error {
	if !s.storyRepository.SeeIfExists(context.Background(), dto.UserId, dto.StoryId) {
		return fmt.Errorf("no cuch story")
	}

	return s.storyRepository.RemoveStory(context.Background(), dto.UserId, dto.StoryId)
}

func (s storyUseCase) GetAllStoriesForOneUser(ctx context.Context, userId string) ([]dto.StoryDTO, error) {
	userFollowing, err := gateway.GetAllUserFollowing(context.Background(), userId)
	if err != nil {
		fmt.Println(err)
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
			story.User = domain.Profile{Id: u.Id, ProfilePhoto: encoded, Username: "aaa"}


			retVal = append(retVal, story)
		}

	}

	return retVal, nil

}

func NewStoryUseCase(storyRepository repository.StoryRepo, useCase RedisUseCase) StoryUseCase {
	return &storyUseCase{
		storyRepository: storyRepository,
		redisUseCase: useCase,
	}
}
