import Box from '@material-ui/core/Box';
import { withStyles } from '@material-ui/core/styles';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import Typography from '@material-ui/core/Typography';
import React from 'react';
import { BrowserPolicy, MobileDevicePolicy, WorkstationPolicy } from './DeviceHygiene';


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

interface styledTabsProps {
  value: number;
  centered: boolean;
  onChange: (event: React.ChangeEvent<{}>, newValue: number) => void;
}


interface styledTabProps {
  label: string;
}

const StyledTabs = withStyles({
  root: {
    maxHeight: 10,
    borderBottom: '1px solid #e8e8e8'
  }
})((props: styledTabsProps) => <Tabs {...props} TabIndicatorProps={{ children: <div /> }} />);

const StyledTab = withStyles((theme) => ({
  root: {
    textTransform: 'none',
    fontSize: '16px',
    fontFamily: 'Open Sans, Rajdhani',
    color: 'black',
    marginRight: theme.spacing(1),
    '&:focus': {
      opacity: 1
    }
  }
}))((props: styledTabProps) => <Tab disableRipple {...props} />);

// PolicyTab is main tab which wraps basic TRASA policy configs and device hygiene policy config
export default function PolicyTab() {
  const [tabValue, setTabValue] = React.useState(1);

  const handleChange = (event: React.ChangeEvent<{}>, newValue: number) => {
    setTabValue(newValue);
  };
  return (
    <div>
      <StyledTabs
        value={tabValue}
        centered
        onChange={handleChange}
        aria-label="styled tabs example"
      >
        <StyledTab label="Mobile Device" />
        <StyledTab label="Workstation" />
        <StyledTab label="Web Browser" />
      </StyledTabs>
      <TabPanel value={tabValue} index={0}>
        <MobileDevicePolicy />
      </TabPanel>
      <TabPanel value={tabValue} index={1}>
        <WorkstationPolicy />
      </TabPanel>
      <TabPanel value={tabValue} index={2}>
        <BrowserPolicy />
      </TabPanel>
    </div>
  );
}
