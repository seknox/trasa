import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import axios from 'axios';
import React from 'react';
import ProgressHOC from '../../../utils/Components/Progressbar';
import Constants from '../../../Constants';
import BackupTable from './BackupTable';

export default function Backup() {
  const [backups, setbackups] = React.useState([]);
  const [loader, setloader] = React.useState(false);

  const initBackup = () => {
    //mixpanel.track('system-bnr-takebackup');
    setloader(true);

    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/system/backup/create`)
      .then((response) => {
        setloader(false);
        let dataArr = [];
        const data = response.data.data[0];
        dataArr = data.map(function (n: any) {
          const date = new Date(n.createdAt * 1000);
          return [n.backupName, n.backupType, date.toDateString(), n.backupID];
        });
        setbackups(dataArr);
      })
      .catch((error) => {
        setloader(false);
        console.log(error);
      });
  };

  const getBackupData = () => {
    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/system/backups`)
      .then((response) => {
        let dataArr = [];
        const data = response.data.data[0];
        dataArr = data.map(function (n: any) {
          const date = new Date(n.createdAt * 1000);
          return [n.backupName, n.backupType, date.toLocaleString(), n.backupID];
        });
        setbackups(dataArr);
      })
      .catch((error) => {
        console.log(error);
      });
  };

  React.useEffect(() => {
    getBackupData();
  }, []);

  return (
    <Grid container spacing={2} alignItems="center" justify="center">
      <Grid item xs={12}>
        <Button variant="contained" color="secondary" size="small" onClick={initBackup}>
          Take Backup Now
        </Button>

        {/* <Typography variant="h4"> Last Backup taken at <b></b>  hours. </Typography>  */}

        <br />
        {loader ? <ProgressHOC /> : ''}
      </Grid>

      <Grid item xs={12} sm={12} md={9}>
        <BackupTable backups={backups} />
      </Grid>
    </Grid>
  );
}
