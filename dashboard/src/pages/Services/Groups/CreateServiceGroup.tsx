import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';
import LinearProgress from '@material-ui/core/LinearProgress';
import { makeStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
import Typography from '@material-ui/core/Typography';
// import TextField from '@material-ui/core/TextField';
import axios from 'axios';
import React, { useState } from 'react';
import Constants from '../../../Constants';

const useStyles = makeStyles((theme) => ({
  textFieldRoot: {
    padding: 1,
    'label + &': {
      marginTop: theme.spacing(3),
    },
  },

  textFieldInput: {
    borderRadius: 4,
    backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    fontSize: 16,
    padding: '10px 12px',
    width: 'calc(100% - 4px)',
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
    },
  },
  textFieldInputBig: {
    borderRadius: 4,
    backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    fontSize: 16,
    //    padding: '10px 100px',
    //     width: 'calc(100% - 4px)',
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
    },
  },
  textFieldFormLabel: {
    fontSize: 18,
  },
}));

export default function CreateServicegroup(props: any) {
  const classes = useStyles();

  const [groupName, setGroupName] = useState('');
  const [loader, setLoader] = useState(false);

  const submitCreateGroup = (event: any) => {
    setLoader(true);

    if (groupName.length === 0) {
      setGroupName(props.groupMeta.groupName);
    }
    const reqData = {
      groupName,
      groupType: 'servicegroup',
      groupID: props.update ? props.groupMeta.groupID : '',
    };
    event.preventDefault();
    let url = `${Constants.TRASA_HOSTNAME}/api/v1/groups/create`;
    if (props.update) {
      url = `${Constants.TRASA_HOSTNAME}/api/v1/groups/update`;
    }

    axios
      .post(url, reqData)
      .then((response) => {
        setLoader(false);
        if (response.data.status === 'success') {
          window.location.reload();
        }
      })
      .catch(function (error) {
        console.log(error);
      });
  };

  return (
    <div>
      <Dialog
        fullWidth
        maxWidth="md"
        open={props.open}
        onClose={props.handleClose}
        aria-labelledby="form-dialog-title"
      >
        <DialogContent>
          <DialogContentText>
            {props.update ? (
              <Typography variant="h2"> Update Group </Typography>
            ) : (
              <Typography variant="h2"> Create service Group </Typography>
            )}
          </DialogContentText>

          <Divider light />
          <Grid container spacing={2}>
            <br />

            <Grid item xs={4} sm={4} md={4}>
              <Typography variant="h3">Name :</Typography>
            </Grid>
            <Grid item xs={8} sm={8} md={8}>
              <TextField
                fullWidth
                defaultValue={props.update ? props.groupMeta.groupName : ''}
                onChange={(e) => setGroupName(e.target.value)}
                name="groupName"
                InputProps={{
                  disableUnderline: true,
                  classes: {
                    root: classes.textFieldRoot,
                    input: classes.textFieldInputBig,
                  },
                }}
                InputLabelProps={{
                  shrink: true,
                  className: classes.textFieldFormLabel,
                }}
              />
            </Grid>
          </Grid>

          <br />
          <Button id='createGroupSubmitBtn' variant="contained" color="secondary" onClick={submitCreateGroup}>
            Submit
          </Button>

          <Divider light />
          {loader ? <LinearProgress /> : null}
        </DialogContent>
      </Dialog>
    </div>
  );
}
