import Grid from '@material-ui/core/Grid';
import axios from 'axios';
import echarts from 'echarts';
import ReactEcharts from 'echarts-for-react'; // or var ReactEcharts = require('echarts-for-react');
import Moment from 'moment';
import React, { useEffect, useState } from 'react';
import Constants from '../../../Constants';
import { AccessStatsFilterProps } from '../../../types/analytics';
import { EchartDivHeight, EchartElementHeight } from '../../../utils/Responsive';
import useStyles from '../../../utils/styles';

// //'#1B2948', '#000066' ,  '#37A2DA', '#67E0E3', '#dd6b66', '#759aa0',  '#8dc1a9', '#ea7e53' //'#1B2948','#32C5E9',
const colorPalette = ['#000066', '#37A2DA', '#030417', '#67E0E3', '#1B2948', '#32C5E9', '#000086'];

echarts.registerTheme('trasaThemeTE', {
  color: colorPalette,
  // height : window.innerHeight < 750 ? '100': '520',
  height: EchartElementHeight(),
  // height: '500',
  // backgroundColor: '#030417'
  textStyle: {
    fontFamily: 'Open Sans, Rajdhani',
  },
});

const onEvents = {
  // 'click': this.onChartClick,
  // 'legendselectchanged': this.onChartLegendselectchanged
};

function OptionBuilder(sourceData: any) {
  const option = {
    title: {
      text: "Today's Authentication Events",
      // link: "let's refer to our document site here"
    },
    legend: {},
    tooltip: {},
    trigger: 'item',

    grid: {
      left: 2,
      bottom: 5,
      right: 5,
      containLabel: true,
    },

    xAxis: [{ type: 'value', gridIndex: 0, name: 'Hour' }],
    yAxis: [{ type: 'value', gridIndex: 0, name: 'Minute' }],
    dataset: {
      dimensions: [
        'User',
        'Service Name',
        'Hour',
        'Minute',
        'Status',
        { name: 'Status', type: 'ordinal' },
      ],
      source: sourceData,
    },
    series: [
      {
        name: '',
        type: 'scatter',
        symbolSize: 10,
        symbol: 'roundRect',
        xAxisIndex: 0,
        yAxisIndex: 0,
        encode: {
          x: 'Hour',
          y: 'Minute',
          tooltip: [0, 1, 2, 3, 4],
        },
        // layout: 'force',
        // symbolSize: 20,
        // data: datArr,
      },
    ],
  };

  return option;
}

const status = ['All', 'Failed', 'Success'];

export default function TodaysEvent(props: AccessStatsFilterProps) {
  const [opt, setOpt] = useState({});
  const [total, setTotal] = useState(0);
  const [filter, setFilter] = useState('All');
  const classes = useStyles();

  useEffect(() => {
    // let reqPath = Constants.TRASA_HOSTNAME+'/api/v1/events/' + props.entityType + '/' + props.entityID + '/hexaevents'
    const reqPath = `${Constants.TRASA_HOSTNAME}/api/v1/stats/todayauths/${props.entityType}/${props.entityID}/${filter}`;

    axios
      .get(reqPath)
      .then((response) => {
        const data = response.data.status === 'success' && response.data.data[0];
        const totalCount = data && data.length;

        const datArr = data.map(function (val: any) {
          const d = Moment.unix(val.loginTime / 1000000000);
          val.hour = d.hour();
          val.minutes = d.minute();
          return [
            val.user,
            val.appName ? val.appName : 'Dashboard',
            val.hour,
            val.minutes,
            val.status,
          ];
        });

        const lopt = OptionBuilder(datArr);

        setOpt(lopt);
        setTotal(totalCount);
      })
      .catch((error) => {
        console.log('catched ', error);
      });
  }, [props.entityType, props.entityID, filter]);

  function handleChange(event: React.ChangeEvent<HTMLSelectElement>) {
    setFilter(event.target.value);
  }
  return (
    <div>
      <ReactEcharts
        option={opt}
        notMerge
        theme="trasaThemeTE"
        onEvents={onEvents}
        style={EchartDivHeight()}
      />

      <Grid container>
        <Grid item xs={6}>
          <select
            className={classes.selectAnalytics}
            value={filter}
            name="time"
            onChange={(e) => handleChange(e)}
          >
            {status.map((name) => (
              <option key={name} value={name}>
                {name}
              </option>
            ))}
          </select>
        </Grid>
        <Grid item xs={6}>
          <div className={classes.analyticsText}>Total: {total}</div>
        </Grid>
      </Grid>
    </div>
  );
}
