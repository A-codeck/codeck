package comment

type CommentModel interface {
	GetCommentByID(id string) (Comment, bool)
	GetCommentsByActivityID(activityID string) []Comment
	CreateComment(comment Comment) Comment
	DeleteComment(id string) bool
}

var DefaultCommentModel CommentModel = NewInMemoryComment()
