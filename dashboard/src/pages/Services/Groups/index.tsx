import { createMuiTheme, MuiThemeProvider } from '@material-ui/core';
import AppBar from '@material-ui/core/AppBar';
import Button from '@material-ui/core/Button';
import Drawer from '@material-ui/core/Drawer';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import { makeStyles } from '@material-ui/core/styles';
import Toolbar from '@material-ui/core/Toolbar';
import AuthServicegroupIcon from '@material-ui/icons/DeviceHubSharp';
import axios from 'axios';
import MUIDataTable, { MUIDataTableOptions } from 'mui-datatables';
import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import Constants from '../../../Constants';
import Servicesetting from '../Service/Settings/ServiceSetting';
import CreateServicegroup from './CreateServiceGroup';

const theme = createMuiTheme({
  typography: { fontFamily: 'Open Sans, Rajdhani' },
  palette: {
    type: 'light',
    primary: { 500: '#000080' },
    secondary: { A400: '#000080' }, // '#000080' },
  },
});

const useStyles = makeStyles((th) => ({
  paper: {
    backgroundColor: '#fdfdfd',
    padding: th.spacing(2),
    //  textAlign: 'center',
    color: th.palette.text.secondary,
  },
}));

type Anchor = 'top' | 'left' | 'bottom' | 'right';

export default function Servicegroup() {
  const [creatUserGroupDlgState, setCreatUserGroupDlg] = useState(false);
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

  const [groups, setgroups] = useState([]);

  useEffect(() => {
    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/groups/service`)
      .then((response) => {
        let dataArr = [];
        const data = response.data.data[0];

        dataArr = data.map(function (n: any) {
          const cdate = new Date(n.createdAt * 1000);
          const udate = new Date(n.updatedAt * 1000);
          return [
            n.groupName,
            n.memberCount,
            n.status ? 'Active' : 'Disabled',
            cdate.toDateString(),
            udate.toDateString(),
            n.groupID,
          ];
        });
        setgroups(dataArr);
      })
      .catch((error) => {
        console.log(error);
      });
  }, []);

  const handleCreateGroupDlgState = () => {
    setCreatUserGroupDlg(!creatUserGroupDlgState);
  };

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
      name: ' View Details',
      options: {
        filter: false,
        customBodyRender: (value: any) => {
          return (
            <Button
              component={Link}
              to={`/services/groups/group/${value}`}
              variant="outlined"
              color="secondary"
            >
              View / Edit
            </Button>
          );
        },
      },
    },
  ];

  const options = {
    filter: true,
    responsive: 'scrollMaxHeight',
    selectableRows: 'none',
  };

  const [configDrawerState, setConfigDrawerState] = useState({ right: false });
  const classes = useStyles();

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

  return (
    <div>
      <CreateServicegroup
        open={creatUserGroupDlgState}
        handleClose={handleCreateGroupDlgState}
        update={false}
      />

      <Paper className={classes.paper}>
        <AppBar position="static" color="default" elevation={0}>
          <Toolbar>
            <Grid container spacing={2}>
              <Grid item xs={2}>
                <Button variant="contained" size="small" onClick={handleCreateGroupDlgState}>
                  <AuthServicegroupIcon fontSize="small" />
                  Create group
                </Button>
              </Grid>

              {/* <Grid item xs={2}>
                <Button
                  variant="contained"
                  size="small"
                  onClick={toggleConfigDrawer('right', true)}
                >
                  Create new Service
                </Button>
              </Grid> */}
            </Grid>
          </Toolbar>
        </AppBar>
        <div>
          <MuiThemeProvider theme={theme}>
            <MUIDataTable
              title="Service Group"
              data={groups}
              columns={columns}
              options={options as MUIDataTableOptions}
            />
          </MuiThemeProvider>
        </div>
      </Paper>
      <Drawer
        anchor="right"
        open={configDrawerState.right}
        onClose={toggleConfigDrawer('right', false)}
      >
        <Paper className={classes.paper}>
          <Servicesetting newApp serviceDetail={serviceDetail} />
        </Paper>
      </Drawer>
    </div>
  );
}
