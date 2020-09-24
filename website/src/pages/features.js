import Layout from '@theme/Layout';
import React from 'react';
import FeaturesComp from '../components/features';
import styles from './styles.module.css';
import { Typography } from '@material-ui/core';
import clsx from 'clsx';
import ThemeBase from '../components/muiTheme';

export default function FeaturesPage() {
  return (
    <Layout title={`TRASA features`} description="TRASA features">
      <header className={clsx('hero hero--primary', styles.heroBanner)}>
        <div className="container">
          <ThemeBase>
            <Typography variant="h1">Features</Typography>
          </ThemeBase>
        </div>
      </header>
      <main>
        <section className={styles.features}>
          <div className="container">
            <FeaturesComp />
          </div>
        </section>
      </main>
    </Layout>
  );
}
