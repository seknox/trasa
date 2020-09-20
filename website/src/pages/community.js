import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import clsx from 'clsx';
import React from 'react';
import { Typography, Grid } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import styles from './styles.module.css';

import ThemeBase from '../components/muiTheme';

const useStyles = makeStyles(() => ({
  ctaPad: {
    marginTop: 50,
  },
}));

function Home() {
  const context = useDocusaurusContext();
  const classes = useStyles();
  const { siteConfig = {} } = context;
  return (
    <Layout
      title={`${siteConfig.title}`}
      description="Description will go into a meta tag in <head />"
    >
      <header className={clsx('hero hero--primary', styles.heroBanner)}>
        <div className="container">
          <ThemeBase>
            <Grid container spacing={2} className={classes.background}>
              <Grid item xs={12} sm={12}>
                <div className={classes.ctaPad}>
                  <Typography variant="h1">Community</Typography>
                  <Typography variant="body1">
                    TRASA is an open source project with a growing community. <br /> There are
                    active, dedicated users willing to help you through various mediums.
                  </Typography>
                </div>
              </Grid>
            </Grid>
          </ThemeBase>
        </div>
      </header>

      <main />
    </Layout>
  );
}

export default Home;
