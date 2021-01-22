import { createMuiTheme, MuiThemeProvider } from '@material-ui/core';
import IconButton from '@material-ui/core/IconButton';
import DownloadIcon from '@material-ui/icons/CloudDownloadOutlined';
import axios from 'axios';
import MUIDataTable, { MUIDataTableColumn, MUIDataTableOptions } from 'mui-datatables';
import React from 'react';
import Constants from '../../../Constants';
import {LogtableV2Theme} from '../../../utils/styles/themes'
const fileDownload = require('js-file-download');


export default function BackupTable(props: any) {
  const downloadBackup = (val: any) => {
    const config = {
      responseType: 'blob',
      headers: {
        'X-SESSION': localStorage.getItem('X-SESSION'),
        'X-CSRF': localStorage.getItem('X-CSRF'),
      },
    };

    // since the default instance has middlewares to parse json response,
    //we are using new axios instance for downloading files.
    axios.create(config)
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/system/backup/${val}`)
      .then((response) => {
        fileDownload(response.data, 'trasa-system-backup.zip');
      })
      .catch((error) => {
        console.log(error);
      });
  };

  const columns = [
    {
      name: 'Backup Name',
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
      <MuiThemeProvider theme={LogtableV2Theme}>
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
