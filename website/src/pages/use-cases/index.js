import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import clsx from 'clsx';
import React from 'react';
import { Typography, Grid, Paper } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import styles from '../styles.module.css';
import Link from '@docusaurus/Link';

import ThemeBase from '../../components/muiTheme';

const useStyles = makeStyles(() => ({
  ctaPad: {
    marginTop: 50,
  },
  link: {
    fontSize: 30,
  },
}));

export default function UseCases() {
  const context = useDocusaurusContext();
  const classes = useStyles();
  const { siteConfig = {} } = context;
  return (
    <Layout title="TRASA Use Cases" description="TRASA use cases. Know why do you need TRASA">
      <ThemeBase>
        <Grid
          container
          spacing={2}
          direction="column"
          alignItems="center"
          justify="center"
          className={classes.background}
        >
          <Grid item xs={12} sm={12}>
            <div className={classes.ctaPad}>
              <Typography variant="h1">Why do you need TRASA ?</Typography>
              <Typography variant="body1" style={{ textAlign: 'center' }}>
                TRASA is an open source project with a growing community. <br /> There are active,
                dedicated users willing to help you through various mediums.
              </Typography>
            </div>
          </Grid>
          <Grid item xs={12} sm={12} md={4}>
            <Paper>
              <Typography variant="h3">Two Factor Authentication</Typography>
            </Paper>
          </Grid>
        </Grid>
      </ThemeBase>
    </Layout>
  );
}
