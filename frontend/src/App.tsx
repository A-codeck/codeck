import React, { useState } from 'react';
import { ThemeProvider } from '@mui/material/styles';
import { CssBaseline } from '@mui/material';
import { BrowserRouter as Router } from 'react-router-dom';
import { theme } from './theme/theme';
import { AuthProvider, useAuth } from './contexts/AuthContext';
import LoginPage from './components/LoginPage';
import HomePage from './components/HomePage';

const AppContent: React.FC = () => {
  const { isAuthenticated } = useAuth();
  const [isRegisterMode, setIsRegisterMode] = useState(false);

  const handleToggleMode = () => {
    setIsRegisterMode(!isRegisterMode);
  };

  if (!isAuthenticated) {
    return (
      <LoginPage 
        onToggleMode={handleToggleMode}
        isRegisterMode={isRegisterMode}
      />
    );
  }

  return <HomePage />;
};

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Router>
        <AuthProvider>
          <AppContent />
        </AuthProvider>
      </Router>
    </ThemeProvider>
  );
}

export default App;
