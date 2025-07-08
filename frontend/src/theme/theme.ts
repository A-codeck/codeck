import { createTheme } from '@mui/material/styles';

export const theme = createTheme({
  palette: {
    mode: 'dark',
    primary: {
      main: '#057A55',
      contrastText: '#ffffff',
    },
    secondary: {
      main: '#057A55',
      contrastText: '#ffffff',
    },
    background: {
      default: '#111928',
      paper: '#111928',
    },
    text: {
      primary: '#ffffff',
      secondary: '#a0a0a0',
    },
    divider: '#4B5563',
  },
  typography: {
    fontFamily: '"Inter", "Roboto", "Helvetica", "Arial", sans-serif',
    h1: {
      fontSize: '2rem',
      fontWeight: 600,
    },
    h2: {
      fontSize: '1.5rem',
      fontWeight: 600,
    },
    h3: {
      fontSize: '1.25rem',
      fontWeight: 600,
    },
    body1: {
      fontSize: '0.875rem',
    },
    body2: {
      fontSize: '0.75rem',
    },
  },
  components: {
    MuiPaper: {
      defaultProps: {
        elevation: 0,
      },
      styleOverrides: {
        root: {
          border: '1px solid #4B5563',
          boxShadow: 'none',
          backgroundColor: '#111928',
        },
      },
    },
    MuiButton: {
      styleOverrides: {
        root: {
          textTransform: 'none',
          borderRadius: '6px',
          fontWeight: 500,
        },
        contained: {
          boxShadow: 'none',
          '&:hover': {
            boxShadow: 'none',
          },
        },
      },
    },
    MuiTextField: {
      styleOverrides: {
        root: {
          '& .MuiOutlinedInput-root': {
            '& fieldset': {
              borderColor: '#4D545D',
            },
            '&:hover fieldset': {
              borderColor: '#057A55',
            },
            '&.Mui-focused fieldset': {
              borderColor: '#057A55',
            },
          },
        },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: {
          border: '1px solid #374151',
          borderRadius: '8px',
          boxShadow: 'none',
          backgroundColor: '#111928',
        },
      },
    },
    MuiDialog: {
      styleOverrides: {
        paper: {
          border: '1px solid #374151',
          borderRadius: '8px',
          boxShadow: 'none',
          backgroundColor: '#111928',
        },
      },
    },
  },
});
