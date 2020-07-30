import { makeStyles } from '@material-ui/core';
import AppBar from '@material-ui/core/AppBar';
import Grid from '@material-ui/core/Grid';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import axios from 'axios';
import QueryString from 'query-string';
import React, { useEffect, useState } from 'react';
import ReactPlayer from 'react-player';
import Constants from '../../../Constants';
import SSHLogXterm from './TrasaSSHGWLog';

function TabContainer(props: any) {
  return (
    <Typography component="div" style={{ padding: 8 * 3 }}>
      {props.children}
    </Typography>
  );
}

const lightColor = 'rgba(255, 255, 255, 0.7)'; // 'rgba(255, 255, 255, 0.7)'; // '#030417';

const useStyles = makeStyles((theme) => ({
  root: {
    flexgrow: 1,
  },
  paper: {
    backgroundColor: '#fdfdfd',
    padding: theme.spacing(2),
    color: theme.palette.text.secondary,
  },
  mainContent: {
    flex: 1,
    padding: '48px 36px 0',
    background: 'white', // '#eaeff1',
  },
  searchBar: {
    borderBottom: '1px solid rgba(0, 0, 0, 0.12)',
  },
  searchInput: {
    fontSize: theme.typography.fontSize,
  },
  block: {
    display: 'block',
  },
  addUser: {
    marginRight: theme.spacing(1),
  },
  contentWrapper: {
    margin: '40px 16px',
  },
  secondaryBar: {
    zIndex: 0,
  },
  button: {
    borderColor: lightColor,
  },
}));

function RecordedSession(props: any) {
  const classes = useStyles();
  //const [value, setvalue] = useState(props.sessionType === 'guac' ? 1 : 0);
  // const [httpRaw, sethttpRaw] = useState('');
  const [sskey, setsskey] = useState('');
  const [sessionLog, setsessionLog] = useState('');
  const [tabIndex, setTabIndex] = useState(0);

  useEffect(()=>{
    setTabIndex(props.sessionType == 'guac' || props.sessionType == 'rdp'?1:0)
  },[props.sessionType])

  useEffect(() => {
    axios
      .get(
        `${Constants.TRASA_HOSTNAME}/api/v1/logs/sessionlog?sessionID=${props.sessionID}&type=${props.sessionType}&day=${props.day}&month=${props.month}&year=${props.year}`,
      )
      .then((response) => {
        setsskey(response.headers.sskey);
        if (response.data === '') {
          setsessionLog('Log Not Found');
        } else {
          setsessionLog(response.data);
        }
      })
      .catch((error) => {
        console.log(error);
      });
  }, [props.sessionID, props.orgID, props.sessionType, props.day, props.month, props.year]);

  const handleChange = (event: React.ChangeEvent<{}>, tvalue: number) => {
    console.log('handleChange', tvalue)
    setTabIndex(tvalue);
    //setvalue(tvalue)
  };

  return (
    <div>
      <AppBar
        component="div"
        className={classes.secondaryBar}
        color="primary"
        position="static"
        elevation={0}
      >
        <Toolbar>
          <Grid container alignItems="center" spacing={8}>
            <Grid item xs>
              <Typography color="inherit" variant="h5">
                View Recorded Session
              </Typography>
            </Grid>
          </Grid>
        </Toolbar>
      </AppBar>
      <AppBar
        component="div"
        className={classes.secondaryBar}
        color="primary"
        position="static"
        elevation={0}
      >
        <Tabs value={tabIndex} onChange={handleChange} textColor="inherit">
          {props.sessionType === 'http' ||
          props.sessionType === 'ssh' ||
          props.sessionType === 'guac-ssh' ? (
            <Tab textColor="inherit" label="Raw Log" value={0} />
            ) : null}
          {props.sessionType === 'http' ||
          props.sessionType === 'guac' ||
          props.sessionType === 'rdp' ? (
            <Tab textColor="inherit" label="Video" value={1} />
            ) : null}
        </Tabs>
      </AppBar>

      <div className={classes.mainContent}>
        {tabIndex === 0 &&
          (props.sessionType === 'http' ||
            props.sessionType === 'ssh' ||
            props.sessionType === 'db' ||
            props.sessionType === 'guac-ssh') && (
            <TabContainer>
              {' '}
              <RecordedSessionRaw data={sessionLog} sessionType={props.sessionType} />{' '}
            </TabContainer>
        )}
        {tabIndex === 1 &&
          (props.sessionType === 'http' ||
            props.sessionType === 'guac' ||
            props.sessionType === 'rdp') &&
          sskey && (
            <TabContainer>
              {' '}
              <RecordedSessionVideo
                sessionID={props.sessionID}
                sskey={sskey}
                sessionType={props.sessionType}
                orgID={props.orgID}
                day={props.day}
                month={props.month}
                year={props.year}
              />
            </TabContainer>
        )}
      </div>
    </div>
  );
}

function RecordedSessionRaw(props: any) {
  const loading = props.data === '';
  return (
    <div style={{ whiteSpace: 'pre-wrap' }}>
      {loading ? 'Loading ...' : null}
      {(props.sessionType === 'ssh' || props.sessionType === 'guac-ssh') && !loading ? (
        <SSHLogXterm sessionLog={props.data} />
      ) : (
        props.data
      )}
    </div>
  );
}

// "http://192.168.0.100:3339/api/v1/events/sessionlog?sessionID=f82585772437fbddeaa802322582bc770a&type=http"
function RecordedSessionVideo(props: any) {
  const reqVal = `${Constants.TRASA_HOSTNAME}/api/v1/logs/vsessionlog?sessionID=${props.sessionID}&type=${props.sessionType}&ssKey=${props.sskey}&orgID=${props.orgID}&day=${props.day}&month=${props.month}&year=${props.year}`;

  return (
    <ReactPlayer
      url={reqVal}
      playing
      controls
      width="100%"
      height="100%"
      playbackRate={props.sessionType === 'http' ? 0.1 : 1}
    />
  );
}

// type TParams = {
//   sessionID: string;
//   type: string;
//   orgID: string;
//   day: string;
//   month: string;
//   year: string;
// };

type TParams = {
  sessionID: string | string[] | null | undefined;
  type: string | string[] | null | undefined;
  orgID: string | string[] | null | undefined;
  day: string | string[] | null | undefined;
  month: string | string[] | null | undefined;
  year: string | string[] | null | undefined;
};


const RecordedSessionViewer = (props: any) => {
  const [data, setData] = React.useState<any>({sessionID:'',type: '', day: '', month: '', year: '' })
  React.useEffect(() => {
    const hashed = QueryString.parse(props.location.hash);
    console.log('ProxyDomain: ', hashed);
    if (hashed) {
      setData({sessionID: hashed.sessionID,type: hashed.type, day: hashed.day, month: hashed.month, year: hashed.year})
      console.log('ProxyDomain: ', hashed);
    }
  }, [props.location.hash])
  return (
    <div>
      <RecordedSession
        sessionID={data.sessionID}
        sessionType={data.type}
        day={data.day}
        month={data.month}
        year={data.year}
      />
    </div>
  );
};


// const RecordedSessionViewer = ({ match }: RouteComponentProps<TParams>) => {
//   return (
//     <div>
//       <RecordedSession
//         sessionID={match.params.sessionID}
//         sessionType={match.params.type}
//         // orgID={match.params.orgID}
//         day={match.params.day}
//         month={match.params.month}
//         year={match.params.year}
//       />
//     </div>
//   );
// };

// userarray={location.state.userdetail}   userdetail userarray={location.state.user}
export default RecordedSessionViewer;
