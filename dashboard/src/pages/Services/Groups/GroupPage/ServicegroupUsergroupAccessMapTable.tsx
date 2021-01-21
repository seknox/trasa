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
import Constants from '../../../../Constants';

const ttheme = createMuiTheme({
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
  links: {
    backgroundColor: 'Transparent',
    borderRadius: 3,
    border: 0,
    // color: 'black',
    fontColor: 'blue',
    fontSize: '12px',
    fontFamily: 'Open Sans, Rajdhani',
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
}));

const links = {
  backgroundColor: 'Transparent',
  borderRadius: 3,
  border: 0,
  // color: 'black',
  fontColor: 'blue',
  fontSize: '12px',
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

export default function ServicegroupUsergroupMapTable(props: any) {
  const classes = useStyles();
  const [addedGroups, setaddedGroups] = useState([]);
  const [updateRow, setupdateRow] = useState(false);
  const [updateValue, setupdateValue] = useState({ privilege: '', mapID: '' });
  const [loader, setLoader] = useState(false);
  const [updateRowIndex, setupdateRowIndex] = useState(0);

  const fetchAssignedUserGroups = (groupID: string) => {
    const url = `${Constants.TRASA_HOSTNAME}/api/v1/accessmap/servicegroup/addedusergroups/${groupID}`;
    axios.get(url).then((response) => {
      if (response.data.status === 'success') {
        let dataArr = [];
        const data = response.data.data[0];
        dataArr = data.map(function (n: any) {
          //  console.log(n)
          const cdate = new Date(n.addedAt * 1000);
          // WARNING:should the sequence of returned element be altered, the index value should also be updated in deleteuserGroupmap function.
          return [
            n.userGroupName,
            n.privilege,
            n.policyName,
            cdate.toDateString(),
            n.mapID,
            n.usergroupID,
            n.policyID,
          ];
        });
        setaddedGroups(dataArr);
      }
    });
  };

  useEffect(() => {
    fetchAssignedUserGroups(props.groupID);
  }, [props.groupID]);

  const deleteGroupMaps = (rowsDeleted: any) => {
    const mapID = rowsDeleted.data.map((i: any) => {
      return addedGroups[i.index][4];
    });

    // const val = JSON.stringify({serviceID, permissionID})
    const req = { mapID };

    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/accessmap/servicegroup/usergroup/delete`, req)
      .then(() => {})
      .catch((error) => {
        console.log(error);
      });
  };

  const updateUsername = (e: any) => {
    const row = addedGroups[updateRowIndex];
    const updateVal = { privilege: e.target.value, mapID: row[4] };
    setupdateValue(updateVal);
  };

  const updateEditRowState = (tableMeta: any) => {
    setupdateRow(true);
    setupdateRowIndex(tableMeta.rowIndex);
  };

  const udpateAppusername = () => {
    setLoader(true);
    const val = updateValue;

    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/accessmap/servicegroup/usergroup/update`, val)
      .then((r) => {
        setLoader(false);
        if (r.data.status === 'success') {
          fetchAssignedUserGroups(props.groupID);
          setupdateRow(false);
          setupdateRowIndex(0);
        }
      })
      .catch((error) => {
        console.log(error);
      });
  };

  const columns = [
    {
      name: 'User Group Name',
      options: {
        filter: true,
        customBodyRender: (value: any, tableMeta: MUIDataTableMeta) => {
          return (
            <a style={links} href={`/users/groups/group/${tableMeta.rowData[5]}`}>
              {value}
            </a>
          );
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
      name: 'Assgned Policy',
      options: {
        filter: true,
        customBodyRender: (value: any, tableMeta: MUIDataTableMeta) => {
          return (
            <a style={links} href="/policies#Access Policies">
              {value}
            </a>
          );
        },
      },
    },
    {
      name: 'Assigned On',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return <div style={tableBodyFont}>{value}</div>;
        },
      },
    },
    {
      name: 'Edit Privilege',
      options: {
        filter: true,
        customBodyRender: (value: any, tableMeta: MUIDataTableMeta) => {
          return (
            <div>
              {updateRow && updateRowIndex === tableMeta.rowIndex ? (
                <div>
                  <IconButton id="savePrivilegeBtn" color="primary" onClick={udpateAppusername}>
                    {loader && <CircularProgress size={48} className={classes.fabProgress} />}
                    {loader ? <CheckIcon /> : <SaveIcon />}
                  </IconButton>
                </div>
              ) : (
                <IconButton
                  id="editGroupPrivilegeBtn"
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
    {
      name: 'usergroupID',
      options: {
        filter: true,
        display: false,
      },
    },
    {
      name: 'policyID',
      options: {
        filter: true,
        display: false,
      },
    },
    {
      name: 'mapID',
      options: {
        filter: true,
        display: false,
      },
    },
  ];
  return (
    <MuiThemeProvider theme={ttheme}>
      <MUIDataTable
        title="Group Access Map"
        data={addedGroups}
        columns={columns as MUIDataTableColumn[]}
        options={{
          filter: true,
          responsive: 'scrollMaxHeight',
          // resizableColumns: true,
          onRowsDelete: deleteGroupMaps,
        }}
      />
    </MuiThemeProvider>
  );
}
