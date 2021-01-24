import { makeStyles, MuiThemeProvider } from '@material-ui/core';
// import Cities from "./cities";
import AppBar from '@material-ui/core/AppBar';
import Button from '@material-ui/core/Button';
import FormControl from '@material-ui/core/FormControl';
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import InputLabel from '@material-ui/core/InputLabel';
import LoadingBar from '@material-ui/core/LinearProgress';
import MenuItem from '@material-ui/core/MenuItem';
import Paper from '@material-ui/core/Paper';
import Select from '@material-ui/core/Select';
import TextField from '@material-ui/core/TextField';
import Toolbar from '@material-ui/core/Toolbar';
import Tooltip from '@material-ui/core/Tooltip';
import LeftIcon from '@material-ui/icons/ChevronLeft';
import RightIcon from '@material-ui/icons/ChevronRight';
import axios from 'axios';
import Moment from 'moment-timezone';
import MUIDataTable, {
  MUIDataTableColumn,
  MUIDataTableMeta,
  MUIDataTableOptions,
} from 'mui-datatables';
import React, { useEffect, useState } from 'react';
import { Route, Switch, withRouter } from 'react-router-dom';
import Constants from '../../../Constants';
import { MuiDataTableTheme } from '../../../utils/styles/themes';
import HTTPSession from './RecordedSession';

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
    marginTop: 50,
    // overflow: 'hidden',
  },
  formControl: {
    margin: theme.spacing(1),
    minWidth: 120,
  },
  selectEmpty: {
    marginTop: theme.spacing(2),
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
  success: {
    paddingLeft: 5,
    paddingRight: 0,
    maxWidth: 50,
    background: 'green',
    color: 'white',
  },
  failed: {
    paddingLeft: 10,
    // paddingRight: 5,
    maxWidth: 50,
    background: 'maroon',
    color: 'white',
  },
  redDot: {
    maxWidth: 15,
    marginLeft: 40,
    paddingLeft: 5,
    backgroundColor: 'tomato',
    color: 'white',
  },
  viewButton: {
    marginLeft: 20,
  },
  privDiv: {
    padding: theme.spacing(),
    borderRadius: '30px',
    // maxWidth: 100,
    background: 'gainsboro',
    color: 'black',
    textAlign: 'center',
  },
  textField: {
    borderRadius: 4,
    backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    color: 'black',
    fontSize: '12px',
    fontFamily: 'Open Sans, Rajdhani',
    padding: 'theme.spacing(2)',
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
    },
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

type logtableProps = {
  entityType: string;
  entityID: string;
};

export function LogTableV2(props: logtableProps) {
  const classes = useStyles();
  const [eventData, setEventData] = useState([]);
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(50);
  const [loading, setLoading] = useState(true);
  const [dateFrom, setDateFrom] = useState('');
  const [dateTo, setDateTo] = useState('');

  const custumFooter = (
    count: any,
    page: any,
    rowsPerPage: any,
    changeRowsPerPage: any,
    changePage: any,
  ) => {
    return (
      <div>
        {loading ? <LoadingBar /> : null}
        <Grid container>
          <Grid item lg={6}>
            <FormControl className={classes.formControl}>
              <InputLabel htmlFor="age-simple">Rows Per Page</InputLabel>
              <Select
                value={rowsPerPage}
                onChange={onNumberOfRowsChange(changeRowsPerPage)}
                inputProps={{
                  name: 'rowsPerPage',
                  id: 'rowsPerPage',
                }}
              >
                <MenuItem value={10}>10</MenuItem>
                <MenuItem value={20}>20</MenuItem>
                <MenuItem value={30}>30</MenuItem>
                <MenuItem value={50}>50</MenuItem>
                <MenuItem value={100}>100</MenuItem>
              </Select>
            </FormControl>
          </Grid>
          <Grid item lg={6}>
            <IconButton aria-label="Delete" onClick={onPageLeft}>
              <LeftIcon />
            </IconButton>
            {page + 1}
            <IconButton aria-label="Delete" onClick={onPageRight}>
              <RightIcon />
            </IconButton>
          </Grid>
        </Grid>
      </div>
    );
  };

  const options = {
    filter: true,
    rowsPerPage,
    filterType: 'textField',
    customFooter: custumFooter,
  };

  useEffect(() => {
    const reqPath = `${Constants.TRASA_HOSTNAME}/api/v1/logs/auth/${props.entityType}/${props.entityID}`;

    axios
      .get(reqPath)
      .then((response) => {
        // console.log(response.data);
        // userDataMain = response.data;
        let dataArr = [];
        let data = [];
        if (!response.data.data) {
          data = [];
        } else {
          data = response.data.data[0];
        }

        dataArr = data.map(function (n: any) {
          return [
            n.eventID,
            n.email,
            n.privilege,
            !n.serviceName ? 'Dashboard' : n.serviceName,
            !n.serviceType ? 'Dashboard' : n.serviceType,
            n.userIP,
            n.userAgent,
            n.status,
            n.failedReason,
            n.loginTime,
            n.logoutTime,
            n.guests,
            n.sessionID,
            n.orgID,
            n.sessionRecord,
          ];
        });

        setEventData(dataArr);
        setLoading(false);
      })
      .catch((error) => {
        if (error.response) {
          console.log(error.response.data);
          // commented out for local debug   window.location.href = '/login'
        } else {
          // Something happened in setting up the request that triggered an Error
          setLoading(false);
          console.log('Error-###', error);
        }
      });
  }, [props.entityType, props.entityID]);

  const getPaginatedLogs = (page: any, size: any) => {
    setLoading(true);

    let reqPath = `${Constants.TRASA_HOSTNAME}/api/v1/logs/auth/${props.entityType}/${
      props.entityID
    }/${page * size}/${size}`;
    if (dateFrom && dateTo) {
      reqPath = `${reqPath}/${dateFrom}/${dateTo}`;
    }

    axios
      .get(reqPath || '')
      .then((response) => {
        //        console.log(response.data);
        // userDataMain = response.data;
        let dataArr = [];
        const { data } = response;

        if (!data.data) {
          setEventData([]);
        } else {
          dataArr = data.data[0].map(function (n: any) {
            return [
              n.eventID,
              n.email,
              n.privilege,
              !n.serviceName ? 'Dashboard' : n.serviceName,
              !n.serviceType ? 'Dashboard' : n.serviceType,
              n.userIP,
              n.userAgent,
              n.status,
              n.failedReason,
              n.loginTime,
              n.logoutTime,
              n.guests,
              n.sessionID,
              n.orgID,
              n.sessionRecord,
            ];
          });
        }

        setEventData(dataArr);
        setPage(page);
        setLoading(false);
      })
      .catch((error) => {
        if (error.response) {
          console.log(error.response.data);
          // commented out for local debug   window.location.href = '/login'
        } else {
          // Something happened in setting up the request that triggered an Error
          setLoading(false);
          console.log('Error-###', error);
        }
      });
  };

  const onPageLeft = () => {
    let lpage = page;
    lpage = page > 0 ? page - 1 : page;

    getPaginatedLogs(lpage, rowsPerPage);
  };
  const onPageRight = () => {
    let lpage = page;
    lpage = page + 1;

    getPaginatedLogs(lpage, rowsPerPage);
  };
  const onNumberOfRowsChange = (changeRowsPerPage: any) => (e: any) => {
    const numRows = e.target.value;

    getPaginatedLogs(0, numRows);
    setRowsPerPage(numRows);
    changeRowsPerPage(numRows);
  };

  // const onTableChange = (action: any, tableState: any) => {
  //   console.log(action, tableState);
  // };

  const onDateChange = (e: any) => {
    // console.log(e.target.id,e.target.value)

    switch (e.target.id) {
      case 'dateFrom':
        setDateFrom(e.target.value);
        break;
      case 'dateTo':
        setDateTo(e.target.value);
        break;
      default:
        break;
    }
    // setState({[e.target.id]:e.target.value || undefined})
  };

  function statusDiv(val: boolean) {
    if (val) {
      return <div className={classes.success}>success</div>;
    }
    return <div className={classes.failed}>failed</div>;
  }

  const columns = [
    {
      name: 'Event ID',
      options: {
        filter: true,
        display: false,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'User',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'Privilege',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return <div className={classes.privDiv}>{value}</div>;
        },
      },
    },
    {
      name: 'Service Name',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'Service Type',
      options: {
        filter: false,
        display: false,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'IP address',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'User Agent',
      options: {
        filter: false,
        display: false,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'Status',
      options: {
        filter: true,
        filterOptions: [false, true],
        customBodyRender: (value: any) => {
          return statusDiv(value);
        },
      },
    },
    {
      name: 'Failed Reason',
      options: {
        filter: true,
        filterOptions: [false, true],
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },

    {
      name: 'Logged In Time',
      options: {
        filter: false,
        customBodyRender: (value: any) => {
          const d = Moment.unix(value / 1000000000);
          return <div style={tableBodyFont}> {d.format()} </div>;
        },
      },
    },
    {
      name: 'Logged Out Time',
      options: {
        filter: false,
        display: false,
        customBodyRender: (value: any) => {
          const d = Moment.unix(value / 1000000000);
          return <div style={tableBodyFont}> {d.format()} </div>;
        },
      },
    },
    {
      name: 'Guests Joined',
      options: {
        filter: false,
        display: false,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'Recorded Session',
      options: {
        filter: false,
        customBodyRender(
          value: any,
          tableMeta: MUIDataTableMeta,
          updateValue: (value: string) => void,
        ) {
          // let d=Moment(tableMeta.rowData[8],"YYYY-MM-DDTHH:mm:ssZ")
          const d = Moment.unix(tableMeta.rowData[9] / 1000000000).tz('UTC');
          //  d.month(d.month()+1)
          const month = d.month() + 1;
          return (
            <div>
              {tableMeta.rowData[14] ? (
                <Button
                  className={classes.viewButton}
                  variant="outlined"
                  color="secondary"
                  onClick={() => {
                    //  console.log(tableMeta.rowData[4])

                    // window.location.href="/overview/session/"+sessionID
                    if (tableMeta.rowData[4] === 'ssh') {
                      window.open(
                        `/monitor/sessions/view#type=ssh&year=${d.year()}&month=${month}&day=${d.date()}&sessionID=${value}`,
                        '_blank',
                      );
                    } else if (tableMeta.rowData[4] === 'db') {
                      window.open(
                        `/monitor/sessions/view#type=db&year=${d.year()}&month=${month}&day=${d.date()}&sessionID=${value}`,
                        '_blank',
                      );
                    } else if (tableMeta.rowData[4] === 'guac-ssh') {
                      window.open(
                        `/monitor/sessions/view#type=guac-ssh&year=${d.year()}&month=${month}&day=${d.date()}&sessionID=${value}`,
                        '_blank',
                      );
                    } else if (
                      tableMeta.rowData[4] === 'guac-rdp' ||
                      tableMeta.rowData[4] === 'rdp'
                    ) {
                      window.open(
                        `/monitor/sessions/view#type=guac&year=${d.year()}&month=${month}&day=${d.date()}&sessionID=${value}`,
                        '_blank',
                      );
                    } else if (tableMeta.rowData[4] === 'guac-vnc') {
                      window.open(
                        `/monitor/sessions/view#type=guac&year=${d.year()}&month=${month}&day=${d.date()}&sessionID=${value}`,
                        '_blank',
                      );
                    } else if (tableMeta.rowData[4] === 'guac-vnc') {
                      window.open(
                        `/monitor/sessions/view#type=guac&year=${d.year()}&month=${month}&day=${d.date()}&sessionID=${value}`,
                        '_blank',
                      );
                    } else if (tableMeta.rowData[4] === 'http') {
                      window.open(
                        `/monitor/sessions/view#type=http&year=${d.year()}&month=${month}&day=${d.date()}&sessionID=${value}`,
                        '_blank',
                      );
                    }
                  }}
                >
                  View
                </Button>
              ) : (
                <div className={classes.redDot}> - </div>
              )}
            </div>
          );
        },
      },
    },
    {
      name: 'orgID',
      options: {
        filter: false,
        display: false,
        searchable: false,
        print: false,
        download: false,
      },
    },
    {
      name: 'Session Recorded',
      options: {
        filter: false,
        display: false,
        searchable: false,
        print: false,
        download: false,
      },
    },
  ];

  return (
    <div>
      {/* <Content> */}
      <Paper className={classes.paper}>
        <AppBar className={classes.searchBar} position="static" color="default" elevation={0}>
          <Toolbar>
            <Grid container spacing={2} alignItems="center">
              <Grid item />
              <Grid item xs>
                <div style={tableBodyFont}>
                  {' '}
                  Tip - You can search and sort log data using search bar or filter function{' '}
                </div>
              </Grid>
              <Grid item>
                <TextField
                  id="dateFrom"
                  label="From"
                  type="date"
                  value={dateFrom}
                  // defaultValue="2017-05-24"
                  onChange={onDateChange}
                  className={classes.textField}
                  InputLabelProps={{
                    shrink: true,
                  }}
                />
              </Grid>

              <Grid item>
                <TextField
                  id="dateTo"
                  label="To"
                  type="date"
                  value={dateTo}
                  // defaultValue="2017-05-24"
                  onChange={onDateChange}
                  className={classes.textField}
                  InputLabelProps={{
                    shrink: true,
                  }}
                />
              </Grid>

              <Grid item>
                <Button
                  variant="contained"
                  size="small"
                  color="secondary"
                  onClick={() => {
                    getPaginatedLogs(0, rowsPerPage);
                  }}
                >
                  Filter By Date
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
              title="Authentication Event Logs"
              data={eventData}
              columns={columns as MUIDataTableColumn[]}
              options={options as MUIDataTableOptions}
            />
          </MuiThemeProvider>
        </div>
      </Paper>
      {/* </Content> */}
    </div>
  );
}

const SessionLog = (props: any) => {
  return (
    <Switch>
      <Route
        exact
        path="/monitor/sessions"
        render={(routeProps) => <LogTableV2 {...routeProps} {...props} />}
      />
      {/* <Route path="/monitor/sessions/history/:type/:orgID/:sessionID" component={HTTPSession} />
      <Route path="/monitor/sessions/history/:type/:orgID/:sessionID" component={HTTPSession} />
      <Route path="/monitor/sessions/history/:type/:orgID/:sessionID" component={HTTPSession} /> */}
    </Switch>
  );
};

// export default withStyles(styles)(Overview)
export default withRouter(SessionLog);

// {() => {
//   //  console.log(tableMeta.rowData[4])

//   // window.location.href="/overview/session/"+sessionID
//   if (tableMeta.rowData[4] === 'ssh') {
//     window.open(
//       `/monitor/sessions/view#type=ssh&year=${d.year()}&month=${month}&day=${d.date()}&sessionID=${value}`,
//       '_blank'
//     );
//   } else if (tableMeta.rowData[4] === 'db') {
//     window.open(
//       `/monitor/sessions/view#type=db&year=${d.year()}&month=${month}&day=${d.date()}&sessionID=${value}`,
//       '_blank'
//     );
//   } else if (tableMeta.rowData[4] === 'guac-ssh') {
//     window.open(
//       `/monitor/sessions/view#type=guac-ssh&year=${d.year()}&month=${month}&day=${d.date()}&sessionID=${value}`,
//       '_blank'
//     );
//   } else if (
//     tableMeta.rowData[4] === 'guac-rdp' ||
//     tableMeta.rowData[4] === 'rdp'
//   ) {
//     window.open(
//       `/monitor/sessions/view#type=guac&year=${d.year()}&month=${month}&day=${d.date()}&sessionID=${value}`,
//       '_blank'
//     );
//   } else if (tableMeta.rowData[4] === 'guac-vnc') {
//     window.open(
//       `/monitor/sessions/view#type=guac&year=${d.year()}&month=${month}&day=${d.date()}&sessionID=${value}`,
//       '_blank'
//     );
//   } else if (tableMeta.rowData[4] === 'guac-vnc') {
//     window.open(
//       `/monitor/sessions/view#type=guac&year=${d.year()}&month=${month}&day=${d.date()}&sessionID=${value}`,
//       '_blank'
//     );
//   } else if (tableMeta.rowData[4] === 'http') {
//     window.open(
//       `/monitor/sessions/view#type=http&year=${d.year()}&month=${month}&day=${d.date()}&sessionID=${value}`,
//       '_blank'
//     );
//   }
// }}
