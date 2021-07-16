package mapper

import (
	"story-service/domain"
	"story-service/dto"
)

func MapDTOToStory(dto dto.StoryDTO) domain.Story {
	var mentions []domain.Profile
	for _, s := range dto.Mentions {
		mentions = append(mentions, domain.Profile{Id: s})
	}
	return domain.Story{
		Id: dto.StoryId,
		Profile: domain.Profile{Id: dto.UserId},
		Timestamp: dto.Timestamp,
		Location: dto.Location,
		Mentions: mentions,
		CloseFriends: dto.CloseFriends,
		Banned: false,
		StoryType: domain.StoryType{Type: dto.Type},
		Media: dto.MediaPath,
		IsCampaign: dto.IsCampaign,
		CampaignId: dto.CampaignId,
		Link: dto.Link,
	}
}
