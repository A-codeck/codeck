import React, { useState, useEffect, forwardRef, useImperativeHandle } from 'react';
import {
  Box,
  Paper,
  Typography,
  List,
  ListItem,
  ListItemButton,
  ListItemText,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Avatar,
  Divider,
  Chip,
} from '@mui/material';
import { 
  Add as AddIcon, 
  Group as GroupIcon, 
  Home as HomeIcon 
} from '@mui/icons-material';
import { Group, GroupCreateRequest, UserStats } from '../types/api';
import { useAuth } from '../contexts/AuthContext';
import { apiService } from '../services/api';

interface GroupsSidebarProps {
  selectedGroupId?: string;
  onGroupSelect: (groupId: string | undefined) => void;
  onGroupsChange: () => void;
  onGroupsLoaded?: (groups: Group[]) => void; // New callback to pass groups data
}

export interface GroupsSidebarRef {
  openCreateDialog: () => void;
}

const GroupsSidebar = forwardRef<GroupsSidebarRef, GroupsSidebarProps>(({
  selectedGroupId,
  onGroupSelect,
  onGroupsChange,
  onGroupsLoaded,
}, ref) => {
  const { user } = useAuth();
  const [groups, setGroups] = useState<Group[]>([]);
  const [groupRankings, setGroupRankings] = useState<{ [groupId: string]: UserStats[] }>({});
  const [loading, setLoading] = useState(true);
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [newGroupData, setNewGroupData] = useState<GroupCreateRequest>({
    name: '',
    description: '',
    end_date: '',
    creator_id: '',
  });

  // Expose functions to parent component
  useImperativeHandle(ref, () => ({
    openCreateDialog: () => setCreateDialogOpen(true),
  }));

  // Load user's groups using the new API endpoint
  useEffect(() => {
    loadGroups();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [user]);

  const loadGroups = async () => {
    if (!user) return;

    try {
      setLoading(true);
      const userGroups = await apiService.getUserGroups(user.id);
      setGroups(userGroups);
      
      // Load rankings for each group
      await loadGroupRankings(userGroups);
      
      // Call the callback to pass groups data to parent
      if (onGroupsLoaded) {
        onGroupsLoaded(userGroups);
      }
    } catch (error) {
      console.error('Error loading groups:', error);
      // Even if there's an error, we should call the callback with empty array
      if (onGroupsLoaded) {
        onGroupsLoaded([]);
      }
    } finally {
      setLoading(false);
    }
  };

  const loadGroupRankings = async (groupList: Group[]) => {
    if (!user) return;

    const rankings: { [groupId: string]: UserStats[] } = {};

    for (const group of groupList) {
      try {
        // Load group members
        const membersData = await apiService.getGroupMembers(group.id, user.id);
        
        // Calculate rankings by getting activities for each member
        const userStats: UserStats[] = [];
        
        for (const member of membersData.members) {
          try {
            // Get user info
            const userData = await apiService.getUser(member.user_id);
            
            // Get user's activities in this group
            const groupActivities = await apiService.getGroupActivities(group.id, user.id);
            const userActivities = groupActivities.filter(activity => activity.creator_id === member.user_id);
            
            userStats.push({
              user_id: member.user_id,
              user_name: userData.name,
              activity_count: userActivities.length,
            });
          } catch (error) {
            console.error(`Error loading data for user ${member.user_id}:`, error);
          }
        }
        
        // Sort by activity count (descending) and take top 3
        userStats.sort((a, b) => b.activity_count - a.activity_count);
        rankings[group.id] = userStats.slice(0, 3);
        
      } catch (error) {
        console.error(`Error loading ranking for group ${group.id}:`, error);
        rankings[group.id] = [];
      }
    }
    
    setGroupRankings(rankings);
  };

  const handleCreateGroup = async () => {
    if (!user) return;

    try {
      // Add creator_id before sending the request
      const groupDataWithCreator = {
        ...newGroupData,
        creator_id: user.id,
      };
      
      const newGroup = await apiService.createGroup(groupDataWithCreator);
      const updatedGroups = [...groups, newGroup];
      setGroups(updatedGroups);
      
      // Load rankings for the updated groups list
      await loadGroupRankings(updatedGroups);
      
      // Call the callback to pass updated groups data to parent
      if (onGroupsLoaded) {
        onGroupsLoaded(updatedGroups);
      }
      setCreateDialogOpen(false);
      setNewGroupData({
        name: '',
        description: '',
        end_date: '',
        creator_id: '',
      });
      onGroupsChange();
    } catch (error) {
      console.error('Error creating group:', error);
    }
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setNewGroupData({
      ...newGroupData,
      [e.target.name]: e.target.value,
    });
  };

  const getMedalEmoji = (position: number) => {
    switch (position) {
      case 0: return '🥇';
      case 1: return '🥈';
      case 2: return '🥉';
      default: return '';
    }
  };

  return (
    <Paper
      sx={{
        width: 280,
        height: '100%',
        display: 'flex',
        flexDirection: 'column',
        borderRight: '1px solid',
        borderColor: 'divider',
        borderRadius: 0,
      }}
    >
      <Box sx={{ flex: 1, overflow: 'auto', p: 2 }}>
        {/* Home Section */}
        <ListItem disablePadding sx={{ mb: 1 }}>
          <ListItemButton
            selected={!selectedGroupId}
            onClick={() => onGroupSelect(undefined)}
            sx={{ borderRadius: 1 }}
          >
            <Avatar sx={{ mr: 2, width: 32, height: 32, bgcolor: 'primary.main' }}>
              <HomeIcon />
            </Avatar>
            <ListItemText 
              primary="Home" 
              secondary="See all activities"
              primaryTypographyProps={{
                variant: 'body2',
                fontWeight: 'medium',
              }}
              secondaryTypographyProps={{
                variant: 'caption',
              }}
            />
          </ListItemButton>
        </ListItem>

        <Divider sx={{ my: 2 }} />

        {/* Groups Section */}
        <Typography variant="caption" color="text.secondary" sx={{ px: 1, mb: 1, display: 'block', textTransform: 'uppercase', letterSpacing: 1 }}>
          Groups
        </Typography>

        {groups.map((group, groupIndex) => (
          <Box key={group.id}>
            <ListItem disablePadding sx={{ mb: 1 }}>
              <ListItemButton
                selected={selectedGroupId === group.id}
                onClick={() => onGroupSelect(group.id)}
                sx={{ borderRadius: 1, flexDirection: 'column', alignItems: 'flex-start', py: 1.5 }}
              >
                <Box sx={{ display: 'flex', alignItems: 'center', width: '100%', mb: 1 }}>
                  <Avatar sx={{ mr: 2, width: 32, height: 32, bgcolor: 'secondary.main' }}>
                    <GroupIcon />
                  </Avatar>
                  <Box sx={{ flex: 1 }}>
                    <Typography variant="body2" fontWeight="medium" noWrap>
                      {group.name}
                    </Typography>
                    <Typography variant="caption" color="text.secondary" noWrap>
                      {group.description}
                    </Typography>
                  </Box>
                </Box>
                
                {/* Top 3 Users */}
                {groupRankings[group.id] && groupRankings[group.id].length > 0 && (
                  <Box sx={{ width: '100%', pl: 5 }}>
                    {groupRankings[group.id].map((userStat, index) => (
                      <Box key={userStat.user_id} sx={{ display: 'flex', alignItems: 'center', py: 0.25 }}>
                        <Typography variant="caption" sx={{ mr: 0.5 }}>
                          {getMedalEmoji(index)}
                        </Typography>
                        <Typography variant="caption" color="text.secondary" noWrap sx={{ flex: 1 }}>
                          {userStat.user_name}
                        </Typography>
                        <Chip 
                          label={userStat.activity_count} 
                          size="small" 
                          sx={{ 
                            height: 16, 
                            fontSize: '0.625rem',
                            '& .MuiChip-label': { px: 0.75 }
                          }}
                        />
                      </Box>
                    ))}
                  </Box>
                )}
              </ListItemButton>
            </ListItem>
            {groupIndex < groups.length - 1 && <Divider sx={{ my: 1 }} />}
          </Box>
        ))}

        {groups.length === 0 && !loading && (
          <Box sx={{ textAlign: 'center', py: 4 }}>
            <Typography variant="body2" color="text.secondary" gutterBottom>
              No groups yet
            </Typography>
            <Typography variant="caption" color="text.secondary">
              Create your first group to get started
            </Typography>
          </Box>
        )}

        <Divider sx={{ my: 2 }} />

        {/* Add Group Button */}
        <Button
          variant="outlined"
          color="secondary"
          startIcon={<AddIcon />}
          fullWidth
          onClick={() => setCreateDialogOpen(true)}
          sx={{ borderRadius: 1 }}
        >
          Add Group
        </Button>
      </Box>

      {/* Create Group Dialog */}
      <Dialog
        open={createDialogOpen}
        onClose={() => setCreateDialogOpen(false)}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>Create New Group</DialogTitle>
        <DialogContent>
          <TextField
            fullWidth
            name="name"
            label="Group Name"
            value={newGroupData.name}
            onChange={handleInputChange}
            required
            sx={{ mb: 2, mt: 1 }}
          />
          <TextField
            fullWidth
            name="description"
            label="Description"
            value={newGroupData.description}
            onChange={handleInputChange}
            multiline
            rows={3}
            sx={{ mb: 2 }}
          />
          <TextField
            fullWidth
            name="end_date"
            label="End Date"
            type="date"
            value={newGroupData.end_date}
            onChange={handleInputChange}
            required
            InputLabelProps={{
              shrink: true,
            }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setCreateDialogOpen(false)}>
            Cancel
          </Button>
          <Button
            onClick={handleCreateGroup}
            variant="contained"
            color="secondary"
            disabled={!newGroupData.name || !newGroupData.end_date}
          >
            Create Group
          </Button>
        </DialogActions>
      </Dialog>
    </Paper>
  );
});

GroupsSidebar.displayName = 'GroupsSidebar';

export default GroupsSidebar;
