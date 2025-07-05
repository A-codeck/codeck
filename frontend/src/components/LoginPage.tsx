import React, { useState } from 'react';
import {
  Box,
  Card,
  CardContent,
  TextField,
  Button,
  Typography,
  Alert,
  Link,
  Container,
  Divider,
} from '@mui/material';
import { useAuth } from '../contexts/AuthContext';
import { apiService } from '../services/api';

interface LoginPageProps {
  onToggleMode: () => void;
  isRegisterMode: boolean;
}

const LoginPage: React.FC<LoginPageProps> = ({ onToggleMode, isRegisterMode }) => {
  const { login } = useAuth();
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    password: '',
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      if (isRegisterMode) {
        // Register new user
        await apiService.register({
          name: formData.name,
          email: formData.email,
          password: formData.password,
        });
        
        // Small delay to ensure database transaction is complete
        await new Promise(resolve => setTimeout(resolve, 100));
        
        // After successful registration, login
        const loginResponse = await apiService.login({
          email: formData.email,
          password: formData.password,
        });
        
        login(loginResponse.user, loginResponse.token);
      } else {
        // Login existing user
        const loginResponse = await apiService.login({
          email: formData.email,
          password: formData.password,
        });
        
        login(loginResponse.user, loginResponse.token);
      }
    } catch (err: any) {
      setError(err.response?.data?.message || err.response?.data || err.message || 'An error occurred');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box
      sx={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        bgcolor: 'background.default',
      }}
    >
      <Container maxWidth="sm">
        <Card sx={{ p: 3 }}>
          <CardContent>
            <Typography variant="h1" align="center" gutterBottom>
              CODECK
            </Typography>
            <Typography variant="body1" align="center" color="text.secondary" sx={{ mb: 4 }}>
              COmpetição DE Código do Kim
            </Typography>

            <Typography variant="h2" align="center" gutterBottom>
              {isRegisterMode ? 'Create Account' : 'Sign In'}
            </Typography>

            {error && (
              <Alert severity="error" sx={{ mb: 2 }}>
                {error}
              </Alert>
            )}

            <Box component="form" onSubmit={handleSubmit} sx={{ mt: 2 }}>
              {isRegisterMode && (
                <TextField
                  fullWidth
                  name="name"
                  label="Full Name"
                  value={formData.name}
                  onChange={handleChange}
                  required
                  sx={{ mb: 2 }}
                />
              )}
              
              <TextField
                fullWidth
                name="email"
                label="Email"
                type="email"
                value={formData.email}
                onChange={handleChange}
                required
                sx={{ mb: 2 }}
              />
              
              <TextField
                fullWidth
                name="password"
                label="Password"
                type="password"
                value={formData.password}
                onChange={handleChange}
                required
                sx={{ mb: 3 }}
              />

              <Button
                type="submit"
                fullWidth
                variant="contained"
                size="large"
                disabled={loading}
                sx={{ mb: 2 }}
              >
                {loading 
                  ? (isRegisterMode ? 'Creating Account...' : 'Signing In...') 
                  : (isRegisterMode ? 'Create Account' : 'Sign In')
                }
              </Button>

              <Divider sx={{ my: 2 }} />

              <Box textAlign="center">
                <Typography variant="body2" color="text.secondary">
                  {isRegisterMode ? 'Already have an account?' : "Don't have an account?"}
                  {' '}
                  <Link
                    component="button"
                    type="button"
                    onClick={onToggleMode}
                    sx={{ 
                      color: 'secondary.main',
                      textDecoration: 'none',
                      '&:hover': {
                        textDecoration: 'underline',
                      },
                    }}
                  >
                    {isRegisterMode ? 'Sign In' : 'Create Account'}
                  </Link>
                </Typography>
              </Box>
            </Box>
          </CardContent>
        </Card>
      </Container>
    </Box>
  );
};

export default LoginPage;
