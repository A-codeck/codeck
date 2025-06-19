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

type GroupMember struct {
	UserID   string  `json:"user_id"`
	GroupID  string  `json:"group_id"`
	Nickname *string `json:"nickname,omitempty"`
}

type GroupInvite struct {
	InviteCode string  `json:"invite_code"`
	GroupID    string  `json:"group_id"`
	CreatedBy  string  `json:"created_by"`
	CreatedAt  string  `json:"created_at"`
	ExpiresAt  *string `json:"expires_at,omitempty"`
	IsActive   bool    `json:"is_active"`
}
