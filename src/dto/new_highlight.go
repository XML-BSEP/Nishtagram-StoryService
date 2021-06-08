package dto


type NewHighlight struct {
	Stories []string `json:"stories"`
	HighlightName string `json:"highlightName"`
	Id string `json:"id"`
	UserId string
}
