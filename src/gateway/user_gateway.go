package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	logger "github.com/jelena-vlajkov/logger/logger"
	"os"
	"story-service/dto"
)

func GetUser(ctx context.Context, userId string, logger *logger.Logger) (dto.ProfileUsernameImageDTO, error) {
	client := resty.New()
	domain := os.Getenv("USER_DOMAIN")
	if domain == "" {
		domain = "127.0.0.1"
	}
	fmt.Println(domain)
	if os.Getenv("DOCKER_ENV") == "" {
		resp, _ := client.R().
			EnableTrace().
			Get("https://" + domain + ":8082/getProfileUsernameImageById?userId=" + userId)

		var responseDTO dto.ProfileUsernameImageDTO
		err := json.Unmarshal(resp.Body(), &responseDTO)
		if err != nil {
			logger.Logger.Errorf("error while getting profile info for user %v, error: %v\n", userId, err)
		}

		return responseDTO, nil
	} else {
		resp, _ := client.R().
			EnableTrace().
			Get("http://" + domain + ":8082/getProfileUsernameImageById?userId=" + userId)

		var responseDTO dto.ProfileUsernameImageDTO
		err := json.Unmarshal(resp.Body(), &responseDTO)
		if err != nil {
			logger.Logger.Errorf("error while getting profile info for user %v, error: %v\n", userId, err)
		}

		return responseDTO, nil
	}

}
func IsProfilePrivate(ctx context.Context, userId string) (bool, error) {
	client := resty.New()
	domain := os.Getenv("USER_DOMAIN")
	if domain == "" {
		domain = "127.0.0.1"
	}

	if os.Getenv("DOCKER_ENV") == "" {
		resp, err := client.R().
			SetBody(gin.H{"id" : userId}).
			SetContext(ctx).
			EnableTrace().
			Post("https://" + domain + ":8082/getPrivacyAndTagging?userId=" + userId)

		if err != nil {
			return false, err
		}

		if resp.StatusCode() != 200 {
			return false, fmt.Errorf("Err")
		}

		var privacyCheckResponseDto dto.PrivacyTaggingDTO
		if err := json.Unmarshal(resp.Body(), &privacyCheckResponseDto); err != nil {
			return false, err
		}
		if privacyCheckResponseDto.PrivacyPermission == "Private" {
			return false, err
		}
		return true, err
	} else {
		resp, err := client.R().
			SetBody(gin.H{"id" : userId}).
			SetContext(ctx).
			EnableTrace().
			Post("http://" + domain + ":8082/getPrivacyAndTagging?userId=" + userId)

		if err != nil {
			return false, err
		}

		if resp.StatusCode() != 200 {
			return false, fmt.Errorf("Err")
		}

		var privacyCheckResponseDto dto.PrivacyTaggingDTO
		if err := json.Unmarshal(resp.Body(), &privacyCheckResponseDto); err != nil {
			return false, err
		}
		if privacyCheckResponseDto.PrivacyPermission == "Private" {
			return false, err
		}
		return true, err
	}

}
