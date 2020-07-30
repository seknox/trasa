import Button from '@material-ui/core/Button';
import Divider from '@material-ui/core/Divider';
import Drawer from '@material-ui/core/Drawer';
import Grid from '@material-ui/core/Grid';
import LinearProgress from '@material-ui/core/LinearProgress';
import Paper from '@material-ui/core/Paper';
import { makeStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
import PropTypes from 'prop-types';
import React from 'react';

const useStyles = makeStyles((theme) => ({
  menu: {
    marginLeft: '38%',
  },
  fullList: {
    marginTop: '3%',
    padding: theme.spacing(2),
    // maxWidt: '1000px',
  },
  paper: {
    padding: theme.spacing(2),
    textAlign: 'center',
    marginRight: '2%',
  },
  textField: {
    marginLeft: theme.spacing(1),
    marginRight: theme.spacing(1),
  },
  ftMargin: {
    // padding: theme.spacing(5),
    textAlign: 'center',
    marginLeft: theme.spacing(1) * 2,
    marginRight: theme.spacing(1) * 2,
  },
  headings: {
    color: 'black',
    fontSize: '15px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  subHeadings: {
    color: 'black',
    fontSize: '10px',
    fontFamily: 'Open Sans, Rajdhani',
  },
}));

type ConsoleMenuProps = {
  progress: number;
  openMenu: boolean;
  toggleDrawer: () => void;
  onClipboardChange: (e: any) => void;
  CtrlAltDel: () => void;
  onFileOpen: (event: React.ChangeEvent<HTMLInputElement>) => void;
  onUpload: () => void;
  clipboardVal: string;
  appType: string;
};

export default function ConsoleMenu(props: ConsoleMenuProps) {
  // const [top, setTop] = useState(false);
  // const [value, setValue] = useState('cp');

  // const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
  //   setValue(event.target.value);
  // };

  const classes = useStyles();

  const menuItems = (
    <div className={classes.fullList}>
      <Divider variant="middle" /> <br />
      <Grid container spacing={0}>
        <Grid item xs={4}>
          <Paper className={classes.paper}>
            <Grid container>
              <Grid item xs={12}>
                <div className={classes.headings}>Copy / Paste Value</div>
                <Divider variant="middle" />
              </Grid>
              <Grid item xs={12}>
                <TextField
                  id="outlined-multiline-static"
                  label="Value"
                  multiline
                  rows="4"
                  className={classes.textField}
                  margin="normal"
                  variant="outlined"
                  onChange={props.onClipboardChange}
                  value={props.clipboardVal}
                  fullWidth
                />
              </Grid>
            </Grid>

            {/* <input  placeholder={"PASTE HERE"} /> */}
          </Paper>
        </Grid>
        <Grid item xs={4}>
          <Paper className={classes.paper}>
            <div className={classes.ftMargin}>
              <Grid container>
                <Grid item xs={12}>
                  <div className={classes.headings}>File Transfer</div>
                  <div className={classes.subHeadings}>(only for rdp session) </div>
                  <Divider variant="middle" /> <br />
                </Grid>
                <Grid item xs={12}>
                  <div>
                    <input type="file" id="file" onChange={props.onFileOpen} /> <br /> <br />
                  </div>
                </Grid>
                <Grid item xs={6}>
                  <Button variant="contained" color="primary" onClick={props.onUpload}>
                    Upload
                    {/* <CloudUploadIcon/> */}
                  </Button>
                </Grid>
                <Grid item xs={6}>
                  <Button
                    variant="contained"
                    color="primary"
                    onClick={() => {
                      window.open('/my/file/browse');
                    }}
                  >
                    {' '}
                    Download{' '}
                  </Button>
                </Grid>
              </Grid>

              <div>
                <br />
                {props.progress ? (
                  <LinearProgress variant="determinate" value={props.progress || 0} />
                ) : (
                  ''
                )}
              </div>
            </div>
          </Paper>
        </Grid>

        <Grid item xs={4}>
          <Paper className={classes.paper}>
            <Grid container>
              <Grid item xs={12}>
                <div className={classes.headings}>Keyboard Events</div>
                <Divider variant="middle" /> <br />
              </Grid>
              <Grid item xs={12}>
                <Button variant="contained" color="primary" onClick={props.CtrlAltDel}>
                  {' '}
                  Ctrl+alt-Del{' '}
                </Button>
              </Grid>
            </Grid>
          </Paper>
        </Grid>
      </Grid>
    </div>
  );

  return (
    <div>
      <Drawer anchor="top" open={props.openMenu} onClose={props.toggleDrawer}>
        <div>{menuItems}</div>
      </Drawer>
    </div>
  );
}

ConsoleMenu.propTypes = {
  classes: PropTypes.object.isRequired,
};
