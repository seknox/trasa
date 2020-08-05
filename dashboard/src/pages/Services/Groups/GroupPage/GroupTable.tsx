import { createMuiTheme, MuiThemeProvider } from '@material-ui/core';
import Button from '@material-ui/core/Button';
import MUIDataTable, {
  MUIDataTableColumn,
  MUIDataTableMeta,
  MUIDataTableOptions,
} from 'mui-datatables';
import React, { useEffect } from 'react';
import { Link } from 'react-router-dom';

const theme = createMuiTheme({
  typography: { fontFamily: 'Open Sans, Rajdhani' },
  palette: {
    type: 'light',
    primary: { 500: '#000080' },
    secondary: { A400: '#000080' }, // '#000080' },
  },
});

export default function AuthServicegroupTable(props: any) {
  useEffect(() => {}, []);

  const columns = [
    {
      name: 'Service name',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'Hostname / IP address',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'Service type',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'View App',
      options: {
        filter: false,
        customBodyRender: (
          value: any,
          tableMeta: MUIDataTableMeta,
          updateValue: (value: string) => void,
        ) => {
          return (
            <Button
              component={Link}
              to={`/services/service/${value}`}
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

  const options = {
    filter: true,
    responsive: 'scrollMaxHeight',
    onRowsDelete: props.removeServices,
  };

  return (
    <div>
      <MuiThemeProvider theme={theme}>
        <MUIDataTable
          title="Added Services"
          data={props.allServicesArray}
          columns={columns as MUIDataTableColumn[]}
          options={options as MUIDataTableOptions}
        />
      </MuiThemeProvider>
    </div>
  );
}
