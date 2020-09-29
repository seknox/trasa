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
    <Layout title={`${siteConfig.title}`} description="zero trust service access">
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
                    Zero Trust Service Access
                  </Typography>
                  <Typography variant="body1" style={{ textAlign: 'center' }}>
                    Verify user identity, device security hygiene, and security policies, <br /> and
                    verify it in every access.
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

                <Grid item xs={12} sm={12} md={8}>
                  <Typography variant="h2">Identity Aware</Typography>

                  <Typography variant="body1">
                    Control access based on JIT (just in time) risk calculated by user's profile and
                    user devices.
                  </Typography>

                  <Typography variant="h2">Application Layer conrol</Typography>

                  <Typography variant="body1">
                    Being Layer 7 proxy, Trasa IAP understands user's sessions in more detail. This
                    is powerfull to identify malicious user actions even though was authenticated
                    within common trust.
                  </Typography>

                  <Typography variant="h2">Privilege Access Management features</Typography>

                  <Typography variant="body1">
                    TRASA IAP is PAM capable. Meaning it can manage Privilege sessions, record it,
                    control it. Hackers and malicious insiders are known to exploit trust and
                    perform malicious actions once authorized. Our system performs deep action
                    audits to identify malicious actions performed within authorized session.
                  </Typography>

                  <Typography variant="h2"> Proactive Monitoring</Typography>

                  <Typography variant="body1">
                    We continuously monitor threats to your organization's internal services and
                    compromised user accounts in both surface and deep web.
                  </Typography>
                </Grid>
              </Grid>
            </ThemeBase>
          </div>
        </section>
      </main>
    </Layout>
  );
}
