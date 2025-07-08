package responses

import (
	"backend/models/activity"
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
	GroupImage  *string `json:"group_image,omitempty" example:"https://example.com/image.jpg"` // Optional
	Description *string `json:"description,omitempty" example:"A group for studying algorithms"`
}

type GroupUpdateRequest struct {
	GroupImage  *string `json:"group_image,omitempty" example:"https://example.com/image.jpg"`
	Description *string `json:"description,omitempty" example:"Updated description"`
}

type GroupDeleteRequest struct {
	CreatorID int `json:"creator_id" example:"1"`
}

type AddUserToGroupRequest struct {
	UserID int `json:"user_id" example:"1"`
}

type AddUserByEmailRequest struct {
	Email       string `json:"email" example:"user@example.com"`
	RequesterID int    `json:"requester_id" example:"1"`
}

type RemoveUserFromGroupRequest struct {
	UserID      int `json:"user_id" example:"1"`
	RequesterID int `json:"requester_id" example:"2"`
}

type LeaveGroupRequest struct {
	UserID int `json:"user_id" example:"1"`
}

type AddUserToGroupResponse struct {
	Message string `json:"message" example:"User added to group successfully"`
	GroupID int    `json:"group_id" example:"1"`
	UserID  int    `json:"user_id" example:"1"`
}

type GroupMembersResponse struct {
	GroupID     int                 `json:"group_id" example:"1"`
	Members     []group.GroupMember `json:"members"`
	MemberCount int                 `json:"member_count" example:"3"`
}

type GroupMemberWithUser struct {
	UserID    int     `json:"user_id" example:"1"`
	GroupID   int     `json:"group_id" example:"1"`
	Nickname  *string `json:"nickname,omitempty" example:"Cool Coder"`
	UserName  string  `json:"user_name" example:"John Doe"`
	UserEmail string  `json:"user_email" example:"john@example.com"`
}

type GroupMembersWithUsersResponse struct {
	GroupID     int                   `json:"group_id" example:"1"`
	Members     []GroupMemberWithUser `json:"members"`
	MemberCount int                   `json:"member_count" example:"3"`
}

type GroupActivitiesResponse struct {
	GroupID       int                 `json:"group_id" example:"1"`
	Activities    []activity.Activity `json:"activities"`
	ActivityCount int                 `json:"activity_count" example:"5"`
}

type CreateInviteRequest struct {
	CreatorID int     `json:"creator_id" example:"1"`
	ExpiresAt *string `json:"expires_at,omitempty" example:"2025-12-31T23:59:59Z"`
}

type JoinGroupRequest struct {
	UserID int `json:"user_id" example:"1"`
}

type DeactivateInviteRequest struct {
	RequesterID int `json:"requester_id" example:"1"`
}

type SetNicknameRequest struct {
	UserID      int    `json:"user_id" example:"1"`
	RequesterID int    `json:"requester_id" example:"2"`
	Nickname    string `json:"nickname" example:"Cool Coder"`
}

type DeleteNicknameRequest struct {
	UserID      int `json:"user_id" example:"1"`
	RequesterID int `json:"requester_id" example:"2"`
}

type ActivityCreateRequest struct {
	CreatorID   int    `json:"creator_id" example:"1"`
	GroupID     int    `json:"group_id" example:"1"`
	Title       string `json:"title" example:"Algorithm Contest"`
	Date        string `json:"date" example:"2025-12-31"`
	Description string `json:"description,omitempty" example:"A competitive programming contest"`
	// Note: Image file is handled as multipart form data, not JSON
}

type ActivityUpdateRequest struct {
	Date          *string `json:"date,omitempty" example:"2025-12-31"`
	ActivityImage *string `json:"activity_image,omitempty" example:"https://example.com/image.jpg"`
	Description   *string `json:"description,omitempty" example:"Updated description"`
}

type ActivityDeleteRequest struct {
	CreatorID int `json:"creator_id" example:"1"`
}

type CommentCreateRequest struct {
	UserID  int    `json:"user_id" example:"1"`
	Content string `json:"content" example:"Great activity!"`
}

type CommentDeleteRequest struct {
	RequesterID int `json:"requester_id" example:"1"`
}

type CommentsResponse struct {
	ActivityID   int               `json:"activity_id" example:"1"`
	Comments     []comment.Comment `json:"comments"`
	CommentCount int               `json:"comment_count" example:"5"`
}

type CommentDeleteResponse struct {
	Message   string `json:"message" example:"Comment deleted successfully"`
	CommentID int    `json:"comment_id" example:"1"`
}

// Enhanced comment structure with user information
type CommentWithUser struct {
	ID         int    `json:"id"`
	Content    string `json:"content"`
	UserID     int    `json:"user_id"`
	ActivityID int    `json:"activity_id"`
	CreatedAt  string `json:"created_at"`
	UserName   string `json:"user_name"`
	UserEmail  string `json:"user_email"`
}

type CommentsWithUsersResponse struct {
	ActivityID   int               `json:"activity_id" example:"1"`
	Comments     []CommentWithUser `json:"comments"`
	CommentCount int               `json:"comment_count" example:"5"`
}
