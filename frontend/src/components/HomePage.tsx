import React, { useState } from 'react';
import {
  Box,
  AppBar,
  Toolbar,
  Typography,
  IconButton,
  Menu,
  MenuItem,
  Avatar,
} from '@mui/material';
import { useAuth } from '../contexts/AuthContext';
import GroupsSidebar from './GroupsSidebar';
import ActivityFeed from './ActivityFeed';
import GroupRanking from './GroupRanking';

const HomePage: React.FC = () => {
  const { user, logout } = useAuth();
  const [selectedGroupId, setSelectedGroupId] = useState<string | undefined>();
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

  const handleMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const handleLogout = () => {
    logout();
    handleMenuClose();
  };

  const handleGroupSelect = (groupId: string | undefined) => {
    setSelectedGroupId(groupId);
  };

  const handleGroupsChange = () => {
    // This can be used to refresh data when groups change
    // For now, we'll keep it simple
  };

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', height: '100vh' }}>
      {/* Header */}
      <AppBar position="static" sx={{ bgcolor: 'primary.main' }}>
        <Toolbar>
          <Typography variant="h2" component="div" sx={{ flexGrow: 1 }}>
            CODECK
          </Typography>
          
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
            <Typography variant="body1">
              Welcome, {user?.name}
            </Typography>
            
            <IconButton
              size="large"
              edge="end"
              aria-label="account menu"
              onClick={handleMenuOpen}
              color="inherit"
            >
              <Avatar sx={{ width: 32, height: 32 }}>
                {user?.name?.charAt(0)}
              </Avatar>
            </IconButton>
          </Box>

          <Menu
            anchorEl={anchorEl}
            open={Boolean(anchorEl)}
            onClose={handleMenuClose}
            anchorOrigin={{
              vertical: 'bottom',
              horizontal: 'right',
            }}
            transformOrigin={{
              vertical: 'top',
              horizontal: 'right',
            }}
          >
            <MenuItem onClick={handleLogout}>
              Logout
            </MenuItem>
          </Menu>
        </Toolbar>
      </AppBar>

      {/* Main Content */}
      <Box sx={{ display: 'flex', flex: 1, overflow: 'hidden' }}>
        {/* Left Sidebar - Groups */}
        <GroupsSidebar
          selectedGroupId={selectedGroupId}
          onGroupSelect={handleGroupSelect}
          onGroupsChange={handleGroupsChange}
        />

        {/* Center - Activity Feed */}
        <ActivityFeed selectedGroupId={selectedGroupId} />

        {/* Right Sidebar - Group Ranking (only show when group is selected) */}
        {selectedGroupId && (
          <GroupRanking groupId={selectedGroupId} />
        )}
      </Box>
    </Box>
  );
};

export default HomePage;
