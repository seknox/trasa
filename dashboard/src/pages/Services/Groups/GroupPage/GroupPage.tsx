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
import React, { useState } from 'react';
import UpdateGroup from '../CreateServiceGroup';
import AddservicesToGroup from './AddServicesToGroup';
import GroupServiceTable from './GroupTable';

const useStyles = makeStyles((theme) => ({
  root: {
    flexGrow: 1,
    marginTop: 20,
    justifyContent: 'space-between',
  },
  //   paper: {
  //     height: 200,
  //     maxWidth: 200,
  //   },
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
  textFieldInput: {
    borderRadius: 4,
    backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    fontSize: 46,
    padding: '5px 9px',
    marginLeft: 100,
    width: 'calc(100% - 200px)',
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
    // padding: '10px 100px',
    // width: 'calc(100% - 4px)',
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
    },
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
}));

export default function GroupPage(props: any) {
  const classes = useStyles();

  const [open, setopen] = useState(false);
  const [editGroup, seteditGroup] = useState(false);
  const [groupDeleteDlgState, setgroupDeleteDlgState] = useState(false);

  const changeAddServiceDlgState = () => {
    setopen(!open);
  };

  const changeEditGroupDlgState = () => {
    seteditGroup(!editGroup);
  };

  const changeGroupDeleteDlgState = () => {
    setgroupDeleteDlgState(!groupDeleteDlgState);
  };

  const closeGroupDlg = () => {
    seteditGroup(false);
  };

  return (
    <div>
      <Paper className={classes.paper}>
        {/* <AppBar className={classes.searchBar} position="static" color="default" elevation={0}> */}
        <Toolbar>
          <Grid container>
            <Grid item xs={8}>
              <Tooltip title="Add users">
                <Button variant="contained" color="secondary" onClick={changeAddServiceDlgState}>
                  <AddIcon />
                  {/* <AddServiceIcon className={classes.buttonIcons} /> */}
                  Add Service
                </Button>
              </Tooltip>
            </Grid>

            <Grid item xs={1}>
              <Tooltip title="edit">
                <Button
                  variant="contained"
                  color="secondary"
                  size="small"
                  onClick={changeEditGroupDlgState}
                >
                  <EditIcon />
                  Edit
                </Button>
              </Tooltip>
            </Grid>

            <Grid item xs={2}>
              <Tooltip title="delete">
                <Button
                  variant="contained"
                  color="secondary"
                  size="small"
                  onClick={changeGroupDeleteDlgState}
                >
                  <DeleteIcon />
                  Delete Group
                </Button>
              </Tooltip>
            </Grid>
          </Grid>
        </Toolbar>
      </Paper>
      <br />
      <div>
        {props.dispayList ? (
          <GroupServiceTable
            allServicesArray={props.addedservicesArray}
            removeServices={props.removeServices}
          />
        ) : (
          ''
        )}
      </div>

      <AddservicesToGroup
        ServicesThatCanBeAdded={props.unaddedservices}
        groupID={props.groupID}
        open={open}
        handleClose={changeAddServiceDlgState}
      />

      <UpdateGroup
        open={editGroup}
        handleClose={closeGroupDlg}
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
                changeGroupDeleteDlgState();
                props.DeleteAuthServicegroup(props.groupID);
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
        </Dialog>
      </div>
    </div>
  );
}
