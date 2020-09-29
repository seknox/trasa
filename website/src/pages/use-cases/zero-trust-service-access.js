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

                <Typography variant="body1">
                  Seknox zero trust access controller controls access based on risk rather than
                  static trusted policies. Stolen credentials? Compromised network? Compromised user
                  devices? Zero trust access controller can detect malicious access hidden between
                  trusted parameters and blocks it near real-time. <br />
                  We profile users, user devices and health of service to dynamically score risk for
                  access being requested. This makes any stolen credentials or comprmised accounts
                  useless for hackers to get unauthorized access.
                </Typography>

                <Typography variant="h2">The perimiter we know has changed.</Typography>

                <Typography variant="body1">
                  Access control based on VPN leads to security design based on trust. While the
                  benefits up VPN has helped users securely connect to protected services, it has
                  also let attackers misuse this trust and compromise more easlity. Afterall, all
                  that is needed to compromise is one trusted credential and than everything is
                  over.
                </Typography>

                <Typography variant="h2">Identity Aware</Typography>

                <Typography variant="body1">
                  Control access based on JIT (just in time) risk calculated by user's profile and
                  user devices.
                </Typography>

                <Typography variant="h2">Application Layer conrol</Typography>

                <Typography variant="body1">
                  Being Layer 7 proxy, Trasa IAP understands user's sessions in more detail. This is
                  powerfull to identify malicious user actions even though was authenticated within
                  common trust.
                </Typography>

                <Typography variant="h2">Privilege Access Management features</Typography>

                <Typography variant="body1">
                  TRASA IAP is PAM capable. Meaning it can manage Privilege sessions, record it,
                  control it. Hackers and malicious insiders are known to exploit trust and perform
                  malicious actions once authorized. Our system performs deep action audits to
                  identify malicious actions performed within authorized session.
                </Typography>

                <Typography variant="h2"> Proactive Monitoring</Typography>

                <Typography variant="body1">
                  We continuously monitor threats to your organization's internal services and
                  compromised user accounts in both surface and deep web.
                </Typography>
              </Grid>
            </ThemeBase>
          </div>
        </section>
      </main>
    </Layout>
  );
}
