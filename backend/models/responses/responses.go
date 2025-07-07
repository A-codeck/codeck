package responses

import (
	"backend/models/comment"
	"backend/models/group"
	"backend/models/user"
)

type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request"`
}

type SuccessResponse struct {
	Message string `json:"message" example:"Operation completed successfully"`
}

type UserCreateRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Name     string `json:"name" example:"John Doe"`
	Password string `json:"password" example:"password123"`
}

type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123"`
}

type LoginResponse struct {
	User  user.User `json:"user"`
	Token string    `json:"token" example:"dummy-jwt-token-1"`
}

type GroupCreateRequest struct {
	Name        string  `json:"name" example:"Study Group"`
	EndDate     string  `json:"end_date" example:"2025-12-31"`
	GroupImage  *string `json:"group_image,omitempty" example:"https://example.com/image.jpg"`
	Description *string `json:"description,omitempty" example:"A group for studying algorithms"`
}

type GroupUpdateRequest struct {
	EndDate     *string `json:"end_date,omitempty" example:"2025-12-31"`
	GroupImage  *string `json:"group_image,omitempty" example:"https://example.com/image.jpg"`
	Description *string `json:"description,omitempty" example:"Updated description"`
}

type GroupDeleteRequest struct {
	CreatorID string `json:"creator_id" example:"user123"`
}

type AddUserToGroupRequest struct {
	UserID string `json:"user_id" example:"user123"`
}

type RemoveUserFromGroupRequest struct {
	UserID      string `json:"user_id" example:"user123"`
	RequesterID string `json:"requester_id" example:"user456"`
}

type AddUserToGroupResponse struct {
	Message string `json:"message" example:"User added to group successfully"`
	GroupID string `json:"group_id" example:"group123"`
	UserID  string `json:"user_id" example:"user123"`
}

type GroupMembersResponse struct {
	GroupID     string              `json:"group_id" example:"group123"`
	Members     []group.GroupMember `json:"members"`
	MemberCount int                 `json:"member_count" example:"3"`
}

type CreateInviteRequest struct {
	CreatorID string  `json:"creator_id" example:"user123"`
	ExpiresAt *string `json:"expires_at,omitempty" example:"2025-12-31T23:59:59Z"`
}

type JoinGroupRequest struct {
	UserID string `json:"user_id" example:"user123"`
}

type DeactivateInviteRequest struct {
	RequesterID string `json:"requester_id" example:"user123"`
}

type SetNicknameRequest struct {
	UserID      string `json:"user_id" example:"user123"`
	RequesterID string `json:"requester_id" example:"user456"`
	Nickname    string `json:"nickname" example:"Cool Coder"`
}

type DeleteNicknameRequest struct {
	UserID      string `json:"user_id" example:"user123"`
	RequesterID string `json:"requester_id" example:"user456"`
}

type ActivityCreateRequest struct {
	Title         string  `json:"title" example:"Algorithm Contest"`
	Date          string  `json:"date" example:"2025-12-31"`
	ActivityImage *string `json:"activity_image,omitempty" example:"https://example.com/image.jpg"`
	Description   *string `json:"description,omitempty" example:"A competitive programming contest"`
}

type ActivityUpdateRequest struct {
	Date          *string `json:"date,omitempty" example:"2025-12-31"`
	ActivityImage *string `json:"activity_image,omitempty" example:"https://example.com/image.jpg"`
	Description   *string `json:"description,omitempty" example:"Updated description"`
}

type ActivityDeleteRequest struct {
	CreatorID string `json:"creator_id" example:"user123"`
}

type CommentCreateRequest struct {
	UserID  string `json:"user_id" example:"user123"`
	Content string `json:"content" example:"Great activity!"`
}

type CommentDeleteRequest struct {
	RequesterID string `json:"requester_id" example:"user123"`
}

type CommentsResponse struct {
	ActivityID   string            `json:"activity_id" example:"activity123"`
	Comments     []comment.Comment `json:"comments"`
	CommentCount int               `json:"comment_count" example:"5"`
}

type CommentDeleteResponse struct {
	Message   string `json:"message" example:"Comment deleted successfully"`
	CommentID string `json:"comment_id" example:"comment123"`
}
