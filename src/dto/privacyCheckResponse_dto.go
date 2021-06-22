package dto

type PrivacyTaggingDTO struct {
	PrivacyPermission string `bson:"privacy_permission" json:"privacy_permission"`
	AllowTagging bool `bson:"allow_tagging" json:"allow_tagging"`
	ProfileId string `bson:"profile_id" json:"profile_id"`
}
