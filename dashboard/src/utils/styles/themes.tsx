import { createMuiTheme } from '@material-ui/core';

export const MuiDataTableTheme = () =>
  createMuiTheme({
    overrides: {
      MuiTableHead: {
        root: {
          fontWeight: 'bold',
          fontSize: '14px',
          '&:nth-child(3)': {
            paddingLeft: 50,
            textAlign: 'center',
          },
        },
      },
      MuiButton: {
        root: {
          textTransform: 'none',
          text: 'Open Sans, Rajdhani',
        },
      },
    },
    typography: { fontFamily: 'Open Sans, Rajdhani' },
    palette: {
      type: 'light',
      primary: { 500: '#000080' },
      secondary: { A400: '#000080' }, // '#000080' },
    },

    // #0A2053
  });
