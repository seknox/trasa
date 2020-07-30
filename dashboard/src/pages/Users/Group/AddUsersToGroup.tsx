import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import axios from 'axios';
import React, { useEffect, useState } from 'react';
import MultiSelect from 'react-multi-select-component';
import Constants from '../../../Constants';
import { Usertype } from '../../../types/users';
import ProgressBar from '../../../utils/Components/Progressbar';

const useStyles = makeStyles((theme) => ({
  root: {
    flexGrow: 1,
    height: 250,
  },
  container: {
    flexGrow: 1,
    position: 'relative',
  },
  paper: {
    position: 'absolute',
    zIndex: 1,
    marginTop: theme.spacing(1),
    left: 0,
    right: 0,
    minHeight: '20px',
  },
  chip: {
    margin: `${theme.spacing(1) / 2}px ${theme.spacing(1) / 4}px`,
  },
  inputRoot: {
    flexWrap: 'wrap',
  },
  inputInput: {
    width: 'auto',
    flexGrow: 1,
  },
  divider: {
    height: theme.spacing(2),
  },
  submitRight: {
    marginLeft: '35%',
  },
  dialogueRoot: {
    minHeight: 500,
    // overflow: 'auto',
  },
  selectZindex: {
    overflow: 'visible',
    // zIndex: 99999,
  },
}));

type AddUserToGroupProps = {
  open: boolean;
  groupID: string;
  usersThatCanBeAdded: Usertype[];
  handleClose: () => void;
};

export default function AddUserToGroup(props: AddUserToGroupProps) {
  const classes = useStyles();

  const [selected, setSelected] = useState([]);

  const [listVals, setListVals] = useState([{ label: '', value: '', id: '' }]);

  const [loader, setLoader] = useState(false);

  const submitAddUserRequest = () => {
    setLoader(true);
    const userArr = selected.map((users: any) => {
      console.log('users: ', users);
      return users.id;
    });

    const reqData = { groupID: props.groupID, userIDs: userArr, updateType: 'add' };

    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/groups/user/update`, reqData)
      .then((r) => {
        setLoader(false);
        if (r.data.status === 'success') {
          window.location.href = `/users/groups/group/${props.groupID}`;
        }
      })
      .catch(function (error) {
        setLoader(false);
        console.log('catched: ', error);
      });
  };

  useEffect(() => {
    const vals = props.usersThatCanBeAdded.map((v) => {
      return { label: `${v.firstName} ${v.lastName} (${v.email})`, value: v.email, id: v.ID };
    });

    setListVals(vals);
  }, [props.usersThatCanBeAdded]);

  return (
    <Dialog
      fullWidth
      maxWidth="md"
      open={props.open}
      onClose={props.handleClose}
      aria-labelledby="form-dialog-title"
      classes={{
        paper: classes.dialogueRoot,
      }}
    >
      <DialogTitle disableTypography id="alert-dialog-title">
        Add users to this group{' '}
      </DialogTitle>
      <DialogContent>
        <div className={classes.selectZindex}>
          <br />
          <MultiSelect
            options={listVals}
            value={selected}
            onChange={setSelected}
            labelledBy="Select users"
          />
        </div>

        <br />
        <br />
        <br />
        <Grid container spacing={2} direction="column" alignItems="center" justify="center">
          <Grid item xs={12}>
            <Button variant="contained" color="secondary" onClick={submitAddUserRequest}>
              Add Selected Users
            </Button>
            <br />
          </Grid>
        </Grid>
        {loader ? <ProgressBar /> : null}
      </DialogContent>
      <DialogActions>
        <Button variant="outlined" color="primary" onClick={props.handleClose}>
          Close
        </Button>
      </DialogActions>
    </Dialog>
  );
}
