import Button from '@material-ui/core/Button';
import Paper from '@material-ui/core/Paper';
import { makeStyles } from '@material-ui/core/styles';
import axios from 'axios';
import 'date-fns';
import MUIDataTable, {
  MUIDataTableColumn,
  MUIDataTableMeta,
  MUIDataTableOptions,
} from 'mui-datatables';
import React, { useEffect, useState } from 'react';
import Constants from '../../../../Constants';
import ViewAdhocSession from './ViewAdhocSession';

const lightColor = 'rgba(255, 255, 255, 0.7)'; //'rgba(255, 255, 255, 0.7)'; // '#030417';

const useStyles = makeStyles((theme) => ({
  paper: {
    maxWidth: 1500,
    margin: 'auto',
    marginTop: 50,
    overflow: 'hidden',
    padding: theme.spacing(2),
  },
  contentWrapper: {
    margin: '40px 16px',
  },
}));

type adhocProps = {
  orgID: string;
};

export default function AdhocRequestHistory(props: adhocProps) {
  const [requestHistory, setRequestHistory] = useState([]);
  const [reqDlgOpen, setReqDlgOpen] = useState(false);
  const [reqID, setReqID] = useState('');
  const [reqText, setReqtext] = useState('');
  const [haveData, setHaveData] = useState(false);
  const [sessionID, setSessionID] = useState('');
  const [serviceType, setServiceType] = useState('');
  const [localeAuthorizedOn, setlocaleAuthorizedOn] = useState('');
  const [localeAuthorizedPeriod, setlocaleAuthorizedPeriod] = useState('');
  const [timestamp, setTimestamp] = useState('');

  useEffect(() => {
    axios
      .get(Constants.TRASA_HOSTNAME + '/api/v1/iam/adhoc/requests/all')
      .then((response) => {
        //console.log(response.data);
        let data = response.data.data[0];
        if (data != null) {
          let dataArr = [];

          dataArr = data.map(function (n: any) {
            let reqDate = new Date(n.reqTime * 1000);
            let localeAuthorizedOn = 'not authorized';
            if (n.authorizedOn > 0) {
              let authorizedOn = new Date(n.authorizedOn * 1000);
              localeAuthorizedOn = authorizedOn.toString();
            }

            let localeAuthorizedPeriod = 'not authorized';

            if (n.isAuthorized) {
              let authorizedPeriod = new Date(n.authorizedPeriod);
              localeAuthorizedPeriod = authorizedPeriod.toString();
            }

            return [
              n.requesterEmail,
              n.requesteeEmail,
              n.isAuthorized ? 'granted' : 'rejected',
              n.requestTxt,
              n.serviceID,
              reqDate.toString(),
              n.reqID +
                ':-:' +
                n.requestTxt +
                ':-:' +
                n.sessionID +
                ':-:' +
                localeAuthorizedOn +
                ':-:' +
                localeAuthorizedPeriod +
                ':-:' +
                n.authorizedOn +
                ':-:' +
                n.serviceType,
            ];
          });
          setRequestHistory(dataArr);
          setHaveData(true);
        }
      })
      .catch((error) => {
        console.log(error);
      });
  }, []);

  const handleRequestDialogueOpen = (reqID: string) => {
    let reqData = reqID.split(':-:');

    setReqDlgOpen(true);
    setReqID(reqData[0]);
    setReqtext(reqData[1]);
    setSessionID(reqData[2]);
    setServiceType(reqData[6]);
    setlocaleAuthorizedOn(reqData[3]);
    setlocaleAuthorizedPeriod(reqData[4]);
    setTimestamp(reqData[5]);
  };

  const handleRequestDialogueClose = () => {
    setReqDlgOpen(false);
  };

  const classes = useStyles();

  const columns = [
    {
      name: 'Requested By',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'Sanctioned By',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'Status',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'Access Reason',
      options: {
        filter: true,
        display: false,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'Service Name',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'Requested On',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'View Detail',
      options: {
        filter: false,
        customBodyRender: (value: any) => {
          return (
            <Button
              onClick={() => handleRequestDialogueOpen(value)}
              variant="outlined"
              color="secondary"
            >
              View
            </Button>
          );
        },
      },
    },
  ];

  return (
    <div className={classes.contentWrapper}>
      {haveData ? (
        <MUIDataTable
          title={'Adhoc authentication session logs'}
          data={requestHistory}
          columns={columns as MUIDataTableColumn[]}
          options={options as MUIDataTableOptions}
        />
      ) : (
        <Paper className={classes.paper} elevation={1}>
          <p>There is no adhoc access history to show.</p>
        </Paper>
      )}

      <ViewAdhocSession
        reqText={reqText}
        viewAdhocSessionDlgState={reqDlgOpen}
        handleRequestDialogueClose={handleRequestDialogueClose}
        orgID={props.orgID}
        sessionID={sessionID}
        serviceType={serviceType}
        localeAuthorizedOn={localeAuthorizedOn}
        localeAuthorizedPeriod={localeAuthorizedPeriod}
        timestamp={timestamp}
      />
    </div>
  );
}

const options = {
  filter: true,
  responsive: 'scrollMaxHeight',
};
