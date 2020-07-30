import axios from 'axios';
import echarts from 'echarts';
import ReactEcharts from 'echarts-for-react'; // or var ReactEcharts = require('echarts-for-react');
import React, { useEffect, useState } from 'react';
import Constants from '../../../Constants';
import { AccessStatsFilterProps } from '../../../types/analytics';
import { EchartDivHeight } from '../../../utils/Responsive';
import useStyles from '../../../utils/styles';

const colorPalette = ['#000066', '#37A2DA', '#030417', '#67E0E3', '#1B2948', '#32C5E9', '#000086'];

function getHeight() {
  const h = window.innerHeight;
  switch (true) {
    case h < 750:
      return '220';
    case h > 750 && h < 920:
      return '280';
    case h > 920 && h < 1500:
      return '320';
    default:
      return '520';
  }
}

function OptionBuilder(sourceData: any) {
  const option = {
    title: {
      text: 'Top Authentication Failed Reasons',
      // top: '100'
      // link: "let's refer to our document site here"
    },
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)',
      position: ['10%', '32%'],
    },
    grid: {
      left: 2,
      top: 5,
      bottom: 5,
      right: 5,
      containLabel: true,
    },

    series: [
      {
        // roam: true,
        // smooth: true,
        // layout: 'force',
        // name:"Top Authentication Failed Reasons",
        type: 'pie',
        top: '10',
        // top : '100',
        // selectedMode: 'single',
        radius: ['0%', '10%'],
        height: 'auto',
        width: 'auto',
        // radius: [0, 300],
        // radius: ['40%', '55%'],
        roseType: 'area',
        label: {
          normal: {
            position: 'inner',
          },
        },
        labelLine: {
          normal: {
            show: false,
          },
        },
        data: [{ value: 0, name: '' }],
        // data: sourceData
      },
      {
        top: '30',
        name: 'Failed Reason:',
        type: 'pie',
        radius: ['0%', '40%'],
        roam: true,
        data: sourceData,
        label: {
          color: '#030417',
        },
      },
    ],
  };
  // [ {"name":"root","value":1}, {"name":"yetanotheruniqueemail@mailinator.com","value":9}, {"name":"bhrg3se","value":60}, {"name":"sakshyam@seknox.com","value":16} ]
  return option;
}

const time = ['All', 'Today', 'Yesterday', '7 Days', '30 Days', '90 Days'];

echarts.registerTheme('trasaThemeF', {
  color: colorPalette,
  // height : window.innerHeight < 750 ? '100': '520',
  height: getHeight(),
  // height: '500',
  // backgroundColor: '#030417'
  textStyle: {
    fontFamily: 'Open Sans, Rajdhani',
  },
});

export default function AppAndUserFailedPlot(props: AccessStatsFilterProps) {
  const [opt, setOpt] = useState({});
  const [filter, setFilter] = useState('All');
  useEffect(() => {
    axios
      .get(
        `${Constants.TRASA_HOSTNAME}/api/v1/stats/failedreasons/${props.entityType}/${props.entityID}/${filter}`,
      )
      .then((response) => {
        const optl = OptionBuilder(response.data.data[0]);
        setOpt(optl);
      })
      .catch((error) => {
        if (error.response) {
          console.log(error.response.data);
        } else {
          // Something happened in setting up the request that triggered an Error
          console.log('Error', error.message);
        }
      });
  }, [filter, props.entityType, props.entityID]);

  function handleChange(event: React.ChangeEvent<HTMLSelectElement>) {
    setFilter(event.target.value);
  }

  const classes = useStyles();

  return (
    <div>
      {/* <ShowcaseButton name={"Prev"} onClick={()=>{this.setState({selectedApp:this.state.selectedApp-1})}}>Prev</ShowcaseButton> */}

      <ReactEcharts option={opt} theme="trasaThemeF" style={EchartDivHeight()} />
      <select className={classes.selectAnalytics} name="time" onChange={(e) => handleChange(e)}>
        {time.map((name) => (
          <option key={name} value={name}>
            {name}
          </option>
        ))}
      </select>
    </div>
  );
}
