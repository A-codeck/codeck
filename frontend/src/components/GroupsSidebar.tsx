import React, { useState, useEffect } from 'react';
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
} from '@mui/material';
import { Add as AddIcon, Group as GroupIcon } from '@mui/icons-material';
import { Group, GroupCreateRequest } from '../types/api';
import { useAuth } from '../contexts/AuthContext';
import { apiService } from '../services/api';

interface GroupsSidebarProps {
  selectedGroupId?: string;
  onGroupSelect: (groupId: string | undefined) => void;
  onGroupsChange: () => void;
}

const GroupsSidebar: React.FC<GroupsSidebarProps> = ({
  selectedGroupId,
  onGroupSelect,
  onGroupsChange,
}) => {
  const { user } = useAuth();
  const [groups, setGroups] = useState<Group[]>([]);
  const [loading, setLoading] = useState(true);
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [newGroupData, setNewGroupData] = useState<GroupCreateRequest>({
    name: '',
    description: '',
    end_date: '',
    creator_id: '',
  });

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
    } catch (error) {
      console.error('Error loading groups:', error);
    } finally {
      setLoading(false);
    }
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
      setGroups([...groups, newGroup]);
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
      <Box sx={{ p: 2, borderBottom: '1px solid', borderColor: 'divider' }}>
        <Typography variant="h3" gutterBottom>
          Groups
        </Typography>
        <Button
          variant="contained"
          color="secondary"
          startIcon={<AddIcon />}
          fullWidth
          onClick={() => setCreateDialogOpen(true)}
        >
          New Group
        </Button>
      </Box>

      <Box sx={{ flex: 1, overflow: 'auto' }}>
        <List>
          <ListItem disablePadding>
            <ListItemButton
              selected={!selectedGroupId}
              onClick={() => onGroupSelect(undefined)}
            >
              <ListItemText 
                primary="All Groups" 
                secondary="See all activities"
              />
            </ListItemButton>
          </ListItem>

          {groups.map((group) => (
            <ListItem key={group.id} disablePadding>
              <ListItemButton
                selected={selectedGroupId === group.id}
                onClick={() => onGroupSelect(group.id)}
              >
                <Avatar sx={{ mr: 2, width: 32, height: 32 }}>
                  <GroupIcon />
                </Avatar>
                <ListItemText
                  primary={group.name}
                  secondary={group.description}
                  primaryTypographyProps={{
                    variant: 'body2',
                    noWrap: true,
                  }}
                  secondaryTypographyProps={{
                    variant: 'caption',
                    noWrap: true,
                  }}
                />
              </ListItemButton>
            </ListItem>
          ))}

          {groups.length === 0 && !loading && (
            <ListItem>
              <ListItemText
                primary="No groups yet"
                secondary="Create your first group to get started"
                sx={{ textAlign: 'center' }}
              />
            </ListItem>
          )}
        </List>
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
};

export default GroupsSidebar;
