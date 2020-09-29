import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import clsx from 'clsx';
import React from 'react';

import styles from '../styles.module.css';
import { Typography, Grid, Paper } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import ThemeBase from '../../components/muiTheme';

import useBaseUrl from '@docusaurus/useBaseUrl';
import FeaturesComp from '../../components/features';

const useStyles = makeStyles(() => ({
  ctaPad: {
    marginTop: 50,
  },
  link: {
    fontSize: 30,
  },
  content: {
    backgroundColor: '#ffffff',
  },
}));

export default function ZeroTrustServiceAccess() {
  const context = useDocusaurusContext();
  const { siteConfig = {} } = context;
  const classes = useStyles();
  const zeroTrust = useBaseUrl('arch/audit.svg');
  return (
    <Layout title={`${siteConfig.title}`} description="Privileged Access Management">
      <header className={clsx('hero hero--primary', styles.heroBanner)}>
        <div className="container">
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
                  <Typography variant="h1" component="h1">
                    Privileged Access Management
                  </Typography>
                  <Typography variant="body1" style={{ textAlign: 'center' }}>
                    Monitor, verify and trust user devices.
                  </Typography>
                </div>
              </Grid>
            </Grid>
          </ThemeBase>
        </div>
      </header>

      <main>
        <section>
          <div className="container">
            <ThemeBase>
              <Grid
                container
                spacing={2}
                direction="column"
                alignItems="center"
                justify="center"
                className={classes.content}
              >
                <Grid item xs={12} sm={12}>
                  <FeaturesComp />
                </Grid>
              </Grid>
            </ThemeBase>
          </div>
        </section>
      </main>
    </Layout>
  );
}
