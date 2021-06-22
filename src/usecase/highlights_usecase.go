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
	"story-service/repository"
	"strings"
)

type HighlightUseCase interface {
	AddStoryToHighlight(ctx context.Context, dto dto.HighlightDTO) error
	RemoveStoryFrom(ctx context.Context, dto dto.HighlightDTO) error
	GetHighlights(ctx context.Context, dto dto.HighlightDTO) ([]dto.HighlightsPreviewDTO, error)
	GetHighlightByName(ctx context.Context, userId string, highlightName string) (dto.OneHighlightDTO, error)
	DecodeBase64Str(media string, userId string, ctx context.Context) (string, error)
	EncodeBase64String(media string, userId string, ctx context.Context) (string, error)
	UpdateHighlights(highlightDTO dto.NewHighlight, ctx context.Context) error
}

type highlightUseCase struct {
	highlightRepository repository.HighlightRepo
	storyRepository repository.StoryRepo
	logger *logger.Logger
}

func (h highlightUseCase) UpdateHighlights(dto dto.NewHighlight, ctx context.Context) error {
	h.logger.Logger.Infof("updating highlight %v for user %v\n", dto.HighlightName, dto.UserId)
	userId := dto.UserId
	if !h.highlightRepository.SeeIfHighlightExists(context.Background(), dto.UserId, dto.HighlightName) {
		h.highlightRepository.CreateHighlight(dto.UserId, dto.HighlightName, context.Background())
	}
	var postsToSave []string
	if len(dto.Stories) == 0 {
		h.highlightRepository.DeleteHighlight(userId, dto.HighlightName, context.Background())
		return nil
	}
	for _, highlight := range dto.Stories {
		img, _ := h.storyRepository.GetStoryById( context.Background(), dto.UserId, highlight)
		postsToSave = append(postsToSave, img.StoryId)
	}

	h.highlightRepository.UpdatePostsInHighlight(dto.UserId, dto.HighlightName, postsToSave, context.Background())

	return nil
}

func (h highlightUseCase) AddStoryToHighlight(ctx context.Context, dto dto.HighlightDTO) error {
	h.logger.Logger.Infof("adding story %v to highlight %v for user %v", dto.StoryId, dto.HighlightName, dto.UserId)
	if !h.highlightRepository.SeeIfHighlightExists(context.Background(), dto.UserId, dto.HighlightName) {
		h.logger.Logger.Infof("creating highlight %v for user %v", dto.HighlightName, dto.UserId)
		h.highlightRepository.CreateHighlight(dto.UserId, dto.HighlightName, context.Background())
	}
	var postsToSave []string

	for _, highlight := range dto.Stories {
		img, _ := h.EncodeBase64String(highlight.Story, dto.UserId, context.Background())

		postsToSave = append(postsToSave, img)
	}

	return nil

}

func (h highlightUseCase) EncodeBase64String(media string, userId string, ctx context.Context) (string, error) {
	h.logger.Logger.Infof("encoding image for user %v\n", userId)
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
		h.logger.Logger.Errorf("error while encoding image for user %v, error: %v\n", userId, err)
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
		h.logger.Logger.Errorf("error while encoding image for user %v, error: %v\n", userId, err)
	}
	uuid := uuid.NewString()
	f, err := os.Create(uuid + "." + format[0])

	if err != nil {
		h.logger.Logger.Errorf("error while encoding image for user %v, error: %v\n", userId, err)
	}

	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		h.logger.Logger.Errorf("error while encoding image for user %v, error: %v\n", userId, err)
	}
	if err := f.Sync(); err != nil {
		h.logger.Logger.Errorf("error while encoding image for user %v, error: %v\n", userId, err)
	}

	os.Chdir(workingDirectory)
	return userId + "/" + uuid + "." + format[0], nil
}

func (h highlightUseCase) RemoveStoryFrom(ctx context.Context, dto dto.HighlightDTO) error {
	h.logger.Logger.Infof("removing story %v from highlight %v for user %v\n", dto.StoryId, dto.HighlightName, dto.UserId)
	if len(dto.Stories) == 0 {
		h.logger.Logger.Infof("deleting highlight %v for user %v\n", dto.HighlightName, dto.UserId)
		h.highlightRepository.DeleteHighlight(dto.UserId, dto.HighlightName, context.Background())
	}
	var postsToSave []string

	for _, highlight := range dto.Stories {
		img, _ := h.EncodeBase64String(highlight.Story, dto.UserId, context.Background())

		postsToSave = append(postsToSave, img)
	}

	return nil
}

func (h highlightUseCase) GetHighlights(ctx context.Context, req dto.HighlightDTO) ([]dto.HighlightsPreviewDTO, error) {
	h.logger.Logger.Infof("getting all highlights for user %v\n", req.Id)
	if req.UserId != req.Id {
		if req.Id != req.UserId {
			userFollowing, _ := gateway.GetAllUserFollowing(context.Background(), req.UserId, h.logger)
			isOkay := false
			for  _, u := range userFollowing {
				if u.Id == req.Id {
					isOkay = true
					break
				}
			}
			if !isOkay {
				isPrivate, _ := gateway.IsProfilePrivate(ctx, req.Id)
				if !isPrivate {
					h.logger.Logger.Errorf("error while getting all stories by user %v, error: no followings\n", req.UserId)
					return nil, fmt.Errorf("oh no i hope i don't fall")
				}
			}
		}
	}

	highlights, _ := h.highlightRepository.GetAllHighlightsByUser(context.Background(), req.Id)
	var retVal []dto.HighlightsPreviewDTO
	if highlights != nil {
		for _, highlight := range highlights {
			encoded, _ := h.DecodeBase64Str(highlight.HighlightPhoto, req.Id, context.Background())
			highlight.HighlightPhoto = encoded
			retVal = append(retVal, highlight)
		}
	}
	return retVal, nil
}

func (h highlightUseCase) GetHighlightByName(ctx context.Context, userId string, highlightName string) (dto.OneHighlightDTO, error) {
	h.logger.Logger.Infof("getting highlight %v for user %v", highlightName, userId)
	stories, mainStory, _ := h.highlightRepository.GetHighlightByName(context.Background(), userId, highlightName)

	if len(stories) > 0 {
		var stories2 []string
		var storiesToReturn []dto.StoryDTO
		for _, s := range stories {

			story, err := h.storyRepository.GetStoryById(context.Background(), userId, s)
			img, _ := h.DecodeBase64Str(story.MediaPath.Path, story.UserId, context.Background())
			story.MediaPath.Path = img
			stories2 = append(stories2, story.StoryId)
			if err != nil {
				continue
			}
			isVideo := false
			if story.Type == "VIDEO" {
				isVideo = true
			}
			storiesToReturn = append(storiesToReturn, dto.StoryDTO{Story: img, StoryId: story.StoryId, CloseFriends: story.CloseFriends, IsVideo: isVideo})
		}

		return dto.OneHighlightDTO{
			MainPicture: domain.Media{Path: mainStory},
			UserId: userId,
			StoryId: stories2,
			Stories: storiesToReturn,
		}, nil
	}
	return dto.OneHighlightDTO{}, fmt.Errorf("no such highlight")


}

func (h highlightUseCase) DecodeBase64Str(media string, userId string, ctx context.Context) (string, error) {
	h.logger.Logger.Infof("decoding image %v for user %v\n", media, userId)
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

func NewHighlightUseCase(highlightRepo repository.HighlightRepo, storyRepository repository.StoryRepo, logger *logger.Logger) HighlightUseCase {
	return &highlightUseCase{highlightRepository: highlightRepo, storyRepository: storyRepository, logger: logger}
}