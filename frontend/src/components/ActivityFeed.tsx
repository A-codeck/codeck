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
import { ActivityWithGroups, ActivityWithGroup } from '../types/api';
import { useAuth } from '../contexts/AuthContext';
import { apiService } from '../services/api';
import ActivityCard from './ActivityCard';
import AddActivityDialog from './AddActivityDialog';

interface ActivityFeedProps {
  selectedGroupId?: string;
}

const ActivityFeed: React.FC<ActivityFeedProps> = ({ selectedGroupId }) => {
  const { user } = useAuth();
  const [activities, setActivities] = useState<(ActivityWithGroup | ActivityWithGroups)[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [addDialogOpen, setAddDialogOpen] = useState(false);
  const [refreshTrigger, setRefreshTrigger] = useState(0);

  useEffect(() => {
    loadActivities();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedGroupId, user, refreshTrigger]);

  const loadActivities = async () => {
    if (!user) return;

    try {
      setLoading(true);
      setError(''); // Clear any previous errors
      
      let activitiesData: (ActivityWithGroup | ActivityWithGroups)[] = [];

      if (selectedGroupId) {
        // Load activities for specific group
        const groupActivities = await apiService.getGroupActivities(selectedGroupId, user.id);
        activitiesData = groupActivities.map(activity => ({ ...activity }));
      } else {
        // Load user's feed from all groups with group information
        const feedActivities = await apiService.getUserFeedWithGroups(user.id);
        activitiesData = feedActivities;
      }

      // Sort by date (newest first)
      activitiesData.sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime());
      
      setActivities(activitiesData);
      setError(''); // Clear error on successful load
    } catch (err: any) {
      console.error('Error loading activities:', err);
      console.error('Error response:', err.response?.data);
      console.error('Error status:', err.response?.status);
      
      // Provide more specific error messages
      let errorMessage = 'Failed to load activities';
      if (err.response?.status === 403) {
        errorMessage = 'You do not have permission to view this group\'s activities';
      } else if (err.response?.status === 404) {
        errorMessage = 'Group not found or you are not a member';
      } else if (err.response?.data?.error) {
        errorMessage = err.response.data.error;
      }
      
      setError(errorMessage);
      setActivities([]); // Clear activities on error
    } finally {
      setLoading(false);
    }
  };

  const handleActivityAdded = () => {
    // Trigger a refresh by updating the refresh trigger
    setRefreshTrigger(prev => {
      console.log('ActivityFeed: refreshTrigger updated from', prev, 'to', prev + 1);
      return prev + 1;
    });
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
    <Box 
      sx={{ 
        flex: 1, 
        overflow: 'auto',
        p: 3
      }}
    >
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
        <Alert severity="error" sx={{ mb: 3, backgroundColor: (theme) => theme.palette.background.default }}>
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
        preselectedGroupId={selectedGroupId}
      />
    </Box>
  );
};

export default ActivityFeed;
