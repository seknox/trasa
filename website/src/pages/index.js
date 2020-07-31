import Link from '@docusaurus/Link';
import useBaseUrl from '@docusaurus/useBaseUrl';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import clsx from 'clsx';
import React from 'react';
import styles from './styles.module.css';
// import Button from '../comps/button'
import MainCta from '../comps/main';

function Home() {
  const context = useDocusaurusContext();
  const { siteConfig = {} } = context;
  return (
    <Layout
      title={`${siteConfig.title}`}
      description="Description will go into a meta tag in <head />"
    >
      <header className={clsx('hero hero--primary', styles.heroBanner)}>
        <div className="container">
          {/* <h1 className="hero__title">zero trust access manager </h1>
          <p className="hero__subtitle">manage secure remote access to web apps, SSH, RDP and database services</p>
          <p className="hero__subtitle">zero trust</p>
          <div className={styles.buttons}>
            <Link
              className={clsx(
                'button button--outline button--secondary button--lg',
                styles.getStarted,
              )}
              to={useBaseUrl('docs/')}>
              Get Started
            </Link>
          </div> */}
          <MainCta />
        </div>
      </header>

      <main />
    </Layout>
  );
}

export default Home;
