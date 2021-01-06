import Avatar from '@material-ui/core/Avatar';
import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
// import UserCards from './UserOverviewCard';
import ListItemText from '@material-ui/core/ListItemText';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import Paper from '@material-ui/core/Paper';
import { makeStyles, Theme, withStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import MoreVertIcon from '@material-ui/icons/MoreVert';
import axios from 'axios';
import cx from 'clsx';
import React, { Dispatch, SetStateAction, useEffect, useState } from 'react';
import Constants from '../../Constants';
import DialogueWrapper from '../../utils/Components/DialogueWrapComponent';
import { HeaderFontSize } from '../../utils/Responsive';
import LoginBox from '../Auth/index';
import { SetPasswordComponent } from '../Publicpages/PasswordSetupReset';

const fileDownload = require('js-file-download');

const useStyles = makeStyles((theme) => ({
  root: {
    flexgrow: 1,
    // marginBotton: '5%',
  },

  paper: {
    backgroundColor: '#fdfdfd',

    padding: theme.spacing(2),
    //  textAlign: 'center',
    color: theme.palette.text.secondary,
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
    margin: '-40px 34% 0 43%',
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
}));

export default function UserOverview() {
  const classes = useStyles();

  const [userData, setUserData] = useState<any>({});
  const [userGroups, setUsergroups] = useState([]);
  const [userDevices, setUserdevices] = useState([]);
  const [userAccessMaps, setuserAccessMaps] = useState([]);

  useEffect(() => {
    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/my/account-details`)
      .then((r) => {
        if (r.data.status === 'success') {
          setUserData(r.data.data[0].user);
          setUsergroups(r.data.data[0].userGroups);
          setUserdevices(r.data.data[0].userDevices);
          setuserAccessMaps(r.data.data[0].userAccessMaps);
        }
      })
      .catch((error) => {
        console.log('err: ', error);
      });
  }, []);

  return (
    <div className={classes.root}>
      {/* direction="row" alignItems="center" justify="center" */}
      <Grid container spacing={3} direction="row" justify="center" alignItems="center">
        <Grid item xs={12} sm={12} md={6}>
          <Grid container spacing={3}>
            <Grid item xs={12}>
              <UserCard
                email={userData && userData.email}
                user={userData}
                totalGroups={userGroups.length}
                totalServices={userAccessMaps.length}
                totalUserdevices={userDevices.length}
              />
            </Grid>
          </Grid>
        </Grid>
      </Grid>
    </div>
  );
}

function AggStat(props: any) {
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

function UserCard(props: any) {
  const classes = useStyles();

  function getTimeFromTimestamp(tstr: number) {
    const date = new Date(tstr * 1000);
    return date.toDateString();
  }

  const [updatePassDlgState, setUpdatePassDlgState] = useState(false);

  const changeUpdatePassDlgState = () => {
    setUpdatePassDlgState(!updatePassDlgState);
  };

  const closeUpdatePassDlg = () => {
    setUpdatePassDlgState(false);
  };

  const [hasAuthenticated, setHasAuthentecated] = useState(false);

  function changeHasAuthenticated() {
    setHasAuthentecated(!hasAuthenticated);
  }

  const [changePassToken, setChangePassToken] = useState('');

  return (
    <Grid container spacing={2} direction="column" alignItems="center" justify="center">
      <Paper className={cx(classes.card, classes.shadowRise)}>
        <Grid container spacing={2}>
          <Avatar className={classes.avatar}>
            {' '}
            {props.user.firstName ? props.user.firstName[0] + props.user.lastName[0] : 'U'}{' '}
          </Avatar>
          <SettinMenu changeUpdatePassDlgState={changeUpdatePassDlgState} />
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
            <AggStat
              totalGroups={props.totalGroups}
              totalServices={props.totalServices}
              totalUserdevices={props.totalUserdevices}
            />
          </Grid>
        </Grid>
      </Paper>

      <DialogueWrapper
        open={updatePassDlgState}
        handleClose={closeUpdatePassDlg}
        title="Change Password"
        maxWidth="md"
        fullScreen
      >
        <UpdatePasswordDlg
          hasAuthenticated={hasAuthenticated}
          userData={props.user}
          changeHasAuthenticated={changeHasAuthenticated}
          setToken={setChangePassToken}
          changePassToken={changePassToken}
          closeUpdatePassDlg={closeUpdatePassDlg}
        />
      </DialogueWrapper>
    </Grid>
  );
}

type UpdatePasswordDlgProps = {
  hasAuthenticated: boolean;
  userData: any;
  changeHasAuthenticated: () => void;
  setToken: Dispatch<SetStateAction<string>>;
  changePassToken: string;
  closeUpdatePassDlg: () => void;
};

function UpdatePasswordDlg(props: UpdatePasswordDlgProps) {
  const { closeUpdatePassDlg, changePassToken, userData, changeHasAuthenticated, setToken } = props;

  return (
    <Grid container spacing={2} direction="column">
      <Grid item xs={12}>
        {props.hasAuthenticated === true ? (
          <SetPasswordComponent
            token={changePassToken}
            closeUpdatePassDlg={closeUpdatePassDlg}
            update
          />
        ) : (
          <LoginBox
            autofillEmail
            intent="AUTH_REQ_CHANGE_PASS"
            showForgetPass={false}
            title="Authenticate again to proceed"
            userData={userData}
            proxyDomain={''}
            changeHasAuthenticated={changeHasAuthenticated}
            setData={setToken}
          />
        )}
      </Grid>
    </Grid>
  );
}

const MMenu = (props: any) => (
  <Menu
    elevation={0}
    getContentAnchorEl={null}
    anchorOrigin={{
      vertical: 'bottom',
      horizontal: 'right',
    }}
    transformOrigin={{
      vertical: 'top',
      horizontal: 'center',
    }}
    {...props}
  />
);
const StyledMenu = withStyles((theme: Theme) => ({
  paper: {
    border: '1px solid #d3d4d5',
  },
}))(MMenu) as typeof MMenu;

const StyledMenuItem = withStyles((theme) => ({
  root: {
    '&:hover': {
      backgroundColor: theme.palette.secondary.main,
      '& .MuiListItemIcon-root, & .MuiListItemText-primary': {
        color: theme.palette.common.white,
      },
    },
  },
}))(MenuItem) as typeof MenuItem;

type SettingsMenuProps = {
  changeUpdatePassDlgState: () => void;
};

function SettinMenu(props: SettingsMenuProps) {
  const [anchorEl, setAnchorEl] = React.useState<any>();

  const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const GeneratePrivateKey = () => {
    const config = {
      responseType: 'blob',
      headers: {
        'X-CSRF': localStorage.getItem('X-CSRF'),
      },
    };

    axios.post(`${Constants.TRASA_HOSTNAME}/api/v1/my/generatekey`, null, config).then((resp) => {
      fileDownload(resp.data, 'id_rsa.zip', 'application/zip');
    });
  };

  return (
    <div>
      <IconButton
        aria-label="More"
        aria-owns={anchorEl && 'long-menu'}
        aria-haspopup="true"
        onClick={handleClick}
      >
        <MoreVertIcon />
      </IconButton>
      <StyledMenu
        id="customized-menu"
        anchorEl={anchorEl}
        keepMounted
        open={Boolean(anchorEl)}
        onClose={handleClose}
      >
        <StyledMenuItem onClick={props.changeUpdatePassDlgState}>
          <ListItemText primary="Change Password" />
        </StyledMenuItem>
        <StyledMenuItem onClick={GeneratePrivateKey}>
          <ListItemText primary="Get ssh private key" />
        </StyledMenuItem>
      </StyledMenu>
    </div>
  );
}
