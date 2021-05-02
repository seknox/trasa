import Box from '@material-ui/core/Box';
import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import Typography from '@material-ui/core/Typography';
import React from 'react';


const useStyles =  makeStyles(theme => ({
    toolTip: {
        padding: 5,
       // maxWidth: 220,
       fontSize: 16,
      // fontFamily: 'Open Sans, Rajdhani',
        border: '1px solid #dadde9',
      },
      selectCustom: {
        fontSize: 15,
        fontFamily: 'Open Sans, Rajdhani',
        borderRadius: 4,
        backgroundColor: theme.palette.common.white,
        border: '1px solid #ced4da',
        height: 31,
        marginTop: 5,
       // padding: '10px 100px',
        width: 'calc(100% - 4px)',
        transition: theme.transitions.create(['border-color', 'box-shadow']),
        '&:focus': {
          borderColor: '#80bdff',
          boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
        },
      },
      tabRoot: {
            flexGrow: 1,
            backgroundColor: theme.palette.background.paper,
            display: 'flex',
            //height: 300,
      },
      tabs: {
        borderRight: `1px solid ${theme.palette.divider}`,
      },
}))


function TabPanel(props: any) {
    const { children, value, index, ...other } = props;
  
    return (
      <Typography
        component="div"
        role="tabpanel"
        hidden={value !== index}
        id={`vertical-tabpanel-${index}`}
        aria-labelledby={`vertical-tab-${index}`}
        {...other}
      >
        {value === index && <Box p={3}>{children}</Box>}
      </Typography>
    );
  }


  function a11yProps(index: any) {
    return {
      id: `vertical-tab-${index}`,
      'aria-controls': `vertical-tabpanel-${index}`,
    };
  }

  export function DeviceHygiene() {
    const  classes = useStyles()

    const [tabValue, setValue] = React.useState(5);

    const handleTabChange = (event: React.ChangeEvent<{}>, newValue: number) => {
      setValue(newValue);
    }

    return (         
        <Grid container spacing={2}  justify="center">
        <div className={classes.tabRoot}>
            <Grid item xs={3}>
            <Tabs
            indicatorColor="primary"
                orientation="vertical"
                variant="scrollable"
                value={tabValue}
                onChange={handleTabChange}
                className={classes.tabs}
            >
                
                <Tab label={<Typography variant="h4">Login Security</Typography>} {...a11yProps(0)} />
                <Tab label={<Typography variant="h4">OS Security</Typography>}  {...a11yProps(1)} />
                <Tab label={<Typography variant="h4">Endpoint Security</Typography>}  {...a11yProps(2)} /> 
                <Tab label={<Typography variant="h4">Network Security</Typography>}  {...a11yProps(2)} /> 
                <Tab label={<Typography variant="h4">Version Blocking</Typography>}  {...a11yProps(2)} /> 
            </Tabs>

                </Grid>
            
            <Grid item xs={12} sm={12} md={9}>

            <TabPanel />

            </Grid>
            </div>
        </Grid>
    )
}


// DevicePolicy wraps device policy components
export function MobileDevicePolicy() {
    const  classes = useStyles()

    const [tabValue, setValue] = React.useState(5);

    const handleTabChange = (newValue: any) => {
      setValue(newValue);
    }

    return (         
        <Grid container spacing={2}  justify="center">
        <div className={classes.tabRoot}>
            <Grid item xs={3}>
            <Tabs
            indicatorColor="primary"
                orientation="vertical"
                variant="scrollable"
                value={tabValue}
                onChange={handleTabChange}
                className={classes.tabs}
            >
                
                <Tab label={<Typography variant="h4">Baseline</Typography>} {...a11yProps(0)} />
                <Tab label={<Typography variant="h4">Android</Typography>}  {...a11yProps(1)} />
                <Tab label={<Typography variant="h4">IOS</Typography>}  {...a11yProps(2)} /> 
    
            </Tabs>

                </Grid>
            
            <Grid item xs={12} sm={12} md={9}>

            <TabPanel />

            </Grid>
            </div>
        </Grid>
    )
}


export function WorkstationPolicy() {
    const  classes = useStyles()

    const [tabValue, setValue] = React.useState(5);

    const handleTabChange = (newValue: any) => {
      setValue(newValue);
    }

    return (         
        <Grid container spacing={2}  justify="center">
        <div className={classes.tabRoot}>
            <Grid item xs={3}>
            <Tabs
            indicatorColor="primary"
                orientation="vertical"
                variant="scrollable"
                value={tabValue}
                onChange={handleTabChange}
                className={classes.tabs}
            >
                
                <Tab label={<Typography variant="h4">Baseline</Typography>} {...a11yProps(0)} />
                <Tab label={<Typography variant="h4">Windows</Typography>}  {...a11yProps(1)} />
                <Tab label={<Typography variant="h4">Mac</Typography>}  {...a11yProps(2)} /> 
                <Tab label={<Typography variant="h4">Linux</Typography>}  {...a11yProps(2)} /> 
            </Tabs>

                </Grid>
            
            <Grid item xs={12} sm={12} md={9}>

            <TabPanel />

            </Grid>
            </div>
        </Grid>
    )
}


export function BrowserPolicy() {
    const  classes = useStyles()

    const [tabValue, setValue] = React.useState(5);

    const handleTabChange = (newValue: any) => {
      setValue(newValue);
    }

    return (         
        <Grid container spacing={2}  justify="center">
        <div className={classes.tabRoot}>
            <Grid item xs={3}>
            <Tabs
            indicatorColor="primary"
                orientation="vertical"
                variant="scrollable"
                value={tabValue}
                onChange={handleTabChange}
                className={classes.tabs}
            >
                <Tab label={<Typography variant="h4">Chrome</Typography>}  {...a11yProps(1)} />
                <Tab label={<Typography variant="h4">Firefox</Typography>}  {...a11yProps(2)} /> 
    
            </Tabs>

                </Grid>
            
            <Grid item xs={12} sm={12} md={9}>

            <TabPanel />

            </Grid>
            </div>
        </Grid>
    )
}