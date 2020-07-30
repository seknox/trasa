import { createMuiTheme, MuiThemeProvider } from '@material-ui/core';
import Button from '@material-ui/core/Button';
import MUIDataTable from 'mui-datatables';
import React from 'react';
import { Link } from 'react-router-dom';

const theme = createMuiTheme({
  typography: { fontFamily: 'Open Sans, Rajdhani' },
  palette: {
    type: 'light',
    primary: { 500: '#000080' },
    secondary: { A400: '#000080' },
  },
  // #0A2053
});

export default function GroupsAssignedToUser(props: any) {
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
      name: 'Group Details',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return (
            <Button
              component={Link}
              to={`/users/groups/group/${value}`}
              variant="outlined"
              color="secondary"
            >
              View
            </Button>
          );
        },
      },
    },
  ];

  function getTitle(arr: any) {
    const count = arr.length;
    return `Groups (Total: ${count})`;
  }

  return (
    <div>
      <MuiThemeProvider theme={theme}>
        <MUIDataTable
          title={getTitle(props.userGroups)}
          data={props.userGroups.map(function (n: any) {
            return [n.groupName, new Date(n.createdAt * 1000).toLocaleDateString(), n.groupID];
          })}
          // data={this.props.userGroups}
          columns={columns}
          options={{
            filter: true,
            // responsive: "scroll",
            selectableRows: 'none',
            // resizableColumns: true,
            // onRowsDelete: this.handleDeletePermission
          }}
        />
      </MuiThemeProvider>
    </div>
  );
}
