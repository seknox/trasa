import Grid from '@material-ui/core/Grid';
import React, { useState } from 'react';
import { AccessStatsFilterProps } from '../../../../types/analytics';
import useStyles from '../../../../utils/styles';
import GeoMap from './GeoMap';
import IPMap from './IPPool';

const time = ['All', 'Today', 'Yesterday', '7 Days', '30 Days', '90 Days'];
const status = ['All', 'Success', 'Failed'];
const plot = ['Geo', 'IP'];

export default function AccessLocation(props: AccessStatsFilterProps) {
  // const [opt, setOpt] = useState({});
  const [timeFilter, setTimeFilter] = useState('Today');
  const [statusFilter, setStatusFilter] = useState('All');
  const [plotType, setplotType] = useState('IP');
  const classes = useStyles();

  function handleTimeChange(event: React.ChangeEvent<HTMLSelectElement>) {
    setTimeFilter(event.target.value);
  }

  function handleStatusChange(event: React.ChangeEvent<HTMLSelectElement>) {
    setStatusFilter(event.target.value);
  }

  function changePlotType(event: React.ChangeEvent<HTMLSelectElement>) {
    setplotType(event.target.value);
  }

  return (
    <div>
      {plotType === 'Geo' ? (
        <GeoMap
          statusFilter={statusFilter}
          timeFilter={timeFilter}
          entityType={props.entityType}
          entityID={props.entityID}
        />
      ) : (
        <IPMap
          entityType={props.entityType}
          entityID={props.entityID}
          statusFilter={statusFilter}
          timeFilter={timeFilter}
        />
      )}

      <Grid container spacing={2}>
        <Grid item xs={4}>
          <select
            className={classes.selectAnalytics}
            value={timeFilter}
            name="time"
            onChange={(e) => handleTimeChange(e)}
          >
            {time.map((name) => (
              <option key={name} value={name}>
                {name}
              </option>
            ))}
          </select>
        </Grid>{' '}
        <Grid item xs={4}>
          <select
            className={classes.selectAnalytics}
            value={statusFilter}
            name="time"
            onChange={(e) => handleStatusChange(e)}
          >
            {status.map((name) => (
              <option key={name} value={name}>
                {name}
              </option>
            ))}
          </select>
        </Grid>
        <Grid item xs={4}>
          <select
            className={classes.selectAnalytics}
            value={plotType}
            name="time"
            onChange={(e) => changePlotType(e)}
          >
            {plot.map((name) => (
              <option key={name} value={name}>
                {name}
              </option>
            ))}
          </select>
        </Grid>
      </Grid>
    </div>
  );
}
