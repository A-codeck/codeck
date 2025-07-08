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

  const handleCreateGroup = async (groupData: GroupCreateRequest) => {
    if (!user) return;

    try {
      const newGroup = await apiService.createGroup(groupData);
      const updatedGroups = [...groups, newGroup];
      setGroups(updatedGroups);
      
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
                  <Box sx={{ flex: 1, minWidth: 0 }}>
                    <Typography variant="body2" fontWeight="medium" noWrap>
                      {group.name}
                    </Typography>
                    <Typography 
                      variant="caption" 
                      color="text.secondary" 
                      sx={{
                        display: '-webkit-box',
                        WebkitLineClamp: 2,
                        WebkitBoxOrient: 'vertical',
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                        lineHeight: 1.2,
                        maxHeight: '2.4em'
                      }}
                    >
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
