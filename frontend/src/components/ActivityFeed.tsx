import React, { useState, useEffect } from 'react';
import {
  Box,
  Paper,
  Typography,
  Button,
  CircularProgress,
  Alert,
} from '@mui/material';
import { Add as AddIcon } from '@mui/icons-material';
import { ActivityWithGroup } from '../types/api';
import { useAuth } from '../contexts/AuthContext';
import { apiService } from '../services/api';
import ActivityCard from './ActivityCard';
import AddActivityDialog from './AddActivityDialog';

interface ActivityFeedProps {
  selectedGroupId?: string;
}

const ActivityFeed: React.FC<ActivityFeedProps> = ({ selectedGroupId }) => {
  const { user } = useAuth();
  const [activities, setActivities] = useState<ActivityWithGroup[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [addDialogOpen, setAddDialogOpen] = useState(false);

  useEffect(() => {
    loadActivities();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedGroupId, user]);

  const loadActivities = async () => {
    if (!user) return;

    try {
      setLoading(true);
      setError('');
      
      let activitiesData: ActivityWithGroup[] = [];

      if (selectedGroupId) {
        // Load activities for specific group
        const groupActivities = await apiService.getGroupActivities(selectedGroupId, user.id);
        activitiesData = groupActivities.map(activity => ({ ...activity }));
      } else {
        // Load user's feed from all groups
        const feedActivities = await apiService.getUserFeed(user.id);
        activitiesData = feedActivities.map(activity => ({ ...activity }));
      }

      // Sort by date (newest first)
      activitiesData.sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime());
      
      setActivities(activitiesData);
    } catch (err: any) {
      console.error('Error loading activities:', err);
      setError('Failed to load activities');
    } finally {
      setLoading(false);
    }
  };

  const handleActivityAdded = () => {
    loadActivities();
  };

  const getFeedTitle = () => {
    if (selectedGroupId) {
      return 'Group Activities';
    }
    return 'Activity Feed';
  };

  const getFeedSubtitle = () => {
    if (selectedGroupId) {
      return 'Activities from this group';
    }
    return 'Activities from all your groups';
  };

  if (loading) {
    return (
      <Box
        sx={{
          flex: 1,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          minHeight: 400,
        }}
      >
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box sx={{ flex: 1, p: 3 }}>
      {/* Header */}
      <Paper sx={{ p: 3, mb: 3 }}>
        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          <Box>
            <Typography variant="h2" gutterBottom>
              {getFeedTitle()}
            </Typography>
            <Typography variant="body1" color="text.secondary">
              {getFeedSubtitle()}
            </Typography>
          </Box>
          <Button
            variant="contained"
            color="secondary"
            startIcon={<AddIcon />}
            onClick={() => setAddDialogOpen(true)}
            size="large"
          >
            Log Activity
          </Button>
        </Box>
      </Paper>

      {/* Error State */}
      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      {/* Activities List */}
      <Box>
        {activities.length > 0 ? (
          activities.map((activity) => (
            <ActivityCard
              key={activity.id}
              activity={activity}
              showGroupTag={!selectedGroupId}
            />
          ))
        ) : (
          <Paper sx={{ p: 4, textAlign: 'center' }}>
            <Typography variant="h3" color="text.secondary" gutterBottom>
              No activities yet
            </Typography>
            <Typography variant="body1" color="text.secondary" sx={{ mb: 3 }}>
              {selectedGroupId 
                ? 'No activities have been posted in this group yet.'
                : 'Start by logging your first coding activity!'
              }
            </Typography>
            <Button
              variant="contained"
              color="secondary"
              startIcon={<AddIcon />}
              onClick={() => setAddDialogOpen(true)}
            >
              Log Your First Activity
            </Button>
          </Paper>
        )}
      </Box>

      {/* Add Activity Dialog */}
      <AddActivityDialog
        open={addDialogOpen}
        onClose={() => setAddDialogOpen(false)}
        onActivityAdded={handleActivityAdded}
      />
    </Box>
  );
};

export default ActivityFeed;
