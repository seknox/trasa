import CssBaseline from '@material-ui/core/CssBaseline';
import { createMuiTheme, ThemeProvider } from '@material-ui/core/styles';
import React from 'react';
import 'semantic-ui-css/semantic.min.css';

let theme = createMuiTheme({
  typography: {
    // useNextVariants: true,
    fontFamily: 'Open Sans, Rajdhani',
    h1: {
      fontWeight: 700,
      fontSize: 32,
      letterSpacing: 0.5,
      padding: 10,
      color: '#0A2053',
    },
    h2: {
      fontWeight: 700,
      fontSize: 24,
      letterSpacing: 0.5,
      color: '#0A2053',
      padding: 10,
    },
    h3: {
      fontWeight: 700,
      padding: 10,
      fontSize: 18,
      letterSpacing: 0.5,
      color: '#18202c',
    },
    h4: {
      fontWeight: 600,
      fontSize: 16,
      letterSpacing: 0.5,
      color: '#404854',
      padding: 12,
    },
    h5: {
      fontWeight: 600,
      fontSize: 14,
      letterSpacing: 0.5,
      padding: 13,
    },
    h6: {
      fontWeight: 700,
      fontSize: 12,
      letterSpacing: 0.5,
      padding: 12,
    },
    subtitle1: {
      fontWeight: 600,
      fontSize: 16,
      letterSpacing: 0.5,
      color: '#404854',
    },

    button: {
      textTransform: 'none',
      fontFamily: 'Open Sans, Rajdhani',
    },
  },
  palette: {
    primary: {
      light: '#000080',
      main: '#000066',
      dark: '#03052b', // 'rgba(1,1,35,1)', // '#000080' #1b1b32 #0A2053 #000066}, #0A2053
      contrastText: '#fff',
    },
    secondary: { A400: '#000080' },
  },

  shape: {
    borderRadius: 4,
  },
});

theme = {
  ...theme,
  overrides: {
    MuiDrawer: {
      paper: {
        backgroundColor: 'white', // '#18202c',
      },
    },
    MuiAppBar: {
      root: {
        backgroundColor: '#03052b',
      },
    },

    MuiButton: {
      contained: {
        color: 'white',
        fontWeight: 600,
        fontSize: '14px',
        boxShadow: 'none',
        '&:active': {
          boxShadow: 'none',
          backgroundColor: '#000080',
        },
        '&:hover, &:focus': {
          backgroundColor: theme.palette.primary.light,
          boxShadow: '0 0 10px #030417',
        },
        backgroundColor: '#000080',
      },
      text: {
        color: 'blue',
        boxShadow: 'none',
        '&:active': {
          boxShadow: 'none',
        },
        label: {
          textTransform: 'initial',
          color: 'white',
          fontWeight: 600,
          fontSize: '14px',
          fontFamily: 'Open Sans, Rajdhani',
        },
      },
      // palette: {
      //   light: '#63ccff',
      //   main: '#000080', // Dark '#030417', //' Original - #009be5', //Navy - 000080
      //   dark: '#006db3',
      // },
    },

    MuiTextField: {
      root: {
        padding: 1,
        'label + &': {
          marginTop: theme.spacing(3),
        },
        fontSize: 16,

        transition: theme.transitions.create(['border-color', 'box-shadow']),
        '&:focus': {
          border: '#63ccff',
          borderColor: '#80bdff',
          boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
        },
        borderColor: '#63ccff',
      },
      // outlined: {
      //   border: '#63ccff',
      // },
    },
    MuiDialog: {},
    MuiDialogTitle: {
      root: {
        fontWeight: 700,
        padding: theme.spacing(2),
        fontSize: 18,
        letterSpacing: 0.5,
        color: '#18202c',
      },
    },
    MuiIconButton: {
      root: {
        padding: theme.spacing(1),
      },
    },
    MuiTooltip: {
      tooltip: {
        borderRadius: 4,
        fontSize: 14,
      },
    },
    MuiDivider: {
      root: {
        backgroundColor: '#404854',
      },
    },
    MuiListItemText: {
      primary: {
        fontWeight: theme.typography.fontWeightMedium,
      },
    },
    MuiListItemIcon: {
      root: {
        color: 'inherit',
        marginRight: 0,
        '& svg': {
          fontSize: 20,
        },
      },
    },
    MuiAvatar: {
      root: {
        width: 32,
        height: 32,
      },
    },
  },
  props: {
    MuiTab: {
      disableRipple: true,
    },
  },
  mixins: {
    ...theme.mixins,
    toolbar: {
      minHeight: 48,
    },
  },
};

type DashboardBaseProps = {
  children?: React.ReactNode;
};

export default function ThemeWrapper(props: DashboardBaseProps) {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      {props.children}
    </ThemeProvider>
  );
}
