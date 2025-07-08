import React, { useState } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  Box,
  Typography,
  Alert,
  CircularProgress,
  InputAdornment,
} from '@mui/material';
import {
  Link as LinkIcon,
  Group as GroupIcon,
} from '@mui/icons-material';
import { apiService } from '../services/api';
import { useAuth } from '../contexts/AuthContext';

interface JoinGroupDialogProps {
  open: boolean;
  onClose: () => void;
  onGroupJoined: () => void;
}

const JoinGroupDialog: React.FC<JoinGroupDialogProps> = ({
  open,
  onClose,
  onGroupJoined,
}) => {
  const { user } = useAuth();
  const [inviteLink, setInviteLink] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const extractInviteCode = (link: string): string => {
    // Extract invite code from various possible formats:
    // - https://example.com/join/ABC123
    // - /join/ABC123
    // - ABC123
    
    const trimmed = link.trim();
    
    // If it's a full URL
    if (trimmed.includes('/join/')) {
      const parts = trimmed.split('/join/');
      return parts[1] || '';
    }
    
    // If it's just the code
    return trimmed;
  };

  const handleJoinGroup = async () => {
    if (!user || !inviteLink.trim()) return;
    
    const inviteCode = extractInviteCode(inviteLink);
    
    if (!inviteCode) {
      setError('Please enter a valid invite link or code');
      return;
    }
    
    setLoading(true);
    setError('');
    
    try {
      await apiService.joinGroupByInvite(inviteCode, user.id);
      onGroupJoined();
      onClose();
      setInviteLink('');
    } catch (error: any) {
      setError(
        error.response?.data?.error || 
        'Failed to join group. Please check the invite link and try again.'
      );
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    setInviteLink('');
    setError('');
    onClose();
  };

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle>
        <Box display="flex" alignItems="center" gap={1}>
          <GroupIcon color="primary" />
          <Typography variant="h6">Join a Group</Typography>
        </Box>
      </DialogTitle>

      <DialogContent>
        <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
          Enter an invite link or invite code to join a group. You can get these from group owners.
        </Typography>

        <TextField
          fullWidth
          label="Invite Link or Code"
          placeholder="https://example.com/join/ABC123 or ABC123"
          value={inviteLink}
          onChange={(e) => {
            setInviteLink(e.target.value);
            if (error) setError('');
          }}
          InputProps={{
            startAdornment: (
              <InputAdornment position="start">
                <LinkIcon />
              </InputAdornment>
            ),
          }}
          disabled={loading}
          error={!!error}
          helperText={error}
        />

        {error && (
          <Alert severity="error" sx={{ mt: 2 }}>
            {error}
          </Alert>
        )}
      </DialogContent>

      <DialogActions>
        <Button onClick={handleClose} disabled={loading}>
          Cancel
        </Button>
        <Button
          variant="contained"
          onClick={handleJoinGroup}
          disabled={!inviteLink.trim() || loading}
          startIcon={loading ? <CircularProgress size={20} /> : <GroupIcon />}
        >
          {loading ? 'Joining...' : 'Join Group'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default JoinGroupDialog;
