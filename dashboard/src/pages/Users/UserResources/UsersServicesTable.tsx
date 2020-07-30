import { createMuiTheme, MuiThemeProvider } from '@material-ui/core';
import MUIDataTable from 'mui-datatables';
import React from 'react';

const theme = createMuiTheme({
  typography: { fontFamily: 'Open Sans, Rajdhani' },
  palette: {
    type: 'light',
    primary: { 500: '#000080' },
    secondary: { A400: '#000080' },
  },
});

const tableBodyFont = {
  backgroundColor: 'Transparent',
  borderRadius: 3,
  border: 0,
  fontColor: 'dimgray',
  fontSize: '12px',
  fontFamily: 'Open Sans, Rajdhani',
};

export default function UsersAppTable(props: any) {
  const columns = [
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
      name: 'Assigned Privilege',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'Assigned Policy',
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
  ];

  function getTitle(arr: any) {
    const count = arr.length;
    return `Users authorized apps and services (Total: ${count})`;
  }

  return (
    <div>
      <MuiThemeProvider theme={theme}>
        <MUIDataTable
          title={getTitle(props.AppUsers)}
          data={props.AppUsers.map(function (n: any) {
            return [n.serviceName, n.privilege, n.policy.policyName, n.userAddedAt];
          })}
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
