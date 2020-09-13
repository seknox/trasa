import { makeStyles, MuiThemeProvider } from '@material-ui/core';
import Button from '@material-ui/core/Button';
import axios from 'axios';
import MUIDataTable, {
  MUIDataTableColumn,
  MUIDataTableMeta,
  MUIDataTableOptions,
} from 'mui-datatables';
import React, { useEffect, useState } from 'react';
import Constants from '../../../Constants';
import { LogtableV2Theme } from '../../../utils/styles/themes';

const useStyles = makeStyles((theme) => ({
  mainContent: {
    flex: 1,
    padding: '48px 36px 0',
    background: '#eaeff1', // '#eaeff1',
  },
  paper: {
    maxWidth: 1500,
    margin: 'auto',
    marginTop: 50,
    height: 'auto',
    // overflow: 'hidden',
  },
  formControl: {
    margin: theme.spacing(1),
    minWidth: 120,
  },
  selectEmpty: {
    marginTop: theme.spacing(2),
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
    height: '1000px',
  },
  secondaryBar: {
    zIndex: 0,
  },

  svg: {
    width: 100,
    height: 100,
  },
  polygon: {
    fill: theme.palette.common.white,
    stroke: theme.palette.divider,
    strokeWidth: 1,
  },
}));

const tableBodyFont = {
  backgroundColor: 'Transparent',
  borderRadius: 3,
  border: 0,
  fontColor: 'dimgray',
  fontSize: '12px',
  fontFamily: 'Open Sans, Rajdhani',
};

// LiveSessionTable lists live sessions and lets admins connect to supported live session services.
export default function LiveSessionTable() {
  const classes = useStyles();
  const [eventData, seteventData] = useState([]);
  const [userData, setuserData] = useState({ User: { email: '' } });
  const [loader, setLoader] = useState(false);

  const connect = () => {
    const data = {
      //session: localStorage.getItem('X-SESSION'),
      csrf: localStorage.getItem('X-CSRF'),
    };

    const wskt = new WebSocket(
      `${Constants.TRASA_HOSTNAME_WEBSOCKET}/api/v1/logs/livesessions`,
      'livesessions',
    );

    wskt.onclose = (e: any) => {
      console.log('Socket is closed. Reconn will be attempted in 1 second.', e.reason);
      // setTimeout(() => {
      //   connect();
      // }, 1000);
    };

    wskt.onerror = (err: any) => {
      console.error('Socket encountered error: ', err.message, 'Closing socket');
      wskt.close();
    };

    wskt.onopen = () => {
      wskt.send(JSON.stringify(data));
      setInterval(() => {
        wskt.send('pong');
      }, 5000);
    };
    wskt.onmessage = (evt: any) => {
      //  console.log('--------------- ', evt.data)
      const result = JSON.parse(evt.data) || [];
      // alert("Message is received...");
      let dataArr = [];
      dataArr = result.map(function (n: any) {
        n = JSON.parse(n);

        return [
          n.email ? n.email : n.username,
          n.privilege,
          n.serviceName,
          n.serviceType,
          n.userIP,
          n.serverIP,
          n.loginTime,
          n.connID,
          n.serviceID || '',
        ];
      });

      seteventData(dataArr);
      setLoader(false);
    };
  };

  useEffect(() => {
    connect();

    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/my`)
      .then((response) => {
        response.data.data && setuserData(response.data.data[0]);
      })
      .catch((error) => {
        console.log(error);
      });
  }, []);

  const options = {
    filter: true,
    responsive: 'scrollMaxHeight',
    // count:100,
    //  rowsPerPage: 100,
  };

  const columns = [
    {
      name: 'User',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'Username',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'Service name',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'Endpoint Type',
      options: {
        filter: false,
        display: false,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'Client IP',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },

    {
      name: 'Server IP',
      options: {
        filter: false,
        display: true,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'Logged In Time',
      options: {
        filter: false,
        customBodyRender: (value: any) => {
          return <div style={tableBodyFont}> {new Date(value / 1000000).toString()} </div>;
        },
      },
    },
    {
      name: 'Join Session',
      options: {
        filter: false,
        customBodyRender: (
          value: any,
          tableMeta: MUIDataTableMeta,
          // updateValue: (value: string) => void,
        ) => {
          return (
            <div>
              {tableMeta.rowData[3] !== 'http' ? (
                <Button
                  variant="outlined"
                  color="secondary"
                  onClick={() => {
                    if (tableMeta.rowData[3] === 'ssh') {
                      // serviceID/:username/:email/:connID/:hostname/:totp
                      window.open(
                        `/monitor/sessions/joinxterm#serviceID=${
                          tableMeta.rowData[8]
                        }&username=${encodeURIComponent(
                          encodeURIComponent(tableMeta.rowData[1]),
                        )}&email=${userData.User.email}&connID=${value}&hostname=${
                          tableMeta.rowData[5]
                        }`,
                      );
                    } else {
                      // join/:serviceID/:username/:email/:connID/:hostname
                      window.open(
                        `/monitor/sessions/join#serviceID=${
                          tableMeta.rowData[8]
                        }&username=${encodeURIComponent(
                          encodeURIComponent(tableMeta.rowData[1]),
                        )}&email=${userData.User.email}&connID=${value}&hostname=${
                          tableMeta.rowData[5]
                        }`,
                      );
                    }
                  }}
                >
                  View
                </Button>
              ) : null}
            </div>
          );
        },
      },
    },
    {
      name: 'App ID',
      options: {
        display: false,
      },
    },
  ];

  return (
    <div className={classes.contentWrapper}>
      <MuiThemeProvider theme={LogtableV2Theme}>
        <MUIDataTable
          title="Active Sessions ( LIVE feed )"
          data={eventData}
          columns={columns as MUIDataTableColumn[]}
          options={options as MUIDataTableOptions}
        />
      </MuiThemeProvider>
    </div>
  );
}
