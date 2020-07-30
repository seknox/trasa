import Box from '@material-ui/core/Box';
import FormControl from '@material-ui/core/FormControl';
import Grid from '@material-ui/core/Grid';
import { makeStyles, withStyles } from '@material-ui/core/styles';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import TextField from '@material-ui/core/TextField';
import Tooltip from '@material-ui/core/Tooltip';
import Typography from '@material-ui/core/Typography';
import React, { useState } from 'react';
import EndpointSecurity from './EndpointSecurity';
import LoginSecurity from './LoginSecurity';
import NetworkSecurity from './NetworkSecurity';
import OsSecurity from './OsSecurity';
import TrustedDevices from "./TrustedDevices";

const useStyles =  makeStyles(theme => ({
    expiry: {
        marginLeft: 10,
        width:100,
    },
    formControl: {
        fontSize: 15,
        borderRadius: 4,
        backgroundColor: theme.palette.common.white,
        border: '1px solid #ced4da',
        height: 31,
        marginTop: 5,
        //padding: '10px 100px',
        width: 'calc(100% - 4px)',
        transition: theme.transitions.create(['border-color', 'box-shadow']),
        '&:focus': {
          borderColor: '#80bdff',
          boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
        },
    },
    button: {
        margin: theme.spacing(1),
        marginTop: '5%',
    },
    rightIcon: {
        marginRight: theme.spacing(1),
    },
      textFieldInput: {
        borderRadius: 4,
        backgroundColor: theme.palette.common.white,
        border: '1px solid #ced4da',
        fontSize: 16,
        padding: '10px 12px',
        width: 'calc(100% - 4px)',
        transition: theme.transitions.create(['border-color', 'box-shadow']),
        '&:focus': {
          borderColor: '#80bdff',
          boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
        },
      },
      textFieldInputBig: {
        borderRadius: 4,
        backgroundColor: theme.palette.common.white,
        border: '1px solid #ced4da',
        fontSize: 16,
        transition: theme.transitions.create(['border-color', 'box-shadow']),
        '&:focus': {
          borderColor: '#80bdff',
          boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
        },
      },
      textFieldRoot: {
        padding: 1,
        'label + &': {
          marginTop: theme.spacing(3),
        },
      },
      textFieldFormLabel: {
        fontSize: 18,
      },
      policyName: {
          marginLeft: '20%'
      },
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
                //height: 224,
          },
          tabs: {
            borderRight: `1px solid ${theme.palette.divider}`,
          },
    timePickerTextField : {
        marginLeft: theme.spacing(1),
        marginRight: theme.spacing(1),
        width: 100,
    },

    timePickerContainer : {
        display: 'flex',
        flexWrap: 'wrap',
    },
}))


const ITEM_HEIGHT = 48;
const ITEM_PADDING_TOP = 8;
const MenuProps = {
  PaperProps: {
    style: {
      maxHeight: ITEM_HEIGHT * 4.5 + ITEM_PADDING_TOP,
    },

  },
};

const CTooltip = withStyles(theme => ({
    tooltip: {
     // backgroundColor: '#f5f5f9',
      color: 'rgba(0, 0, 0, 0.87)',
      maxWidth: 420,
      fontSize: 16,
      border: '1px solid #dadde9',
    },
  }))(Tooltip);

export default function TrasaUAC(props: any){

    const [tabValue, setValue] = React.useState(0);

    const handleTabChange = (event: React.ChangeEvent<{}>, newValue: number) => {
      setValue(newValue);
    }



    function renderPolicyComponent() {
        switch (tabValue) {
            case 0:
                return (<TabPanel value={tabValue} index={0}>
                    <TrustedDevices
                        blockUntrustedDevices={props.blockUntrustedDevices}
                        changeBlockUntrustedDevices={props.changeBlockUntrustedDevices}

                    />

                </TabPanel> )

                case 1:
                    return (<TabPanel value={tabValue} index={1}>
                                <LoginSecurity
                                    changeBlockAutologinEnabled={props.changeBlockAutologinEnabled}
                                    blockAutologinEnabled={props.blockAutologinEnabled}
                                    changeBlockTfaNotConfigured={props.changeBlockTfaNotConfigured}
                                    blockTfaNotConfigured={props.blockTfaNotConfigured}
                                    changeBlockIdleScreenLockDisabled={props.changeBlockIdleScreenLockDisabled}
                                    blockIdleScreenLockDisabled={props.blockIdleScreenLockDisabled}
                                    changeBlockRemoteLoginEnabled={props.changeBlockRemoteLoginEnabled}
                                    blockRemoteLoginEnabled={props.blockRemoteLoginEnabled}

                                />

                            </TabPanel> )
                case 2:
                    return ( <TabPanel value={tabValue} index={2}>
                                <OsSecurity
                                    blockCriticalAutoUpdateDisabled={props.blockCriticalAutoUpdateDisabled}
                                    changeBlockCriticalAutoUpdateDisabled={props.changeBlockCriticalAutoUpdateDisabled}

                                    changeBlockJailBroken={props.changeBlockJailBroken}
                                    blockJailBroken={props.blockJailBroken}
                                    changeBlockDebuggingEnabled={props.changeBlockDebuggingEnabled}
                                    blockDebuggingEnabled={props.blockDebuggingEnabled}
                                    changeBlockEmulated={props.changeBlockEmulated}
                                    blockEmulated={props.blockEmulated}

                                />

                            </TabPanel> )
                case 3:
                    return ( <TabPanel value={tabValue} index={3}>
                                <EndpointSecurity
                                    blockAntivirusDisabled={props.blockAntivirusDisabled}
                                    changeBlockAntivirusDisabled={props.changeBlockAntivirusDisabled}
                                    blockFirewallDisabled={props.blockFirewallDisabled}
                                    changeBlockFirewallDisabled={props.changeBlockFirewallDisabled}

                                    changeBlockEncryptionNotSet={props.changeBlockEncryptionNotSet}
                                    blockEncryptionNotSet={props.blockEncryptionNotSet}

                                />

                                </TabPanel> )
                case 4:
                    return (<TabPanel value={tabValue} index={4}>

                            </TabPanel>)

        }
    }
        const  classes  = useStyles();

        return (


                <Grid container spacing={2}  justify="center">
                <div className={classes.tabRoot}>
                    <Grid item xs={3}>
                    <Tabs
                    indicatorColor="primary"
                        orientation="vertical"
                       // variant="scrollable"
                        value={tabValue}
                        onChange={handleTabChange}
                        className={classes.tabs}
                    >

                        <Tab label={<Typography variant="h4">Trusted Devices</Typography>} {...a11yProps(0)} />
                        <Tab label={<Typography variant="h4">Login Security</Typography>} {...a11yProps(1)} />
                        <Tab label={<Typography variant="h4">OS Security</Typography>}  {...a11yProps(2)} />
                        <Tab label={<Typography variant="h4">Endpoint Security</Typography>}  {...a11yProps(3)} />
                        {/* <Tab label={<Typography variant="h4">Network Security</Typography>}  {...a11yProps(2)} /> */}
                        {/* <Tab label={<Typography variant="h4">Version Blocking</Typography>}  {...a11yProps(2)} />  */}
                    </Tabs>

                        </Grid>

                    <Grid item xs={12} sm={12} md={9}>

                    {renderPolicyComponent()}

                    </Grid>
                    </div>
                </Grid>

        )
}



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



