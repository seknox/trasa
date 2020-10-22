import { Typography } from '@material-ui/core';
import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import React from 'react';
import useBaseUrl from '@docusaurus/useBaseUrl';
import ThemeBase from '../muiTheme';

const useStyles = makeStyles(() => ({
  ctaPad: {
    marginTop: 150,
    textAlign: 'center',
  },
  contained: {
    color: '#000080',
    backgroundColor: 'white',
    borderColorL: '#000080',
    fontWeight: 600,
    //  fontSize: '14px',
    boxShadow: 'none',
  },
  image: {
    boxShadow: '0 0 20px 0 rgba(0,0,0,0.12)',
    backgroundColor: 'white',
  },
}));

export default function Enterprise() {
  const classes = useStyles();
  const imgUrl = useBaseUrl('arch/usecase-reference.png');
  return (
    <ThemeBase>
      <Grid container spacing={2} direction="column" justify="center" alignItems="center">
        <Grid item xs={8}>
          <div className={classes.ctaPad}>
            <Typography variant="h1">Why you need TRASA ?</Typography>
            <Typography variant="body1" component="span" style={{ textAlign: 'center' }}>
              Data center or dynamic cloud infrastructure, dedicated servers or ephemeral
              applications and <br /> services, access by internal team or managed service provider;
            </Typography>
            <Typography variant="subtitle1" style={{ textAlign: 'center' }}>
              <b>TRASA</b> is a free and open source project that provides modern security features
              and enables best practice security to protect internal infrastructure ( Web, SSH, RDP,
              and Database services) from unauthorized or malicious access.
            </Typography>
          </div>
        </Grid>
        <Grid item xs={12}>
          <div className={classes.image}>
            <img src={imgUrl} alt="protect internal infrastructure" />
          </div>
        </Grid>
        {/* <Grid item xs={12}>
          <Link className={clsx('button  button--lg', classes.contained)} to={useBaseUrl('docs/')}>
            Learn more about features
          </Link>
        </Grid> */}
      </Grid>
    </ThemeBase>
  );
}
