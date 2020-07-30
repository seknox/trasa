import { createMuiTheme, MuiThemeProvider } from '@material-ui/core';
import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import { makeStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import axios from 'axios';
import MUIDataTable from 'mui-datatables';
import React, { useEffect, useState } from 'react';
import Constants from '../../../Constants';

const useStyles = makeStyles((theme) => ({
  root: {
    flexgrow: 1,
    // marginBotton: '5%',
  },

  paperTrans: {
    backgroundColor: 'transparent',
    // padding: theme.spacing(1),
    textAlign: 'center',
  },
  paperTrans1: {
    backgroundColor: 'transparent',
  },
  paper: {
    backgroundColor: '#fdfdfd',
    padding: theme.spacing(2),
    //  textAlign: 'center',
    color: theme.palette.text.secondary,
  },

  aggHeadersBig: {
    color: '#000080', // '#1b1b32',
    fontSize: '30px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  aggHeaders: {
    color: '#1b1b32',
    fontSize: '16px',
    fontFamily: 'Open Sans, Rajdhani',
  },
}));

export default function AssignedGroups(props: any) {
  const [userGroups, setUserGroups] = useState([]);

  useEffect(() => {
    const config = {
      headers: {
        'X-SESSION': localStorage.getItem('X-SESSION'),
        'X-CSRF': localStorage.getItem('X-CSRF'),
      },
    };
    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/user/assignedgroups/${props.userID}`, config)
      .then((resp) => {
        if (resp.data.status == 'success' && resp.data.data) setUserGroups(resp.data.data[0]);
      });
  }, []);

  const classes = useStyles();

  return (
    <div className={classes.root}>
      <Paper>
        <Grid container spacing={2}>
          <Grid item xs={12}>
            <Typography variant="h3"> Groups</Typography>
            <Divider light />
          </Grid>

          <Grid item xs={2}>
            <Paper className={classes.paperTrans} elevation={0}>
              <div className={classes.aggHeaders}>
                {' '}
                <b> Total </b>
              </div>
              <div className={classes.aggHeadersBig}>
                {' '}
                <b> 4 </b>
              </div>
            </Paper>
          </Grid>

          <Grid item xs={12}>
            <AssignedGroupTable userGroups={userGroups} />
          </Grid>

          <Grid item xs={12}>
            <Divider light />
          </Grid>
        </Grid>
      </Paper>
    </div>
  );
}

const theme = createMuiTheme({
  typography: { fontFamily: 'Open Sans, Rajdhani' },
  palette: {
    type: 'light',
    primary: { 500: '#000080' },
    secondary: { A400: '#000080' }, // '#000080' },
  },
});

const tableHeader = {
  fontColor: 'black',
  fontWeight: 'bold',
  fontSize: '17px',
  fontFamily: 'Open Sans, Rajdhani',
};

const tableBodyFont = {
  backgroundColor: 'Transparent',
  borderRadius: 3,
  border: 0,
  fontColor: 'dimgray',
  fontSize: '12px',
  fontFamily: 'Open Sans, Rajdhani',
};

const tableBodyheader = {
  backgroundColor: 'Transparent',
  fontWeight: 'bold',
  border: 0,
  color: 'dimgray',
  fontSize: '14px',
  fontFamily: 'Open Sans, Rajdhani',
};

function AssignedGroupTable(props: any) {
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
      name: 'Assigned On',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'View Details',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return (
            <div style={tableBodyFont}>
              {' '}
              <a href={`/users/groups/group/${value}`}>View Details</a>{' '}
            </div>
          );
        },
      },
    },
  ];

  return (
    <div>
      <MuiThemeProvider theme={theme}>
        <MUIDataTable
          title="Groups"
          data={props.userGroups.map(function (n: any) {
            return [n.groupName, new Date(n.createdAt * 1000).toLocaleDateString(), n.groupID];
          })}
          columns={columns}
          options={{
            filter: true,
            responsive: 'scrollMaxHeight',
            selectableRows: 'none',
            // resizableColumns: true,
            // onRowsDelete: this.handleDeletePermission
          }}
        />
      </MuiThemeProvider>
    </div>
  );
}
