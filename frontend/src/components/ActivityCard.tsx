import React, { useState, useEffect } from 'react';
import {
  Card,
  CardContent,
  Typography,
  Box,
  Button,
  Chip,
  Avatar,
  Collapse,
  TextField,
  IconButton,
  List,
  ListItem,
  ListItemAvatar,
  ListItemText,
  Divider,
} from '@mui/material';
import {
  ExpandMore as ExpandMoreIcon,
  ExpandLess as ExpandLessIcon,
  Comment as CommentIcon,
  Send as SendIcon,
} from '@mui/icons-material';
import { ActivityWithGroup, Comment, User } from '../types/api';
import { useAuth } from '../contexts/AuthContext';
import { apiService } from '../services/api';

interface ActivityCardProps {
  activity: ActivityWithGroup;
  showGroupTag?: boolean;
}

const ActivityCard: React.FC<ActivityCardProps> = ({ 
  activity, 
  showGroupTag = true 
}) => {
  const { user } = useAuth();
  const [commentsOpen, setCommentsOpen] = useState(false);
  const [comments, setComments] = useState<Comment[]>([]);
  const [newComment, setNewComment] = useState('');
  const [loading, setLoading] = useState(false);
  const [creator, setCreator] = useState<User | null>(null);

  useEffect(() => {
    loadCreator();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [activity.creator_id]);

  const loadCreator = async () => {
    try {
      const creatorData = await apiService.getUser(activity.creator_id);
      setCreator(creatorData);
    } catch (error) {
      console.error('Error loading creator:', error);
    }
  };

  const loadComments = async () => {
    try {
      setLoading(true);
      const commentsData = await apiService.getActivityComments(activity.id);
      setComments(commentsData.comments);
    } catch (error) {
      console.error('Error loading comments:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleCommentsToggle = () => {
    if (!commentsOpen) {
      loadComments();
    }
    setCommentsOpen(!commentsOpen);
  };

  const handleAddComment = async () => {
    if (!user || !newComment.trim()) return;

    try {
      const comment = await apiService.createComment(activity.id, {
        content: newComment.trim(),
        user_id: user.id,
      });
      setComments([...comments, comment]);
      setNewComment('');
    } catch (error) {
      console.error('Error adding comment:', error);
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  };

  return (
    <Card sx={{ mb: 2 }}>
      <CardContent>
        {/* Header */}
        <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
          <Avatar sx={{ mr: 2 }}>
            {creator?.name?.charAt(0) || 'U'}
          </Avatar>
          <Box sx={{ flex: 1 }}>
            <Typography variant="body2" fontWeight={600}>
              {creator?.name || 'Unknown User'}
            </Typography>
            <Typography variant="caption" color="text.secondary">
              {formatDate(activity.date)}
            </Typography>
          </Box>
          {showGroupTag && activity.group && (
            <Chip
              label={activity.group.name}
              size="small"
              variant="outlined"
              color="secondary"
            />
          )}
        </Box>

        {/* Content */}
        <Typography variant="h3" gutterBottom>
          {activity.title}
        </Typography>
        
        <Typography variant="body1" color="text.secondary" sx={{ mb: 2 }}>
          {activity.description}
        </Typography>

        {activity.activity_image && (
          <Box
            component="img"
            src={activity.activity_image}
            alt={activity.title}
            sx={{
              width: '100%',
              maxHeight: 300,
              objectFit: 'cover',
              borderRadius: 1,
              mb: 2,
            }}
          />
        )}

        {/* Actions */}
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <Button
            startIcon={<CommentIcon />}
            endIcon={commentsOpen ? <ExpandLessIcon /> : <ExpandMoreIcon />}
            onClick={handleCommentsToggle}
            size="small"
            variant="text"
          >
            {comments.length} Comments
          </Button>
        </Box>

        {/* Comments Section */}
        <Collapse in={commentsOpen}>
          <Box sx={{ mt: 2 }}>
            <Divider sx={{ mb: 2 }} />
            
            {/* Add Comment */}
            <Box sx={{ display: 'flex', gap: 1, mb: 2 }}>
              <Avatar sx={{ width: 32, height: 32 }}>
                {user?.name?.charAt(0)}
              </Avatar>
              <TextField
                placeholder="Add a comment..."
                variant="outlined"
                size="small"
                fullWidth
                value={newComment}
                onChange={(e) => setNewComment(e.target.value)}
                onKeyPress={(e) => {
                  if (e.key === 'Enter' && !e.shiftKey) {
                    e.preventDefault();
                    handleAddComment();
                  }
                }}
              />
              <IconButton 
                onClick={handleAddComment}
                disabled={!newComment.trim()}
                color="secondary"
              >
                <SendIcon />
              </IconButton>
            </Box>

            {/* Comments List */}
            {loading ? (
              <Typography variant="body2" color="text.secondary" sx={{ textAlign: 'center', py: 2 }}>
                Loading comments...
              </Typography>
            ) : (
              <List sx={{ maxHeight: 300, overflow: 'auto' }}>
                {comments.map((comment) => (
                  <ListItem key={comment.id} alignItems="flex-start" sx={{ px: 0 }}>
                    <ListItemAvatar>
                      <Avatar sx={{ width: 32, height: 32 }}>
                        {/* TODO: Get user name from comment.user_id */}
                        U
                      </Avatar>
                    </ListItemAvatar>
                    <ListItemText
                      primary={
                        <Typography variant="body2">
                          {comment.content}
                        </Typography>
                      }
                      secondary={
                        <Typography variant="caption" color="text.secondary">
                          {new Date(comment.created_at).toLocaleString()}
                        </Typography>
                      }
                    />
                  </ListItem>
                ))}
                {comments.length === 0 && (
                  <Typography variant="body2" color="text.secondary" sx={{ textAlign: 'center', py: 2 }}>
                    No comments yet. Be the first to comment!
                  </Typography>
                )}
              </List>
            )}
          </Box>
        </Collapse>
      </CardContent>
    </Card>
  );
};

export default ActivityCard;
