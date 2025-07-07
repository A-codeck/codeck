import React, { useState, useEffect } from 'react';
import {
  Box,
  Paper,
  Typography,
  Grid,
  Button,
  Divider,
  Chip,
  Card,
  CardContent,
  Link,
  IconButton,
} from '@mui/material';
import {
  ChevronLeft,
  ChevronRight,
  Launch as LaunchIcon,
} from '@mui/icons-material';
import { useAuth } from '../contexts/AuthContext';
import { apiService } from '../services/api';

interface UserDashboardProps {}

interface ActivityDate {
  date: string;
  count: number;
}

const UserDashboard: React.FC<UserDashboardProps> = () => {
  const { user } = useAuth();
  const [currentDate, setCurrentDate] = useState(new Date());
  const [activityDates, setActivityDates] = useState<ActivityDate[]>([]);
  const [loading, setLoading] = useState(true);

  const problemSites = [
    {
      name: 'LeetCode',
      description: 'Popular platform with coding interview questions and algorithm challenges. Great for preparing for technical interviews.',
      url: 'https://leetcode.com',
      difficulty: 'Beginner to Advanced',
      focus: 'Interview Prep, Algorithms'
    },
    {
      name: 'HackerRank',
      description: 'Comprehensive coding challenges covering algorithms, data structures, and domain-specific problems.',
      url: 'https://hackerrank.com',
      difficulty: 'Beginner to Advanced', 
      focus: 'Algorithms, Data Structures'
    },
    {
      name: 'Codeforces',
      description: 'Competitive programming platform with regular contests and a vast archive of algorithmic problems.',
      url: 'https://codeforces.com',
      difficulty: 'Intermediate to Expert',
      focus: 'Competitive Programming'
    },
    {
      name: 'AtCoder',
      description: 'Japanese competitive programming platform known for high-quality problems and beginner-friendly contests.',
      url: 'https://atcoder.jp',
      difficulty: 'Beginner to Expert',
      focus: 'Competitive Programming'
    },
    {
      name: 'CodeChef',
      description: 'Global competitive programming platform with monthly contests and practice problems.',
      url: 'https://codechef.com',
      difficulty: 'Beginner to Advanced',
      focus: 'Competitive Programming'
    }
  ];

  useEffect(() => {
    loadUserActivityDates();
  }, [user, currentDate]);

  const loadUserActivityDates = async () => {
    if (!user) return;

    try {
      setLoading(true);
      // Get the first and last day of the current month
      const firstDay = new Date(currentDate.getFullYear(), currentDate.getMonth(), 1);
      const lastDay = new Date(currentDate.getFullYear(), currentDate.getMonth() + 1, 0);
      
      // Get user's activities for the current month
      const activities = await apiService.getUserFeedWithGroups(user.id);
      
      // Filter activities for current month and count by date
      const activityMap = new Map<string, number>();
      
      activities.forEach(activity => {
        const activityDate = new Date(activity.date);
        if (activityDate >= firstDay && activityDate <= lastDay) {
          const dateKey = activityDate.toISOString().split('T')[0];
          activityMap.set(dateKey, (activityMap.get(dateKey) || 0) + 1);
        }
      });

      const activityDatesArray: ActivityDate[] = Array.from(activityMap.entries()).map(([date, count]) => ({
        date,
        count
      }));

      setActivityDates(activityDatesArray);
    } catch (error) {
      console.error('Error loading activity dates:', error);
    } finally {
      setLoading(false);
    }
  };

  const getDaysInMonth = (date: Date) => {
    return new Date(date.getFullYear(), date.getMonth() + 1, 0).getDate();
  };

  const getFirstDayOfMonth = (date: Date) => {
    return new Date(date.getFullYear(), date.getMonth(), 1).getDay();
  };

  const isActivityDay = (day: number) => {
    const dateStr = new Date(currentDate.getFullYear(), currentDate.getMonth(), day)
      .toISOString().split('T')[0];
    return activityDates.some(activity => activity.date === dateStr);
  };

  const getActivityCount = (day: number) => {
    const dateStr = new Date(currentDate.getFullYear(), currentDate.getMonth(), day)
      .toISOString().split('T')[0];
    const activity = activityDates.find(activity => activity.date === dateStr);
    return activity ? activity.count : 0;
  };

  const navigateMonth = (direction: 'prev' | 'next') => {
    setCurrentDate(prev => {
      const newDate = new Date(prev);
      if (direction === 'prev') {
        newDate.setMonth(newDate.getMonth() - 1);
      } else {
        newDate.setMonth(newDate.getMonth() + 1);
      }
      return newDate;
    });
  };

  const renderCalendar = () => {
    const daysInMonth = getDaysInMonth(currentDate);
    const firstDay = getFirstDayOfMonth(currentDate);
    const days = [];
    const dayNames = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];

    // Add day names header
    dayNames.forEach(dayName => {
      days.push(
        <Box key={dayName} sx={{ 
          p: 0.5, 
          textAlign: 'center', 
          fontWeight: 'bold', 
          fontSize: '0.75rem',
          minHeight: 24,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center'
        }}>
          {dayName}
        </Box>
      );
    });

    // Add empty cells for days before the first day of the month
    for (let i = 0; i < firstDay; i++) {
      days.push(<Box key={`empty-${i}`} sx={{ p: 1 }} />);
    }

    // Add days of the month
    for (let day = 1; day <= daysInMonth; day++) {
      const hasActivity = isActivityDay(day);
      const activityCount = getActivityCount(day);
      const isToday = new Date().toDateString() === new Date(currentDate.getFullYear(), currentDate.getMonth(), day).toDateString();

      days.push(
        <Box
          key={day}
          sx={{
            p: 0.5,
            textAlign: 'center',
            borderRadius: 1,
            cursor: 'pointer',
            backgroundColor: hasActivity ? 'success.main' : 'transparent',
            color: hasActivity ? 'success.contrastText' : 'text.primary',
            border: isToday ? 2 : 0,
            borderColor: 'primary.main',
            '&:hover': {
              backgroundColor: hasActivity ? 'success.dark' : 'action.hover',
            },
            fontSize: '0.75rem',
            minHeight: 28,
            maxWidth: '100%',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            position: 'relative',
            aspectRatio: '1',
          }}
          title={hasActivity ? `${activityCount} activity(s)` : 'No activities'}
        >
          {day}
          {hasActivity && activityCount > 1 && (
            <Box
              sx={{
                position: 'absolute',
                top: 2,
                right: 2,
                width: 6,
                height: 6,
                borderRadius: '50%',
                backgroundColor: 'warning.main',
              }}
            />
          )}
        </Box>
      );
    }

    return (
      <Box sx={{ 
        display: 'grid', 
        gridTemplateColumns: 'repeat(7, 1fr)', 
        gap: 0.5,
        width: '100%',
        maxWidth: '100%'
      }}>
        {days}
      </Box>
    );
  };

  return (
    <Box 
      sx={{ 
        width: 350, 
        height: '100%',
        overflow: 'auto',
        p: 3 
      }}
    >
      {/* Activity Calendar */}
      <Paper sx={{ p: 3, mb: 3 }}>
        <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'space-between', mb: 0 }}>
          <Typography variant="h3" component="h2" sx={{mb: 2}}>
            Activity Streak
          </Typography>
          <Box sx={{ display: 'flex', alignItems: 'center' }}>
            <IconButton onClick={() => navigateMonth('prev')} size="small">
              <ChevronLeft />
            </IconButton>
            <Typography variant="body1" sx={{ minWidth: 120, textAlign: 'center' }}>
              {currentDate.toLocaleDateString('en-US', { month: 'long', year: 'numeric' })}
            </Typography>
            <IconButton onClick={() => navigateMonth('next')} size="small">
              <ChevronRight />
            </IconButton>
          </Box>
        </Box>
        
        {renderCalendar()}
        
        <Box sx={{ mt: 2, display: 'flex', alignItems: 'center', gap: 1, justifyContent: 'center' }}>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
            <Box sx={{ width: 12, height: 12, backgroundColor: 'success.main', borderRadius: 1 }} />
            <Typography variant="caption">Activity</Typography>
          </Box>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
            <Box sx={{ width: 12, height: 12, backgroundColor: 'action.hover', borderRadius: 1 }} />
            <Typography variant="caption">No Activity</Typography>
          </Box>
        </Box>
      </Paper>

      <Divider sx={{ mb: 3 }} />

      {/* Programming Problem Sites */}
      <Paper sx={{ p: 3 }}>
        <Typography variant="h3" component="h2" gutterBottom>
          Problem Solving Sites
        </Typography>
        <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
          Discover new programming challenges to improve your skills
        </Typography>

        {problemSites.map((site, index) => (
          <Card key={site.name} sx={{ mb: 2, '&:last-child': { mb: 0 } }}>
            <CardContent sx={{ p: 2, '&:last-child': { pb: 2 } }}>
              <Box sx={{ display: 'flex', alignItems: 'flex-start', justifyContent: 'space-between', mb: 1 }}>
                <Typography variant="h4" component="h3">
                  {site.name}
                </Typography>
                <IconButton
                  component={Link}
                  href={site.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  size="small"
                  sx={{ p: 0.5 }}
                >
                  <LaunchIcon fontSize="small" />
                </IconButton>
              </Box>
              
              <Typography variant="body2" color="text.secondary" sx={{ mb: 1.5, lineHeight: 1.4 }}>
                {site.description}
              </Typography>
              
              <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap' }}>
                <Chip 
                  label={site.difficulty} 
                  size="small" 
                  variant="outlined"
                  color="primary"
                />
                <Chip 
                  label={site.focus} 
                  size="small" 
                  variant="outlined"
                  color="secondary"
                />
              </Box>
            </CardContent>
          </Card>
        ))}
      </Paper>
    </Box>
  );
};

export default UserDashboard;
