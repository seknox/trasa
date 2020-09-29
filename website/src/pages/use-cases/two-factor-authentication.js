import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import clsx from 'clsx';
import React from 'react';

import styles from '../styles.module.css';
import { Typography, Grid, Paper } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import ThemeBase from '../../components/muiTheme';

import useBaseUrl from '@docusaurus/useBaseUrl';

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
    <Layout title={`${siteConfig.title}`} description="Two Factor Authentication">
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
                    Two Factor Authentication
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
                  <img src={zeroTrust} width={500} alt="zero trust access" />
                </Grid>
                <Grid item xs={9}>
                  <Typography component="p" variant="body1" align="center" gutterBottom>
                    Multifactor authentication for everything. <br /> TRASA enables native
                    agent-based (most-secure) and agentless (controlled by TRASA access-proxy)
                    two-factor authentication in web, RDP, SSH, and Database services. TRASA
                    supports TOTP, Push U2F and Yubikey as authentorization method.
                  </Typography>
                  <div className={classes.aligner}>
                    <ul>
                      <li>Two factor for windows server and workstation</li>
                      <li>Two factor for linux/unix server and workstation</li>
                      <li>
                        Two factor for hardware devices (router, firewall switches. any vendor){' '}
                      </li>
                      <li>Two factor for web and email services. </li>
                    </ul>
                  </div>
                </Grid>
              </Grid>
            </ThemeBase>
          </div>
        </section>
      </main>
    </Layout>
  );
}
