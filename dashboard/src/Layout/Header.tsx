import AppBar from '@material-ui/core/AppBar';
import Badge from '@material-ui/core/Badge';
import Button from '@material-ui/core/Button';
import CircularProgress from '@material-ui/core/CircularProgress';
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import Menu, { MenuProps } from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import Paper from '@material-ui/core/Paper';
import Popover from '@material-ui/core/Popover';
import { makeStyles, withStyles } from '@material-ui/core/styles';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import AccountBox from '@material-ui/icons/AccountBox';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import MenuIcon from '@material-ui/icons/Menu';
import NotificationsIcon from '@material-ui/icons/Notifications';
import OpenInNew from '@material-ui/icons/OpenInNewOutlined';
import PowerSetting from '@material-ui/icons/PowerSettingsNew';
import axios from 'axios';
import React from 'react';
import Constants from '../Constants';
import packagejson from '../../package.json'

const lightColor = 'rgba(255, 255, 255, 0.7)'; // 'rgba(255, 255, 255, 0.7)'; // '#030417';

const useStyles = makeStyles((theme) => ({
  appBar: {
    backgroundColor: theme.palette.primary.dark,
    maxHeight: 40,
  },

  menuButton: {
    marginLeft: -theme.spacing(1),
  },
  iconButtonAvatar: {
    padding: 4,
  },
  link: {
    textDecoration: 'none',
    color: lightColor,
    '&:hover': {
      color: theme.palette.common.white,
    },
  },
  button: {
    borderColor: lightColor,
    color: 'white',
    fontSize: '15px',
    fontFamily: 'Open Sans, Rajdhani',
    maxHeight: 30,
    padding: 0,
  },
  paper: {
    // borderColor: 'black',
    // textAlign: 'center',
    minWidth: 400,
    maxWidth: 400,
    minHeight: 400,
    maxHeight: 400,
    padding: theme.spacing(2),
  },
  paperItems: {
    //  textAlign: 'center',
    padding: theme.spacing(2),
  },
  notifContent: {
    fontSize: '12px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  loader: {
    display: 'flex',
    '& > * + *': {
      marginLeft: theme.spacing(2),
    },
    menuItemText: {
      color: 'black',
      fontSize: '14px',
      fontFamily: 'Open Sans, Rajdhani',
    },
  },
  version: {
    paddingLeft: 5,
    paddingRight: 5,
    background: '#1b1b32',
    color: 'black',
  },
}));

// const MMenu = (props: any) => (
//   <Menu
//     open={props.open}
//     elevation={0}
//     getContentAnchorEl={null}
//     anchorOrigin={{
//       vertical: 'bottom',
//       horizontal: 'right',
//     }}
//     transformOrigin={{
//       vertical: 'top',
//       horizontal: 'right',
//     }}
//   />
// );

// const StyledMenu = withStyles(() => ({
//   paper: {
//     border: '1px solid #d3d4d5',
//   },
// }))(MMenu) as typeof MMenu;

const StyledMenu = withStyles({
  paper: {
    border: '1px solid #d3d4d5',
  },
})((props: MenuProps) => (
  <Menu
    elevation={0}
    getContentAnchorEl={null}
    anchorOrigin={{
      vertical: 'bottom',
      horizontal: 'center',
    }}
    transformOrigin={{
      vertical: 'top',
      horizontal: 'center',
    }}
    {...props}
  />
));

const StyledMenuItem = withStyles(() => ({
  root: {
    // '&:hover': {
    //    backgroundColor: theme.palette.secondary.main,
    //   '& .MuiListItemIcon-root, & .MuiListItemText-primary': {
    //     color: theme.palette.common.white,
    //   },
    // },
  },
}))(MenuItem) as typeof MenuItem;

const StyledBadge = withStyles(() => ({
  badge: {
    right: -3,
    // top: 13,
    // border: `2px solid ${theme.palette.background.paper}`,
    // padding: '0 4px',
    maxHeight: 10,
    fontSize: '10px',
  },
}))(Badge);

export default function MainHeaderBar(props: any) {
  const classes = useStyles();
  // const [value, setValue] = React.useState(0);
  const [anchorEl, setAnchorEl] = React.useState<HTMLElement | null>(null);
  const [alertAnchorEl, setalertAnchorEl] = React.useState<HTMLButtonElement | null>(null);
  const [notifData, setnotifData] = React.useState([]);
  const [loader, setLoader] = React.useState(false);

  const changeAlertMenuState = (e: any) => {
    setalertAnchorEl(e.currentTarget);
  };

  const handleAlertMenuClose = () => {
    setalertAnchorEl(null);
  };

  const onUserAvatarClicked = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const handleUserMenuClicked = (name: any) => () => {
    setAnchorEl(null);
    if (name === 'logout') {
      localStorage.clear();
      axios.delete(`${Constants.TRASA_HOSTNAME}/api/v1/user/logout`);
      window.location.href = '/login';
    } else if (name === 'my') {
      window.location.href = '/my';
    }
  };

  const getMyPendingNotifs = () => {
    const url = `${Constants.TRASA_HOSTNAME}/api/v1/my/notifs`;
    axios
      .get(url)
      .then((response) => {
        if (response.data.status === 'success') {
          const data = response.data.data[0];
          setnotifData(data);
        }
      })
      .catch((error) => {
        console.log(error);
      });
  };

  React.useEffect(() => {
    getMyPendingNotifs();
  }, []);

  const ResolveNotif = (notificationID: any) => {
    setLoader(true);

    const req = { notifID: notificationID };

    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/my/notif/resolve`, req)
      .then((response) => {
        setLoader(false);
        if (response.data.status === 'success') {
          getMyPendingNotifs();
          setLoader(false);
        }
      })
      .catch((error) => {
        console.log(error);
      });
  };

  const showNotifData = () => {
    if (notifData.length === 0) {
      return <Typography variant="h4">You do not have pending notifications...</Typography>;
    }
    return notifData.map((n, i) => (
      <Grid item xs={12} key={i}>
        <NotifCard notif={n} ResolveNotif={ResolveNotif} />
      </Grid>
    ));
  };

  return (
    <AppBar className={classes.appBar} position="sticky" elevation={0}>
      <Toolbar>
        <Grid container spacing={4} alignItems="center">
          <Grid item>
            <IconButton
              color="inherit"
              size="small"
              aria-label="Open drawer"
              onClick={props.onDrawerToggle}
              className={classes.menuButton}
            >
              <MenuIcon />
            </IconButton>
          </Grid>
          <Grid item xs />

          <Grid item>
            <div className={classes.version}>
              <Typography className={classes.link} component="a" href="#" variant="h6">
                v{packagejson.version}
              </Typography>
            </div>
          </Grid>

          <Grid item>
            <IconButton id="notif-menu" color="inherit" size="small" onClick={changeAlertMenuState}>
              <StyledBadge color="secondary" badgeContent={notifData.length}>
                <NotificationsIcon style={{ maxHeight: 17 }} />
              </StyledBadge>
            </IconButton>

            <Popover
              id="notif-menu"
              open={Boolean(alertAnchorEl)}
              anchorEl={alertAnchorEl}
              anchorOrigin={{
                vertical: 'bottom',
                horizontal: 'right',
              }}
              transformOrigin={{
                vertical: 'top',
                horizontal: 'right',
              }}
              onClose={handleAlertMenuClose}
            >
              <div className={classes.paper}>
                <Grid container spacing={2}>
                  {loader ? (
                    <div className={classes.loader}>
                      <CircularProgress />
                    </div>
                  ) : (
                    showNotifData()
                  )}
                </Grid>
              </div>
            </Popover>
          </Grid>
          <Grid item>
            <Button
              id="user-menu"
              variant="contained"
              size="small"
              className={classes.button}
              onClick={onUserAvatarClicked}
            >
              {(props.userData.firstName[0] || '') + (props.userData.lastName[0] || '')}

              <ExpandMoreIcon />
            </Button>

            <StyledMenu
              id="user-menu"
              aria-haspopup="true"
              keepMounted
              anchorEl={anchorEl}
              open={Boolean(anchorEl)}
              onClose={handleMenuClose}
            >
              <StyledMenuItem onClick={handleUserMenuClicked('my')}>
                <ListItemIcon>
                  <AccountBox fontSize="small" />
                </ListItemIcon>
                <div>My Account</div>
              </StyledMenuItem>

              <StyledMenuItem
                href="https://www.trasa.io/docs"
                target="_blank"
                onClick={() => {
                  window.location.href = 'https://www.trasa.io/docs ';
                }}
              >
                <ListItemIcon>
                  <OpenInNew fontSize="small" />
                </ListItemIcon>
                <div>Go to docs</div>
              </StyledMenuItem>

              <StyledMenuItem onClick={handleUserMenuClicked('logout')}>
                <ListItemIcon>
                  <PowerSetting fontSize="small" />
                </ListItemIcon>
                <div>Sign Out</div>
              </StyledMenuItem>
            </StyledMenu>
          </Grid>
        </Grid>
      </Toolbar>
    </AppBar>
  );
}

type NotifProp = {
  notificationLabel: string;
  notificationID: string;
};

const NotifCard = (props: any) => {
  const { notif } = props;

  function notifTitle() {
    switch (notif.notificationLabel) {
      case 'access-request':
        return 'Access Request';
      default:
        return 'Security Alert';
    }
  }

  function cancelButtonText() {
    switch (notif.notificationLabel) {
      case 'access-request':
        return 'Reject';
      default:
        return 'Dismiss';
    }
  }
  // const preventDefault = event => event.preventDefault();

  function cancelButtonAction(notif: NotifProp) {
    switch (notif.notificationLabel) {
      case 'access-request':
        return '';
      default:
        return (
          <Button
            variant="text"
            onClick={() => {
              props.ResolveNotif(notif.notificationID);
            }}
          >
            {cancelButtonText()}
          </Button>
        );
    }
  }

  function ReviewButtonAction(notif: NotifProp) {
    switch (notif.notificationLabel) {
      case 'access-request':
        return (
          <Button variant="text" href="/control#AdHoc%20Requests">
            Review
          </Button>
        );
      default:
        return '';
      // <Button variant="text" href="/system/security">
      //   Review
      // </Button>
    }
  }

  const classes = useStyles();
  return (
    <div>
      <Paper className={classes.paperItems} elevation={5}>
        {/* <Box className={classes.paperItems} border={1} borderRadius={5} borderColor="navt"> */}
        <Grid container spacing={2} direction="row">
          <Grid item xs={12}>
            <Typography component="h4" variant="h4">
              {notifTitle()}
            </Typography>

            <div className={classes.notifContent}>{notif.notificationText}</div>
          </Grid>
          <Grid container direction="row" justify="flex-end" alignItems="flex-end">
            <Grid item xs={3}>
              {ReviewButtonAction(notif)}
            </Grid>
            <Grid item xs={3}>
              {cancelButtonAction(notif)}
            </Grid>
          </Grid>
        </Grid>
        {/* </Box> */}
      </Paper>
    </div>
  );
};
