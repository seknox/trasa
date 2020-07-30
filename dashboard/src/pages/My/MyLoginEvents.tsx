// import Cities from "./cities";
import AppBar from '@material-ui/core/AppBar';
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import Paper from '@material-ui/core/Paper';
import { makeStyles } from '@material-ui/core/styles';
import Toolbar from '@material-ui/core/Toolbar';
import Tooltip from '@material-ui/core/Tooltip';
import RefreshIcon from '@material-ui/icons/Refresh';
import axios from 'axios';
import MUIDataTable, { MUIDataTableColumn, MUIDataTableOptions } from 'mui-datatables';
import React from 'react';
import { MuiThemeProvider, createMuiTheme } from '@material-ui/core';
import Constants from '../../Constants';

const theme = createMuiTheme({
  typography: { fontFamily: 'Open Sans, Rajdhani' },
  palette: {
    type: 'light',
    primary: { 500: '#000080' },
    secondary: { A400: '#000080' }, // '#000080' },
  },
});

const lightColor = 'rgba(255, 255, 255, 0.7)'; // 'rgba(255, 255, 255, 0.7)'; // '#030417';

const useStyles = makeStyles((theme) => ({
  mainContent: {
    flex: 1,
    padding: '48px 36px 0',
    background: '#eaeff1', // '#eaeff1',
  },
  paper: {
    maxWidth: 1500,
    margin: 'auto',
    marginTop: 100,
    overflow: 'hidden',
  },
  searchBar: {
    borderBottom: '1px solid rgba(0, 0, 0, 0.12)',
  },
  searchInput: {
    fontSize: theme.typography.fontSize,
  },
  block: {
    display: 'block',
  },
  addUser: {
    marginRight: theme.spacing(1),
  },
  contentWrapper: {
    margin: '40px 16px',
  },
  secondaryBar: {
    zIndex: 0,
  },
  button: {
    borderColor: lightColor,
  },
  svg: {
    width: 100,
    height: 100,
  },
  polygon: {
    fill: theme.palette.common.white,
    stroke: theme.palette.divider,
    strokeWidth: 1,
  },
}));

const tableBodyFont = {
  backgroundColor: 'Transparent',
  borderRadius: 3,
  border: 0,
  fontColor: 'dimgray',
  fontSize: '12px',
  fontFamily: 'Open Sans, Rajdhani',
};

export default function GrouopTable() {
  const classes = useStyles();
  const [eventData, seteventData] = React.useState([]);

  React.useEffect(() => {
    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/my/auth/log`)
      .then((response) => {
        let dataArr = [];
        const { data } = response;
        if(!data.data){
          return
        }
        dataArr = data.data[0].map(function (n: any) {
          return [
            n.eventID,
            n.privilege,
            !n.serviceName ? 'Dashboard' : n.serviceName,
            !n.serviceType ? 'Dashboard' : n.serviceType,
            n.userIP,
            n.userAgent,
            n.status.toString(),
            n.loginTime,
            n.logoutTime,
          ];
        });
        seteventData(dataArr);
      })
      .catch((error) => {
        console.log(error);
      });
  }, []);

  const columns = [
    {
      name: 'EventID',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return <div style={tableBodyFont}>{value}</div>;
        },
      },
    },
    {
      name: 'User',
      options: {
        filter: true,
        display: false,
        customBodyRender: (value: any) => {
          return <div style={tableBodyFont}>{value}</div>;
        },
      },
    },
    {
      name: 'Service Name',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return <div style={tableBodyFont}>{value}</div>;
        },
      },
    },
    {
      name: 'Service Type',
      options: {
        filter: true,
        display: false,
        customBodyRender: (value: any) => {
          return <div style={tableBodyFont}>{value}</div>;
        },
      },
    },
    {
      name: 'IP address',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return <div style={tableBodyFont}>{value}</div>;
        },
      },
    },

    {
      name: 'User Agent',
      options: {
        filter: true,
        display: false,
        customBodyRender: (value: any) => {
          return <div style={tableBodyFont}>{value}</div>;
        },
      },
    },
    {
      name: 'Status',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return <div style={tableBodyFont}>{value}</div>;
        },
      },
    },
    {
      name: 'Logged In Time',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return <div style={tableBodyFont}>{new Date(value / 1000000).toString()}</div>;
        },
      },
    },
    {
      name: 'Logged Out Time',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return <div style={tableBodyFont}>{new Date(value / 1000000).toString()}</div>;
        },
      },
    },
  ];

  return (
    <Paper className={classes.paper}>
      <AppBar className={classes.searchBar} position="static" color="default" elevation={0}>
        <Toolbar>
          <Grid container spacing={2} alignItems="center">
            <Grid item />
            <Grid item xs>
              Tip - You can search and sort log data using search bar or filter function
            </Grid>
            <Grid item>
              <Tooltip title="Reload">
                <IconButton>
                  <RefreshIcon className={classes.block} color="inherit" />
                </IconButton>
              </Tooltip>
            </Grid>
          </Grid>
        </Toolbar>
      </AppBar>
      <div className={classes.contentWrapper}>
        <MuiThemeProvider theme={theme}>
          <MUIDataTable
            title="Your Authentication Event"
            data={eventData}
            columns={columns as MUIDataTableColumn[]}
            options={options as MUIDataTableOptions}
          />
        </MuiThemeProvider>
      </div>
    </Paper>
  );
}

const options = {
  filter: true,
  responsive: 'scrollMaxHeight',
};
