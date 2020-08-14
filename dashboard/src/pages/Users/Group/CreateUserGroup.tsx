import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';
import LinearProgress from '@material-ui/core/LinearProgress';
import TextField from '@material-ui/core/TextField';
import Typography from '@material-ui/core/Typography';
import { Box } from '@material-ui/core';
// import TextField from '@material-ui/core/TextField';
import axios from 'axios';
import React, { useState } from 'react';
import Constants from '../../../Constants';

export default function CreateUserGroup(props: any) {
  // const [open, setOpen] = useState(false);
  const [loader, setLoader] = useState(false);
  const [groupName, setGroupName] = useState('');

  const submitCreateGroup = (event: any) => {
    setLoader(true);

    let lgroupName = groupName;
    if (groupName.length === 0) {
      lgroupName = props.groupMeta.groupName;
    }
    const reqData = {
      groupName: lgroupName,
      groupType: 'usergroup',
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
        window.location.reload();
      });
  };

  function handleGroupNameChange(event: React.ChangeEvent<HTMLInputElement>) {
    const val = event.target.value;
    setGroupName(val);
  }

  return (
    <div>
      <Dialog
        fullWidth
        maxWidth="md"
        open={props.open}
        onClose={props.handleClose}
        // aria-labelledby="form-dialog-title"
      >
        <DialogContent>
          <DialogContentText>
            {props.update ? (
              <Typography variant="h2"> Update Group </Typography>
            ) : (
              <Typography variant="h2"> Create User Group </Typography>
            )}
          </DialogContentText>

          <Divider light />
          <br />
          <Grid container spacing={2}>
            <Grid item xs={4} sm={4} md={4}>
              <Typography variant="h3">Name :</Typography>
            </Grid>
            <Grid item xs={8} sm={8} md={8}>
              <TextField
                fullWidth
                variant="outlined"
                size="small"
                defaultValue={props.update ? props.groupMeta.groupName : ''}
                onChange={handleGroupNameChange}
                name="groupName"
              />
            </Grid>
          </Grid>

          <br />
          <Box display="flex" alignItems="center" justifyContent="center">
            <Button variant="contained" color="secondary" onClick={(e) => submitCreateGroup(e)} id='createGroupSubmitBtn'>
              Submit
            </Button>
          </Box>

          {loader ? <LinearProgress /> : null}
        </DialogContent>
      </Dialog>
    </div>
  );
}
