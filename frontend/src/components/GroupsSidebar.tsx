import React, { useState, useEffect, forwardRef, useImperativeHandle } from 'react';
import {
  Box,
  Paper,
  Typography,
  ListItem,
  ListItemButton,
  ListItemText,
  Button,
  Avatar,
  Divider,
  Chip,
  IconButton,
  Tooltip,
} from '@mui/material';
import { 
  Add as AddIcon, 
  Home as HomeIcon,
  Settings as SettingsIcon,
} from '@mui/icons-material';
import { Group, GroupCreateRequest, UserStats } from '../types/api';
import { useAuth } from '../contexts/AuthContext';
import { apiService } from '../services/api';
import CreateGroupDialog from './CreateGroupDialog';
import GroupManagementDialog from './GroupManagementDialog';

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
  const [managementDialogOpen, setManagementDialogOpen] = useState(false);
  const [selectedGroupForManagement, setSelectedGroupForManagement] = useState<Group | null>(null);

  // Expose functions to parent component
  useImperativeHandle(ref, () => ({
    openCreateDialog: () => {
      setCreateDialogOpen(true);
    },
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

  const handleCreateGroup = async (groupData: GroupCreateRequest) => {
    if (!user) return;

    try {
      const newGroup = await apiService.createGroup(groupData);
      const updatedGroups = [...groups, newGroup];
      setGroups(updatedGroups);
      
      // Load rankings for the updated groups list
      await loadGroupRankings(updatedGroups);
      
      // Call the callback to pass updated groups data to parent
      if (onGroupsLoaded) {
        onGroupsLoaded(updatedGroups);
      }
      onGroupsChange();
    } catch (error) {
      console.error('Error creating group:', error);
      throw error; // Re-throw to let the dialog handle it
    }
  };

  const getMedalEmoji = (position: number) => {
    switch (position) {
      case 0: return '🥇';
      case 1: return '🥈';
      case 2: return '🥉';
      default: return '';
    }
  };

  const handleManageGroup = (group: Group, event: React.MouseEvent) => {
    event.stopPropagation();
    setSelectedGroupForManagement(group);
    setManagementDialogOpen(true);
  };

  const handleGroupManagementUpdate = () => {
    loadGroups();
    onGroupsChange();
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
                  <Avatar 
                    sx={{ mr: 2, width: 32, height: 32 }}
                    src={group.group_image ? (group.group_image.startsWith('http') ? group.group_image : `${process.env.REACT_APP_API_URL || 'http://localhost:8080'}${group.group_image}`) : undefined}
                    alt={group.name}
                  >
                    {!group.group_image ? group.name.charAt(0).toUpperCase() : null}
                  </Avatar>
                  <Box sx={{ flex: 1 }}>
                    <Typography variant="body2" fontWeight="medium" noWrap>
                      {group.name}
                    </Typography>
                    <Typography variant="caption" color="text.secondary" noWrap>
                      {group.description}
                    </Typography>
                  </Box>
                  {/* Show manage button for group owners */}
                  {user && group.creator_id === user.id && (
                    <Tooltip title="Manage Group">
                      <IconButton
                        size="small"
                        onClick={(event) => handleManageGroup(group, event)}
                        sx={{ ml: 1 }}
                      >
                        <SettingsIcon fontSize="small" />
                      </IconButton>
                    </Tooltip>
                  )}
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
      <CreateGroupDialog
        open={createDialogOpen}
        onClose={() => setCreateDialogOpen(false)}
        onCreateGroup={handleCreateGroup}
        creatorId={user?.id || ''}
      />

      {/* Group Management Dialog */}
      {selectedGroupForManagement && (
        <GroupManagementDialog
          open={managementDialogOpen}
          onClose={() => {
            setManagementDialogOpen(false);
            setSelectedGroupForManagement(null);
          }}
          group={selectedGroupForManagement}
          onGroupUpdated={handleGroupManagementUpdate}
        />
      )}
    </Paper>
  );
});

GroupsSidebar.displayName = 'GroupsSidebar';

export default GroupsSidebar;
