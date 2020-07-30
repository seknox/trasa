import { createMuiTheme, MuiThemeProvider } from '@material-ui/core';
import Button from '@material-ui/core/Button';
// import UserCreateUpdateDrawer from './UserCreateUpdateDrawer';
// import DeleteConfirmDialogue from '../../../components/ui/confirms'
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import { makeStyles } from '@material-ui/core/styles';
import axios from 'axios';
import MUIDataTable, { MUIDataTableColumn, MUIDataTableOptions } from 'mui-datatables';
import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import Constants from '../../../Constants';
import ProgressHOC from '../../../utils/Components/Progressbar';
import CreateUser from './UserCreateDrawer';

const MuiDataTableTheme = createMuiTheme({
  typography: { fontFamily: 'Open Sans, Rajdhani' },
  palette: {
    type: 'light',
    primary: { 500: '#000080' },
    secondary: { A400: '#000080' }, // '#000080' },
  },
});

const useStyles = makeStyles((theme) => ({
  mainContent: {
    flex: 1,
    padding: '48px 36px 0',
    background: '#eaeff1', // '#eaeff1',
  },
  success: {
    paddingLeft: 5,
    paddingRight: 0,
    maxWidth: 50,
    background: 'green',
    color: 'white',
  },
  failed: {
    paddingLeft: 10,
    // paddingRight: 5,
    maxWidth: 50,
    background: 'maroon',
    color: 'white',
  },
  paper: {
    maxWidth: 1500,
    margin: 'auto',
    //  marginTop: 100,
    overflow: 'hidden',
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
  // button: {
  //     borderColor: lightColor,
  //     color: 'white'
  // },
  svg: {
    width: 100,
    height: 100,
  },
  polygon: {
    fill: theme.palette.common.white,
    stroke: theme.palette.divider,
    strokeWidth: 1,
  },
  headers: {
    color: 'black',
    fontSize: '12px',
    fontFamily: 'Open Sans, Rajdhani',
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

type UsertableProps = {
  groupID: string;
  isGroup: boolean;
  allUsers: any;
  deleteUsers: () => void;
  // updateUserTable: () => void;
};

export default function UserTable(props: UsertableProps) {
  const classes = useStyles();
  const [users, setusers] = useState<any>([]);
  // const [usersObj, setusersObj] = useState({});
  // const [deleteDialogueOpen, setdeleteDialogueOpen] = useState(false);
  const [dlgOpen, setdlgOpen] = useState(false);
  const [loader, setLoader] = useState(false);

  const handleDlgClose = () => {
    setdlgOpen(false);
  };

  // const handleDlgOpen = () => {
  //   setdlgOpen(true);
  // };

  useEffect(() => {
    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/user/all`)
      .then((response) => {
        let dataArr = [];
        const data = response.data.data[0];
        dataArr = data.map(function (n: any) {
          const date = new Date(n.CreatedAt * 1000);
          return [
            n.email.toString(),
            n.firstName,
            n.lastName,
            n.userName.toString(),
            date.toDateString(),
            date.toDateString(),
            n.idpName,
            n.status,
            n.ID,
          ];
        });
        setusers(dataArr);
        // setusersObj(response.data.data[0]);
      })
      .catch((error) => {
        console.log(error);
      });
  }, []);

  const updateUserTable = (n: any) => {
    setLoader(true);
    const userArr = users;
    const date = new Date(n.CreatedAt * 1000);
    userArr.unshift([
      n.email.toString(),
      n.firstName,
      n.lastName,
      n.userName.toString(),
      n.userRole.toString(),
      date.toDateString(),
      n.idpName,
      n.status,
      n.ID,
    ]);

    setusers(userArr);
    setLoader(false);
  };

  const statusDiv = (val: boolean) => {
    if (val) {
      return <div className={classes.success}>enabled</div>;
    }
    return <div className={classes.failed}>disabled</div>;
  };

  const columns = [
    {
      name: 'Email',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'First Name',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'Last Name ',
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
        filter: false,
        customBodyRender: (value: any) => {
          return value;
        },
      },
    },
    {
      name: 'User Role',
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
        filter: false,
        customBodyRender: (value: any) => {
          return <div style={tableBodyFont}> {value} </div>;
        },
      },
    },
    {
      name: 'Identity Provider',
      options: {
        filter: false,
        customBodyRender: (value: any) => {
          return <div style={tableBodyFont}> {value} </div>;
        },
      },
    },
    {
      name: 'Status',
      options: {
        filter: true,
        filterOptions: [false, true],
        customBodyRender: (value: any) => {
          return statusDiv(value);
        },
      },
    },
    {
      name: 'View Profile',
      options: {
        filter: false,
        customBodyRender: (value: any) => {
          return (
            <Button component={Link} to={`/users/user/${value}`} variant="outlined" color="primary">
              View
            </Button>
          );
        },
      },
    },
  ];

  const options = {
    filter: true,
    responsive: 'scrollMaxHeight',
    onRowsDelete: props.deleteUsers,
    // isRowSelectable:()=>false,
    selectableRows: props.isGroup ? 'multiple' : 'none',

    // onTableChange:(a,b)=>{console.log(a,b);console.log(this.state.users)},
  };

  return (
    <div>
      {props.isGroup ? '' : <CreateUser updateUserTable={updateUserTable} />}
      <MuiThemeProvider theme={MuiDataTableTheme}>
        <MUIDataTable
          title="Users"
          data={props.allUsers}
          columns={columns as MUIDataTableColumn[]}
          options={options as MUIDataTableOptions}
        />
      </MuiThemeProvider>
      <Dialog
        open={dlgOpen}
        fullWidth
        onClose={handleDlgClose}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
      >
        <DialogTitle id="alert-dialog-title">Deleting Users</DialogTitle>
        <DialogContent>
          <DialogContentText id="alert-dialog-description">
            {loader ? <ProgressHOC /> : ''}
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button variant="contained" onClick={handleDlgClose} color="primary" autoFocus>
            Close
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
}
