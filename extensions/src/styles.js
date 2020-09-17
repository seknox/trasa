import { makeStyles } from '@material-ui/core/styles';

const useStyles = makeStyles((theme) => ({
  root: {
    //display: 'flex',
  },
  card: {
    textAlign: 'center',
    backgroundColor: 'white', //'#d0d3d4', // //rgba(1,1,35,1)
  },
  icon: {
    marginLeft: '40%',
    color: 'navy',
  },
  paper: {
    //raised: true,
    minWidth: 350,
    minHeight: 350,
    backgroundColor: 'white', //'#d0d3d4', // //rgba(1,1,35,1)
    padding: theme.spacing(2),
    fontFamily: 'Open Sans, Rajdhani',
    textAlign: 'center',
  },

  title: {
    marginBottom: 16,
    fontSize: 20,
    color: '#1A237E',
    fontFamily: 'Open Sans, Rajdhani',
  },
  input: {
    color: 'teal',
  },
  button: {
    backgroundColor: '#1A237E',
    color: 'white',
  },
  errorText: {
    color: 'white',
    fontSize: '15px',
    fontFamily: 'Open Sans, Rajdhani',
    background: 'maroon',
  },
  list: {
    fontSize: '14px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  listItem: {
    fontSize: '14px',
    fontFamily: 'Open Sans, Rajdhani',
  },
}));

export default useStyles;
