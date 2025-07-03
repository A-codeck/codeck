package comment

type CommentModel interface {
	GetCommentByID(id int) (Comment, bool)
	GetCommentsByActivityID(activityID int) []Comment
	CreateComment(comment Comment) Comment
	DeleteComment(id int) bool
}

// DefaultCommentModel must be set in main.go after DB initialization
var DefaultCommentModel CommentModel
