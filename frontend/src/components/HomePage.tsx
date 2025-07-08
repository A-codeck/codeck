import React, { useState, useRef } from 'react';
import {
  Box,
  AppBar,
  Toolbar,
  Typography,
  IconButton,
  Menu,
  MenuItem,
  Avatar,
  Paper,
  Button,
  CircularProgress,
} from '@mui/material';
import { Add as AddIcon, Group as GroupIcon } from '@mui/icons-material';
import { Group } from '../types/api';
import { useAuth } from '../contexts/AuthContext';
import GroupsSidebar, { GroupsSidebarRef } from './GroupsSidebar';
import ActivityFeed from './ActivityFeed';
import GroupRanking from './GroupRanking';
import UserDashboard from './UserDashboard';
import JoinGroupDialog from './JoinGroupDialog';

const HomePage: React.FC = () => {
  const { user, logout } = useAuth();
  const [selectedGroupId, setSelectedGroupId] = useState<string | undefined>();
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [userGroups, setUserGroups] = useState<Group[]>([]);
  const [groupsLoaded, setGroupsLoaded] = useState(false);
  const [joinGroupDialogOpen, setJoinGroupDialogOpen] = useState(false);
  const groupsSidebarRef = useRef<GroupsSidebarRef>(null);

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

  const handleGroupsLoaded = (groups: Group[]) => {
    setUserGroups(groups);
    setGroupsLoaded(true);
  };

  const handleGroupJoined = () => {
    // Refresh the groups list
    window.location.reload(); // Simple approach, could be improved with state management
  };

  const renderNoGroupsMessage = () => (
    <Box sx={{ flex: 1, display: 'flex', alignItems: 'center', justifyContent: 'center', p: 3 }}>
      <Paper sx={{ p: 6, textAlign: 'center', maxWidth: 600 }}>
        <GroupIcon sx={{ fontSize: 80, color: 'text.secondary', mb: 3 }} />
        <Typography variant="h3" gutterBottom>
          Bem-vindo ao CODECK!
        </Typography>
        <Typography variant="h4" color="text.secondary" gutterBottom>
          Você ainda não faz parte de nenhum grupo
        </Typography>
        <Typography variant="body1" color="text.secondary" sx={{ mb: 4, lineHeight: 1.6 }}>
          Para começar a compartilhar suas atividades de programação, você precisa criar seu primeiro grupo ou entrar em um já existente. 
          Os grupos ajudam você a se conectar com outros programadores e acompanhar seu progresso juntos.
        </Typography>
        <Box sx={{ display: 'flex', gap: 2, justifyContent: 'center', flexWrap: 'wrap' }}>
          <Button
            variant="contained"
            color="secondary"
            startIcon={<AddIcon />}
            size="large"
            onClick={() => {
              groupsSidebarRef.current?.openCreateDialog();
            }}
          >
            Criar Seu Primeiro Grupo
          </Button>
          <Button
            variant="outlined"
            color="secondary"
            size="large"
            onClick={() => setJoinGroupDialogOpen(true)}
          >
            Entrar em um Grupo
          </Button>
        </Box>
        <Typography variant="caption" color="text.secondary" sx={{ mt: 3, display: 'block' }}>
          Assim que você fizer parte de um grupo, poderá registrar atividades e ver no que outros estão trabalhando!
        </Typography>
      </Paper>
    </Box>
  );

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', height: '100vh' }}>
      {/* Header */}
      <AppBar position="static" elevation={0} sx={{ backgroundColor: (theme) => theme.palette.background.default }}>
        <Toolbar>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, flexGrow: 1 }}>
            {/* Logo Placeholder */}
            <Box
              component="img"
              src="/logo4x.png"
              alt="CODECK Logo"
              sx={{
                width: 42,
                height: 50,
                borderRadius: 2,
                objectFit: 'cover',
                mr: -1,
              }}
            />
            <Typography variant="h2" color="#3DC959" component="div">
              CODECK
            </Typography>
          </Box>
          
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
        <Typography variant="body1">
          Bem-vindo, {user?.name}
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
              Sair
            </MenuItem>
          </Menu>
        </Toolbar>
      </AppBar>

      {/* Main Content */}
      <Box sx={{ display: 'flex', flex: 1, overflow: 'hidden' }}>
        {/* Always show sidebar for group management */}
        <GroupsSidebar
          ref={groupsSidebarRef}
          selectedGroupId={selectedGroupId}
          onGroupSelect={handleGroupSelect}
          onGroupsChange={handleGroupsChange}
          onGroupsLoaded={handleGroupsLoaded}
        />

        {/* Conditional Content Based on Groups */}
        {(() => {
          console.log('HomePage render: groupsLoaded =', groupsLoaded, 'userGroups.length =', userGroups.length);
          return groupsLoaded && userGroups.length === 0 ? (
            renderNoGroupsMessage()
          ) : groupsLoaded ? (
            <>
              {/* Center - Activity Feed */}
              <ActivityFeed selectedGroupId={selectedGroupId} />

              {/* Right Sidebar - Conditional based on group selection */}
              {selectedGroupId ? (
                <GroupRanking groupId={selectedGroupId} />
              ) : (
                <UserDashboard />
              )}
            </>
          ) : (
            // Show loading while groups are being loaded
            <Box sx={{ flex: 1, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
              <CircularProgress />
            </Box>
          );
        })()}
      </Box>

      {/* Join Group Dialog */}
      <JoinGroupDialog
        open={joinGroupDialogOpen}
        onClose={() => setJoinGroupDialogOpen(false)}
        onGroupJoined={handleGroupJoined}
      />
    </Box>
  );
};

export default HomePage;
