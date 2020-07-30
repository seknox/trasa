import echarts, { EChartOption } from 'echarts';
import ReactEcharts from 'echarts-for-react'; // or var ReactEcharts = require('echarts-for-react');
import React from 'react';

const colorPalette = [
  '#000066', // '#1B2948','#32C5E9',
];

echarts.registerTheme('trasaThemeB', {
  top: 20,
  color: colorPalette,
  height: '150',
  // backgroundColor: '#030417',
  textStyle: {
    fontFamily: 'Open Sans, Rajdhani',
    color: '#030417',
    fontSize: '34px',
  },
});

export default function BrowserTypes(props: any) {
  function getserviceTypes(v: any) {
    if (v === null || v === undefined) {
      return '';
    }
    return v.map((b: any) => {
      return `${b.name} ${b.version}`;
    });
  }

  const option = {
    title: {
      text: 'Browsers',
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        axisPointer: {
          type: 'shadow',
        },
      },
    },
    grid: {
      left: 2,
      bottom: 5,
      right: 5,
      containLabel: true,
    },
    xAxis: [
      {
        type: 'value',
      },
    ],
    yAxis: [
      {
        // gridIndex: 1,
        type: 'category',
        data: getserviceTypes(props.browsersByType), // this.props.ServiceByType, // Object.keys(failed),
      },
    ],
    series: [
      {
        name: 'Browser',
        type: 'bar',
        stack: 'component',
        data: props.browsersByType,
        label: {
          show: true,
          position: 'inside',
          color: 'white',
        },
      },
    ],
  };

  return (
    <div>
      <ReactEcharts
        option={option as EChartOption}
        notMerge
        lazyUpdate
        theme="trasaThemeB"
        style={{ height: 250 }}
      />
    </div>
  );
}
