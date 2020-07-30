import axios from 'axios';
import echarts, { EChartOption } from 'echarts';
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

echarts.registerTheme('trasaThemeIP', {
  top: 20,
  color: colorPalette,
  height: EchartMapElementHeight(),
  // height: '520',
  // backgroundColor: '#030417'
  textStyle: {
    fontFamily: 'Open Sans, Rajdhani',
  },
});

type AggIp = {
  key: string;
  name: string;
  value: number;
};

export default function IPRadialPlot(props: AccessStatsFilterProps) {
  // const [opt, setOpt] = useState({});
  const [ipPool, setIPPool] = useState<AggIp[]>([]);

  useEffect(() => {
    const reqPath = `${Constants.TRASA_HOSTNAME}/api/v1/stats/ips/${props.entityType}/${props.entityID}/${props.timeFilter}/${props.statusFilter}`;

    axios
      .get(reqPath)
      .then((response) => {
        const resp = response.data.data[0];
        if (resp.children.length === 0) {
          setIPPool([]);
        } else {
          setIPPool([resp]);
        }
      })
      .catch((error) => {
        console.log('Error', error.message);
      });
  }, [props.timeFilter, props.statusFilter, props.entityType, props.entityID]);

  const option = {
    title: {
      text: 'Accessed IP Breakdown',
      // left: 'center'
    },

    tooltip: {
      trigger: 'item', // as "none" | "item" | "axis" | undefined,
      triggerOn: 'mousemove', // as "none" | "mousemove" | "click" | "mousemove|click" | undefined
    },

    series: [
      {
        type: 'tree',
        roam: true,
        // layout: 'force',
        data: ipPool, // : [],
        // top: '18%',
        // bottom: '100%',
        layout: 'radial',
        // orient: 'vertical',
        symbol: 'diamond', // 'path://M111.88095,129.17857,48.24294,129.0783,28.673076,68.523967,80.216246,31.199603,131.64154,68.686207Z',
        // color: 'navy',
        symbolSize: 10,
        initialTreeDepth: 2,
        itemStyle: {
          color: '#000066',
          borderColor: '#000066',
        },
        animationDurationUpdate: 750,
      },
    ],
  };
  return (
    <div>
      <ReactEcharts
        option={option as EChartOption}
        notMerge
        theme="trasaThemeIP"
        style={EchartMapDivHeight()}
        // style={{minHeight: window.innerHeight < 750 ? 500: 690}}
        // style={{minHeight: 690}}
      />
    </div>
  );
}
