import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import { makeStyles } from '@material-ui/core/styles';
import Toolbar from '@material-ui/core/Toolbar';
import Tooltip from '@material-ui/core/Tooltip';
import AddIcon from '@material-ui/icons/AddCircleSharp';
import DeleteIcon from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import axios from 'axios';
import React, { useEffect, useState } from 'react';
import { Route, Switch, withRouter, RouteComponentProps } from 'react-router-dom';
import Constants from '../../../Constants';
import Layout from '../../../Layout/DashboardBase';
import Headers from '../../../Layout/HeaderAndContent';
import GroupUsersTable from '../User/userListTable';
import AddUsersToGroup from './AddUsersToGroup';
import UpdateGroup from './CreateUserGroup';
import ProgressBar from '../../../utils/Components/Progressbar';

const useStyles = makeStyles((theme) => ({
  root: {
    flexGrow: 1,
    marginTop: 20,
    justifyContent: 'space-between',
  },
  appRoot: {
    flexgrow: 1,
    padding: '28px 16px 0',
    // marginBotton: '5%',
  },
  paper: {
    backgroundColor: '#fdfdfd',
    padding: theme.spacing(1),
    color: theme.palette.text.secondary,
  },
  demo: {
    //  marginTop: 100,
    height: 70,
    // marginRight: 70,
  },
  card: {
    maxWidth: 250,
    //  width: 50,
    padding: theme.spacing(5),
    height: 200,
  },
  passthru: {
    marginTop: 30,
    marginLeft: 350,
  },
  tableroot: {
    width: '100%',
    marginTop: theme.spacing(3),
    //  overflowX: 'auto',
  },
  table: {
    minWidth: 700,
  },
  tablerow: {
    '&:nth-of-type(odd)': {
      backgroundColor: theme.palette.background.default,
    },
  },
  errorText: {
    fontSize: 15,
    color: 'red',
  },
  stepperButtonLeft: {
    // marginTop:'15%',
    marginLeft: theme.spacing(1),
    marginTop: theme.spacing(1),
    marginBottom: theme.spacing(1),
  },

  stepperButtonRight: {
    marginTop: '15%',
    marginRight: '10%',
  },
  newPaper: {
    backgroundColor: '#fdfdfd',
    padding: theme.spacing(2),

    color: theme.palette.text.secondary,
  },
  fab: {
    margin: theme.spacing(1),
  },
  Warning: {
    color: 'maroon',
    fontSize: '44px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  WarningButton: {
    color: 'white',
    //  fontSize: '44px',
    fontFamily: 'Open Sans, Rajdhani',
    background: 'maroon',
    '&:hover': {
      color: theme.palette.common.white,
      background: 'maroon',
    },
  },
  settingHeader: {
    color: 'black',
    fontSize: '15px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  settingSHeader: {
    color: 'black',
    fontSize: '12px',
    fontFamily: 'Open Sans, Rajdhani',
  },
}));

function UserGroupPage(props: any) {
  const [groupMeta, setgroupMeta] = useState<any>({});
  const [allAddedUsersObj, setallAddedUsersObj] = useState({});
  const [addedUsersArray, setaddedUsersArray] = useState([]);
  const [addUsersDlg, setaddUsersDlg] = useState(false);
  const [addedUsersIndex, setaddedUsersIndex] = useState([]);
  // const [addedUsers, setaddedUsers] = useState([]);
  const [unaddedUsers, setunaddedUsers] = useState([]);

  const onserviceselected = (addedUsersIndex: any, ServiceData: any) => {
    // const addedUsers = ServiceData
    //   .filter((u: any, i: any) => addedUsersIndex.indexOf(i) > -1)
    //   .map((app: any) => ({ appName: app.appName, serviceID: app.serviceID }));
    setaddedUsersIndex(addedUsersIndex);
    // setaddedUsers(addedUsers);
  };

  const openaddUsersDlg = () => {
    setaddUsersDlg(!addUsersDlg);
  };

  const fetchUsers = (isGroupIDOrAllUsers: string) => {
    let url = '';
    if (isGroupIDOrAllUsers === 'allusers') {
      url = `${Constants.TRASA_HOSTNAME}/api/v1/users/all`;
    } else {
      url = `${Constants.TRASA_HOSTNAME}/api/v1/groups/user/${isGroupIDOrAllUsers}`;
    }
    axios
      .get(url)
      .then((response) => {
        const resp = response.data.data[0];
        setallAddedUsersObj(resp.addedUsers);
        setunaddedUsers(resp.unaddedUsers);
        setgroupMeta(resp.groupMeta);

        const ddataArr = resp.addedUsers.map(function (n: any) {
          const date = new Date(n.CreatedAt * 1000);
          return [
            n.email.toString(),
            n.firstName,
            n.lastName,
            n.userName.toString(),
            n.userRole.toString(),
            date.toDateString(),
            n.idpName,
            n.status,
            n.ID,
          ];
        });
        setaddedUsersArray(ddataArr);
      })
      .catch((error) => {
        console.log(error);
      });
  };

  const deleteUsers = (rowsDeleted: any) => {
    setaddedUsersArray([]);
    const users = rowsDeleted.data.map((v: any) => {
      return addedUsersArray[v.index][8];
    });

    const reqData = { groupID: groupMeta.groupID, userIDs: users, updateType: 'remove' };
    const url = `${Constants.TRASA_HOSTNAME}/api/v1/groups/user/update`;
    axios
      .post(url, reqData)
      .then(() => {
        // TODO : BUG: IMP: test and fix folowing filter function. should return new array without elements that are removed
        const laddedUsersArray = addedUsersArray.filter((v, i) => users.includes(v[8]));
        setaddedUsersArray(laddedUsersArray);

        // fetchUsers(props.groupid);
      })
      .catch((error) => {
        console.log(error);
      });
  };

  useEffect(() => {
    fetchUsers(props.groupid);
  }, [props.groupid]);

  const getRoute = (staticVal: any, dynamicVal: any) => {
    return staticVal + dynamicVal;
  };

  return (
    <Layout>
      <Headers
        pageName={[
          { name: 'User Group', route: '/users/groups' },
          {
            name: groupMeta.groupName,
            route: getRoute('/users/groups/group/', props.groupid),
          },
        ]}
        tabHeaders={['Overview']}
        Components={[
          <GroupUsersTableWrapper
            groupID={props.groupid}
            groupMeta={groupMeta}
            allAddedUsersObj={allAddedUsersObj}
            addedUsersArray={addedUsersArray}
            openaddUsersDlg={openaddUsersDlg}
            open={addUsersDlg}
            onserviceselected={onserviceselected}
            addedUsersIndex={addedUsersIndex}
            unaddedUsers={unaddedUsers}
            deleteUsers={deleteUsers}
          />,
        ]}
      />
    </Layout>
  );
}

function GroupUsersTableWrapper(props: any) {
  const classes = useStyles();
  const [loader, setLoader] = useState(false);
  // const [submit, setsubmit] = useState(false);
  const [open, setopen] = useState(false);
  const [editGroupDlgState, seteditGroupDlgState] = useState(false);
  const [groupDeleteDlgState, setgroupDeleteDlgState] = useState(false);

  const changeAddUserDlgState = () => {
    setopen(!open);
  };

  const changeEditGroupDlgState = () => {
    seteditGroupDlgState(!editGroupDlgState);
  };

  const changeGroupDeleteDlgState = () => {
    setgroupDeleteDlgState(!groupDeleteDlgState);
  };

  const DeleteUserGroup = () => {
    setLoader(true);
    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/groups/delete/${props.groupMeta.groupID}`)
      .then((r) => {
        setLoader(false);
        if (r.data.status === 'success') {
          changeGroupDeleteDlgState();
          window.location.href = '/users/groups';
        }
      })
      .catch((error) => {
        setLoader(false);
        console.log(error);
      });
  };

  // const deleteUsersFromGroup = (rowsDeleted: any) => {
  //   let users = rowsDeleted.data.map((v: any) => {
  //     return addedUsersArray[v.index][8];
  //   });

  //   axios
  //     .post(`${Constants.TRASA_HOSTNAME}/api/v1/groups/user/update/${props.groupMeta.groupID}`)
  //     .then((response) => {  })
  //     .catch((error) => {
  //       console.log(error);
  //     });
  // };

  return (
    <div>
      <Paper className={classes.paper}>
        {/* <AppBar className={classes.searchBar} position="static" color="default" elevation={0}> */}
        <Toolbar>
          <Grid container>
            <Grid item xs={9}>
              <Tooltip title="Add users">
                <Button variant="contained" onClick={changeAddUserDlgState}>
                  <AddIcon />
                  Add Users
                </Button>
              </Tooltip>
            </Grid>

            <Grid item xs={1}>
              <Tooltip title="edit">
                <Button variant="contained" size="small" onClick={changeEditGroupDlgState}>
                  <EditIcon />
                  Edit
                </Button>
              </Tooltip>
            </Grid>

            <Grid item xs={2}>
              <Tooltip title="delete">
                <Button variant="contained" size="small" onClick={changeGroupDeleteDlgState}>
                  <DeleteIcon />
                  Delete Group
                </Button>
              </Tooltip>
            </Grid>
          </Grid>
        </Toolbar>
        {/* </AppBar> */}
        <div>
          <GroupUsersTable
            allUsers={props.addedUsersArray}
            isGroup
            groupID={props.groupID}
            deleteUsers={props.deleteUsers}
          />
        </div>
      </Paper>
      <br />
      <div />

      <AddUsersToGroup
        usersThatCanBeAdded={props.unaddedUsers}
        groupID={props.groupID}
        open={open}
        handleClose={changeAddUserDlgState}
      />

      <UpdateGroup
        open={editGroupDlgState}
        handleClose={changeEditGroupDlgState}
        allUsers={props.allUsers}
        update
        groupMeta={props.groupMeta}
      />

      <div>
        <Dialog
          open={groupDeleteDlgState}
          onClose={changeGroupDeleteDlgState}
          aria-labelledby="alert-dialog-title"
          aria-describedby="alert-dialog-description"
        >
          <DialogTitle id="alert-dialog-title">
            {' '}
            <div className={classes.Warning}> !!! WARNING !!! </div>
          </DialogTitle>
          <DialogContent>
            <DialogContentText id="alert-dialog-description">
              Deleting this group will cause all group users access to services removed. User may be
              locked out of system. Make sure you know what you are doing.
            </DialogContentText>
          </DialogContent>
          <DialogActions>
            <Button
              onClick={() => {
                DeleteUserGroup();
              }}
              className={classes.WarningButton}
            >
              Yes, Delete this group.
            </Button>
            <Button
              onClick={changeGroupDeleteDlgState}
              variant="contained"
              color="primary"
              autoFocus
            >
              No
            </Button>
          </DialogActions>
          {loader ? <ProgressBar /> : null}
          <br />
        </Dialog>
      </div>
    </div>
  );
}

type TParams = { groupid: string };

const GroupPageWithProps = ({ match }: RouteComponentProps<TParams>) => {
  return (
    <div>
      <UserGroupPage groupid={match.params.groupid} />
    </div>
  );
};

const UserGroupPageRoute = () => {
  return (
    <Switch>
      <Route exact path="/users/groups/group/:groupid" component={GroupPageWithProps} />
    </Switch>
  );
};

export default withRouter(UserGroupPageRoute);
