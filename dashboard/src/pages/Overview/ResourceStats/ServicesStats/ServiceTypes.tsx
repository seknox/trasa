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

export default function serviceTypes(props: any) {
  function getserviceTypes(v: any) {
    return v.map((a: any) => {
      return a.name;
    });
  }

  const option = {
    title: {
      text: 'Service Types',
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
        type: 'value', // as "value" | "category" | "time" | "log" | undefined
      },
    ],
    yAxis: [
      {
        // gridIndex: 1,
        type: 'category',
        data: getserviceTypes(props.ServiceByType), // this.props.ServiceByType, // Object.keys(failed),
      },
    ],
    series: [
      {
        name: 'Total services',
        type: 'bar',
        stack: 'component',
        data: props.ServiceByType,
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
