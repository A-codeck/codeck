package group

type Group struct {
	ID          string  `json:"id"`
	CreatorID   string  `json:"creator_id"`
	Name        string  `json:"name"`
	StartDate   string  `json:"start_date"`
	EndDate     string  `json:"end_date"`
	GroupImage  *string `json:"group_image,omitempty"`
	Description *string `json:"description,omitempty"`
}
