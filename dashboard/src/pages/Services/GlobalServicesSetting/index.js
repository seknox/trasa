import Grid from '@material-ui/core/Grid';
import axios from 'axios';
import React from 'react';
import Constants from '../../../Constants';
import TRASASshAuth from './SshAuth';
import DeviceHygieneCheck from './DeviceHygieneCheck';
import DynamicAccess from './DynamicAccess';

// TODO @bhrg3se move this global service settings api to relative service database table in backend.
export default function GlobalServieSetting() {
  const [sshCertEnforced, setSshCertEnforced] = React.useState(false);
  const [deviceHygieneEnabled, setDeviceHygieneEnabled] = React.useState(false);
  const [dynamicAccessSettings, setDynamicAccessSettings] = React.useState({});

  React.useEffect(() => {
    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/system/settings/all`)
      .then((r) => {
        if (r.data.status === 'success') {
          setSshCertEnforced(r.data.data?.[0]?.sshCertSetting?.status);
          setDeviceHygieneEnabled(r.data?.data?.[0].deviceHygiene?.status);
          setDynamicAccessSettings(r.data?.data?.[0]?.dynamicAccess);
        }
      })
      .catch((error) => {
        console.log(error);
      });
  }, []);

  return (
    <Grid container spacing={2} direction="row" alignItems="center" justify="center">
      <Grid item xs={9}>
        <TRASASshAuth status={sshCertEnforced} />
      </Grid>
      <Grid item xs={9}>
        <DeviceHygieneCheck status={deviceHygieneEnabled} />
      </Grid>
      <Grid item xs={9}>
        <DynamicAccess settings={dynamicAccessSettings} />
      </Grid>
    </Grid>
  );
}
