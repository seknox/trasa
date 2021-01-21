import Button from '@material-ui/core/Button';
import green from '@material-ui/core/colors/green';
import purple from '@material-ui/core/colors/purple';
// import DialogTitle from '@material-ui/core/DialogTitle';
import Dialog from '@material-ui/core/Dialog';
import MuiDialogActions from '@material-ui/core/DialogActions';
import MuiDialogContent from '@material-ui/core/DialogContent';
import MuiDialogTitle from '@material-ui/core/DialogTitle';
import Divider from '@material-ui/core/Divider';
import FormControl from '@material-ui/core/FormControl';
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import InputBase from '@material-ui/core/InputBase';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import Paper from '@material-ui/core/Paper';
import Select from '@material-ui/core/Select';
import { makeStyles, withStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
import Tooltip from '@material-ui/core/Tooltip';
import CastIcon from '@material-ui/icons/Cast';
import SearchIcon from '@material-ui/icons/Search';
import Security from '@material-ui/icons/Security';
import axios from 'axios';
import React, { useEffect, useState } from 'react';
import DatabaseIcon from '../../assets/database.png';
import RdpIcon from '../../assets/rdp.png';
import Service from '../../assets/services.png';
import SshIcon from '../../assets/ssh.png';
import Constants from '../../Constants';
import ProgressHOC from '../../utils/Components/Progressbar';
import NewConnDlg from './NewConn';
import Typography from '@material-ui/core/Typography';

const useStyles = makeStyles((theme) => ({
  root: {
    flexgrow: 1,
    // marginBotton: '5%',
  },

  paper: {
    backgroundColor: '#fdfdfd',

    padding: theme.spacing(2),
    textAlign: 'center',
    color: theme.palette.text.secondary,
  },
  extendedIcon: {
    marginRight: '10px',
  },
  searchRoot: {
    marginLeft: '35%',
    padding: '2px 4px',
    display: 'flex',
    alignItems: 'center',
    width: 400,
  },
  searchInput: {
    marginLeft: 8,
    flex: 1,
  },
  iconButton: {
    padding: 10,
  },
  servicesDemiter: {
    marginBottom: 20,
  },
  lightTooltip: {
    backgroundColor: theme.palette.common.white,
    color: 'rgba(0, 0, 0, 0.87)',
    boxShadow: theme.shadows[1],
    fontSize: 11,
  },

  buttonSpace: {
    flexgrow: 1,
    justifyContent: 'space-between',
  },
  formControl: {
    margin: theme.spacing(1),
  },
}));

export default function MyservicesList() {
  const classes = useStyles();
  const [newconDlgOpen, setNewconDlgOpen] = useState(false);
  const [reqOpen, setReqOpen] = useState(false);
  const [serviceName, setServiceName] = useState('');
  const [serviceID, setserviceID] = useState('');
  const [query, setQuery] = useState('');
  const [assignedservices, setAssignedservices] = useState<any[]>([]);
  const [user, setUser] = useState({ email: '' });
  const [selectedServiceIndex, setSelectedServiceIndex] = useState(0);
  const [admins, setAdmins] = useState([]);

  const handleNewconDlgState = () => {
    setNewconDlgOpen(!newconDlgOpen);
  };

  const handleClickOpen = (
    lserviceID: string,
    serviceType: string,
    hostname: string,
    userName: string,
    tfaRequired: string,
  ) => {
    if (serviceType === 'http') {
      // var messenger = document.getElementById("trasextmsngr");

      // messenger.addEventListener("click", messageContentScript);

      // function messageContentScript() {
      window.postMessage(
        {
          direction: 'trasextmsngr',
          message: userName,
        },
        '*',
      );
      // }

      window.open(`https://${hostname}`);
    } else if (serviceType === 'ssh') {
      window.open(
        `/my/service/connectssh#username=${encodeURIComponent(
          userName,
        )}&serviceID=${lserviceID}&hostname=${hostname}&tfaRequired=${tfaRequired}`,
      );
    } else if (serviceType === 'db') {
      // window.open("sql://"+ encodeURIComponent(userName) + "@" + serviceID);setState({ open: false })
    } else if (serviceType === 'radius') {
      // window.open("sql://"+ encodeURIComponent(userName) + "@" + serviceID);setState({ open: false })
    } else {
      window.open(
        `/my/service/connectrdp#username=${encodeURIComponent(
          userName,
        )}&serviceID=${lserviceID}&hostname=${hostname}`,
      );

      // setState({ open: true , serviceID: serviceID, userName: userName})
    }
  };

  // handleClose = () => {
  //   setState({ open: false, menuOpen: false });
  // };

  useEffect(() => {
    const url = `${Constants.TRASA_HOSTNAME}/api/v1/my/services`;

    axios
      .get(url)
      .then((response) => {
        console.log(response.data);
        setUser(response.data.User);
        setAssignedservices(response.data?.data?.[0]?.myServices);
        console.log(response.data?.data?.[0]?.myServices);
        // setState({ user: response.data.User, services: response.data.UserService });
      })
      .catch((error) => {
        console.error(error);
      });

    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/my/services/adhoc/getadmins`)
      .then((r) => {
        if (r.data.status === 'success') {
          setAdmins(r.data?.data?.[0]);
        }
      })
      .catch((error) => {
        console.log(error);
      });
  }, []);

  const handleRequestDialogueOpen = (lserviceID: string, serviceName: string) => {
    setReqOpen(true);
    setserviceID(lserviceID);
    setServiceName(serviceName);
  };

  const handleRequestDialogueClose = () => {
    setReqOpen(false);
  };

  const sendAccessRequest = () => {
    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/my/services/adhoc/request`)
      .then((response) => {})
      .catch((error) => {
        console.error(error);
      });
  };

  const searchService = (e: any) => {
    setQuery(e.target.value);
  };

  const onConnectClicked = (index: any) => {
    // setMenuOpen(true);
    setSelectedServiceIndex(index);
  };

  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);

  const openPrivMenu = (event: React.MouseEvent<HTMLButtonElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const services = assignedservices.filter((a) =>
    JSON.stringify(a).toUpperCase().includes(query.toUpperCase().trim()),
  );
  return (
    <div className={classes.root}>
      <Button variant="contained" onClick={handleNewconDlgState}>
        <CastIcon className={classes.extendedIcon} />
        new connection
      </Button>
      <Paper className={classes.searchRoot}>
        <IconButton className={classes.iconButton} aria-label="Search">
          <SearchIcon />
        </IconButton>
        <InputBase
          className={classes.searchInput}
          onChange={searchService}
          placeholder="Search services by name or hostname"
          inputProps={{ 'aria-label': 'Search  service s' }}
        />
      </Paper>

      <NewConnDlg
        handleNewconDlgState={handleNewconDlgState}
        open={newconDlgOpen}
        close={handleNewconDlgState}
        services={services}
        //  email={user.email}
      />

      <br />
      <br />
      <br />
      <br />
      <div className={classes.servicesDemiter}>
        <p>Pre assigned Services</p>
        <Divider light />{' '}
      </div>

      <Grid container spacing={2}>
        {services.map((value: any, index) => (
          <Grid key={value.id} item xs={6} sm={4} md={3} lg={2}>
            <Paper className={classes.paper}>
              {value.isAuthorised ? (
                <Tooltip
                  title="Request if your policy does not authorize you at this time."
                  placement="top-end"
                  classes={{ tooltip: classes.lightTooltip }}
                >
                  <div style={{ marginLeft: '90%', color: 'navy' }}>
                    {' '}
                    <Security style={{ fontSize: 20 }} />{' '}
                  </div>
                </Tooltip>
              ) : (
                <Tooltip
                  title="Requires Adhoc Permission"
                  placement="top-end"
                  classes={{ tooltip: classes.lightTooltip }}
                >
                  <div style={{ marginLeft: '90%', color: '#b71c1c' }}>
                    {' '}
                    <Security style={{ fontSize: 20 }} />{' '}
                  </div>
                </Tooltip>
              )}

              <img
                alt="ServiceIcon"
                src={returnServiceIcon(value.serviceType)}
                style={{
                  height: 40,
                  marginTop: 1,
                }}
              />

              <div />
              <Typography variant="subtitle1" component="p">
                {' '}
                {value.serviceName}{' '}
              </Typography>

              <div className={classes.buttonSpace}>
                <br />
                {value.isAuthorised ? (
                  <Button
                    id="connect-privilege-menu"
                    name={value.serviceName}
                    variant="outlined"
                    color="secondary"
                    onClick={(e) => {
                      openPrivMenu(e);
                      onConnectClicked(index);
                    }}
                  >
                    Connect
                  </Button>
                ) : (
                  <Button
                    variant="contained"
                    onClick={() => {
                      handleRequestDialogueOpen(value.serviceID, value.serviceName);
                    }}
                    name={value.serviceName}
                  >
                    Request Access
                  </Button>
                )}
              </div>
              <br />
            </Paper>
          </Grid>
        ))}

        <Menu
          id="connect-privilege-menu"
          open={Boolean(anchorEl)}
          // keepMounted
          anchorEl={anchorEl}
          onClose={() => setAnchorEl(null)}
          anchorOrigin={{
            vertical: 'top',
            horizontal: 'right',
          }}
          transformOrigin={{
            vertical: 'top',
            horizontal: 'left',
          }}
        >
          {services[selectedServiceIndex] &&
            services[selectedServiceIndex].usernames.map((v: string) => (
              <MenuItem
                id={v}
                // name={v}
                onClick={() => {
                  handleClickOpen(
                    services[selectedServiceIndex].serviceID,
                    services[selectedServiceIndex].serviceType,
                    services[selectedServiceIndex].hostname,
                    v,
                    services[selectedServiceIndex].policy.tfaRequired,
                  );
                }}
              >
                {v}
              </MenuItem>
            ))}
        </Menu>

        <RequestAccess
          admins={admins}
          serviceID={serviceID}
          serviceName={serviceName}
          reqOpen={reqOpen}
          handleRequestDialogueClose={handleRequestDialogueClose}
        />
      </Grid>

      {/* <TfaPageHOC sendTfa={sendTfa} loading={loading} openTotpDlg={openTotpDlg} tfaFailed={tfaFailed} tfaFailedReason={tfaFailedReason}  /> */}
    </div>
  );
}

const returnServiceIcon = (val: any) => {
  if (val === 'ssh') {
    return SshIcon;
  }
  if (val === 'rdp') {
    return RdpIcon;
  }
  if (val === 'http') {
    return Service;
  }
  if (val === 'radius') {
    return Service;
  }

  if (val === 'db') {
    return DatabaseIcon;
  }
};

// function RenderIfNill(props: any) {
//   if (props.constructor === Array) {
//     return <div />;
//   }
//   return (
//     <h6>
//       looks like you have not created any services yet. Users wont be able to authenticate to protected
//       hosts unless configured.
//     </h6>
//   );
// }

type RequestAccessProps = {
  serviceID: string;
  serviceName: string;
  handleRequestDialogueClose: () => void;
  reqOpen: boolean;
  admins: any[];
};

function RequestAccess(props: RequestAccessProps) {
  const classes = useStyles();
  const [data, setData] = useState({ serviceID: '', requesteeID: '', requestTxt: '' });
  const [progress, setProgress] = useState(false);

  const handleChange = (event: any) => {
    setData({ ...data, [event.target.name]: event.target.value });
  };

  const handleSubmit = (event: any) => {
    setProgress(true);

    const d = data;
    // d.serviceID = `${props.serviceID}:${props.serviceName}`;
    d.serviceID = props.serviceID;
    // data['ServiceName'] = props.ServiceName

    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/my/services/adhoc/request`, data)
      .then((response) => {
        setProgress(false);
        if (response.data.status === 'success') {
        }

        // console.log(response.data);
      })
      .catch((error) => {
        console.log(error);
      });
  };

  return (
    <div>
      <Dialog
        onClose={props.handleRequestDialogueClose}
        aria-labelledby="customized-dialog-title"
        open={props.reqOpen}
        fullWidth
        maxWidth="sm"
      >
        <DialogTitle disableTypography>
          <Typography variant="h2">Request Access to {props.serviceName}</Typography>
        </DialogTitle>
        <DialogContent>
          <Grid container spacing={2}>
            <Grid item xs={5} sm={5} md={5}>
              <Typography variant="h3">Select admin</Typography>
            </Grid>
            <Grid item xs={7} sm={7} md={7}>
              <FormControl className={classes.formControl} fullWidth>
                <Select
                  name="requesteeID"
                  onChange={handleChange}
                  value={data.requesteeID}
                  variant="outlined"
                >
                  {props.admins.map((user: any) => (
                    <MenuItem value={user.ID}>{`${user.firstName} ${user.lastName}`} </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
          </Grid>

          <Grid container spacing={2}>
            <Grid item xs={5} sm={5} md={5}>
              <Typography variant="h3">Specify your intent for access</Typography>
            </Grid>
            <Grid item xs={7} sm={7} md={7}>
              <TextField
                fullWidth
                id="requestTxt"
                multiline
                // rowsMax="4"
                name="requestTxt"
                onChange={handleChange}
                value={data.requestTxt}
                margin="normal"
                helperText="write your reason for access"
                variant="outlined"
              />
            </Grid>
          </Grid>
        </DialogContent>
        {progress ? <ProgressHOC /> : ''}
        <DialogActions>
          <Button onClick={handleSubmit} variant="contained" color="primary">
            Submit Request
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
}

const DialogTitle = withStyles((theme) => ({
  root: {
    borderBottom: `1px solid ${theme.palette.divider}`,
    margin: 0,
    padding: theme.spacing(2),
  },
  closeButton: {
    position: 'absolute',
    right: theme.spacing(1),
    top: theme.spacing(1),
    color: theme.palette.grey[500],
  },
}))(MuiDialogTitle);

const DialogContent = withStyles((theme) => ({
  root: {
    margin: 0,
    padding: theme.spacing(2),
  },
}))(MuiDialogContent);

const DialogActions = withStyles((theme) => ({
  root: {
    borderTop: `1px solid ${theme.palette.divider}`,
    margin: 0,
    padding: theme.spacing(1),
  },
}))(MuiDialogActions);
