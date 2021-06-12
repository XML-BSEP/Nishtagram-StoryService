package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	logger "github.com/jelena-vlajkov/logger/logger"
	"story-service/dto"
)

func GetUser(ctx context.Context, userId string, logger *logger.Logger) (dto.ProfileUsernameImageDTO, error) {
	client := resty.New()
	resp, _ := client.R().
		EnableTrace().
		Get("https://127.0.0.1:8082/getProfileUsernameImageById?userId=" + userId)

	var responseDTO dto.ProfileUsernameImageDTO
	err := json.Unmarshal(resp.Body(), &responseDTO)
	if err != nil {
		logger.Logger.Errorf("error while getting profile info for user %v, error: %v\n", userId, err)
	}

	return responseDTO, nil
}


func IsProfilePrivate(ctx context.Context, userId string, logger *logger.Logger) (bool, error) {
	client := resty.New()

	resp, err := client.R().
		SetBody(gin.H{"id" : userId}).
		SetContext(ctx).
		EnableTrace().
		Post("https://localhost:8082/isPrivate")

	if err != nil {
		logger.Logger.Errorf("error while getting profile info for user %v, error: %v\n", userId, err)
		return false, err
	}

	if resp.StatusCode() != 200 {
		logger.Logger.Errorf("error while getting profile info for user %v, error: %v\n", userId, err)
		return false, fmt.Errorf("Err")
	}

	var privacyCheckResponseDto dto.PrivacyCheckResponseDto
	if err := json.Unmarshal(resp.Body(), &privacyCheckResponseDto); err != nil {
		logger.Logger.Errorf("error while getting profile info for user %v, error: %v\n", userId, err)
		return false, err
	}

	return privacyCheckResponseDto.IsPrivate, err
}

