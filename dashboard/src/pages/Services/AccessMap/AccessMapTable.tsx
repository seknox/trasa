import { createMuiTheme, MuiThemeProvider } from '@material-ui/core';
import CircularProgress from '@material-ui/core/CircularProgress';
import green from '@material-ui/core/colors/green';
import IconButton from '@material-ui/core/IconButton';
import { makeStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
import CheckIcon from '@material-ui/icons/Check';
import EditIcon from '@material-ui/icons/Edit';
import SaveIcon from '@material-ui/icons/Save';
import axios from 'axios';
import MUIDataTable, { MUIDataTableColumn, MUIDataTableMeta } from 'mui-datatables';
import React, { useEffect, useState } from 'react';
import Constants from '../../../Constants';

const theme = createMuiTheme({
  typography: { fontFamily: 'Open Sans, Rajdhani' },
  palette: {
    type: 'light',
    primary: { 500: '#000080' },
    secondary: { A400: '#000080' }, // '#000080' },
  },

  // #0A2053
});

const useStyles = makeStyles(() => ({
  fabProgress: {
    color: green[500],
    position: 'absolute',
    top: -6,
    left: -6,
    zIndex: 1,
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

export default function AppuserTable(props: any) {
  const classes = useStyles();

  const [updateRow, setUpdateRow] = useState(false);
  const [appUsersArr, setappUsersArr] = useState([]);
  // const [updateUsername, setupdateUsername] = useState({});
  const [loader, setLoader] = useState(false);

  const [updateRowIndex, setupdateRowIndex] = useState(0);
  const [updateValue, setUpdateValue] = useState({ privilege: '', mapID: '' });

  const assignedUsers = (ID: string) => {
    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/accessmap/service/user/${ID}`)
      .then((response) => {
        const resp = response.data.data[0];
        if (response.status === 403) {
          window.location.href = '/login';
        }
        const dataArr = resp.map(function (n: any) {
          return [n.trasaID, n.privilege, n.policy.policyName, n.userAddedAt, n.mapID];
        });
        // setAppData(resp.App);
        setappUsersArr(dataArr);
      })
      .catch((error) => {
        console.log(error);
      });
  };
  useEffect(() => {
    assignedUsers(props.ID);
  }, [props.ID]);

  const handleDeletePermission = (rowsDeleted: any) => {
    const { ID } = props;
    const mapIDs = rowsDeleted.data.map((i: any) => {
      setappUsersArr(appUsersArr[i.index][4]);
      return appUsersArr[i.index][4];
    });

    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/accessmap/service/user/delete`, { mapIDs })
      .then(() => {
        window.location.reload();
      })
      .catch((error) => {
        console.log(error);
      });
  };

  const updateUsername = (e: any) => {
    const row = appUsersArr[updateRowIndex];
    const updateVal = { privilege: e.target.value, mapID: row[4] };

    setUpdateValue(updateVal);
  };

  const updateEditRowState = (tableMeta: any) => {
    setUpdateRow(true);
    setupdateRowIndex(tableMeta.rowIndex);
  };

  const updateUserPolicy = () => {
    setLoader(true);
    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/accessmap/service/user/update`, updateValue)
      .then((response) => {
        setLoader(false);
        assignedUsers(props.ID);
        setUpdateRow(false);
        setupdateRowIndex(0);
      })
      .catch((error) => {
        console.log(error);
      });
  };

  const columns = [
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
      name: 'Assigned Privilege',
      options: {
        filter: true,
        customBodyRender: (value: any, tableMeta: MUIDataTableMeta) => {
          return (
            <div>
              {updateRow && updateRowIndex === tableMeta.rowIndex ? (
                <TextField
                  fullWidth
                  id="username"
                  // value={this.props.Expiry}
                  onChange={updateUsername}
                  label="update username"
                  type="text"
                  defaultValue={value}
                  InputLabelProps={{
                    shrink: true,
                  }}
                />
              ) : (
                value
              )}
            </div>
          );
        },
      },
    },
    {
      name: ' Assigned Policy',
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
          return <div style={tableBodyFont}> {new Date(value * 1000).toDateString()} </div>;
        },
      },
    },
    {
      name: 'Edit Privilege',
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
                  <IconButton id="savePrivilegeBtn" color="primary" onClick={updateUserPolicy}>
                    {loader && <CircularProgress size={58} className={classes.fabProgress} />}
                    {loader ? <CheckIcon /> : <SaveIcon />}
                  </IconButton>
                </div>
              ) : (
                <IconButton
                  id="editUserPrivilegeBtn"
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
        title="User Access Map"
        data={appUsersArr}
        columns={columns as MUIDataTableColumn[]}
        options={{
          filter: true,
          responsive: 'scrollMaxHeight',
          // resizableColumns: true,
          onRowsDelete: handleDeletePermission,
        }}
      />
    </MuiThemeProvider>
  );
}
