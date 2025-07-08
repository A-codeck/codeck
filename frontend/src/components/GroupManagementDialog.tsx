import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  Box,
  Typography,
  List,
  ListItem,
  ListItemText,
  IconButton,
  Tab,
  Tabs,
  Chip,
  Paper,
  InputAdornment,
  Snackbar,
  Alert,
  CircularProgress,
} from '@mui/material';
import {
  PersonAdd as PersonAddIcon,
  Link as LinkIcon,
  ContentCopy as ContentCopyIcon,
  Email as EmailIcon,
  ExitToApp as ExitToAppIcon,
  Delete as DeleteIcon,
} from '@mui/icons-material';
import { Group, GroupMember, GroupInvite } from '../types/api';
import { apiService } from '../services/api';
import { useAuth } from '../contexts/AuthContext';

interface GroupManagementDialogProps {
  open: boolean;
  onClose: () => void;
  group: Group;
  onGroupUpdated: () => void;
}

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel({ children, value, index }: TabPanelProps) {
  return (
    <div hidden={value !== index} style={{ paddingTop: 16 }}>
      {value === index && children}
    </div>
  );
}

const GroupManagementDialog: React.FC<GroupManagementDialogProps> = ({
  open,
  onClose,
  group,
  onGroupUpdated,
}) => {
  const { user } = useAuth();
  const [tabValue, setTabValue] = useState(0);
  const [members, setMembers] = useState<GroupMember[]>([]);
  const [invites, setInvites] = useState<GroupInvite[]>([]);
  const [emailToAdd, setEmailToAdd] = useState('');
  const [loading, setLoading] = useState(false);
  const [snackbar, setSnackbar] = useState<{ open: boolean; message: string; severity: 'success' | 'error' }>({
    open: false,
    message: '',
    severity: 'success',
  });

  const isOwner = user && group.creator_id === user.id;

  useEffect(() => {
    if (open && user) {
      loadMembers();
      loadInvites();
    }
  }, [open, user, group.id]);

  const loadMembers = async () => {
    if (!user) return;
    try {
      const response = await apiService.getGroupMembers(group.id, user.id);
      setMembers(response.members);
    } catch (error) {
      console.error('Error loading members:', error);
    }
  };

  const loadInvites = async () => {
    if (!user) return;
    try {
      const invites = await apiService.getGroupInvites(group.id);
      setInvites(invites);
    } catch (error) {
      console.error('Error loading invites:', error);
    }
  };

  const handleAddUserByEmail = async () => {
    if (!user || !emailToAdd.trim()) return;
    
    setLoading(true);
    try {
      await apiService.addUserToGroupByEmail(group.id, emailToAdd.trim(), user.id);
      setEmailToAdd('');
      await loadMembers();
      setSnackbar({
        open: true,
        message: 'User added successfully!',
        severity: 'success',
      });
    } catch (error: any) {
      setSnackbar({
        open: true,
        message: error.response?.data?.error || 'Failed to add user',
        severity: 'error',
      });
    } finally {
      setLoading(false);
    }
  };

  const handleCreateInviteLink = async () => {
    if (!user) return;
    
    setLoading(true);
    try {
      const newInvite = await apiService.createGroupInvite(group.id, user.id); // No expiry = permanent
      
      // Automatically copy the invite link to clipboard
      const inviteUrl = `${window.location.origin}/join/${newInvite.invite_code}`;
      await navigator.clipboard.writeText(inviteUrl);
      
      await loadInvites();
      setSnackbar({
        open: true,
        message: `Invite link created and copied! Share: ${inviteUrl}`,
        severity: 'success',
      });
    } catch (error: any) {
      setSnackbar({
        open: true,
        message: error.response?.data?.error || 'Failed to create invite link',
        severity: 'error',
      });
    } finally {
      setLoading(false);
    }
  };

  const handleCopyInviteLink = (inviteCode: string) => {
    const inviteUrl = `${window.location.origin}/join/${inviteCode}`;
    navigator.clipboard.writeText(inviteUrl);
    setSnackbar({
      open: true,
      message: 'Invite link copied to clipboard!',
      severity: 'success',
    });
  };

  const handleDeleteInvite = async (inviteCode: string) => {
    if (!user) return;
    
    if (!window.confirm('Are you sure you want to delete this invite link? It will no longer work for joining the group.')) {
      return;
    }
    
    setLoading(true);
    try {
      await apiService.deactivateInvite(inviteCode, user.id);
      await loadInvites();
      setSnackbar({
        open: true,
        message: 'Invite link deleted successfully!',
        severity: 'success',
      });
    } catch (error: any) {
      setSnackbar({
        open: true,
        message: error.response?.data?.error || 'Failed to delete invite link',
        severity: 'error',
      });
    } finally {
      setLoading(false);
    }
  };

  const handleLeaveGroup = async () => {
    if (!user) return;
    
    const confirmMessage = isOwner 
      ? 'Are you sure you want to leave this group? As the owner, this will delete the group and remove all members.'
      : 'Are you sure you want to leave this group?';
    
    if (!window.confirm(confirmMessage)) return;
    
    setLoading(true);
    try {
      const result = await apiService.leaveGroup(group.id, user.id);
      setSnackbar({
        open: true,
        message: result.group_deleted ? 'Group deleted successfully!' : 'Left group successfully!',
        severity: 'success',
      });
      onGroupUpdated();
      onClose();
    } catch (error: any) {
      setSnackbar({
        open: true,
        message: error.response?.data?.error || 'Failed to leave group',
        severity: 'error',
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
        <DialogTitle>
          <Box display="flex" justifyContent="space-between" alignItems="center">
            <Typography variant="h6">Manage Group: {group.name}</Typography>
            <Button
              startIcon={<ExitToAppIcon />}
              color="error"
              variant="outlined"
              onClick={handleLeaveGroup}
              disabled={loading}
            >
              {isOwner ? 'Delete Group' : 'Leave Group'}
            </Button>
          </Box>
        </DialogTitle>

        <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
          <Tabs value={tabValue} onChange={(_, newValue) => setTabValue(newValue)}>
            <Tab label="Members" />
            {isOwner && <Tab label="Add Members" />}
            {isOwner && <Tab label="Invite Links" />}
          </Tabs>
        </Box>

        <DialogContent>
          <TabPanel value={tabValue} index={0}>
            <Typography variant="h6" gutterBottom>
              Group Members ({members.length})
            </Typography>
            <List>
              {members.map((member) => (
                <ListItem key={member.user_id}>
                  <ListItemText
                    primary={member.nickname || `User ${member.user_id}`}
                    secondary={member.user_id === group.creator_id ? 'Owner' : 'Member'}
                  />
                  {member.user_id === group.creator_id && (
                    <Chip label="Owner" color="primary" size="small" />
                  )}
                </ListItem>
              ))}
            </List>
          </TabPanel>

          {isOwner && (
            <TabPanel value={tabValue} index={1}>
              <Typography variant="h6" gutterBottom>
                Add User by Email
              </Typography>
              <Box sx={{ display: 'flex', gap: 2, mb: 3 }}>
                <TextField
                  fullWidth
                  label="Email Address"
                  value={emailToAdd}
                  onChange={(e) => setEmailToAdd(e.target.value)}
                  InputProps={{
                    startAdornment: (
                      <InputAdornment position="start">
                        <EmailIcon />
                      </InputAdornment>
                    ),
                  }}
                  disabled={loading}
                />
                <Button
                  variant="contained"
                  onClick={handleAddUserByEmail}
                  disabled={!emailToAdd.trim() || loading}
                  startIcon={loading ? <CircularProgress size={20} /> : <PersonAddIcon />}
                >
                  Add User
                </Button>
              </Box>
            </TabPanel>
          )}

          {isOwner && (
            <TabPanel value={tabValue} index={2}>
              <Typography variant="h6" gutterBottom>
                Invite Links
              </Typography>
              <Box sx={{ mb: 3 }}>
                <Button
                  variant="contained"
                  onClick={handleCreateInviteLink}
                  disabled={loading}
                  startIcon={loading ? <CircularProgress size={20} /> : <LinkIcon />}
                >
                  Generate New Invite Link
                </Button>
              </Box>
              
              {invites.length > 0 && (
                <List>
                  {invites.filter(invite => invite.is_active).map((invite) => {
                    const inviteUrl = `${window.location.origin}/join/${invite.invite_code}`;
                    return (
                      <ListItem key={invite.invite_code}>
                        <Paper sx={{ p: 3, width: '100%' }}>
                          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                            <Box sx={{ flex: 1, mr: 2 }}>
                              <Typography variant="body2" color="text.secondary" gutterBottom>
                                Invite Code: {invite.invite_code}
                              </Typography>
                              <Typography variant="body2" sx={{ 
                                fontFamily: 'monospace', 
                                backgroundColor: 'grey.100', 
                                p: 1, 
                                borderRadius: 1,
                                wordBreak: 'break-all',
                                mb: 1
                              }}>
                                {inviteUrl}
                              </Typography>
                              <Typography variant="caption" color="text.secondary">
                                Created: {new Date(invite.created_at).toLocaleDateString()} • Permanent Link
                              </Typography>
                            </Box>
                            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
                              <IconButton
                                onClick={() => handleCopyInviteLink(invite.invite_code)}
                                color="primary"
                                size="small"
                                title="Copy Link"
                              >
                                <ContentCopyIcon />
                              </IconButton>
                              <IconButton
                                onClick={() => handleDeleteInvite(invite.invite_code)}
                                color="error"
                                size="small"
                                title="Delete Link"
                                disabled={loading}
                              >
                                <DeleteIcon />
                              </IconButton>
                            </Box>
                          </Box>
                        </Paper>
                      </ListItem>
                    );
                  })}
                </List>
              )}
            </TabPanel>
          )}
        </DialogContent>

        <DialogActions>
          <Button onClick={onClose}>Close</Button>
        </DialogActions>
      </Dialog>

      <Snackbar
        open={snackbar.open}
        autoHideDuration={6000}
        onClose={() => setSnackbar({ ...snackbar, open: false })}
      >
        <Alert
          onClose={() => setSnackbar({ ...snackbar, open: false })}
          severity={snackbar.severity}
        >
          {snackbar.message}
        </Alert>
      </Snackbar>
    </>
  );
};

export default GroupManagementDialog;
