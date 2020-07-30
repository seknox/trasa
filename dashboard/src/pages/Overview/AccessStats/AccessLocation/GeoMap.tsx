import axios from 'axios';
import echarts from 'echarts';
import ReactEcharts from 'echarts-for-react'; // or var ReactEcharts = require('echarts-for-react');
import React, { useEffect, useState } from 'react';
import Constants from '../../../../Constants';
import { AccessStatsFilterProps } from '../../../../types/analytics';
import { EchartMapDivHeight, EchartMapElementHeight } from '../../../../utils/Responsive';

const colorPalette = [
  '#000066',
  '#37A2DA',
  '#67E0E3',
  '#dd6b66',
  '#759aa0',
  '#8dc1a9',
  '#ea7e53', // '#1B2948','#32C5E9',
];

echarts.registerTheme('trasaThemeW', {
  // top: 20,
  color: colorPalette,
  // height: '320',
  height: EchartMapElementHeight(),
  // height : window.innerHeight < 750 ? '350': '520',
  // height: '520',
  // backgroundColor: '#030417'
  textStyle: {
    fontFamily: 'Open Sans, Rajdhani',
  },
});

require('echarts/map/js/world.js');

const onEvents = {
  // 'click': this.onChartClick,
  // 'legendselectchanged': this.onChartLegendselectchanged
};

function OptionBuilder(sourceData: any) {
  const option = {
    title: {
      text: 'Access By Location',
      top: 'top',
    },
    grid: {
      left: 2,
      // bottom: 5,
      right: 5,
      containLabel: true,
    },

    tooltip: {
      trigger: 'item',
      formatter(params: any) {
        return `${params.seriesName}<br/>${params.name} : ${params.value}`;
      },
    },
    toolbox: {
      show: true,
      orient: 'vertical',
      left: 'right',
      top: 'center',
      feature: {
        dataView: { readOnly: false },
        restore: {},
        saveAsImage: {},
      },
    },
    visualMap: {
      min: 0,
      max: 1000,
      text: ['High', 'Low'],
      realtime: false,
      calculable: true,
      inRange: {
        color: ['#000066', '#1B2948', '#800000'], // '#000066']
      },
    },
    series: [
      {
        name: "Event's By Location",
        type: 'map',
        mapType: 'world',
        roam: true,
        itemStyle: {
          normal: {
            borderColor: 'white',
            color: '#1B2948',
          },
        },
        data: sourceData,
      },
    ],
  };

  return option;
}

export default function GeoMap(props: AccessStatsFilterProps) {
  const [opt, setOpt] = useState({});

  useEffect(() => {
    axios
      .get(
        `${Constants.TRASA_HOSTNAME}/api/v1/stats/mapplot/${props.entityType}/${props.entityID}/${props.timeFilter}/${props.statusFilter}`,
      )
      .then((response) => {
        if (response.data.status === 'success') {
          const optl = OptionBuilder(response.data.data[0]);
          setOpt(optl);
        }
      })
      .catch((error) => {
        if (error.response.status === 403) {
          window.location.href = '/login';
        }

        if (error.response) {
          console.log(error.response.data);
        } else {
          console.log('Error', error.message);
        }
      });
  }, [props.statusFilter, props.timeFilter, props.entityType, props.entityID]);

  return (
    <div>
      <ReactEcharts
        option={opt}
        notMerge
        lazyUpdate
        theme="trasaThemeW"
        onEvents={onEvents}
        style={EchartMapDivHeight()}
        // style={{minHeight: window.innerHeight < 750 ? 500: 690}}
        // style={{minHeight: 690}}
      />
    </div>
  );
}
