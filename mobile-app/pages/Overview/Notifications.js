import React, {Component} from 'react';
import {
  Body,
  Button,
  Card,
  CardItem,
  Container,
  Content,
  DatePicker,
  Header,
  Icon,
  Left,
  List,
  ListItem,
  Right,
  Text,
  Title,
} from 'native-base';
import {FlatList, ScrollView, TouchableOpacity} from 'react-native';
import AsyncStorage from '@react-native-community/async-storage';
import {getStoredNotifications} from './NotificationStorage';

export class ViewNotification extends Component {
  constructor(props) {
    super(props);
    this.state = {chosenDate: new Date()};
    this.setDate = this.setDate.bind(this);
  }
  setDate(newDate) {
    this.setState({chosenDate: newDate});
  }
  render() {
    const appName = this.props.route.params?.appName || 'Test App';
    const requester = this.props.route.params?.requester;
    const reason = this.props.route.params?.reason;
    const time = this.props.route.params?.time;
    return (
      <Container>
        <Header>
          <Left>
            <Button transparent>
              <Icon
                name={'arrow-back'}
                onPress={() => {
                  this.props.navigation.navigate('Home');
                }}
              />
            </Button>
          </Left>
          <Body>
            <Title>Notifications</Title>
          </Body>
          <Right />
        </Header>
        <ScrollView>
          {/*<Card>
                        <CardItem header>
                            <Text>AdHoc request from {requester} in {appName} app.</Text>
                        </CardItem>
                        <CardItem >
                            <Text>{reason}</Text>
                        </CardItem>
                        <CardItem cardHeader>
                            <Text>{time}</Text>
                        </CardItem>
                        <CardItem cardHeader>
                            <Text>Go to dashboard to grant or deny the request.</Text>
                        </CardItem>

                    </Card>
*/}

          <Card>
            <CardItem header>
              <Text>Request from</Text>
            </CardItem>
            <CardItem>
              <Text>{requester}</Text>
            </CardItem>

            <CardItem header>
              <Text>App Name</Text>
            </CardItem>
            <CardItem>
              <Text>{appName}</Text>
            </CardItem>

            <CardItem header>
              <Text>Reason</Text>
            </CardItem>
            <CardItem>
              <Text>{reason}</Text>
            </CardItem>

            <CardItem header>
              <Text>Time</Text>
            </CardItem>
            <CardItem>
              <Text>{time}</Text>
            </CardItem>
          </Card>
        </ScrollView>
      </Container>
    );
  }
}

export class Notifications extends React.Component {
  constructor(props) {
    super(props);
    this.state = {notifs: []};
  }

  componentDidMount() {
    getStoredNotifications((notifs) => {
      this.setState({notifs: notifs});
    });
  }

  showNotifiacationDetails = (value) => {
    value.isRead = true;
    if (value.ASkey) {
      AsyncStorage.mergeItem(value.ASkey, JSON.stringify(value));
    }

    if (value.type === 'GRANT' || value.type === 'REQUEST') {
      this.props.navigation.navigate('ViewNotification', value);
    } else {
      this.props.navigation.navigate('TwoFA', value);
    }
  };

  render() {
    return (
      <Container>
        <Header>
          <Left>
            <Button transparent>
              <Icon
                name={'arrow-back'}
                onPress={() => {
                  this.props.navigation.navigate('Home');
                }}
              />
            </Button>
          </Left>
          <Body>
            <Title style={{color:'white'}}>Notifications</Title>
          </Body>
          <Right />
        </Header>
        <ScrollView>
          {this.state?.notifs?.length?null:
              <Card transparent>
                <CardItem >
                  <Text>You don't have any notifications yet</Text>

                </CardItem>
              </Card>
          }
          {this.state.notifs.map((val, index) => {
            return (
              <TouchableOpacity
                  key={index}
                onPress={() => {
                  this.showNotifiacationDetails(val);
                }}>
                <Card button key={index}>
                  <CardItem header>
                    {val.type == 'GRANT' && (
                      <Text>
                        AdHOC granted by {val.requestee} in {val.appName} app
                      </Text>
                    )}
                    {val.type == 'REQUEST' && (
                      <Text>
                        AdHOC requested by {val.requester} in {val.appName} app
                      </Text>
                    )}
                    {(val.type != 'REQUEST' || val.type != 'GRANT') && (
                      <Text>
                        U2F notification from {val.ipAddr} in {val.appName}
                      </Text>
                    )}
                  </CardItem>
                  {/*<CardItem>*/}
                  {/*    <Text>{val.isRead?"READ":"UNREAD"}</Text>*/}
                  {/*    <Icon name={val.isRead?"star-outline":"star"}/>*/}
                  {/*</CardItem>*/}

                  <CardItem>
                    <Text note>{val.time}</Text>
                  </CardItem>
                </Card>
              </TouchableOpacity>
            );
          })}
        </ScrollView>
      </Container>
    );
  }
}
