import Avatar from '@material-ui/core/Avatar';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import Divider from '@material-ui/core/Divider';
import Drawer from '@material-ui/core/Drawer';
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import Paper from '@material-ui/core/Paper';
import { makeStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import Delete from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import MoreVertIcon from '@material-ui/icons/MoreVert';
import axios from 'axios';
import cx from 'clsx';
import React, { useState } from 'react';
import Constants from '../../../Constants';
import { Usertype } from '../../../types/users';
// import UserCards from './UserOverviewCard';
import { HeaderFontSize } from '../../../utils/Responsive';
import UserCrud from './crud';

const useStyles = makeStyles((theme) => ({
  root: {
    flexgrow: 1,
    // marginBotton: '5%',
  },

  paperSmaller: {
    backgroundColor: '#fdfdfd',
    // backgroundColor: '#00001aff', //'#E0E0E0', 'rgba(10,34,52,1)'
    marginTop: '20%',

    // marginBotton: '5%',
    padding: theme.spacing(2),
    textAlign: 'center',
    color: theme.palette.text.secondary,
  },
  paper: {
    backgroundColor: '#fdfdfd',
    // backgroundColor: '#00001aff', //'#E0E0E0', 'rgba(10,34,52,1)' // #011019
    // minWidth: 400,
    // minHeight: 300,
    padding: theme.spacing(2),
    //  textAlign: 'center',
    color: theme.palette.text.secondary,
  },
  paperLarge: {
    backgroundColor: '#fdfdfd',
    // backgroundColor: '#00001aff', //'#E0E0E0', 'rgba(10,34,52,1)'
    minWidth: 400,
    minHeight: 300,
    marginTop: '5%',
    marginBotton: '5%',
    padding: theme.spacing(2),
    textAlign: 'center',
    color: theme.palette.text.secondary,
  },
  paperHeighted: {
    backgroundColor: '#fdfdfd',
    minWidth: 800,
    minHeight: 300,
    marginTop: '5%',
    marginBotton: '5%',
    padding: theme.spacing(2),
    textAlign: 'center',
    color: theme.palette.text.secondary,
  },

  // avatar: {
  //   backgroundColor: '#000080',
  //   width: theme.spacing(5),
  //   height: theme.spacing(5)
  // },

  heading: {
    fontWeight: 'bold',
    marginLeft: 50,
    marginBottom: 20,
    color: 'black',
    fontSize: '17px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  Heading: {
    fontWeight: 'bold',
    color: 'black',
    fontSize: '17px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  drawer: {
    background: 'white',
    width: 456,
  },
  drawerContent: {
    marginLeft: 30,
  },
  deviceHeader: {
    fontWeight: 'bold',
    background: 'WhiteSmoke',
    color: 'black',
    fontSize: '16px',
    fontFamily: 'Open Sans, Rajdhani',
    marginBottom: '1%',
    // marginLeft: 10,
  },
  deviceDetail: {
    padding: theme.spacing(1),
    color: 'black',
    fontSize: '14px',
    // padding: theme.spacing(1),
    fontFamily: 'Open Sans, Rajdhani',
    // marginLeft: 10,
  },
  paperTrans: {
    backgroundColor: 'transparent',
  },
  aggHeadersBig: {
    color: '#000066',
    fontSize: '21px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  aggHeaders: {
    color: '#1b1b32',
    fontSize: '19px',
    fontFamily: 'Open Sans, Rajdhani',
  },

  card: {
    marginTop: 40,
    borderRadius: 0.5, // theme.spacing(0.5),
    transition: '0.3s',
    // width: '90%',
    minWidth: 400,
    overflow: 'initial',
    background: '#ffffff',
    textAlign: 'center',
    padding: theme.spacing(2),
  },
  content: {
    textAlign: 'center',
    padding: theme.spacing(2),
  },
  shadowRise: {
    boxShadow: '0 0 20px 0 rgba(0,0,0,0.12)',
  },
  shadowFaded: {
    boxShadow: '0 2px 4px -2px rgba(0,0,0,0.24), 0 4px 24px -2px rgba(0, 0, 0, 0.2)',
  },
  cardHeader: {
    background: 'navy',
    borderRadius: 8,
    margin: '-20px auto 0',
    width: '88%',
    color: 'white',
    fontSize: '18px',
    fontWeight: 'bold',
    minHeight: 50,
  },
  title: {
    color: 'white',
    fontWeight: 'bold',
  },
  subheader: {
    color: 'rgba(255, 255, 255, 0.76)',
  },
  avatar: {
    width: 80,
    height: 80,
    background: 'navy',
    margin: '-40px 34% 0 41%',
    // color: 'white',
    fontSize: '24px',
    fontWeight: 'bold',
    transition: '0.3s',
    '&:hover': {
      transform: 'translateY(-3px)',
      boxShadow: '0 4px 20px 0 rgba(0,0,0,0.12)',
    },
  },
  subHeading: {
    color: '#000066',
    textAlign: 'left',
    fontSize: HeaderFontSize(), // 14,
    fontFamily: 'Open Sans, Rajdhani',
  },
  tHeading: {
    textAlign: 'right',
    color: '#1b1b32',
    fontSize: HeaderFontSize(), // window.innerHeight < 750 ? '14px':'18px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  warningButton: {
    color: 'maroon',
    //  fontSize: '44px',
    fontFamily: 'Open Sans, Rajdhani',
    '&:hover': {
      color: theme.palette.common.white,
      background: 'maroon',
    },
  },
}));

type UserOverviewProps = {
  userData: Usertype;
  userDevices: any;
  userGroups: any;
  userAccessMaps: any;
};
type Anchor = 'top' | 'left' | 'bottom' | 'right';

function getTimeFromTimestamp(tstr: any) {
  const date = new Date(tstr * 1000);
  return date.toDateString();
}

export default function UserOverview(props: UserOverviewProps) {
  const classes = useStyles();
  const [deleteDialogueOpen, setdeleteDialogueOpen] = useState(false);
  // const [deviceInfoDialogueOpen, setdeviceInfoDialogueOpen] = useState(false);
  // const [toDeleteVal, settoDeleteVal] = useState({});
  // const [updateDrawerOpen, setupdateDrawerOpen] = useState(false);
  // const [users, setusers] = useState(false);
  // const [userData, setuserData] = useState({});
  const [right, setright] = useState(false);
  // const [deviceDeleteState, setdeviceDeleteState] = useState(false);
  // const [selDeviceFinger, setselDeviceFinger] = useState({});
  // const [selDevice, setselDevice] = useState({});

  // useEffect(() => {
  //   axios
  //     .get(`${Constants.TRASA_HOSTNAME}/api/v1/my`)
  //     .then((response) => {
  //       setuserData(response.data.data[0]);
  //     })
  //     .catch((error) => {
  //       console.log(error);
  //     });
  // }, []);

  const handleDeleteDialogueClose = () => {
    setdeleteDialogueOpen(false);
  };

  const handleDeleteDialogueOpen = () => {
    setdeleteDialogueOpen(true);
  };

  // const handleDeviceDeleteDialogueClose = () => {
  //   setdeviceDeleteState(false);
  // };

  // const handleDeviceDeleteDialogueOpen = (device: any) => {
  //   setdeviceDeleteState(true);
  //   settoDeleteVal(device);
  // };

  // const handleUpdateDrawerOpen = () => {
  //   setupdateDrawerOpen(true);
  // };

  // const handleUpdateDrawerClose = () => {
  //   setupdateDrawerOpen(false);
  // };

  const openDrawer = (anchor: Anchor, open: boolean) => (
    event: React.KeyboardEvent | React.MouseEvent,
  ) => {
    if (
      event.type === 'keydown' &&
      ((event as React.KeyboardEvent).key === 'Tab' ||
        (event as React.KeyboardEvent).key === 'Shift')
    ) {
      return;
    }

    setright(open);
  };

  const deleteSingleUser = () => {
    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/user/delete/${props.userData.ID}`)
      .then((response) => {
        handleDeleteDialogueClose();
        if (response.data.status === 'success') {
          window.location.href = '/users';
        }
      })
      .catch((error) => {
        console.log(error);
      });
  };

  // const showDeviceDetail = (dev: any, devFinger: any) => {
  //   setselDevice(dev);
  //   setselDeviceFinger(devFinger);
  //   setdeviceInfoDialogueOpen(true);
  // };

  const sideList = (
    <div>
      <UserCrud
        userData={props.userData}
        update
        updateUserTable={() => {}}
        handleVerifyLinkChange={() => {}}
        handleDrawerClose={() => {}}
      />
    </div>
  );

  return (
    <div className={classes.root}>
      {/* direction="row" alignItems="center" justify="center" */}
      <Grid container spacing={3} direction="row" justify="center" alignItems="center">
        <Grid item xs={12} sm={12} md={6}>
          <Grid container spacing={3}>
            <Grid item xs={12}>
              <UserCardH
                email={props.userData.email}
                user={props.userData}
                deleteUser={handleDeleteDialogueOpen}
                handleUpdateDrawerOpen={openDrawer('right', true)}
                totalGroups={props.userGroups.length}
                totalServices={props.userAccessMaps.length}
                totalUserdevices={props.userDevices.length}
              />
            </Grid>
          </Grid>
        </Grid>

        <Dialog
          open={deleteDialogueOpen}
          onClose={handleDeleteDialogueClose}
          aria-labelledby="alert-dialog-title"
          aria-describedby="alert-dialog-description"
        >
          <DialogTitle id="alert-dialog-title">Confirm Delete?</DialogTitle>
          <DialogContent>
            <DialogContentText id="alert-dialog-description">
              Are You Sure You Want To Delete?
            </DialogContentText>
          </DialogContent>
          <DialogActions>
            <Button onClick={deleteSingleUser} className={classes.warningButton}>
              Yes, Delete
            </Button>
            <Button
              onClick={handleDeleteDialogueClose}
              color="primary"
              variant="contained"
              autoFocus
            >
              No
            </Button>
          </DialogActions>
        </Dialog>
        <Drawer anchor="right" open={right} onClose={openDrawer('right', false)}>
          <div className={classes.drawer}>
            <div className={classes.drawerContent}>{sideList}</div>
          </div>
        </Drawer>
      </Grid>
    </div>
  );
}

function AggStats(props: any) {
  const classes = useStyles();
  return (
    <div className={classes.root}>
      <br />
      <Grid container spacing={2}>
        <Grid item xs={4}>
          <Paper className={classes.paperTrans} elevation={0}>
            <div className={classes.aggHeaders}>
              {' '}
              <b> Groups </b>
            </div>
            <div className={classes.aggHeadersBig}>
              {' '}
              <b> {props.totalGroups} </b>{' '}
            </div>
          </Paper>
        </Grid>
        <Grid item xs={4}>
          <Paper className={classes.paperTrans} elevation={0}>
            <div className={classes.aggHeaders}>
              {' '}
              <b> Services </b>
            </div>
            <div className={classes.aggHeadersBig}>
              {' '}
              <b> {props.totalServices} </b>{' '}
            </div>
          </Paper>
        </Grid>
        <Grid item xs={4}>
          <Paper className={classes.paperTrans} elevation={0}>
            <div className={classes.aggHeaders}>
              {' '}
              <b> Device </b>
            </div>
            <div className={classes.aggHeadersBig}>
              {' '}
              <b> {props.totalUserdevices} </b>{' '}
            </div>
          </Paper>
        </Grid>
      </Grid>
    </div>
  );
}

function UserCardH(props: any) {
  const classes = useStyles();

  return (
    <Grid container spacing={2} direction="column" alignItems="center" justify="center">
      <Grid item xs={12}>
        <Paper className={cx(classes.card, classes.shadowRise)}>
          <Grid container spacing={2}>
            <Avatar className={classes.avatar}>
              {' '}
              {props.user.firstName ? props.user.firstName[0] + props.user.lastName[0] : 'U'}{' '}
            </Avatar>
            <DeleteMenu
              deleteUser={props.deleteUser}
              handleUpdateDrawerOpen={props.handleUpdateDrawerOpen}
            />
            <Grid item xs={12}>
              <Typography variant="h3">
                {`${props.user.firstName} ${props.user.middleName} ${props.user.lastName}`}
              </Typography>
              <Divider light variant="middle" />
            </Grid>

            <Grid item xs={12}>
              <Grid container spacing={1}>
                <Grid item xs={12}>
                  <Grid container spacing={4}>
                    <Grid item xs={6}>
                      <div className={classes.tHeading}> Identity Provider :</div>
                    </Grid>
                    <Grid item xs={6}>
                      <div className={classes.subHeading}>{props.user.idpName}</div>
                    </Grid>
                  </Grid>
                </Grid>
                <Grid item xs={12}>
                  <Grid container spacing={4}>
                    <Grid item xs={6}>
                      <div className={classes.tHeading}>Email Address :</div>
                    </Grid>
                    <Grid item xs={6}>
                      <div className={classes.subHeading}>{props.user.email}</div>
                    </Grid>
                  </Grid>
                </Grid>
                <Grid item xs={12}>
                  <Grid container spacing={4}>
                    <Grid item xs={6}>
                      <div className={classes.tHeading}>Username :</div>
                    </Grid>
                    <Grid item xs={6}>
                      <div className={classes.subHeading}>{props.user.userName}</div>
                    </Grid>
                  </Grid>
                </Grid>
                <Grid item xs={12}>
                  <Grid container spacing={4}>
                    <Grid item xs={6}>
                      <div className={classes.tHeading}>Role :</div>
                    </Grid>
                    <Grid item xs={6}>
                      <div className={classes.subHeading}>{props.user.userRole}</div>
                    </Grid>
                  </Grid>
                </Grid>
                <Grid item xs={12}>
                  <Grid container spacing={4}>
                    <Grid item xs={6}>
                      <div className={classes.tHeading}>Status :</div>
                    </Grid>
                    <Grid item xs={6}>
                      <div className={classes.subHeading}>
                        {props.user.status ? 'active' : 'disabled'}
                      </div>
                    </Grid>
                  </Grid>
                </Grid>
                <Grid item xs={12}>
                  <Grid container spacing={4}>
                    <Grid item xs={6}>
                      <div className={classes.tHeading}>Created at :</div>
                    </Grid>
                    <Grid item xs={6}>
                      <div className={classes.subHeading}>
                        {getTimeFromTimestamp(props.user.CreatedAt)}
                      </div>
                    </Grid>
                  </Grid>
                </Grid>
                <Grid item xs={12}>
                  <Grid container spacing={4}>
                    <Grid item xs={6}>
                      <div className={classes.tHeading}>Updated at :</div>
                    </Grid>
                    <Grid item xs={6}>
                      <div className={classes.subHeading}>
                        {getTimeFromTimestamp(props.user.UpdatedAt)}
                      </div>
                    </Grid>
                  </Grid>
                </Grid>
                <br />
                <br /> <br />
                <br />
              </Grid>
              {/* </Paper> */}
            </Grid>

            <Grid item xs={12}>
              <Divider light variant="middle" />
              <AggStats
                totalGroups={props.totalGroups}
                totalServices={props.totalServices}
                totalUserdevices={props.totalUserdevices}
              />
            </Grid>
          </Grid>
        </Paper>
      </Grid>
    </Grid>
  );
}

function DeleteMenu(props: any) {
  const [anchorEl, setanchorEl] = useState(null);

  const handleClick = (event: any) => {
    setanchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setanchorEl(null);
  };

  return (
    <div>
      <IconButton
        aria-label="More"
        aria-owns={anchorEl ? 'long-menu' : ''}
        aria-haspopup="true"
        onClick={handleClick}
      >
        <MoreVertIcon />
      </IconButton>
      <Menu
        id="long-menu"
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={handleClose}
        PaperProps={{}}
      >
        <MenuItem onClick={props.handleUpdateDrawerOpen}>
          <Button variant="contained" color="secondary">
            Update <EditIcon />
          </Button>
        </MenuItem>
        <MenuItem onClick={props.deleteUser}>
          <Button color="secondary" variant="contained">
            Delete
            <Delete />
          </Button>
        </MenuItem>
      </Menu>
    </div>
  );
}
