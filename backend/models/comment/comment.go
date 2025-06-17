package comment

type Comment struct {
	ID         string `json:"id"`
	ActivityID string `json:"activity_id"`
	UserID     string `json:"user_id"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
}
