import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import { Typography } from '@material-ui/core';
import Layout from '@theme/Layout';
import clsx from 'clsx';
import React from 'react';
import styles from './styles.module.css';
import Disclosure from '../components/security/disclosure';

function Home() {
  const context = useDocusaurusContext();
  const { siteConfig = {} } = context;
  return (
    <Layout title={`${siteConfig.title}`} description="TRASA security">
      <header className={clsx('hero hero--primary', styles.heroBanner)}>
        <div className="container">
          <Disclosure />
        </div>
      </header>

      <main />
    </Layout>
  );
}

export default Home;
