import Snackbar from '@material-ui/core/Snackbar';
import MuiAlert from '@material-ui/lab/Alert';
import React from 'react';

function Alert(props: any) {
  return <MuiAlert elevation={6} variant="filled" {...props} />;
}

export default function Snack(props: any) {
  return (
    <Snackbar
      open={props.open}
      autoHideDuration={6000}
      onClose={props.snackClose}
      anchorOrigin={{ vertical: 'top', horizontal: 'center' }}
    >
      <Alert severity={props.success ? 'success' : 'error'}>{props.statusMsg}</Alert>
    </Snackbar>
  );
}
