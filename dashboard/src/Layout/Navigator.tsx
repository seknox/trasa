import Divider from '@material-ui/core/Divider';
import Drawer from '@material-ui/core/Drawer';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import { makeStyles } from '@material-ui/core/styles';
import DashbordIcon from '@material-ui/icons/Assessment';
import ServiceIcon from '@material-ui/icons/Apps';
import ManageIcon from '@material-ui/icons/AccountTree';
import SessionsIcon from '@material-ui/icons/MissedVideoCallSharp';
import ControlIcon from '@material-ui/icons/SecuritySharp';
import SettingsIcon from '@material-ui/icons/Settings';
import Account from '@material-ui/icons/SupervisedUserCircle';
import classNames from 'classnames';
import React from 'react';
import { Link } from 'react-router-dom';
// import PhonelinkSetupIcon from '@material-ui/icons/PhonelinkSetup';
import logo from '../assets/trasa-white.svg';
// import ChevronRightIcon from '@material-ui/icons/ChevronRight';

const categories = () => {
  const ar = [];
  ar.push({
    id: 'Overview',
    path: 'overview',
    icon: <DashbordIcon />,
    children: [],
  });

    ar.push({
    id: 'Providers',
    path: 'providers',
    icon: <ManageIcon />,
    children: [],
    // children: [
    //   { id: 'Users', path: 'users/groups', icon: <Account /> },
    //   { id: 'Services', path: 'services/groups', icon: <ServiceIcon /> },
    // ],
  });

  ar.push({
    id: 'Users',
    path: 'users/groups',
    icon: <Account />,
    children: [],
  });

  ar.push({
    id: 'Services',
    path: 'services/groups',
    icon: <ServiceIcon />,
    children: [],
  });

  ar.push({
    id: 'Control',
    path: 'control',
    icon: <ControlIcon />,
    children: [],
  });

  ar.push({
    id: 'Monitor',
    path: 'monitor/sessions',
    icon: <SessionsIcon />,
    children: [],
  });



  ar.push({
    id: 'System',
    path: 'system',
    icon: <SettingsIcon />,
    children: [],
  });

  // }

  return ar;
};

const NormalUsersMenu = [
  {
    id: 'My',
    path: 'my',
    icon: <Account />,
    children: [
      // { id: 'Group', path: 'groups', icon: <SettingsIcon /> },
    ],
  },
];

const useStyles = makeStyles((theme) => ({
  categoryHeader: {
    color: 'black', // 'black',
    fontSize: '15px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  categoryHeaderPrimary: {
    color: 'black', // 'black',
    fontSize: '18px',
    fontWeight: 500,
    fontFamily: 'Open Sans, Rajdhani',
  },
  item: {
    color: 'black', // '#404854',
  },
  itemCategory: {
    boxShadow: '0 -1px 0 #404854 inset',
    color: 'white',
  },
  firebase: {
    fontSize: 24,
    backgroundColor: '#eeeeee', // 'rgb(3,4,27)',   //'#030417', // rgb(3,4,27)
    fontFamily: theme.typography.fontFamily,
    color: theme.palette.common.white,
  },
  itemActionable: {
    '&:hover': {
      backgroundColor: 'rgba(255, 255, 255, 0.08)',
    },
  },
  itemActiveItem: {
    color: '#4fc3f7',
  },
  itemPrimary: {
    color: 'black', // 'black',
    fontSize: '15px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  textDense: {},
  divider: {
    //   marginTop: theme.spacing(2),
  },
  logoImage: {
    width: 30,
    height: 20,
  },
  toolbar: {
    color: 'white',
    backgrounColor: 'white',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'flex-end',
    padding: '0 8px',
    ...theme.mixins.toolbar,
  },
  trasaLogo: {
    backgroundColor: 'rgba(1,1,35,1)', // 'rgba(1,1,35,1)','#000066'
    paddingLeft: '15px',
    paddingRight: '15px',
    paddingTop: '35px',
    paddingBottom: '20px',
  },
}));

export default function Navigator(props: any) {
  const classes = useStyles();
  // const [open, setopen] = React.useState(false);

  return (
    <Drawer anchor={props.anchor} open={props.open} onClose={props.onClose}>
      <div className={classes.trasaLogo}>
        <img src={logo} alt="logo" width={180} />
      </div>
      <Divider className={classes.divider} />
      <List disablePadding>
        {props.userrole === 'selfUser'
          ? NormalUsersMenu.map(({ id, children, path, icon }) => (
              <React.Fragment key={id}>
                <ListItem
                  className={classes.categoryHeader}
                  button
                  key={id}
                  component={Link}
                  to={`/${path}`}
                >
                  <ListItemIcon className={classNames(classes.item, classes.itemCategory)}>
                    {icon}
                  </ListItemIcon>
                  <ListItemText
                    classes={{
                      primary: classes.categoryHeaderPrimary,
                    }}
                  >
                    {id}
                  </ListItemText>
                </ListItem>
                {children.map(({ id: childId, path: childPath, icon }) => (
                  <ListItem
                    button
                    dense
                    key={childId}
                    className={classNames(
                      classes.item,
                      classes.itemActionable,
                      // active && classes.itemActiveItem,
                    )}
                    component={Link}
                    to={`/${path}/${childPath}`}
                  >
                    <ListItemIcon>{icon}</ListItemIcon>
                    <ListItemText
                      classes={{
                        primary: classes.itemPrimary,
                        // textDense: classes.textDense,
                      }}
                    >
                      {childId}
                    </ListItemText>
                  </ListItem>
                ))}
                {/* <Divider className={classes.divider} /> */}
              </React.Fragment>
            ))
          : categories().map(({ id, children, path, icon }) => (
              <React.Fragment key={id}>
                <ListItem
                  className={classes.categoryHeader}
                  button
                  key={id}
                  component={Link}
                  to={`/${path}`}
                >
                  <ListItemIcon>{icon}</ListItemIcon>
                  <ListItemText
                    classes={{
                      primary: classes.categoryHeaderPrimary,
                    }}
                  >
                    <b>{id} </b>
                  </ListItemText>
                </ListItem>
                {children.map(({ id: childId, path: childPath, icon }) => (
                  <ListItem
                    button
                    dense
                    key={childId}
                    className={classNames(
                      classes.item,
                      classes.itemActionable,
                      // active && classes.itemActiveItem,
                    )}
                    component={Link}
                    to={`/${path}/${childPath}`}
                  >
                    <ListItemIcon>{icon}</ListItemIcon>
                    <ListItemText
                      classes={{
                        primary: classes.itemPrimary,
                        // textDense: classes.textDense,
                      }}
                    >
                      {childId}
                    </ListItemText>
                  </ListItem>
                ))}
                {/* <Divider light /> */}
              </React.Fragment>
            ))}
      </List>
    </Drawer>
  );
}
