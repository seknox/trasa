import { createMuiTheme, makeStyles, MuiThemeProvider } from '@material-ui/core';
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
import Moment from 'moment';
import MUIDataTable, { MUIDataTableColumn } from 'mui-datatables';
import React from 'react';
import Constants from '../../../Constants';

const theme = createMuiTheme({
  typography: { fontFamily: 'Open Sans, Rajdhani' },
  palette: {
    type: 'light',
    primary: { 500: '#000080' },
    secondary: { A400: '#000080' }, // '#000080' },
  },
});

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

  svg: {
    width: 100,
    height: 100,
  },
  polygon: {
    fill: theme.palette.common.white,
    stroke: theme.palette.divider,
    strokeWidth: 1,
  },
  textField: {
    fontSize: 14,
    fontFamily: 'Open Sans, Rajdhani',
    width: 'calc(100% - 4px)',
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
  fontSize: '14px',
  fontFamily: 'Open Sans, Rajdhani',
};

export default function InAppTrailsTable() {
  const classes = useStyles();

  const [eventData, seteventData] = React.useState([]);
  const [loader, setLoader] = React.useState(false);
  const [page, setpage] = React.useState(0);
  const [rowsPerPage, setrowsPerPage] = React.useState(50);
  const [date, setDate] = React.useState({ dateFrom: '', dateTo: '' });

  React.useEffect(() => {
    const reqPath = `${Constants.TRASA_HOSTNAME}/api/v1/logs/inapptrails/org`;

    axios
      .get(reqPath)
      .then((response) => {
        let dataArr = [];
        const data = response.data.data[0];
        dataArr = data.map(function (n: any) {

          return [
            n.eventID,
            n.email,
            n.userAgent,
            n.clientIP,
            n.status.toString(),
            n.eventTime,
            n.description,
          ];
        });

        seteventData(dataArr);
      })
      .catch((error) => {
        console.log(error);
      });
  }, []);

  const getPaginatedLogs = (page: any, size: any) => {
    setLoader(true);

    const reqPath = `${Constants.TRASA_HOSTNAME}/api/v1/logs/inapptrails/org/${
      page * size
    }/${size}/${date.dateFrom}/${date.dateTo}`;

    axios
      .get(reqPath || '')
      .then((response) => {
        let dataArr = [];
        const data = response.data.data[0];
        dataArr = data.map(function (n: any) {

          return [
            n.eventID,
            n.email,
            n.userAgent,
            n.clientIP,
            n.status.toString(),
            n.eventTime,
            n.description,
          ];
        });

        seteventData(dataArr);
        setLoader(false);
        setpage(page);
      })
      .catch((error) => {
        setLoader(false);
        console.log(error);
      });
  };

  const onPageLeft = () => {
    let lpage = page;
    lpage = lpage > 0 ? lpage - 1 : lpage;

    getPaginatedLogs(lpage, rowsPerPage);
  };

  const onPageRight = () => {
    let lpage = page;
    lpage += 1;

    getPaginatedLogs(lpage, rowsPerPage);
  };

  const onNumberOfRowsChange = (changeRowsPerPage: any) => (e: any) => {
    const numRows = e.target.value;

    getPaginatedLogs(0, numRows);
    setrowsPerPage(numRows);
    changeRowsPerPage(numRows);
  };

  const custumFooter = (rowsPerPage: any, changeRowsPerPage: any, changePage: any) => {
    return (
      <div>
        {loader ? <LoadingBar /> : null}
        <Grid container>
          <Grid item lg={6}>
            <FormControl>
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

  const onDateChange = (e: any) => {
    setDate({ ...date, [e.target.id]: e.target.value });
  };

  const columns = [
    {
      name: 'Event ID',
      options: {
        // filter: true,
        display: true,
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
      name: 'User Agent',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'User IP',
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
          return value;
        },
      },
    },
    {
      name: 'Event Time',
      options: {
        filter: false,
        customBodyRender: (value: any) => {
          const d = Moment.unix(value / 1000000000);
          return d.format();
        },
      },
    },
    {
      name: 'Details',
      options: {
        filter: false,
        customBodyRender: (value: any) => {
          return value;
        },
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
                  value={date.dateFrom}
                  // defaultValue="2017-05-24"
                  onChange={onDateChange}
                  className={classes.textField}
                />
              </Grid>

              <Grid item>
                <TextField
                  id="dateTo"
                  label="To"
                  type="date"
                  value={date.dateTo}
                  // defaultValue="2017-05-24"
                  onChange={onDateChange}
                />
              </Grid>

              <Grid item>
                <Button
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
          <MuiThemeProvider theme={theme}>
            <MUIDataTable
              title="In-app Audit Trails"
              data={eventData}
              columns={columns as MUIDataTableColumn[]}
              // options={this.options}
            />
          </MuiThemeProvider>
          ;
        </div>
      </Paper>
    </div>
  );
}
