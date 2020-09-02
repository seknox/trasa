import React, {Component} from 'react';
import {StyleSheet, View, Clipboard, TouchableOpacity} from 'react-native';
import {Text, Toast} from 'native-base';
import {AnimatedCircularProgress} from 'react-native-circular-progress';
import {Circle, Text as SVGText, Svg, TSpan} from 'react-native-svg';
import generateTOTP from '../../../utils/totpGenerate';

const styles = StyleSheet.create({
  circle: {
    overflow: 'hidden',
    position: 'relative',
    justifyContent: 'center',
    alignItems: 'center',
  },
  loader: {
    position: 'absolute',
    top: 0,
  },
});

//Toast.show({text:"Copied",buttonText: 'Okay'})

export default class CircularProgress extends React.Component {
  render():
    | React.ReactElement<any>
    | string
    | number
    | {}
    | React.ReactNodeArray
    | React.ReactPortal
    | boolean
    | null
    | undefined {
    return (
      <AnimatedCircularProgress
        size={300}
        width={16}
        fill={this.props.percent}
        tintColor="#00e0ff"
        backgroundColor="#3d5875"
        padding={0}
        children={() => (
          <TouchableOpacity
            onPress={() => {
              Clipboard.setString(this.props.otp);
            }}>
            <Text style={{fontSize: 50}}>{this.props.otp}</Text>
          </TouchableOpacity>
        )}
        renderCap={({center}) => (
          <Svg>
            <Circle
              cx={center.x}
              cy={center.y}
              r="8"
              fill="#00e0ff"
              height={'10'}
            />
          </Svg>
        )}
        // dashedBackground={{ width: 10, gap: 20 }}
      />
    );
  }
}

export class CircleProgress extends Component {
  // propTypes: {
  //   color: React.PropTypes.string,
  //   bgcolor: React.PropTypes.string,
  //   radius: React.PropTypes.number,
  //   percent: React.PropTypes.number
  // }
  // constructor(props) {
  //   super(props);
  //   //this.state = this.compute(this.props.percent);
  // }
  // componentWillReceiveProps(nextProps) {
  //   this.setState(this.compute(nextProps.percent));
  // }
  compute(percent) {
    let degree = percent * 3.6 + 'deg';
    let color = this.props.color;
    if (percent >= 50) {
      color = this.props.bgcolor;
      degree = (percent - 50) * 3.6 + 'deg';
    }
    return {percent, degree, color};
  }
  render() {
    let {percent, degree, color} = this.compute(this.props.percent);
    return (
      <View
        style={[
          styles.circle,
          {
            width: this.props.radius * 2,
            height: this.props.radius * 2,
            borderRadius: this.props.radius,
            backgroundColor: this.props.bgcolor,
          },
        ]}>
        <View
          style={[
            styles.loader,
            {
              left: 0,
              width: this.props.radius,
              height: this.props.radius * 2,
              backgroundColor: this.props.color,
              borderTopLeftRadius: this.props.radius,
              borderBottomLeftRadius: this.props.radius,
            },
          ]}
        />
        <View
          style={[
            styles.loader,
            {
              left: this.props.radius,
              width: this.props.radius,
              height: this.props.radius * 2,
              backgroundColor: color,
              borderTopRightRadius: this.props.radius,
              borderBottomRightRadius: this.props.radius,
              transform: [
                {
                  translateX: -this.props.radius / 2,
                },
                {
                  rotate: degree,
                },
                {
                  translateX: this.props.radius / 2,
                },
              ],
            },
          ]}
        />
      </View>
    );
  }
}

CircleProgress.defaultProps = {
  color: '#3c80f7',
  bgcolor: '#e3e3e3',
};
