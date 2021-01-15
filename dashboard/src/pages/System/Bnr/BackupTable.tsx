import { createMuiTheme, MuiThemeProvider } from '@material-ui/core';
import IconButton from '@material-ui/core/IconButton';
import DownloadIcon from '@material-ui/icons/CloudDownloadOutlined';
import axios from 'axios';
import MUIDataTable, { MUIDataTableColumn, MUIDataTableOptions } from 'mui-datatables';
import React from 'react';
import Constants from '../../../Constants';

const fileDownload = require('js-file-download');

const theme = createMuiTheme({
  typography: { fontFamily: 'Open Sans, Rajdhani' },
  palette: {
    type: 'light',
    primary: { 500: '#000080' },
    secondary: { A400: '#000080' }, // '#000080' },
  },
});



export default function BackupTable(props: any) {
  const downloadBackup = (val: any) => {
    const config = {
      responseType: 'blob',
      headers: {
        'X-SESSION': localStorage.getItem('X-SESSION'),
        'X-CSRF': localStorage.getItem('X-CSRF'),
      },
    };

    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/system/backup/${val}`, config)
      .then((response) => {
        fileDownload(response.data, 'trasa-system-backup.zip');
      })
      .catch((error) => {
        console.log(error);
      });
  };

  const columns = [
    {
      name: 'Bacup Name',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'Backup Type',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'Backup Taken At',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'Download',
      options: {
        filter: false,
        customBodyRender: (value: any) => {
          return (
            <IconButton
              color="secondary"
              aria-label="dlownload backup file"
              onClick={() => downloadBackup(value)}
            >
              <DownloadIcon />
            </IconButton>
          );
        },
      },
    },
  ];

  return (
    <div>
      <MuiThemeProvider theme={theme}>
        <MUIDataTable
          title="Backups"
          data={props.backups}
          columns={columns as MUIDataTableColumn[]}
          options={options as MUIDataTableOptions}
        />
      </MuiThemeProvider>
    </div>
  );
}

const options = {
  filter: true,
  filterType: 'textField',
};
