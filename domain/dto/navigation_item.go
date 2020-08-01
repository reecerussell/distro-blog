package dto

type NavigationItem struct {
	ID string `json:"id"`
	Text string `json:"text"`
	Target string `json:"target"`
	URL *string `json: "url"`
	PageID *string `json:"pageID"`
	IsHidden bool `json:"isHidden"`
	IsBrand bool `json:"isBrand"`
}
