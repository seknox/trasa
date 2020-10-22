import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import clsx from 'clsx';
import React from 'react';
import { Typography, Grid } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import Link from '@docusaurus/Link';
import styles from './styles.module.css';

import ThemeBase from '../components/muiTheme';

const useStyles = makeStyles(() => ({
  ctaPad: {
    marginTop: 50,
  },
  link: {
    fontSize: 24,
    textDecoration: 'underline',
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
                  <Typography variant="h2">Discussion Forum</Typography>
                  <span>
                    <Typography variant="body1" component="span">
                      Ask or answer on topics related to TRASA:
                    </Typography>
                    <Link className={classes.link} to="https://discuss.trasa.io">
                      TRASA Community Forum
                    </Link>
                  </span>
                </div>
              </Grid>

              <Grid item xs={12} sm={12}>
                <div className={classes.ctaPad}>
                  <Typography variant="h2">Realtime communication</Typography>
                  <span>
                    <Typography variant="body1" component="span">
                      Chat with others:
                    </Typography>
                    <Link className={classes.link} to="https://discord.gg/4wRmuv9">
                      Join TRASA discord chat
                    </Link>
                  </span>
                </div>
              </Grid>

              <Grid item xs={12} sm={12}>
                <div className={classes.ctaPad}>
                  <Typography variant="h2">Blog</Typography>
                  <span>
                    <Typography variant="body1" component="span">
                      Read and share our blog:
                    </Typography>
                    <Link className={classes.link} to="/blog">
                      TRASA blog
                    </Link>
                  </span>
                </div>
              </Grid>

              <Grid item xs={12} sm={12}>
                <div className={classes.ctaPad}>
                  <Typography variant="h2">Contribute</Typography>
                  <span>
                    <Typography variant="body1" component="span">
                      Contribute to project:
                    </Typography>
                    <Link className={classes.link} to="https://github.com/seknox/trasa/issues">
                      Contribute in Github
                    </Link>
                  </span>
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
