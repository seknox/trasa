import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import Save from '@material-ui/icons/PersonAdd';
import axios from 'axios';
import React, { useEffect, useState } from 'react';
import Constants from '../../../Constants';
import ServicegroupUsergroupMapTable from '../Groups/GroupPage/ServicegroupUsergroupAccessMapTable';
import AppUsersTable from './AccessMapTable';
import NAddUsergroupToServicegroup from './NAddUsergroupToServicegroup';
// import AssignUsersAndGroupsToApp from './AssignUsersAndGroupsToApp';

const useStyles = makeStyles((theme) => ({
  root: {
    flexGrow: 1,
    marginTop: 20,
    justifyContent: 'space-between',
  },
  appRoot: {
    // flexgrow: 1,
    // padding: '28px 16px 0',
    // marginBotton: '5%',
  },
  paper: {
    backgroundColor: '#fdfdfd',
    padding: theme.spacing(1),
    color: theme.palette.text.secondary,
  },
  demo: {
    height: 70,
  },
  card: {
    maxWidth: 250,
    padding: theme.spacing(5),
    height: 200,
  },
  textFieldFormLabel: {
    fontSize: 18,
  },
  passthru: {
    marginTop: 30,
    marginLeft: 350,
  },
  tableroot: {
    width: '100%',
    // marginTop: theme.spacing(3),
    // overflow: 'auto',
  },
  table: {
    // minWidth: 700,
  },
  tablerow: {
    '&:nth-of-type(odd)': {
      backgroundColor: theme.palette.background.default,
    },
  },
  appValueBHText: {
    color: 'white',
    fontSize: '25px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  appValueHText: {
    color: 'black',
    fontSize: '15px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  appValueHMText: {
    marginTop: '5%',
    color: 'black',
    fontSize: '12px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  appValueSHText: {
    color: 'black',
    fontSize: '12px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  fab: {
    // margin: theme.spacing(2),
    // marginLeft: '20%',
    fontSize: '10px',
  },
}));

/// ////////////////////////////////////////////////////////

export default function Appuser(props: any) {
  const [appData, setAppData] = useState({
    ID: '',
    serviceName: '',
    serviceType: '',
    rdpProtocol: '',
    remoteserviceName: '',
    hostname: '',
    adhoc: false,
    passthru: false,
    nativeLog: false,
  });
  const [appUsers, setAppusers] = useState([]);
  const [allUsers, setAllusers] = useState([]);
  const [assignuserDlgState, setAssignuserDlgState] = useState(false);
  const [assignGroupDlgState, setAssignGroupDlgState] = useState(false);
  // const [updateAppUserState, setupdateAppUserState] = useState(false);
  // const [update, setUpdate] = useState(false);
  // const [singleAppUser, setSingleAppuser] = useState({});
  const [userGroups, setuserGroups] = useState([]);
  const [policies, setPolicies] = useState([]);

  const fetchUserGroups = () => {
    const url = `${Constants.TRASA_HOSTNAME}/api/v1/accessmap/servicegroup/usergroupstoadd`;
    axios
      .get(url)
      .then((response) => {
        const resp = response.data.data[0];
        setuserGroups(resp);
      })
      .catch((error) => {
        console.log(error);
      });
  };

  const fetchPolicies = () => {
    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/groups/policy/all`)
      .then((response) => {
        // this.setState({allUsers: response.data})
        const resp = response.data.data[0];
        setPolicies(resp);
      })
      .catch((error) => {
        console.log(error);
      });
  };

  const dataforAssignUserToApp = (entityID: string) => {
    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/accessmap/service/user/${entityID}`)
      .then((response) => {
        const resp = response.data.data[0];
        props.UpdateserviceName(resp.App.serviceName);
        if (response.status === 403) {
          window.location.href = '/login';
        }
        setAppData(resp.App);
        setAppusers(resp.AppUsers);
      })
      .catch((error) => {
        console.log(error);
      });

    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/user/all`)
      .then((response) => {
        if (response.status === 403) {
          window.location.href = '/login';
        }
        setAllusers(response.data.data[0]);
      })
      .catch((error) => {
        console.log(error);
      });
  };

  const changeAssignuserDlgState = () => {
    setAssignuserDlgState(!assignuserDlgState);
  };

  const changeAssignGroupDlgState = () => {
    setAssignGroupDlgState(!assignGroupDlgState);
  };

  const classes = useStyles();
  useEffect(() => {
    fetchUserGroups();
    fetchPolicies();
    dataforAssignUserToApp(props.entityID);
  }, [props.entityID]);

  return (
    <div className={classes.appRoot}>
      <NAddUsergroupToServicegroup
        open={assignuserDlgState}
        serviceName={appData.serviceName}
        handleClose={changeAssignuserDlgState}
        userGroups={allUsers}
        policies={policies}
        groupMeta={props.groupMeta}
        assignuser
        ID={props.entityID}
        renderFor="assignUserToApp"
      />

      <NAddUsergroupToServicegroup
        open={assignGroupDlgState}
        serviceName={appData.serviceName}
        handleClose={changeAssignGroupDlgState}
        userGroups={userGroups}
        policies={policies}
        groupMeta={props.groupMeta}
        assignuser={false}
        ID={props.entityID}
        renderFor="assignUserGroup"
      />
      <br />

      <Grid container>
        <Grid item xs={2}>
          <Button id='assignUserBtn' variant="contained" size="small" onClick={changeAssignuserDlgState}>
            <Save />
            Assign User
          </Button>
        </Grid>
        <Grid item xs={2}>
          <Button id='assignUserGroupBtn' variant="contained" size="small" onClick={changeAssignGroupDlgState}>
            <Save />
            Assign User Groups
          </Button>
        </Grid>
      </Grid>

      <ServicegroupUsergroupMapTable groupID={props.entityID} />

      <br />
      <AppUsersTable AppUsers={appUsers} ID={props.entityID} />
    </div>
  );
}
