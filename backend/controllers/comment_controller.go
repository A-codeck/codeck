package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"backend/models/activity"
	"backend/models/comment"
	"backend/models/group"
	"backend/models/responses"

	"github.com/gorilla/mux"
)

type CommentController struct {
	CommentModel  comment.CommentModel
	ActivityModel activity.ActivityModel
	GroupModel    group.GroupModel
}

// swagger imports (used in annotations)
var (
	_ = responses.ErrorResponse{}
)

func NewCommentController(commentModel comment.CommentModel, activityModel activity.ActivityModel, groupModel group.GroupModel) *CommentController {
	return &CommentController{
		CommentModel:  commentModel,
		ActivityModel: activityModel,
		GroupModel:    groupModel,
	}
}

// GetCommentsByActivity godoc
// @Summary Get comments by activity ID
// @Description Get all comments for a specific activity
// @Tags comments
// @Accept json
// @Produce json
// @Param activity_id path string true "Activity ID"
// @Success 200 {object} responses.CommentsResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /activities/{activity_id}/comments [get]
func (cc *CommentController) GetCommentsByActivity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	activityIDStr := vars["activity_id"]
	activityID, err := strconv.Atoi(activityIDStr)
	if err != nil {
		log.Printf("Invalid activity_id: %v", err)
		http.Error(w, "Invalid activity_id", http.StatusBadRequest)
		return
	}

	_, exists := cc.ActivityModel.GetActivityByID(activityID)
	if !exists {
		log.Printf("Activity not found: id=%d", activityID)
		http.Error(w, "Activity not found", http.StatusNotFound)
		return
	}

	comments := cc.CommentModel.GetCommentsByActivityID(activityID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"activity_id":   activityID,
		"comments":      comments,
		"comment_count": len(comments),
	})
}

// CreateComment godoc
// @Summary Create a new comment
// @Description Create a new comment on an activity
// @Tags comments
// @Accept json
// @Produce json
// @Param activity_id path string true "Activity ID"
// @Param comment body responses.CommentCreateRequest true "Comment creation data"
// @Success 201 {object} comment.Comment
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /activities/{activity_id}/comments [post]
func (cc *CommentController) CreateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	activityIDStr := vars["activity_id"]
	activityID, err := strconv.Atoi(activityIDStr)
	if err != nil {
		log.Printf("Invalid activity_id: %v", err)
		http.Error(w, "Invalid activity_id", http.StatusBadRequest)
		return
	}

	_, exists := cc.ActivityModel.GetActivityByID(activityID)
	if !exists {
		log.Printf("Activity not found: id=%d", activityID)
		http.Error(w, "Activity not found", http.StatusNotFound)
		return
	}

	var request struct {
		UserID  int    `json:"user_id"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Failed to decode request payload: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if request.UserID == 0 || request.Content == "" {
		log.Println("Missing required fields (user_id, content)")
		http.Error(w, "Missing required fields (user_id, content)", http.StatusBadRequest)
		return
	}

	newComment := comment.Comment{
		ActivityID: activityID,
		UserID:     request.UserID,
		Content:    request.Content,
	}

	createdComment := cc.CommentModel.CreateComment(newComment)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdComment)
}

// DeleteComment godoc
// @Summary Delete a comment
// @Description Delete a comment (only comment author or activity creator can delete)
// @Tags comments
// @Accept json
// @Produce json
// @Param comment_id path string true "Comment ID"
// @Param request body responses.CommentDeleteRequest true "Delete request with requester_id"
// @Success 200 {object} responses.CommentDeleteResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 403 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /comments/{comment_id} [delete]
func (cc *CommentController) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentIDStr := vars["comment_id"]
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		log.Printf("Invalid comment_id: %v", err)
		http.Error(w, "Invalid comment_id", http.StatusBadRequest)
		return
	}

	existingComment, exists := cc.CommentModel.GetCommentByID(commentID)
	if !exists {
		log.Printf("Comment not found: id=%d", commentID)
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}
	targetActivity, exists := cc.ActivityModel.GetActivityByID(existingComment.ActivityID)
	if !exists {
		http.Error(w, "Activity not found", http.StatusNotFound)
		return
	}

	var request struct {
		RequesterID int `json:"requester_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if request.RequesterID == 0 {
		http.Error(w, "Missing requester_id", http.StatusBadRequest)
		return
	}

	if request.RequesterID != existingComment.UserID && request.RequesterID != targetActivity.CreatorID {
		http.Error(w, "Forbidden: Only comment author or activity creator can delete comments", http.StatusForbidden)
		return
	}

	success := cc.CommentModel.DeleteComment(commentID)
	if !success {
		http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Comment deleted successfully",
		"comment_id": commentID,
	})
}
