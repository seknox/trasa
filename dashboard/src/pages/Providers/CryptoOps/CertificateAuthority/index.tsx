import Button from '@material-ui/core/Button';
import FormControl from '@material-ui/core/FormControl';
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import MenuItem from '@material-ui/core/MenuItem';
import Paper from '@material-ui/core/Paper';
import Select from '@material-ui/core/Select';
import { makeStyles } from '@material-ui/core/styles';
import Toolbar from '@material-ui/core/Toolbar';
import Tooltip from '@material-ui/core/Tooltip';
import Typography from '@material-ui/core/Typography';
import DownloadIcon from '@material-ui/icons/CloudDownloadOutlined';
import axios from 'axios';
import MUIDataTable, { MUIDataTableColumn, MUIDataTableMeta } from 'mui-datatables';
import React, { useState } from 'react';
import CAIcon from '../../../../assets/ca.svg';
import Constants from '../../../../Constants';
import DialogueWrapper from '../../../../utils/Components/DialogueWrapComponent';

const fileDownload = require('js-file-download');

const useStyles = makeStyles((theme) => ({
  root: {
    // flexgrow: 1,
    // marginBotton: '5%',
  },

  paper: {
    backgroundColor: '#fdfdfd',
    padding: theme.spacing(2),
    //  textAlign: 'center',
    color: theme.palette.text.secondary,
  },

  initButton: {
    background: '#0A2053', // '#0A2053',
    borderRadius: 3,
    border: 0,
    color: 'white',
    height: 38,
    padding: '0 30px',
    '&:hover': {
      color: theme.palette.common.white,
      background: 'skyblue',
    },
    // boxShadow: '0 3px 5px 2px rgba(255, 105, 135, .3)',
  },
  initedButton: {
    background: '#025531',
    borderRadius: 3,
    border: 0,
    color: 'white',
    height: 38,
    padding: '0 30px',
  },
  label: {
    textTransform: 'capitalize',
  },
  heading: {
    color: '#025531',
    fontSize: '24px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  keys: {
    color: 'navy',
    fontSize: '16px',
  },
  initStatusText: {
    color: 'grey',
    fontSize: '14px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  Warning: {
    color: 'maroon',
    fontSize: '44px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  WarningButton: {
    color: 'maroon',
    //  fontSize: '44px',
    fontFamily: 'Open Sans, Rajdhani',
    '&:hover': {
      color: theme.palette.common.white,
      background: 'maroon',
    },
  },
  selectCustom: {
    fontSize: 15,
    borderRadius: 4,
    backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    height: 17,
    // marginTop: 5,
    // padding: '10px 100px',
    // width: 100,
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
    },
  },
  settingHeader: {
    color: 'black',
    fontSize: '15px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  settingSHeader: {
    color: 'black',
    fontSize: '12px',
    fontFamily: 'Open Sans, Rajdhani',
  },
}));

export default function CA(props: any) {
  const classes = useStyles();
  const [caDetail, setCADetail] = React.useState({ subject: {}, issuer: {} });
  const [caArr, setCaArr] = useState([]);
  const [viewDlg, openViewDlg] = useState(false);

  const changeViewDlgState = (rowIndex: any) => {
    openViewDlg(!viewDlg);
    const row = caArr[rowIndex];
    // console.log(row)
    setCADetail({ subject: row[6], issuer: row[7] });
  };

  const closeViewDlg = () => {
    openViewDlg(false);
  };

  const generateSSHCA = (type: any) => () => {
    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/system/sshca/init/${type}`)
      .then((response) => {});
  };

  React.useEffect(() => {
    // code to run on component mount
    const config = {
      headers: {
        'X-SESSION': localStorage.getItem('X-SESSION'),
        'X-CSRF': localStorage.getItem('X-CSRF'),
      },
    };
    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/system/ca/all`, config)
      .then((response) => {
        if (response.data.status === 'success') {
          // console.log(response.data.data[0].cert)
          const data = response.data && response.data.data && response.data.data[0];
          //  setCaArr(data.map(ca=>([ca.name,ca.name,ca.createdAt,ca.lastUpdated,ca.certID])));
          let dataArr = [];
          dataArr = data.map(function (n: any) {
            const cdate = new Date(n.createdAt * 1000);
            const udate = new Date(n.lastUpdated * 1000);
            // console.log('cert data: ', n)

            return [
              n.cert ? n.cert.issuer.common_name : n.entityID,
              n.cert && n.cert.issuer.common_name,
              n.certType,
              cdate.toDateString(),
              udate.toDateString(),
              8,
              n.certID,
              n.cert ? n.cert.subject : {},
              n.cert ? n.cert.issuer : {},
              n.certID,
            ];
          });
          setCaArr(dataArr);
        }
      })
      .catch((err) => {
        console.error(err);
      });
  }, []);

  const [certCrudDlgState, changeCertCrudDlgState] = useState(false);
  const [selectedCAOp, setSelectedCAOp] = React.useState('genSSHHostCA');

  const certOpChange = (event: any) => {
    setSelectedCAOp(event.target.value);
  };

  function closeCertCridDlg() {
    changeCertCrudDlgState(false);
  }

  return (
    <Paper className={classes.paper} elevation={0}>
        <Toolbar>
          <Grid container spacing={3} alignItems="center">
            <Grid item xs={2}>
              <Button variant="contained" onClick={() => changeCertCrudDlgState(!certCrudDlgState)}>
                Generate certs
              </Button>
              <br />
            </Grid>
            <Grid item>
              <Tooltip title="Reload">
                <IconButton>

                </IconButton>
              </Tooltip>
            </Grid>
          </Grid>
        </Toolbar>


      <CATable caArr={caArr} changeViewDlgState={changeViewDlgState} />

      <DialogueWrapper
        open={certCrudDlgState}
        handleClose={closeCertCridDlg}
        title="CA cert ops"
        maxWidth="md"
        fullScreen={false}
      >
        <CAcrud
          selectedCAOp={selectedCAOp}
          certOpChange={certOpChange}
          generateSSHCA={generateSSHCA}
        />
      </DialogueWrapper>

      <DialogueWrapper
        open={viewDlg}
        handleClose={closeViewDlg}
        title="CA detail"
        maxWidth="md"
        fullScreen={false}
      >
        <CADetail caDetail={caDetail} />
      </DialogueWrapper>
    </Paper>
  );
}

function CAcrud(props: any) {
  const classes = useStyles();

  function renderComp() {
    switch (true) {
      // case (props.selectedCAOp === 'genRootCA'):
      //   return (<Grid container spacing={2}  direction="column" justify="center" alignItems="center">
      //           <Grid item xs={12}>
      //           <Typography variant="h1" component="h1"><b>Generate root CA.</b></Typography>
      //           </Grid>

      //           <Grid item xs={12}>
      //           <br />
      //           <CreateCA />
      //           </Grid>
      //            </Grid>)
      // case (props.selectedCAOp === 'uploadRootCA'):
      // return (<Grid container spacing={2}  direction="column" justify="center" alignItems="center">
      //         <Grid item xs={12}>
      //         <Typography variant="h1" component="h1"><b>Upload root CA files</b></Typography>
      //         </Grid>

      //         <Grid item xs={12}>
      //         <br />
      //         <UploadCA />
      //         </Grid>
      //          </Grid>)
      case props.selectedCAOp === 'genSSHHostCA':
        return (
          <Grid container spacing={2} direction="column" justify="center" alignItems="center">
            <Grid item xs={12}>
              <Typography variant="h1" component="h1">
                <b>Generate SSH Host CA</b>
              </Typography>
            </Grid>

            <Grid item xs={12}>
              <Button variant="contained" onClick={props.generateSSHCA('host')}>
                Generate SSH Host CA
              </Button>
            </Grid>
          </Grid>
        );
      case props.selectedCAOp === 'genSSHUserCA':
        return (
          <Grid container spacing={2} direction="column" justify="center" alignItems="center">
            <Grid item xs={12}>
              <Typography variant="h1" component="h1">
                <b>Generate SSH User CA</b>
              </Typography>
            </Grid>

            <Grid item xs={12}>
              <br />
              <Button variant="contained" onClick={props.generateSSHCA('user')}>
                Generate SSH User CA
              </Button>
            </Grid>
          </Grid>
        );
      default:
        return null;
    }
  }

  return (
    <Grid container spacing={2} direction="row" justify="center" alignItems="center">
      <Grid item xs={12}>
        <Typography variant="h3"> Select cert operation :</Typography>
      </Grid>
      <Grid item xs={12}>
        <FormControl fullWidth>
          <Select
            name="certOpType"
            defaultValue={props.selectedCAOp}
            onChange={props.certOpChange}
            inputProps={{
              classes: {
                root: classes.selectCustom,
              },
            }}
          >
            {/* <MenuItem value={'genRootCA'}><div className={classes.settingSHeader}>Generate root CA </div></MenuItem>
               <MenuItem value={'uploadRootCA'}><div className={classes.settingSHeader}>Upload root CA</div></MenuItem> */}
            <MenuItem value="genSSHHostCA">
              <div className={classes.settingSHeader}>Generate SSH host CA</div>
            </MenuItem>
            <MenuItem value="genSSHUserCA">
              <div className={classes.settingSHeader}>Generate SSH user CA</div>
            </MenuItem>
          </Select>
        </FormControl>
      </Grid>

      <Grid item xs={12}>
        {renderComp()}
      </Grid>
    </Grid>
  );
}

function CATable(props: any) {
  const downloadCA = (name: any, type: any) => () => {
    switch (type) {
      case 'SSH_CA':
        axios.get(`${Constants.TRASA_HOSTNAME}/api/v1/system/ca/ssh/${name}`).then((response) => {
          fileDownload(response.data, 'ca-cert.pem', 'application/x-pem-file');
        });
        break;
      case 'HTTP_CA':
        break;
      default:
        break;
      // TODO download http CA
    }
  };
  const downloadable = (name: any, type: any) => {
    if (type !== 'SSH_CA') {
      return false;
    }
    if (name === 'host' || name === 'user') {
      return true;
    }
    return false;
  };
  const columns = [
    {
      name: 'CA Name',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'Issuer',
      options: {
        filter: true,
        display: false,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'CA Type',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'Created At',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'Last updated',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },

    {
      name: 'View Details',
      options: {
        filter: true,
        display: false,
        customBodyRender: (value: any, tableMeta: MUIDataTableMeta, updateValue: any) => {
          // console.log('value', value, 'tableMeta', tableMeta, 'updateValue', updateValue )
          return (
            <Button
              variant="outlined"
              size="small"
              onClick={() => {
                props.changeViewDlgState(tableMeta.rowIndex);
              }}
            >
              View
            </Button>
          );
        },
      },
    },
    {
      // TODO bhrg3se build your certificate download option here.
      name: 'Download',
      options: {
        filter: true,
        customBodyRender: (value: any, tableMeta: MUIDataTableMeta, updateValue: any) => {
          // console.log('value', value, 'tableMeta', tableMeta, 'updateValue', updateValue )
          return (
            <div>
              {downloadable(tableMeta.rowData[0], tableMeta.rowData[2]) && (
                <IconButton color="secondary" aria-label="dlownload certificate">
                  <DownloadIcon onClick={downloadCA(tableMeta.rowData[0], tableMeta.rowData[2])} />
                </IconButton>
              )}
            </div>
          );
        },
      },
    },
  ];

  return (
    <MUIDataTable
      title="Available Certificate Authorities"
      data={props.caArr}
      columns={columns as MUIDataTableColumn[]}
      options={{
        filter: true,
        responsive: 'scrollMaxHeight',
        // resizableColumns: true,
        selectableRows: 'single',
        onRowsDelete() {},
      }}
    />
  );
}

function CADetail(props: any) {
  const classes = useStyles();
  const { caDetail } = props;
  return (
    <div>
      <Grid item xs={12}>
        <Grid container spacing={2}>
          <Grid item xs={5}>
            <img src={CAIcon} alt="caICON" />
          </Grid>

          <Grid item xs={7}>
            <Grid container spacing={2}>
              <Grid item xs={12}>
                <div className={classes.heading}>
                  <Typography variant="h1"> Certificate Authority (CA)</Typography>
                  <Typography variant="h4">
                    {' '}
                    This CA will be used to issue and verify client and server certificates
                  </Typography>
                </div>
              </Grid>
            </Grid>

            <br />
          </Grid>
        </Grid>
      </Grid>
      <Grid container>
        <Grid container item>
          <Grid item xs={12}>
            <Typography variant="h3">Subject</Typography>
          </Grid>
        </Grid>

        <Grid container item>
          <Grid item xs={6}>
            <Typography>CA Name</Typography>
          </Grid>
          <Grid item xs={6}>
            <Typography>{caDetail.subject.common_name} </Typography>
          </Grid>
        </Grid>
        <Grid container item>
          <Grid item xs={6}>
            <Typography>Country</Typography>
          </Grid>
          <Grid item xs={6}>
            <Typography>{caDetail.subject.country} </Typography>
          </Grid>
        </Grid>
        <Grid container item>
          <Grid item xs={6}>
            <Typography>Organization</Typography>
          </Grid>
          <Grid item xs={6}>
            <Typography>{caDetail.subject.organization} </Typography>
          </Grid>
        </Grid>

        <Grid container item>
          <Grid item xs={12}>
            <Typography variant="h3">Issuer</Typography>
          </Grid>
        </Grid>

        <Grid container item>
          <Grid item xs={6}>
            <Typography>CA Name</Typography>
          </Grid>
          <Grid item xs={6}>
            <Typography>{caDetail.issuer.common_name} </Typography>
          </Grid>
        </Grid>
        <Grid container item>
          <Grid item xs={6}>
            <Typography>Country</Typography>
          </Grid>
          <Grid item xs={6}>
            <Typography>{caDetail.issuer.country} </Typography>
          </Grid>
        </Grid>
        <Grid container item>
          <Grid item xs={6}>
            <Typography>Organization</Typography>
          </Grid>
          <Grid item xs={6}>
            <Typography>{caDetail.issuer.organization} </Typography>
          </Grid>
        </Grid>
      </Grid>
    </div>
  );
}
