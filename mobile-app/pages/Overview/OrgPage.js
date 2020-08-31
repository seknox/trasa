import React, {Component} from 'react';
//import RNSecureKeyStore from "react-native-secure-key-store";
import * as SecureStore from 'expo-secure-store';
import axios from 'axios';
import Constants from '../Extra/Constants';
import AsyncStorage from '@react-native-community/async-storage';
import {
  Body,
  Button,
  Card,
  CardItem,
  Container,
  Content,
  Header,
  Icon,
  Left,
  Right,
  Text,
  Title,
} from 'native-base';
import {Alert, Animated, Clipboard, Easing, Image, ScrollView, View} from 'react-native';
import getLogoSource from './Logos/Logos';
import generateTOTP from '../../utils/totpGenerate';
import CircleProgress from '../Extra/Components/CircleProgress';
import OTP from 'otp-client';
import Device from '../Device/Device';
import {CommonActions} from '@react-navigation/native';

export default class OrgPage extends Component {
  constructor(props) {
    super(props);
    this.state = {
      apps: [],
      notifs: [
        {appName: 'win32-A6', username: 'bhrg3se'},
        {appName: 'ununtu45-A6', username: 'haxor66'},
      ],
      info: '{}',
    };
    //this.notifs=[]
  }

  onDelete = async (
    secret, //delete a  item from  public apps
  ) => {
    let tempList = await AsyncStorage.getItem('TOTPlist');
    let totpList = JSON.parse(tempList);
    Alert.alert(
      'Confirm Delete',
      'Are You Sure?',
      [
        {
          text: 'No',
          style: 'cancel',
        },
        {
          text: 'Yes',
          onPress: () => {
            // console.log(totpList);
            totpList = totpList.filter((item) => secret != item.secret);
            AsyncStorage.setItem('TOTPlist', JSON.stringify(totpList)).then(
              (r) => {
                this.props.navigation.navigate('Home',{updated:totpList});

              },
            );
          },
        },
      ],
      {cancelable: false},
    );
  };

  render() {
    const params = this.props.route.params;
    const name = params?.issuer || params?.OrgPage;
    const acc = params?.account;
    const secret = params?.secret;
    // const onDelete =params?.callback;
    const type = params?.type;
    //  const logo=navigation.getParam('logo', 'briefcase')

    return (
      <Container>
        <Header>
          <Left>
            <Button transparent>
              <Icon
                name={'arrow-back'}
                onPress={() => {
                  this.props.navigation.goBack();
                }}
              />
            </Button>
          </Left>
          <Body>
            <Title style={{color:'white'}}>TOTP</Title>
          </Body>
          <Right />
        </Header>
        <ScrollView>
          <Card>
            <CardItem>
              <Left>
                <Image
                  source={getLogoSource(name)}
                  style={{width: 30, height: 30}}
                />
                <Body>
                  <Text> {name}</Text>
                  <Text note>{acc}</Text>
                </Body>
              </Left>
            </CardItem>
            <View>
              <Body>
                {/*{this.props.navigation.getParam("floCodes")?<FlotingCodes secret={secret}/>:<NonFloatingCodes secret={secret}/>}*/}
                <NonFloatingCodes secret={secret} />
              </Body>
            </View>

            <CardItem>
              <Left>
                <Button
                  transparent
                  onPress={() => {
                    Clipboard.setString(generateTOTP(secret));
                    {
                      /*Toast.show({
                                text: "Code copied to clipboard",
                                buttonText: "Okay"
                            })*/
                    }
                  }}>
                  <Icon active name="clipboard" style={{color: '#0e0343'}} />
                  <Text style={{color: '#0e0343'}}>Copy Code</Text>
                </Button>
              </Left>
              <Body />
              <Right>
                <Button
                  transparent
                  onPress={() => {
                    this.onDelete(secret);
                  }}>
                  <Icon active name="trash" style={{color: '#0e0343'}} />
                  <Text style={{color: '#0e0343'}}>Delete</Text>
                </Button>
              </Right>
            </CardItem>
          </Card>

          {/* <Text>
                        {this.state.info}
                    </Text> */}
        </ScrollView>
      </Container>
    );
  }
}

class NonFloatingCodes extends Component {
  constructor(props) {
    super(props);
    this.state = {otp: generateTOTP(props.secret), time: this.second()};
  }
  componentDidMount() {
    setInterval(() => {
      this.setState({
        time: this.second(),
        otp: generateTOTP(this.props.secret),
      });
    }, 500);
  }
  second = () => new Date().getUTCSeconds() % 30;
  render() {
    return (
      <View style={{flex: 1}}>
        <CircleProgress
          radius={30}
          percent={(this.state.time / 30) * 100}
          otp={this.state.otp}
          time={this.state.time}
        />

        {/*<Text style={{fontSize:70,left:3}}>{this.state.otp}</Text>*/}
        {/*<Right>*/}
        {/*  */}
        {/*</Right>*/}
      </View>
    );
  }
}
