import { makeStyles } from '@material-ui/core/styles';
import axios from 'axios';
import echarts from 'echarts';
import React, { useEffect, useState } from 'react';
import Constants from '../Constants';
import Header from './Header';
import Navigator from './Navigator';

const colorPalette = ['#000066', '#37A2DA', '#030417', '#67E0E3', '#1B2948', '#32C5E9', '#000086'];

// #00081C

echarts.registerTheme('trasaTheme', {
  color: colorPalette,
  // height: '500',
  // backgroundColor: '#030417'
  textStyle: {
    fontFamily: 'Open Sans, Rajdhani',
  },
});

const useStyles = makeStyles((theme) => ({
  root: {
    display: 'flex',
    minHeight: '100vh',
  },
  drawer: {
    [theme.breakpoints.up('xl')]: {
      // width: drawerWidth,
      flexShrink: 0,
    },
  },
  appContent: {
    flex: 1,
    display: 'flex',
    flexDirection: 'column',
  },
  mainContent: {
    flex: 1,
    // padding: '48px 36px 0',
    //  background: '#030417', //'#eaeff1',
  },
}));

type DashboardBaseProps = {
  children?: React.ReactNode;
};

export default function DashboardBase(props: DashboardBaseProps) {
  const classes = useStyles();
  const [userData, setuserData] = useState({ firstName: 'U', lastName: '' });
  const [userRole, setuserRole] = useState('');
  const [orgData, setorgData] = useState({ orgName: 'seknox org' });
  const [anchorState, setAnchorState] = React.useState(false);

  useEffect(() => {
    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/my`)
      .then((r) => {
        setuserData(r.data.data[0].User);
        setuserRole(r.data.data[0].User.userRole);
        setorgData(r.data.data[0].Org);
      })
      .catch((error) => {
        console.log(error);
      });
  }, []);

  function toggleDrawer(event: React.KeyboardEvent | React.MouseEvent) {
    if (
      event.type === 'keydown' &&
      ((event as React.KeyboardEvent).key === 'Tab' ||
        (event as React.KeyboardEvent).key === 'Shift')
    ) {
      return;
    }

    setAnchorState(!anchorState);
  }
  return (
    // <TrasaTheme>
    <div className={classes.root}>
      <nav className={classes.drawer}>
        <Navigator
          // PaperProps={{ style: { width: drawerWidth } }}
          open={anchorState}
          anchor="left"
          onClose={toggleDrawer}
          orgData={orgData}
          userrole={userRole}
        />
      </nav>
      <div className={classes.appContent}>
        <Header onDrawerToggle={toggleDrawer} userData={userData} />
        <main className={classes.mainContent}>{props.children}</main>
      </div>
    </div>
    // </TrasaTheme>
  );
}
