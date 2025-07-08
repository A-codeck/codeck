import React, { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Box,
  Paper,
  Typography,
  Button,
  CircularProgress,
  Alert,
  Container,
} from '@mui/material';
import {
  Group as GroupIcon,
  Check as CheckIcon,
} from '@mui/icons-material';
import { useAuth } from '../contexts/AuthContext';
import { apiService } from '../services/api';
import LoginPage from './LoginPage';

const InvitePage: React.FC = () => {
  const { inviteCode } = useParams<{ inviteCode: string }>();
  const navigate = useNavigate();
  const { isAuthenticated, user } = useAuth();
  const [joining, setJoining] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);
  const [isRegisterMode, setIsRegisterMode] = useState(false);

  const handleJoinGroup = async () => {
    if (!user || !inviteCode) return;
    
    setJoining(true);
    setError('');
    
    try {
      await apiService.joinGroupByInvite(inviteCode, user.id);
      setSuccess(true);
      
      // Redirect to home page after a short delay
      setTimeout(() => {
        navigate('/');
      }, 2000);
    } catch (error: any) {
      setError(
        error.response?.data?.error || 
        'Failed to join group. The invite link may be invalid or expired.'
      );
    } finally {
      setJoining(false);
    }
  };

  const handleGoHome = () => {
    navigate('/');
  };

  const handleToggleLoginMode = () => {
    setIsRegisterMode(!isRegisterMode);
  };

  if (!isAuthenticated) {
    return (
      <Container maxWidth="sm" sx={{ mt: 4 }}>
        <Paper sx={{ p: 4, textAlign: 'center', mb: 3 }}>
          <GroupIcon sx={{ fontSize: 60, color: 'primary.main', mb: 2 }} />
          <Typography variant="h4" gutterBottom>
            Join Group Invitation
          </Typography>
          <Typography variant="body1" color="text.secondary" sx={{ mb: 3 }}>
            You need to be logged in to join a group. Please log in or create an account below.
          </Typography>
        </Paper>
        
        <LoginPage 
          onToggleMode={handleToggleLoginMode}
          isRegisterMode={isRegisterMode}
        />
      </Container>
    );
  }

  if (success) {
    return (
      <Container maxWidth="sm" sx={{ mt: 8 }}>
        <Paper sx={{ p: 6, textAlign: 'center' }}>
          <CheckIcon sx={{ fontSize: 80, color: 'success.main', mb: 3 }} />
          <Typography variant="h4" gutterBottom color="success.main">
            Successfully Joined!
          </Typography>
          <Typography variant="body1" color="text.secondary" sx={{ mb: 4 }}>
            You have successfully joined the group. Redirecting to your dashboard...
          </Typography>
          <CircularProgress />
        </Paper>
      </Container>
    );
  }

  return (
    <Container maxWidth="sm" sx={{ mt: 8 }}>
      <Paper sx={{ p: 6, textAlign: 'center' }}>
        <GroupIcon sx={{ fontSize: 80, color: 'primary.main', mb: 3 }} />
        <Typography variant="h4" gutterBottom>
          Join Group
        </Typography>
        <Typography variant="body1" color="text.secondary" sx={{ mb: 4 }}>
          You've been invited to join a group! Click the button below to join.
        </Typography>

        {error && (
          <Alert severity="error" sx={{ mb: 3 }}>
            {error}
          </Alert>
        )}

        <Box sx={{ display: 'flex', gap: 2, justifyContent: 'center', flexWrap: 'wrap' }}>
          <Button
            variant="contained"
            size="large"
            onClick={handleJoinGroup}
            disabled={joining}
            startIcon={joining ? <CircularProgress size={20} /> : <GroupIcon />}
          >
            {joining ? 'Joining...' : 'Join Group'}
          </Button>
          <Button
            variant="outlined"
            size="large"
            onClick={handleGoHome}
            disabled={joining}
          >
            Go to Dashboard
          </Button>
        </Box>

        <Typography variant="caption" color="text.secondary" sx={{ mt: 3, display: 'block' }}>
          Invite Code: {inviteCode}
        </Typography>
      </Paper>
    </Container>
  );
};

export default InvitePage;
