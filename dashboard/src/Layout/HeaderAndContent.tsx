import AppBar from '@material-ui/core/AppBar';
import Breadcrumbs from '@material-ui/core/Breadcrumbs';
import Chip from '@material-ui/core/Chip';
import { makeStyles, Theme, withStyles } from '@material-ui/core/styles';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import HomeIcon from '@material-ui/icons/Home';
import React from 'react';
import { RouteComponentProps, withRouter } from 'react-router-dom';

const StyledChip = withStyles((theme: Theme) => ({
  root: {
    backgroundColor: 'rgba(1,1,35,1)',
    height: theme.spacing(3),
    // padding: theme.spacing(1),
    color: 'silver',
    '&:hover, &:focus': {
      backgroundColor: '#030417',
    },
    '&:active': {
      boxShadow: theme.shadows[1],
      backgroundColor: '#000080',
      fontWeight: 500,
    },
    fontSize: '16px',
    fontWeight: 550,
    fontFamily: 'Open Sans, Rajdhani',
  },
  label: {
    //  marginTop: 5,
  },
  // icon: {
  //   top: 0,
  // }
}))(Chip) as typeof Chip;

const StyledTabs = withStyles((theme: Theme) => ({
  root: {
    marginLeft: theme.spacing(1),
    maxHeight: 10,
  },
  indicator: {
    height: 3,
    borderTopLeftRadius: 3,
    borderTopRightRadius: 3,
    backgroundColor: theme.palette.common.white,
  },
}))(Tabs) as typeof Tabs;

const lightColor = 'rgba(255, 255, 255, 0.7)'; // 'rgba(255, 255, 255, 0.7)'; // '#030417';

const useStyles = makeStyles((theme) => ({
  mainContent: {
    flex: 1,
    background: 'white',
  },
  secondaryBar: { backgroundColor: theme.palette.primary.dark },
  toolbar: {},
  button: {
    borderColor: lightColor,
    color: 'white',
  },
  tabHeader: {
    color: 'white',
    fontSize: '12px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  breadcrumbs: {
    color: 'white',
    flexGrow: 1,
  },
}));

type PageProps = {
  route: string;
  name: string;
};
type SubHeaderAndContentProps = {
  children?: React.ReactNode;
  // history: History;
  tabHeaders: string[];
  pageName: PageProps[];
  Components: React.ReactElement[];
};

type TabContentProps = {
  children?: React.ReactNode;
};

function TabContainer(props: TabContentProps) {
  return <div style={{ padding: 8 * 3 }}>{props.children}</div>;
}

function SubHeaderAndContent(props: SubHeaderAndContentProps & RouteComponentProps) {
  const classes = useStyles();
  const [tabValue, setTabValue] = React.useState(0);

  React.useEffect(() => {
    const hash = window.location.hash.substr(1);
    const tabVal = decodeURIComponent(hash);
    if (tabVal !== '') {
      setTabValue(props.tabHeaders.indexOf(tabVal));
    }
  }, [props.tabHeaders]);

  function handleTabChange(event: React.ChangeEvent<{}>, value: number) {
    setTabValue(value);
    props.history.push(`#${props.tabHeaders[value]}`);
  }

  // const callMixPanel = (val: string) => {
  //   mixpanel.track(val);
  // };

  return (
    <div>
      <AppBar
        component="div"
        className={classes.secondaryBar}
        color="primary"
        position="sticky"
        elevation={1}
      >
        <Toolbar className={classes.toolbar}>
          <Breadcrumbs aria-label="breadcrumb" separator="/" className={classes.breadcrumbs}>
            {props.pageName.map((v: any, i: any) => (
              <StyledChip
                clickable
                component="a"
                href={v.route}
                label={v.name}
                key={v}
                icon={i === 0 ? <HomeIcon style={{ color: 'dimgrey' }} /> : undefined}
              />
            ))}
          </Breadcrumbs>

          {props.tabHeaders.length > 1 ? (
            <StyledTabs value={tabValue} onChange={handleTabChange} textColor="inherit">
              {props.tabHeaders.map((v: any) => (
                <Tab
                  textColor="inherit"
                  key={v}
                  label={
                    <Typography variant="h5" style={{ color: 'white' }}>
                      {v}
                    </Typography>
                  }
                />
              ))}
            </StyledTabs>
          ) : null}
        </Toolbar>
      </AppBar>

      <div className={classes.mainContent}>
        {props.Components.map(
          (v: any, i: any) =>
            tabValue === i && <TabContainer key={v}> {props.Components[i]} </TabContainer>,
        )}
      </div>
    </div>
  );
}

export default withRouter(SubHeaderAndContent);
// export default SubHeaderAndContent;
