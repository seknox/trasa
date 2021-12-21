import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import clsx from 'clsx';
import React from 'react';
// import Button from '../comps/button'
import PrimaryCTA from '../components/landing/primary';
import SecondaryCTA from '../components/landing/secondary';
import styles from './styles.module.css';

import Showcase from '../components/landing/showcase';


function Home() {
  const context = useDocusaurusContext();
  const { siteConfig = {} } = context;
  return (
    <Layout
      title={`${siteConfig.title}`}
      description="Identity aware access proxy for secure remote access to Web, SSH, RDP and Database services"
    >
      <header className={clsx('hero hero--primary', styles.heroBanner)}>
        <div className="container">
          <PrimaryCTA />
        </div>
      </header>

      <main>
        <section className={styles.features}>
          {/* <Features /> */}
          {/* <UseCase /> */}
          <SecondaryCTA />
          <div className="container">
            <Showcase />
          </div>
        </section>
      </main>
    </Layout>
  );
}

export default Home;
