import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import clsx from 'clsx';
import React from 'react';
import { Typography, Grid } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import styles from './styles.module.css';
import Link from '@docusaurus/Link';

import ThemeBase from '../components/muiTheme';

const useStyles = makeStyles(() => ({
  ctaPad: {
    marginTop: 50,
  },
  link: {
    fontSize: 30,
  },
}));

function Home() {
  const context = useDocusaurusContext();
  const classes = useStyles();
  const { siteConfig = {} } = context;
  return (
    <Layout
      title={`${siteConfig.title}`}
      description="TRASA community discussion, contribution and support"
    >
      <header className={clsx('hero hero--primary', styles.heroBanner)}>
        <div className="container">
          <ThemeBase>
            <Grid container spacing={2} className={classes.background}>
              <Grid item xs={12} sm={12}>
                <div className={classes.ctaPad}>
                  <Typography variant="h1">Community</Typography>
                  <Typography variant="body1" style={{ textAlign: 'center' }}>
                    TRASA is an open source project with a growing community. <br /> There are
                    active, dedicated users willing to help you through various mediums.
                  </Typography>
                </div>
              </Grid>
              <Grid item xs={12} sm={12}>
                <div className={classes.ctaPad}>
                  <Typography variant="h2">Discussion List</Typography>
                  <Link className={classes.link} to="https://discuss.seknox.com/c/trasa">
                    TRASA Community Forum
                  </Link>
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
