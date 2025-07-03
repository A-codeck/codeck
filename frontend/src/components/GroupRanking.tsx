import React, { useState, useEffect } from 'react';
import {
  Paper,
  Typography,
  Box,
  List,
  ListItem,
  ListItemAvatar,
  ListItemText,
  Avatar,
  Chip,
  CircularProgress,
} from '@mui/material';
import { EmojiEvents as TrophyIcon, WorkspacePremium as MedalIcon } from '@mui/icons-material';
import { UserStats, Group } from '../types/api';
import { useAuth } from '../contexts/AuthContext';
import { apiService } from '../services/api';

interface GroupRankingProps {
  groupId: string;
}

const GroupRanking: React.FC<GroupRankingProps> = ({ groupId }) => {
  const { user } = useAuth();
  const [rankings, setRankings] = useState<UserStats[]>([]);
  const [group, setGroup] = useState<Group | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadGroupRanking();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [groupId, user]);

  const loadGroupRanking = async () => {
    if (!user) return;

    try {
      setLoading(true);
      
      // Load group info
      const groupData = await apiService.getGroup(groupId, user.id);
      setGroup(groupData);

      // Load group members
      const membersData = await apiService.getGroupMembers(groupId, user.id);
      
      // Calculate rankings by getting activities for each member
      const userStats: UserStats[] = [];
      
      for (const member of membersData.members) {
        try {
          const memberActivities = await apiService.getUserActivities(member.user_id);
          const memberInfo = await apiService.getUser(member.user_id);
          
          userStats.push({
            user_id: member.user_id,
            user_name: member.nickname || memberInfo.name,
            activity_count: memberActivities.length,
          });
        } catch (error) {
          console.error(`Error loading data for user ${member.user_id}:`, error);
        }
      }

      // Sort by activity count (descending)
      userStats.sort((a, b) => b.activity_count - a.activity_count);
      
      setRankings(userStats);
    } catch (error) {
      console.error('Error loading group ranking:', error);
    } finally {
      setLoading(false);
    }
  };

  const getRankIcon = (position: number) => {
    switch (position) {
      case 1:
        return <TrophyIcon sx={{ color: '#FFD700' }} />; // Gold
      case 2:
        return <MedalIcon sx={{ color: '#C0C0C0' }} />; // Silver
      case 3:
        return <MedalIcon sx={{ color: '#CD7F32' }} />; // Bronze
      default:
        return null;
    }
  };

  const getRankColor = (position: number) => {
    switch (position) {
      case 1:
        return '#FFD700';
      case 2:
        return '#C0C0C0';
      case 3:
        return '#CD7F32';
      default:
        return 'text.secondary';
    }
  };

  if (loading) {
    return (
      <Paper
        sx={{
          width: 280,
          height: '100%',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          borderLeft: '1px solid',
          borderColor: 'divider',
          borderRadius: 0,
        }}
      >
        <CircularProgress />
      </Paper>
    );
  }

  return (
    <Paper
      sx={{
        width: 280,
        height: '100%',
        display: 'flex',
        flexDirection: 'column',
        borderLeft: '1px solid',
        borderColor: 'divider',
        borderRadius: 0,
      }}
    >
      <Box sx={{ p: 2, borderBottom: '1px solid', borderColor: 'divider' }}>
        <Typography variant="h3" gutterBottom>
          Group Ranking
        </Typography>
        {group && (
          <Typography variant="body2" color="text.secondary">
            {group.name}
          </Typography>
        )}
      </Box>

      <Box sx={{ flex: 1, overflow: 'auto' }}>
        {rankings.length > 0 ? (
          <List>
            {rankings.map((userStat, index) => {
              const position = index + 1;
              const isCurrentUser = userStat.user_id === user?.id;
              
              return (
                <ListItem
                  key={userStat.user_id}
                  sx={{
                    bgcolor: isCurrentUser ? 'action.selected' : 'transparent',
                  }}
                >
                  <ListItemAvatar>
                    <Box sx={{ position: 'relative' }}>
                      <Avatar sx={{ width: 40, height: 40 }}>
                        {userStat.user_name.charAt(0)}
                      </Avatar>
                      {position <= 3 && (
                        <Box
                          sx={{
                            position: 'absolute',
                            top: -8,
                            right: -8,
                            width: 20,
                            height: 20,
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'center',
                          }}
                        >
                          {getRankIcon(position)}
                        </Box>
                      )}
                    </Box>
                  </ListItemAvatar>
                  
                  <ListItemText
                    primary={
                      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                        <Typography
                          variant="body2"
                          fontWeight={isCurrentUser ? 600 : 400}
                          sx={{ color: getRankColor(position) }}
                        >
                          #{position}
                        </Typography>
                        <Typography
                          variant="body2"
                          fontWeight={isCurrentUser ? 600 : 400}
                          noWrap
                        >
                          {userStat.user_name}
                        </Typography>
                        {isCurrentUser && (
                          <Chip 
                            label="You" 
                            size="small" 
                            color="secondary" 
                            variant="outlined"
                          />
                        )}
                      </Box>
                    }
                    secondary={
                      <Typography variant="caption" color="text.secondary">
                        {userStat.activity_count} activities
                      </Typography>
                    }
                  />
                </ListItem>
              );
            })}
          </List>
        ) : (
          <Box sx={{ p: 3, textAlign: 'center' }}>
            <Typography variant="body2" color="text.secondary">
              No members found
            </Typography>
          </Box>
        )}
      </Box>
    </Paper>
  );
};

export default GroupRanking;
