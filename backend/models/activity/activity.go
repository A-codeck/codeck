package activity

type Activity struct {
	ID            string  `json:"id"`
	CreatorID     string  `json:"creator_id"`
	GroupID       *string `json:"group_id,omitempty"`
	Title         string  `json:"title"`
	Date          string  `json:"date"`
	ActivityImage *string `json:"activity_image,omitempty"`
	Description   *string `json:"description,omitempty"`
}
