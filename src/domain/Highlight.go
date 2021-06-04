package domain

import "time"

type Highlight struct {
	Name string
	Timestamp time.Time
	Profile Profile
	Stories []Story
	MainStory Media
}

func NewHighlight(name string, userId string, stories []string, mainStory string) Highlight {
	var storiesToReturn []Story
	for _, story := range stories {

	}
	return Highlight{
		Name: name,
		Profile: Profile{Id: userId},
		Stories:
	}
}

