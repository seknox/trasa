import LinearProgress from '@material-ui/core/LinearProgress';
import { withStyles } from '@material-ui/core/styles';
import PropTypes from 'prop-types';
import React from 'react';

const styles = (theme) => ({
  root: {
    flexGrow: 1,
    marginLeft: '4%',
    flexDirection: 'column',
  },
});

function Progress(props) {
  const { classes } = props;
  return (
    <div className={classes.root}>
      <LinearProgress />
    </div>
  );
}

Progress.propTypes = {
  classes: PropTypes.object.isRequired,
};

const ProgressHOC = withStyles(styles)(Progress);

export default ProgressHOC;
