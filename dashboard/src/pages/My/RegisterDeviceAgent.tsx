import React, { useEffect, useState } from 'react';
import axios from 'axios';
import Typography from '@material-ui/core/Typography';
import Paper from '@material-ui/core/Paper';
import Grid from '@material-ui/core/Grid';
import Constants from '../../Constants';

export default function (props: any) {
  const [msg, setMsg] = useState('');
  useEffect(
    function () {
      const { deviceAgentToken } = props.match.params;
      if (deviceAgentToken) {
        const data = { token: deviceAgentToken };
        axios
          .post(`${Constants.TRASA_HOSTNAME}/idp/login/deviceAgent/verify`, data)
          .then((resp) => {
            if (resp.data.status === 'success') {
              setMsg('successfully registered device');
            }
          });
      }
    },
    [props.match.params],
  );
  return (
    <div>
      {/* TODO @sshah beautify this view */}
      <Grid alignContent="center">
        <Paper>
          <Typography variant="h3">{msg}</Typography>
        </Paper>
      </Grid>
    </div>
  );
}
