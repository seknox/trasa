import CssBaseline from '@material-ui/core/CssBaseline';
import { createMuiTheme, makeStyles, ThemeProvider } from '@material-ui/core/styles';
import React from 'react';

let theme = createMuiTheme({
  typography: {
    // useNextVariants: true,
    fontFamily: 'Open Sans, Rajdhani',
    h1: {
      fontWeight: 700,
      fontSize: 46,
      letterSpacing: 0.5,
      padding: 10,
      color: '#000066',
    },
    h2: {
      fontWeight: 600,
      fontSize: 36,
      letterSpacing: 0.5,
      color: '#000080',
      padding: 10,
    },
    h3: {
      fontWeight: 600,
      padding: 10,
      fontSize: 32,
      letterSpacing: 0.5,
      color: '#404854',
    },
    h4: {
      fontWeight: 600,
      fontSize: 26,
      letterSpacing: 0.5,
      color: '#404854',
      padding: 12,
    },
    h5: {
      fontWeight: 600,
      fontSize: 18,
      letterSpacing: 0.5,
      padding: 13,
    },
    h6: {
      fontWeight: 700,
      fontSize: 14,
      letterSpacing: 0.5,
      padding: 12,
    },
    // p: {
    //   fontWeight: 700,
    //   fontSize: 14,
    //   letterSpacing: 0.5,
    //   color: '#404854',
    //   padding: 12,
    // },
    button: {
      textTransform: 'none',
      fontFamily: 'Open Sans, Rajdhani',
    },
  },
  palette: {
    type: 'light',
    primary: { 500: 'rgba(1,1,35,1)' },
    secondary: { A400: '#000066' }, // '#000080' #1b1b32 #0A2053 #000066}, #0A2053
    // butprim:  { A400: '#000080'},
    error: { 500: '#9a0036' },
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

    MuiButton: {
      contained: {
        color: 'white',
        fontWeight: 600,
        //  fontSize: '14px',
        boxShadow: 'none',
        '&:active': {
          boxShadow: 'none',
          backgroundColor: '#000080',
        },
        '&:hover, &:focus': {
          backgroundColor: '#000066',
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
          borderColor: '#80bdff',
          boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
        },
      },
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

// const drawerWidth = 250;

const useStyles = makeStyles((theme) => ({
  root: {
    // display: 'flex',
    // minHeight: '100vh',
  },
}));

export default function ThemeBase(props) {
  const classes = useStyles();

  return (
    <ThemeProvider theme={theme}>
      <div className={classes.root}>
        <CssBaseline />

        {props.children}
      </div>
    </ThemeProvider>
  );
}
