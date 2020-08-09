import Button from '@material-ui/core/Button';
import green from '@material-ui/core/colors/green';
import purple from '@material-ui/core/colors/purple';
import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import InputBase from '@material-ui/core/InputBase';
import Paper from '@material-ui/core/Paper';
import { makeStyles } from '@material-ui/core/styles';
import Tooltip from '@material-ui/core/Tooltip';
import SearchIcon from '@material-ui/icons/Search';
import Security from '@material-ui/icons/Security';
import axios from 'axios';
import React, { useEffect, useState } from 'react';
// import {Link } from 'react-router-dom';
import { Link } from 'react-router-dom';
import Drawer from '@material-ui/core/Drawer';
import Service from '../../assets/services.png';
import DatabaseIcon from '../../assets/database.png';
import RdpIcon from '../../assets/rdp.png';
import SshIcon from '../../assets/ssh.png';
import Constants from '../../Constants';
import Servicesetting from './Service/Settings/ServiceSetting';

const useStyles = makeStyles((theme) => ({
  root: {
    flexgrow: 1,
    // marginBotton: '5%',
  },

  paper: {
    backgroundColor: '#fdfdfd',
    // backgroundColor: '#00001aff', //'#E0E0E0', 'rgba(10,34,52,1)' // #011019
    padding: theme.spacing(2),
    textAlign: 'center',
    color: theme.palette.text.secondary,
    // minWidth: 100,
  },

  card: {
    marginLeft: 20,
    height: 200,
  },
  buttonSpace: {
    flexgrow: 1,
    justifyContent: 'space-between',
  },
  Servicebutton: {
    color: 'white',
    '&:hover, &:focus': {
      backgroundColor: '#000066',
      boxShadow: '0 0 10px #030417',
      color: 'white',
    },
    backgroundColor: '#000080',
    // backgroundColor: 'navy', // #1464F4 , #0000CD
  },
  requestButton: {
    color: 'white',
    backgroundColor: '#01579b',
  },
  // form
  formControl: {
    margin: theme.spacing(1),
  },
  inputLabelFocused: {
    color: purple[500],
  },
  inputInkbar: {
    '&:after': {
      backgroundColor: purple[500],
    },
  },
  textFieldRoot: {
    padding: 1,
    'label + &': {
      marginTop: theme.spacing(3),
    },
  },

  textFieldInput: {
    borderRadius: 4,
    backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    fontSize: 16,
    padding: '10px 12px',
    width: 'calc(100% - 4px)',
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
    },
  },
  textFieldInputBig: {
    borderRadius: 4,
    backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    fontSize: 16,
    //    padding: '10px 100px',
    //     width: 'calc(100% - 4px)',
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
    },
  },
  textFieldFormLabel: {
    fontSize: 18,
  },
  buttonProgress: {
    color: green[500],
    position: 'absolute',
    top: '50%',
    left: '50%',
    marginTop: -12,
    marginLeft: -12,
  },
  successText: {
    fontSize: 15,
    color: 'green',
  },
  errorText: {
    fontSize: 15,
    color: 'red',
  },
  fab: {
    margin: theme.spacing(2),
  },
  fab2: {
    margin: theme.spacing(2),
  },
  users: {
    margin: theme.spacing(2),
  },
  dividerInset: {
    margin: `5px 0 0 ${theme.spacing(9)}px`,
  },
  servicesDemiter: {
    marginBottom: 20,
    marginTop: 50,
  },
  lightTooltip: {
    backgroundColor: theme.palette.common.white,
    color: 'rgba(0, 0, 0, 0.87)',
    boxShadow: theme.shadows[1],
    fontSize: 11,
  },
  serviceName: {
    color: 'black',
    fontSize: '12px',
    fontFamily: 'Open Sans, Rajdhani',
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
  divider: {
    width: 1,
    height: 28,
    margin: 4,
  },
}));

type Anchor = 'top' | 'left' | 'bottom' | 'right';

export default function ServiceList() {
  const classes = useStyles();
  const [https, sethttps] = useState([]);
  const [rdp, setrdp] = useState([]);
  const [ssh, setssh] = useState([]);
  const [db, setdb] = useState([]);
  const [radius, setRadius] = useState([]);
  const [other, setothers] = useState([]);
  const [query, setquery] = useState('');

  const fetchservices = (isGroupOrAllservices: string) => {
    let url = '';
    if (isGroupOrAllservices === 'allservices') {
      url = `${Constants.TRASA_HOSTNAME}/api/v1/services/all`;
    } else {
      url = `${Constants.TRASA_HOSTNAME}/api/v1/groups/service/${isGroupOrAllservices}`;
    }
    axios
      .get(url)
      .then((response) => {
        sethttps(response.data.data[0].http);
        setrdp(response.data.data[0].rdp);
        setssh(response.data.data[0].ssh);
        setdb(response.data.data[0].db);
        setRadius(response.data.data[0].radius);
        setothers(response.data.data[0].other);
      })
      .catch((error) => {
        console.log(error);
      });
  };

  useEffect(() => {
    fetchservices('allservices');
  }, []);

  const searchApp = (e: any) => {
    setquery(e.target.value);
  };

  const [configDrawerState, setConfigDrawerState] = useState({ right: false });

  const toggleConfigDrawer = (side: Anchor, open: boolean) => (
    event: React.KeyboardEvent | React.MouseEvent,
  ) => {
    if (
      event.type === 'keydown' &&
      ((event as React.KeyboardEvent).key === 'Tab' ||
        (event as React.KeyboardEvent).key === 'Shift')
    ) {
      return;
    }
    setConfigDrawerState({ ...configDrawerState, [side]: open });
  };

  const serviceDetail = {
    ID: '',
    serviceName: '',
    serviceType: 'db',
    rdpProtocol: '',
    remoteserviceName: '',
    passthru: false,
    hostname: '',
    nativeLog: false,
    adhoc: false,
  };

  return (
    <div className={classes.root}>
      <Button variant="contained" size="small" onClick={toggleConfigDrawer('right', true)}>
        Create new Service
      </Button>
      <Drawer
        anchor="right"
        open={configDrawerState.right}
        onClose={toggleConfigDrawer('right', false)}
      >
        <Paper className={classes.paper}>
          <Servicesetting newApp serviceDetail={serviceDetail} />
        </Paper>
      </Drawer>
      <Paper className={classes.searchRoot}>
        <IconButton className={classes.iconButton} aria-label="Search">
          <SearchIcon />
        </IconButton>
        <InputBase
          className={classes.searchInput}
          onChange={searchApp}
          placeholder="Search Services by name or hostname"
          inputProps={{ 'aria-label': 'Search services' }}
        />
      </Paper>
      <div className={classes.servicesDemiter}>
        <p>HTTPs services</p>
        <Divider light />{' '}
      </div>
      <Renderservices data={https} query={query} serviceType="http" />

      <div className={classes.servicesDemiter}>
        <p>RDP services</p>
        <Divider light />{' '}
      </div>
      <Renderservices data={rdp} serviceType="rdp" query={query} />
      <div className={classes.servicesDemiter}>
        <p>SSH services</p>
        <Divider light />{' '}
      </div>
      <Renderservices data={ssh} serviceType="ssh" query={query} />

      <div className={classes.servicesDemiter}>
        <p>Database services</p>
        <Divider light />{' '}
      </div>
      <Renderservices data={db} serviceType="db" query={query} />
      <div className={classes.servicesDemiter}>
        <p>Radius services</p>
        <Divider light />{' '}
      </div>
      <Renderservices data={radius} serviceType="radius" query={query} />

      <div className={classes.servicesDemiter}>
        <p>Other services</p>
        <Divider light />{' '}
      </div>
      <Renderservices data={other} serviceType="other" query={query} />
    </div>
  );
}

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

// TODO @bhrg3se fix search functionality here!
function Renderservices(props: any) {
  const classes = useStyles();

  const [data, setData] = useState([]);

  useEffect(() => {
    if (!props.data) {
      return;
    }
    const filteredService = props.data.filter((a: any) => {
      return JSON.stringify(a).toUpperCase().includes(props.query.toUpperCase().trim());
    });
    setData(filteredService);
  }, [props.data, props.query]);

  const returnAppIcon = (val: any) => {
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

  // let data = props.data
  //   ? props.data.filter((a: any) => {
  //       JSON.stringify(a).toUpperCase().includes(props.query.toUpperCase().trim());
  //     })
  //   : '';

  return (
    <Grid container spacing={2}>
      {data.map((value: any, k: number) => (
        <Grid key={k} item xs={6} sm={2} md={2} lg={2}>
          <Paper className={classes.paper}>
            {value.adhoc ? (
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
            ) : (
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
            )}

            <img
              alt="app-icon"
              src={returnAppIcon(props.serviceType)}
              style={{
                height: 40,
                marginTop: 1,
              }}
            />

            <div />

            <div className={classes.serviceName}> {value.serviceName} </div>
            <div className={classes.buttonSpace} />
            <br />
            <Button
              variant="outlined"
              color="secondary"
              component={Link}
              to={`/services/service/${value.ID}`}
            >
              Service Setting
            </Button>
          </Paper>
        </Grid>
      ))}
    </Grid>
  );
}
