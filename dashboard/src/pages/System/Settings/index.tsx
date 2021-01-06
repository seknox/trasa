import Grid from '@material-ui/core/Grid';
import axios from 'axios';
import React, { useEffect, useState } from 'react';
import Constants from '../../../Constants';
import EmailSetting from './Emails';
import FcmConfig from './FCMConfig';
import OrgAccountSetting from './OrgAccount';

export type PassPolicyProps = {};

export default function SystemStatus() {
  //  const [passPolicy, setPassPolicy] = useState();
  const [emailSetting, setEmailSetting] = useState({
    integrationType: 'smtp',
    serverAddress: '',
    serverPort: '',
    senderAddress: '',
    authKey: '',
    authPass: '',
  });

  function handleEmailConfigChange(e: any) {
    const { name } = e.target;
    const { value } = e.target;
    setEmailSetting({ ...emailSetting, [name]: value });
  }

  useEffect(() => {
    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/system/settings/all`)
      .then((r) => {
        if (r.data.status === 'success') {
          const parsedSettingVals = JSON.parse(r.data.data[0].emailSettings.settingValue);
          setEmailSetting(parsedSettingVals);
        }
      })
      .catch((error) => {
        console.log(error);
      });
  }, []);

  const [orgData, setOrgSetting] = useState({
    orgName: '',
    domain: '',
    primaryContact: '',
    timezone: '',
  });

  React.useEffect(() => {
    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/org/detail`)
      .then((r) => {
        if (r.data.status === 'success') {
          console.debug('orgData: ', r.data.data[0]);
          setOrgSetting(r.data.data[0]);
        }
      })
      .catch((error) => {
        console.error(error);
      });
  }, []);

  return (
    <Grid container spacing={2} direction="row" alignItems="center" justify="center">
      <Grid item xs={9}>
        <OrgAccountSetting orgData={orgData} />
      </Grid>
      <Grid item xs={9}>
        <EmailSetting
          emailSetting={emailSetting}
          handleEmailConfigChange={handleEmailConfigChange}
        />
      </Grid>
      <Grid item xs={9}>
        <FcmConfig orgData={orgData} />
      </Grid>
    </Grid>
  );
}
