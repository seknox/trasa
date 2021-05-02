import { createMuiTheme, MuiThemeProvider } from '@material-ui/core';
import CircularProgress from '@material-ui/core/CircularProgress';
import green from '@material-ui/core/colors/green';
import IconButton from '@material-ui/core/IconButton';
import { makeStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import CheckIcon from '@material-ui/icons/Check';
import EditIcon from '@material-ui/icons/Edit';
import SaveIcon from '@material-ui/icons/Save';
import axios from 'axios';
import MUIDataTable, { MUIDataTableMeta } from 'mui-datatables';
import React, { useEffect, useState } from 'react';
import Constants from '../../../Constants';
// import SecurityRules from '../../system/Security/SecurityRules';

const theme = createMuiTheme({
  typography: { fontFamily: 'Open Sans, Rajdhani' },
  palette: {
    type: 'light',
    primary: { 500: '#000080' },
    secondary: { A400: '#000080' }, // '#000080' },
    // butprim: '#32C5E9'
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
    //  marginTop: 100,
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
    color: 'blue',
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
  buttonSuccess: {
    backgroundColor: green[500],
    '&:hover': {
      backgroundColor: green[700],
    },
  },
  fabProgress: {
    color: green[500],
    position: 'absolute',
    top: -6,
    left: -6,
    zIndex: 1,
  },
  buttonProgress: {
    color: green[500],
    position: 'absolute',
    top: '50%',
    left: '50%',
    marginTop: -12,
    marginLeft: -12,
  },
  selectCustom: {
    fontSize: 12,
    fontFamily: 'Open Sans, Rajdhani',
    borderRadius: 4,
    backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    height: 25,
    marginTop: 5,
    // padding: '10px 100px',
    width: 'calc(100% - 4px)',
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
    },
  },
}));

const tableHeader = {
  fontColor: 'black',
  fontWeight: 'bold',
  fontSize: '20px',
  fontFamily: 'Open Sans, Rajdhani',
};

const tableBodyFont = {
  backgroundColor: 'Transparent',
  borderRadius: 3,
  border: 0,
  fontColor: 'black',
  fontSize: '14px',
  fontFamily: 'Open Sans, Rajdhani',
};

const tableBodyBrightFont = {
  backgroundColor: 'Transparent',
  borderRadius: 3,
  border: 0,
  fontColor: 'navy',
  color: 'navy',
  fontSize: '14px',
  fontFamily: 'Open Sans, Rajdhani',
};

const tableBodyNameFont = {
  backgroundColor: 'Transparent',
  borderRadius: 3,
  border: 0,
  fontColor: 'black',
  fontSize: '12px',
  fontFamily: 'Open Sans, Rajdhani',
};

const tableBodySubFont = {
  backgroundColor: 'Transparent',
  borderRadius: 3,
  border: 0,
  fontColor: 'dimgray',
  fontSize: '12px',
  fontFamily: 'Open Sans, Rajdhani',
};

const tableBodyheader = {
  backgroundColor: 'Transparent',
  border: 0,
  fontWeight: 'bold',
  color: 'dimgray',
  fontSize: '11px',
  fontFamily: 'Open Sans, Rajdhani',
};

export default function SecurityRules() {
  const [rules, setRules] = useState([]);
  const [loader, setLoader] = useState(false);
  const [updateRow, setUpdateRow] = useState(false);
  const [updateRowIndex, setUpdateRowIndex] = useState(0);
  const [updateValue, setUpdateValue] = useState<any>();

  const classes = useStyles();

  useEffect(() => {
    getSystemSecurityRules();
  }, []);

  function getSystemSecurityRules() {
    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/system/security/rules`)
      .then((r) => {
        if (r.data.status === 'success') {
          const resp = r.data.data[0];

          const dataArr = resp.map(function (n: any) {
            const date = new Date(n.lastModified * 1000);
            return [
              n.name,
              n.status ? 'enabled' : 'disabled',
              'Alert to admins',
              n.createdBy,
              date.toDateString(),
              n.ruleID,
              n.description,
            ];
          });
          setRules(dataArr);
        }
      })
      .catch((error) => {
        console.log(error);
      });
  }

  const updateStatus = (e: any) => {
    const row = rules[updateRowIndex];
    const updateVal = { status: e.target.value, ruleID: row[5] };
    setUpdateValue(updateVal);
    // this.setState({ updateValue: updateVal });
  };

  const updateEditRowState = (tableMeta: any) => {
    setUpdateRow(true);
    setUpdateRowIndex(tableMeta.rowIndex);
  };

  const updateRuleStatus = () => {
    setLoader(true);

    const val = updateValue;
    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/system/security/rule/update`, val)
      .then((response) => {
        getSystemSecurityRules();
        setLoader(false);
        setUpdateRow(false);
        setUpdateRowIndex(0);
        // this.setState({  updateRowIndex: 0, success: false });
      })
      .catch((error) => {
        if (error.response.status === 403) {
          window.location.href = '/login';
        }
        if (error.response) {
        } else {
        }
      });
  };

  const columns = [
    {
      name: 'Name',
      options: {
        filter: true,
        customBodyRender: (value: any, tableMeta: MUIDataTableMeta) => {
          return (
            <div>
              <Typography variant="h6">{value}</Typography>
              <div style={tableBodySubFont}>{rules[tableMeta.rowIndex][6]}</div>
            </div>
          );
        },
      },
    },
    {
      name: 'Status',
      options: {
        filter: true,
        customBodyRender: (value: any, tableMeta: MUIDataTableMeta) => {
          return (
            <div>
              {updateRow && updateRowIndex === tableMeta.rowIndex ? (
                <select className={classes.selectCustom} name="expiry" onChange={updateStatus}>
                  {status.map((name, i) => (
                    <option key={i} value={name} selected={value === name}>
                      {name}
                    </option>
                  ))}
                </select>
              ) : (
                <div style={value === 'enabled' ? tableBodyBrightFont : tableBodyFont}>
                  {' '}
                  {value}{' '}
                </div>
              )}
            </div>
          );
        },
      },
    },
    {
      name: 'Action',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'Created By',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'Last Modified',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'Change Status',
      options: {
        filter: true,
        customBodyRender: (
          value: any,
          tableMeta: MUIDataTableMeta,
          updateValue: (value: string) => void,
        ) => {
          return (
            <div>
              {updateRow && updateRowIndex === tableMeta.rowIndex ? (
                <div>
                  <IconButton color="primary" onClick={updateRuleStatus}>
                    {loader && <CircularProgress size={58} className={classes.fabProgress} />}
                    {loader ? <CheckIcon /> : <SaveIcon />}
                  </IconButton>
                </div>
              ) : (
                <IconButton
                  color="primary"
                  onClick={() => {
                    updateEditRowState(tableMeta);
                  }}
                >
                  <EditIcon />
                </IconButton>
              )}
            </div>
          );
        },
      },
    },
  ];

  return (
    <MuiThemeProvider theme={theme}>
      <MUIDataTable
        title="Security Rules"
        data={rules}
        columns={columns}
        options={{
          filter: true,
          // responsive: "scrollMaxHeight",
          selectableRows: 'none',
          rowsPerPage: 15,
        }}
      />
    </MuiThemeProvider>
  );
}

const status = ['enabled', 'disabled'];
