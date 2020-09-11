import Grid from '@material-ui/core/Grid';
import axios from 'axios';
import echarts from 'echarts';
import ReactEcharts from 'echarts-for-react'; // or var ReactEcharts = require('echarts-for-react');
import React, { useEffect, useState } from 'react';
import Constants from '../../../Constants';
import { AccessStatsFilterProps } from '../../../types/analytics';
import { EchartDivHeight, EchartElementHeight } from '../../../utils/Responsive';
import useStyles from '../../../utils/styles';

const colorPalette = ['#000066', '#37A2DA', '#030417', '#67E0E3', '#1B2948', '#32C5E9', '#000086'];

echarts.registerTheme('trasaThemeLH', {
  color: colorPalette,
  // height : window.innerHeight < 750 ? '100': '520',
  height: EchartElementHeight(),
  // height: '500',
  // backgroundColor: '#030417'
  textStyle: {
    fontFamily: 'Open Sans, Rajdhani',
  },
});

function OptionBuilder(data: any) {
  const option = {
    title: {
      text: 'Login Hours',
      // bottom: 0
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross',
        crossStyle: {
          color: '#999',
        },
      },
    },
    grid: {
      left: 2,
      bottom: 100,
      right: 5,
      containLabel: true,
    },

    xAxis: [
      {
        type: 'category',
        data: [
          0,
          1,
          2,
          3,
          4,
          5,
          6,
          7,
          8,
          9,
          10,
          11,
          12,
          13,
          14,
          15,
          16,
          17,
          18,
          19,
          20,
          21,
          22,
          23,
        ],
        axisPointer: {
          type: 'shadow',
        },
      },
    ],
    yAxis: [
      {
        type: 'value',
        // name: 'Total Logins Events',

        // max: maxy,
        // interval: 1000,
        // axisLabel: {
        //     formatter: '{value}'
        // }
      },
    ],
    series: [
      {
        name: 'total logins',
        type: 'bar',
        data,
      },
    ],
  };

  return option;
}

const time = ['Today', 'Yesterday', '7 Days', '30 Days', '90 Days', 'All time'];
const status = ['All', 'Success', 'Failed'];

export default function LoginPatternByHour(props: AccessStatsFilterProps) {
  const [opt, setOpt] = useState({});
  const [timeFilter, setTimeFilter] = useState('Today');
  const [statusFilter, setStatusFilter] = useState('All');
  const classes = useStyles();

  useEffect(() => {

    const reqPath = `${Constants.TRASA_HOSTNAME}/api/v1/stats/loginhours/${props.entityType}/${props.entityID}/${timeFilter}/${statusFilter}`;

    axios
      .get(reqPath)
      .then((r) => {
        // console.log(r.data.data[0])
        const lopt = OptionBuilder(r.data.data[0]);

        setOpt(lopt);
      })
      .catch((error) => {
        console.log(error);
      });
  }, [timeFilter, statusFilter, props.entityType, props.entityID]);

  function handleChange(event: React.ChangeEvent<HTMLSelectElement>) {
    setTimeFilter(event.target.value);
  }

  function changeStatusFilter(event: React.ChangeEvent<HTMLSelectElement>) {
    setStatusFilter(event.target.value);
  }

  return (
    <div>
      <ReactEcharts
        option={opt}
        notMerge
        lazyUpdate
        theme="trasaThemeLH"
        style={EchartDivHeight()}
      />
      <Grid container spacing={2}>
        <Grid item xs={6}>
          <select className={classes.selectAnalytics} name="time" onChange={(e) => handleChange(e)}>
            {time.map((name) => (
              <option key={name} value={name}>
                {name}
              </option>
            ))}
          </select>
        </Grid>
        <Grid item xs={6}>
          <select
            className={classes.selectAnalytics}
            name="time"
            onChange={(e) => changeStatusFilter(e)}
          >
            {status.map((name) => (
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
