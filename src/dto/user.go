package dto


type ProfileUsernameImageDTO struct {
	Username string `bson:"username" json:"username"`
	ProfilePhoto string `bson:"profile_photo" json:"profile_photo"`
}

