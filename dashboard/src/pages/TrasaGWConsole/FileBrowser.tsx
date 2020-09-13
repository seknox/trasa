import IconButton from '@material-ui/core/IconButton';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemSecondaryAction from '@material-ui/core/ListItemSecondaryAction';
import ListItemText from '@material-ui/core/ListItemText';
import CloudDownloadIcon from '@material-ui/icons/CloudDownload';
import DeleteIcon from '@material-ui/icons/Delete';
import FileIcon from '@material-ui/icons/Description';
import axios from 'axios';
import React, { useEffect, useState } from 'react';
import Constants from '../../Constants';
import Layout from '../../Layout/DashboardBase';

export default function FileBrowser() {
  const [files, setFiles] = useState([]);

  const sendReq = (action: string, name: string, i: number) => {
    const config = {
      headers: {
        'X-CSRF': localStorage.getItem('X-CSRF'),
      },
      url: '',
      method: 'GET',
    };

    switch (action) {
      case 'download':
        config.url = `${Constants.TRASA_HOSTNAME}/api/v1/my/download_file/token`;
        config.method = 'POST';
        break;
      case 'delete':
        if (!window.confirm('Are you sure ?')) {
          return;
        }
        config.url = `${Constants.TRASA_HOSTNAME}/api/v1/my/delete_file/${name}`;
        config.method = 'POST';
        break;
      default:
        return;
    }

    axios(config)
      .then((response) => {
        // console.log(response.data)
        if (action === 'download') {
          if (response.data.status === 'success') {
            const ssKey = response.headers.sskey;
            const link = document.createElement('a');
            link.href = `${Constants.TRASA_HOSTNAME}/api/v1/my/download_file/get/${name}/${ssKey}?a=a`;
            link.setAttribute('download', 'download');
            link.target = '_blank';
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);
          }
        } else if (action === 'delete') {
          const tempFiles = files.splice(i, 1);
          setFiles(tempFiles);
          // window.location.reload()
        }
      })
      .catch((error) => {
        console.log(error);
      });
  };

  useEffect(() => {
    axios.get(`${Constants.TRASA_HOSTNAME}/api/v1/my/download_file/list`).then((res) => {
      // console.log(res.data);
      if (res.data.status === 'success') {
        setFiles(res.data.data[0]);
      }
    });
  }, []);

  return (
    <div>
      <Layout>
        {/* <Headers pageName={'File Browser'} tabHeaders={['Files']} Components={[]}  /> */}

        <List component="nav" aria-label="main mailbox folders">
          {files.map((file, i) => (
            <ListItem>
              <ListItemIcon>
                <FileIcon />
              </ListItemIcon>
              <ListItemText primary={file} />
              <ListItemSecondaryAction>
                <IconButton
                  edge="end"
                  aria-label="download"
                  onClick={() => {
                    sendReq('download', file, i);
                  }}
                >
                  <CloudDownloadIcon />
                </IconButton>
                <IconButton
                  edge="end"
                  aria-label="delete"
                  onClick={() => {
                    sendReq('delete', file, i);
                  }}
                >
                  <DeleteIcon />
                </IconButton>
              </ListItemSecondaryAction>
            </ListItem>
          ))}
        </List>
      </Layout>
    </div>
  );
}
