import Button from '@material-ui/core/Button';
import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
import Typography from '@material-ui/core/Typography';
import axios from 'axios';
import React, {useState} from 'react';
import ProgressHOC from '../../../../utils/Components/Progressbar';
import Constants from '../../../../Constants';

const useStyles = makeStyles((theme) => ({
  container: {
    display: 'flex',
    flexWrap: 'wrap'
  },
  root: {
    flexGrow: 1
  },
  textFieldRoot: {
    padding: 1,
    'label + &': {
      marginTop: theme.spacing(3)
    }
  },
  paper: {
    maxWidth: 500,
    padding: theme.spacing(2),
    textAlign: 'center',
    color: theme.palette.text.secondary
  },

  textFieldInputBig: {
    borderRadius: 4,
    backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    fontSize: 16,
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)'
    }
  }
}));

function CreateCA() {
  const classes = useStyles();
  const [caName, setcaName] = useState('');
  const [caDetails, setcaDetails] = useState({C: '', ST: '', L:'', O: '', OU:''});
  const [loader, setLoader] = useState(false);

  const handleCAName = (event: React.ChangeEvent<HTMLInputElement>) => {
    setcaName(event.target.value);
  };

  const handleCADetial = (event: React.ChangeEvent<HTMLInputElement>) => {
    setcaDetails({ ...caDetails, [event.target.name]: event.target.value });
  };

  const SubmitCARequest = (e:React.FormEvent<HTMLElement>) => {
    e.preventDefault()
    setLoader(true);

    let req = { CN: caName, names: [caDetails] };
    axios
      .post(Constants.TRASA_HOSTNAME + '/api/v1/providers/ca/tsxca/init', req)
      .then((response) => {
        setLoader(false);
      })
      .catch((error) => {
        setLoader(false);
        console.log(error);
      });
  };

  return (
    <div className={classes.paper}>
      <form onSubmit={SubmitCARequest}>
        <Grid container spacing={2} direction="column" alignItems="center" justify="center">
          <Grid item xs={12}>
            <br />
            <Divider light />
            <br /> <br />
            <Typography variant="h3" component="h3">
              <b>CA information</b>
            </Typography>
            <br />
          </Grid>

          <Grid container spacing={2}>
            <Grid item xs={6}>
              <Typography component="h4" variant="h3">
                CA Name :
              </Typography>
            </Grid>
            <Grid item xs={6} sm={6} md={6}>
              <TextField fullWidth onChange={handleCAName} name="CN" value={caName} />
            </Grid>
          </Grid>
          <Grid container spacing={2}>
            <Grid item xs={6}>
              <Typography component="h4" variant="h3">
                Country :
              </Typography>
            </Grid>
            <Grid item xs={6} sm={6} md={6}>
              <TextField fullWidth onChange={handleCADetial} name="C" value={caDetails.C} />
            </Grid>
          </Grid>

          <Grid container spacing={2}>
            <Grid item xs={6}>
              <Typography component="h4" variant="h3">
                State :
              </Typography>
            </Grid>
            <Grid item xs={6} sm={6} md={6}>
              <TextField fullWidth onChange={handleCADetial} name="ST" value={caDetails.ST} />
            </Grid>
          </Grid>

          <Grid container spacing={2}>
            <Grid item xs={6}>
              <Typography component="h4" variant="h3">
                Locality :
              </Typography>
            </Grid>
            <Grid item xs={6} sm={6} md={6}>
              <TextField fullWidth onChange={handleCADetial} name="L" value={caDetails.L} />
            </Grid>
          </Grid>

          <Grid container spacing={2}>
            <Grid item xs={6}>
              <Typography component="h4" variant="h3">
                Organization Name :
              </Typography>
            </Grid>
            <Grid item xs={6} sm={6} md={6}>
              <TextField fullWidth onChange={handleCADetial} name="O" value={caDetails.O} />
            </Grid>
          </Grid>

          <Grid container spacing={2}>
            <Grid item xs={6}>
              <Typography component="h4" variant="h3">
                Organization Unit Name :
              </Typography>
            </Grid>
            <Grid item xs={6} sm={6} md={6}>
              <TextField fullWidth onChange={handleCADetial} name="OU" value={caDetails.OU} />
            </Grid>
          </Grid>

          <Grid container spacing={2}>
            <Grid item xs={6}>
              <Typography component="h4" variant="h3">
                Expiration Time :
              </Typography>
            </Grid>
            <Grid item xs={6} sm={6} md={6}>
              <TextField fullWidth onChange={handleCADetial} name="OU" value={caDetails.OU} />
            </Grid>
          </Grid>

          <Grid item xs={12}>
            <br />
            <Divider light />
            <br /> <br />
            <Typography variant="h3" component="h3">
              <b>Signing profile</b>
            </Typography>
            <br />
          </Grid>

          <Grid container spacing={2}>
            <Grid item xs={6}>
              <Typography component="h4" variant="h3">
                Default expiration Time :
              </Typography>
            </Grid>
            <Grid item xs={6} sm={6} md={6}>
              <TextField fullWidth onChange={handleCADetial} name="OU" value={caDetails.OU} />
            </Grid>
          </Grid>

          <Grid item xs={6}>
            <br />
            <div className={classes.root}>
              <div >
                <Button variant="contained" type="submit">
                  Submit
                </Button>
              </div>
            </div>
            {loader ? <ProgressHOC /> : ''}
          </Grid>
        </Grid>
      </form>
    </div>
  );
}
