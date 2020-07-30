import { makeStyles } from '@material-ui/core/styles';

const useStyles = makeStyles((theme) => ({
  root: {
    flexgrow: 1,
  },

  paperTrans: {
    backgroundColor: 'transparent',
    // padding: theme.spacing(1),
    textAlign: 'center',
  },
  paperTrans1: {
    backgroundColor: 'transparent',
  },
  paper: {
    //  backgroundColor:  '#fdfdfd',
    textAlign: 'center',
    padding: theme.spacing(2),
  },
  paperLarge: {
    backgroundColor: '#fdfdfd',
    // backgroundColor: '#00001aff', //'#E0E0E0', 'rgba(10,34,52,1)'
    minWidth: 400,
    minHeight: 500,
    // padding: theme.spacing(2),
    // textAlign: 'center',
    color: theme.palette.text.secondary,
  },
  paperHeighted: {
    backgroundColor: '#fdfdfd',
    // backgroundColor: '#00001aff', //'#E0E0E0', 'rgba(10,34,52,1)'
    minWidth: 800,
    minHeight: 300,
    marginTop: '5%',
    marginBotton: '5%',
    padding: theme.spacing(2),
    textAlign: 'center',
    color: theme.palette.text.secondary,
  },
  aggHeadersBig: {
    color: '#1b1b32',
    fontSize: '50px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  aggHeaders: {
    color: '#1b1b32',
    fontSize: '18px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  selectLabel: {
    fontSize: 12,
    //  marginBotton: 5,
    padding: '1px',
    color: 'grey',
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
  textPadding: {
    fontSize: 15,
    fontFamily: 'Open Sans, Rajdhani',
    backgroundColor: theme.palette.common.white,
    // height: 31,
    paddingTop: 14,
    paddingLeft: 15,
    // padding: 15,
  },
  selectAnalytics: {
    fontSize: 12,
    fontFamily: 'Open Sans, Rajdhani',
    borderRadius: 4,
    backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    height: 25,
    marginTop: 5,
    // padding: '10px 100px',
    width: 'calc(100% - 4px)',
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
    },
  },
  analyticsText: {
    paddingTop: 10,
    paddingLeft: 15,
    fontSize: 14,
    fontFamily: 'Open Sans, Rajdhani',
  },
  echartStyles: {
    // minHeight: getHeight(),
    height: '180px', // getHeight()
  },
}));

export default useStyles;
