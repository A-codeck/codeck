import React, { useState } from 'react';
import { ThemeProvider } from '@mui/material/styles';
import { CssBaseline } from '@mui/material';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { theme } from './theme/theme';
import { AuthProvider, useAuth } from './contexts/AuthContext';
import LoginPage from './components/LoginPage';
import HomePage from './components/HomePage';
import InvitePage from './components/InvitePage';

const AppContent: React.FC = () => {
  const { isAuthenticated } = useAuth();
  const [isRegisterMode, setIsRegisterMode] = useState(false);

  const handleToggleMode = () => {
    setIsRegisterMode(!isRegisterMode);
  };

  return (
    <Routes>
      <Route 
        path="/join/:inviteCode" 
        element={<InvitePage />} 
      />
      <Route 
        path="/*" 
        element={
          isAuthenticated ? (
            <HomePage />
          ) : (
            <LoginPage 
              onToggleMode={handleToggleMode}
              isRegisterMode={isRegisterMode}
            />
          )
        } 
      />
    </Routes>
  );
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
