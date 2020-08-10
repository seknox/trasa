import { MuiThemeProvider } from '@material-ui/core';
import AppBar from '@material-ui/core/AppBar';
import Button from '@material-ui/core/Button';
// import UserCreateUpdateDrawer from './UserCreateUpdateDrawer';
// import DeleteConfirmDialogue from '../../../components/ui/confirms'
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import Paper from '@material-ui/core/Paper';
import { makeStyles } from '@material-ui/core/styles';
import Toolbar from '@material-ui/core/Toolbar';
import Tooltip from '@material-ui/core/Tooltip';
import GroupIcon from '@material-ui/icons/GroupAddSharp';
import MUIDataTable, { MUIDataTableMeta, MUIDataTableOptions } from 'mui-datatables';
import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { MuiDataTableTheme } from '../../../utils/styles/themes';
import CreateUserGroup from './CreateUserGroup';

const useStyles = makeStyles((theme) => ({
  paper: {
    backgroundColor: '#fdfdfd',
    padding: theme.spacing(2),
    color: theme.palette.text.secondary,
  },
  searchBar: {
    borderBottom: '1px solid rgba(0, 0, 0, 0.12)',
  },
  buttonIcons: {
    marginRight: '5px',
  },
  contentWrapper: {
    margin: '40px 16px',
  },
}));

export default function UserGroupTable(props: any) {
  const classes = useStyles();
  // const [users, setusers] = useState<any>([]);
  // const [usersObj, setusersObj] = useState<any>({});
  const [dlgOpen, setdlgOpen] = useState(false);
  // const [loader, setLoader] = useState(false);
  const [creatUserGroupDlg, setcreatUserGroupDlg] = useState(false);

  const handleDlgClose = () => {
    setdlgOpen(false);
  };

  const handleCreateGroupDlgState = () => {
    setcreatUserGroupDlg(!creatUserGroupDlg);
  };

  // const deleteUsers = (rowsDeleted: any) => {
  //   setLoader(true);
  //   const lusers = rowsDeleted.data.map((v: any) => usersObj[v.index].Id);

  //   const reqData = JSON.stringify({ lusers });

  //   axios
  //     .post(`${Constants.TRASA_HOSTNAME}/api/v1/user/delete`, reqData)
  //     .then(() => {
  //       setLoader(false);
  //     })
  //     .catch((error) => {
  //       setLoader(false);
  //       console.log(error);
  //     });
  // };

  const columns = [
    {
      name: 'Group Name',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'Members',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'Status',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },

    {
      name: 'Created At',
      options: {
        filter: false,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'Last Updated At',
      options: {
        filter: false,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'View Details',
      options: {
        filter: false,
        customBodyRender: (value: any, tableMeta: MUIDataTableMeta) => {
          return (
            <div>
              <Button
                id={tableMeta.rowData[0]}
                component={Link}
                to={`/users/groups/group/${value}`}
                variant="outlined"
                color="secondary"
              >
                View
              </Button>
            </div>
          );
        },
      },
    },
  ];

  const options = {
    filter: true,
    responsive: 'scrollMaxHeight',
    // onRowsDelete: deleteUsers,
    // isRowSelectable:()=>false,
    selectableRows: 'none',

    // onTableChange:(a,b)=>{console.log(a,b);console.log(this.state.users)},
  };

  return (
    <div>
      <CreateUserGroup
        open={creatUserGroupDlg}
        handleClose={handleCreateGroupDlgState}
        allUsers={props.allUsers}
        update={false}
      />

      <Paper className={classes.paper}>
        <AppBar className={classes.searchBar} position="static" color="default" elevation={0}>
          <Toolbar>
            <Grid container spacing={3} alignItems="center">
              <Grid item />
              <Grid item xs>
                <Button variant="contained" onClick={handleCreateGroupDlgState} id='createGroupBtn'>
                  <GroupIcon className={classes.buttonIcons} />
                  Create group
                </Button>
              </Grid>
              <Grid item>
                <Tooltip title="Reload">
                  <IconButton>
                    {/* <RefreshIcon className={classes.block} color="inherit" /> */}
                  </IconButton>
                </Tooltip>
              </Grid>
            </Grid>
          </Toolbar>
        </AppBar>
        <div className={classes.contentWrapper}>
          <MuiThemeProvider theme={MuiDataTableTheme}>
            <MUIDataTable
              title="User Groups"
              data={props.groups}
              columns={columns}
              options={options as MUIDataTableOptions}
            />
          </MuiThemeProvider>
        </div>
      </Paper>

      <Dialog
        open={dlgOpen}
        fullWidth
        onClose={handleDlgClose}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
      >
        <DialogTitle id="alert-dialog-title">Deleting Users</DialogTitle>
        <DialogContent>
          <DialogContentText id="alert-dialog-description">
            {/* {loader ? <ProgressHOC /> : ''} */}
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button variant="contained" onClick={handleDlgClose} color="primary" autoFocus>
            Close
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
}
