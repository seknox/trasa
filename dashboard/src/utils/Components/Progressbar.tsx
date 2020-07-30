import LinearProgress from '@material-ui/core/LinearProgress';
import { createStyles, makeStyles } from '@material-ui/core/styles';
import React from 'react';

const useStyles = makeStyles(() =>
  createStyles({
    root: {
      width: '100%',
      marginTop: 5,
    },
  }),
);

export default function Progress() {
  const classes = useStyles();
  return (
    <div className={classes.root}>
      <LinearProgress />
    </div>
  );
}
