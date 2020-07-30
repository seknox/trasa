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

export default function UsersPerIdp(props: any) {
  function getserviceTypes(v: any) {
    return v.map((i: any) => {
      return i.name;
    });
  }

  const option = {
    title: {
      text: 'Users Per Idp',
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
    // legend: {
    //     data:['l1','l2','l3']
    // },
    xAxis: [
      {
        type: 'value',
      },
    ],
    yAxis: [
      {
        type: 'category',
        data: getserviceTypes(props.totalIdps), // this.props.ServiceByType, // Object.keys(failed),
        // value: this.props.totalIdps[1]
      },
    ],
    series: [
      {
        name: 'Users Per Idp',
        type: 'bar',
        stack: 'component',
        data: props.totalIdps,
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
