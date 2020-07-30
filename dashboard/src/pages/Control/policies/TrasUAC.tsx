import Box from '@material-ui/core/Box';
import Button from '@material-ui/core/Button';
import Checkbox from '@material-ui/core/Checkbox';
import FormControl from '@material-ui/core/FormControl';
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import Input from '@material-ui/core/Input';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemAvatar from '@material-ui/core/ListItemAvatar';
import ListItemSecondaryAction from '@material-ui/core/ListItemSecondaryAction';
import ListItemText from '@material-ui/core/ListItemText';
import MenuItem from '@material-ui/core/MenuItem';
import Select from '@material-ui/core/Select';
import { makeStyles, withStyles } from '@material-ui/core/styles';
import Switch from '@material-ui/core/Switch';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import TextField from '@material-ui/core/TextField';
import Typography from '@material-ui/core/Typography';
import AddIcon from '@material-ui/icons/Add';
import DeleteIcon from '@material-ui/icons/Delete';
import React, { useState } from 'react';
import Tooltip from '@material-ui/core/Tooltip';

const useStyles = makeStyles((theme) => ({
  expiry: {
    marginLeft: 10,
    width: 100
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
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)'
    }
  },
  timePickerTextField: {
    marginLeft: theme.spacing(1),
    marginRight: theme.spacing(1),
    width: 100
  },

  timePickerContainer: {
    display: 'flex',
    flexWrap: 'wrap'
  },
  button: {
    margin: theme.spacing(1),
    marginTop: '5%'
  },
  rightIcon: {
    marginRight: theme.spacing(1)
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
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)'
    }
  },
  textFieldInputBig: {
    borderRadius: 4,
    backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    fontSize: 16,
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)'
    }
  },
  textFieldRoot: {
    padding: 1,
    'label + &': {
      marginTop: theme.spacing(3)
    }
  },
  textFieldFormLabel: {
    fontSize: 18
  },
  policyName: {
    marginLeft: '20%'
  },
  toolTip: {
    padding: 5,
    // maxWidth: 220,
    fontSize: 16,
    // fontFamily: 'Open Sans, Rajdhani',
    border: '1px solid #dadde9'
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
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)'
    }
  },
  tabRoot: {
    flexGrow: 1,
    backgroundColor: theme.palette.background.paper,
    display: 'flex'
    //height: 224,
  },
  tabs: {
    borderRight: `1px solid ${theme.palette.divider}`
  },
}));

const ITEM_HEIGHT = 48;
const ITEM_PADDING_TOP = 8;
const MenuProps = {
  PaperProps: {
    style: {
      maxHeight: ITEM_HEIGHT * 4.5 + ITEM_PADDING_TOP
    }
  }
};

const CTooltip = withStyles((theme) => ({
  tooltip: {
    // backgroundColor: '#f5f5f9',
    color: 'rgba(0, 0, 0, 0.87)',
    maxWidth: 420,
    fontSize: 16,
    border: '1px solid #dadde9'
  }
}))(Tooltip);

export default function TrasaUAC(props: any) {
  const [tabValue, setValue] = React.useState(0);

  const handleTabChange = (event: React.ChangeEvent<{}>, newValue: number) => {
    setValue(newValue);
  };

  const [policyType, setPolicyType] = React.useState('dayAndTime');
  const [days, setDays] = useState<string[]>([]);
  const [fromTime, setFromTime] = useState('');
  const [toTime, setToTime] = useState('');

  function handlerPolicyTypechange(e: React.ChangeEvent<HTMLSelectElement>) {
    setPolicyType(e.target.value);
  }

  // addDayAndTime calls props to update main dayAndTime policy array.
  function addDayAndTime() {
    props.updateDayAndTime(days, fromTime, toTime, 0);
  }

  const updateDays = (event: React.ChangeEvent<{ value: unknown }>) => {
    setDays(event.target.value  as string[]);
  };

  const updateFrom = (event :any) => {
    setFromTime(event.target.value);
  };

  const updateTo = (event: any) => {
    setToTime(event.target.value);
  };

  function renderPolicyComponent() {
    switch (tabValue) {
      case 0:
        return (
          <TabPanel value={tabValue} index={0}>
            <Grid item xs={12} sm={12} md={12}>
              <Grid container spacing={2}>
                <Grid item xs={4}>
                  <Typography variant="h4">Mandatory 2FA: </Typography>
                </Grid>
                <Grid item xs={8}>
                  {/* <FormControlLabel
                                     control={ */}
                  <Switch
                    // tfaRequired, settfaRequired
                    checked={props.tfaRequired}
                    onChange={props.handle2FAChange}
                    //value={props.tfaRequired}
                    color="secondary"
                  />
                </Grid>
                <br />
                {/* <Grid item xs={4} >
                             <Typography variant="h4">Risk Based 2FA: </Typography>
                                 </Grid>
                                 <Grid item xs={8} >
                                     <Switch
                                     // tfaRequired, settfaRequired
                                       checked= {props.tfaRequired}
                                       onChange={props.handle2FAChange}
                                       //value={props.tfaRequired}
                                       color="secondary"
                                       />           
                                 </Grid> */}
              </Grid>
            </Grid>
          </TabPanel>
        );
      case 1:
        return (
          <TabPanel value={tabValue} index={1}>
            <Grid item xs={12} sm={12} md={12}>
              <Grid container spacing={2}>
                <Grid item xs={12} sm={6} md={4} lg={4}>
                  <Typography variant="h4">Record Session: </Typography>
                </Grid>
                <Grid item xs={12} sm={6} md={2} lg={2}>
                  <Switch
                    checked={props.recordSession}
                    onChange={props.handleSessionRecordChange}
                    color="secondary"
                  />
                </Grid>
              </Grid>
            </Grid>
          </TabPanel>
        );
      case 2:
        return (
          <TabPanel value={tabValue} index={2}>
            <Grid item xs={12} sm={12} md={12}>
              <Grid container spacing={2}>
                <Grid item xs={12} sm={12} md={4}>
                  <Typography variant="h4">Allow File Transfers: </Typography>
                </Grid>
                <Grid item xs={12} sm={6} md={2} lg={2}>
                  <Switch
                    checked={props.fileTransfer}
                    onChange={props.handleFileTransferChange}
                    color="secondary"
                  />
                </Grid>
              </Grid>
            </Grid>
          </TabPanel>
        );
      case 3:
        return (
          <TabPanel value={tabValue} index={3}>
            <Grid item xs={12} sm={12} md={12}>
              <Grid container spacing={2}>
                <Grid item xs={12} sm={6} md={4} lg={3}>
                  <Typography variant="h4">IP Source: </Typography>
                </Grid>
                <Grid item xs={12} sm={6} md={8} lg={9}>
                  <TextField
                    fullWidth
                    value={props.ipSource}
                    onChange={props.handleIPSourceChange}
                    name="ipSource"
                    InputProps={{
                      disableUnderline: true,
                      classes: {
                        root: classes.textFieldRoot,
                        input: classes.textFieldInputBig
                      }
                    }}
                    InputLabelProps={{
                      shrink: true,
                      className: classes.textFieldFormLabel
                    }}
                  />
                </Grid>
              </Grid>
            </Grid>
          </TabPanel>
        );
      case 4:
        return (
          <TabPanel value={tabValue} index={4}>
            <Grid container spacing={2}>
              <Grid item xs={12} sm={6} md={1}>
                <Typography variant="h4">Days: </Typography>
              </Grid>
              <Grid item xs={12} sm={12} md={5}>
                <FormControl className={classes.formControl}>
                  <Select
                    fullWidth
                    multiple
                    id="select-multiple-checkbox"
                    value={days}
                    onChange={updateDays}
                    input={<Input />}
                    renderValue={(selected) => (selected as string[]).join(', ')}
                    MenuProps={MenuProps}
                  >
                    {weekDays.map((name) => (
                      <MenuItem key={name} value={name}>
                        <Checkbox checked={days.indexOf(name) > -1} />
                        <ListItemText secondary={name} />
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </Grid>
              <Grid item xs={12} sm={12} md={2}>
                <Typography variant="h3">Time: </Typography>
              </Grid>
              <Grid item xs={6} sm={6} md={2}>
                <TimePicker
                  Time="FROM"
                  value={props.FromTime}
                  name="fromTime"
                  updateTime={updateFrom}
                />
              </Grid>
              <Grid item xs={6} sm={6} md={2}>
                <TimePicker Time="TO" value={props.ToTime} name="toTime" updateTime={updateTo} />
              </Grid>

              <Grid item xs={12}>
                <List>
                  <ListItemText>{props.dayAndTime.length} permissions added</ListItemText>
                  {props.dayAndTime.map((value: any, index: number) => (
                    <ListItem key={index}>
                      {/* <ListItemAvatar /> */}
                      <CTooltip
                        placement="left-start"
                        className={classes.toolTip}
                        title={JSON.stringify(value)}
                      >
                        <Typography variant="h4">{'Permission ' + (index + 1)}</Typography>
                      </CTooltip>
                      <ListItemSecondaryAction>
                        <IconButton
                          aria-label="Delete"
                          style={{ color: 'maroon' }}
                          onClick={() => {
                            props.deleteDayAndTime(index);
                          }}
                        >
                          <DeleteIcon />
                        </IconButton>
                      </ListItemSecondaryAction>
                    </ListItem>
                  ))}
                </List>
              </Grid>

              <Grid item xs={12} sm={12} md={12}>
                <Grid container spacing={2}>
                  <Grid item xs={3}>
                    <Button
                      className={classes.button}
                      variant="contained"
                      color="secondary"
                      size="small"
                      onClick={addDayAndTime}
                    >
                      <AddIcon className={classes.rightIcon} />
                      {props.dayAndTime.length > 0 ? 'Add Another' : 'Add'}
                    </Button>
                  </Grid>
                  {/* <Grid item xs={3}>
                                        <Button className={classes.button} variant="contained" color="secondary" onClick={props.addPermission}>
                                        <SaveIcon className={classes.rightIcon} />
                                            Save
                                        </Button>
                                    </Grid> */}
                </Grid>
              </Grid>
            </Grid>
          </TabPanel>
        );

      case 5:
        return (
          <TabPanel value={tabValue} index={5}>
            <Grid item xs={12} sm={12} md={12}>
              <Grid container spacing={2}>
                <Grid item xs={12} sm={6} md={4} lg={3}>
                  <Typography variant="h3">Expires On: </Typography>
                </Grid>
                <Grid item xs={12} sm={6} md={4} lg={4}>
                  <TextField
                    fullWidth
                    id="Date"
                    name="expiry"
                    value={props.expiry}
                    onChange={props.handleExpiryChange}
                    label="Auto Expiry"
                    type="date"
                    //defaultValue={props.Expiry}
                    // className={classes.textField}
                    InputLabelProps={{
                      shrink: true
                    }}
                  />
                </Grid>
              </Grid>
            </Grid>
          </TabPanel>
        );
    }
  }
  const classes = useStyles();

  return (
    <Grid container spacing={2} justify="center">
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
            <Tab
              label={<Typography variant="h4">Second Factor Auth</Typography>}
              {...a11yProps(0)}
            />
            <Tab
              label={<Typography variant="h4">Session Recording</Typography>}
              {...a11yProps(1)}
            />
            <Tab label={<Typography variant="h4">File Transfers</Typography>} {...a11yProps(2)} />
            <Tab label={<Typography variant="h4">IP Source</Typography>} {...a11yProps(3)} />
            <Tab label={<Typography variant="h4">Day and Time</Typography>} {...a11yProps(4)} />
            <Tab label={<Typography variant="h4">Expiry</Typography>} {...a11yProps(5)} />
          </Tabs>
        </Grid>

        <Grid item xs={12} sm={12} md={9}>
          {renderPolicyComponent()}
        </Grid>
      </div>
    </Grid>
  );
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
    'aria-controls': `vertical-tabpanel-${index}`
  };
}

const weekDays = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];

function TimePicker(props: any) {
  const classes = useStyles();

  return (
    <div>
      <FormControl className={classes.timePickerContainer}>
        <TextField
          id={props.Time}
          onChange={props.updateTime}
          required
          label={props.Time}
          value={props.value}
          type="time"
          //defaultValue="09:00"
          className={classes.timePickerTextField}
          //className={classes.formControl}
          InputLabelProps={{
            shrink: true
          }}
          inputProps={{
            step: 300 // 5 min
          }}
        />
      </FormControl>
    </div>
  );
}
