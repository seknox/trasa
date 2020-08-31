import React from 'react';
import AsyncStorage from '@react-native-community/async-storage';
import {ScrollView, StyleSheet, View} from 'react-native';
import Avatar from 'react-native-badge-avatar';
import getLogoSource from './Logos/Logos';
import {
  Body,
  Button,
  Card,
  CardItem,
  Container,
  Fab,
  H1,
  Header,
  Icon,
  Input,
  Item,
  Left,
  Text,
} from 'native-base';
import Grid from 'react-native-grid-component';
import {SocialIcon} from 'react-native-elements';

// import RNTooltips from "react-native-tooltips";

export function Home(props) {
  const [fabActive, setFabActive] = React.useState(false);
  const [searchQuery, setSearchQuery] = React.useState('');
  const [totpList, setTotpList] = React.useState([]);

  //Here navigation states are defined like totp codes, etc
  React.useEffect(() => {
    let initialList = [];

    AsyncStorage.getItem('TOTPlist')
      .then((tempList) => {
        if (tempList === null) {
          AsyncStorage.setItem('TOTPlist', '[]');
        } else {
          initialList = JSON.parse(tempList);
        }
        //console.log(initialList);

        setTotpList(initialList);
      })
      .catch((e) => {
        alert(e);
      });
  }, [props.route.params?.updated]);

  // console.log(props.route.params)

  const goToOrgPage = (issuer, secret, account, type) => {
    props.navigation.navigate('OrgPage', {
      issuer: issuer,
      secret: secret,
      account: account,
      type: type,
      // callback: delete,
    });
  };

  const goToScanPage = () => {
    props.navigation.navigate('ScanPage');
  };

  const goToCodeInputPage = () => {
    props.navigation.navigate('CodeInput');
  };

  const onSearch = (q) => {
    setSearchQuery(q);
    //   setState({totpList:this.totpList.filter(item=>item.issuer.toUpperCase().includes(q.toUpperCase()))})
  };

  //this is a built in function of react-native-grid-component
  //it returns a single item
  const _renderItem = (data, i) => {
    let name = data.issuer;
    if (
      'angellist, codepen, envelope, etsy, facebook, flickr, foursquare, github-alt, github, gitlab, instagram, linkedin, medium, pinterest, quora, reddit-alien, soundcloud, stack-overflow, steam, stumbleupon, tumblr, twitch, twitter, google, google-plus-official, vimeo, vk, weibo, wordpress, youtube'.includes(
        data.issuer.toLowerCase(),
      )
    ) {
      name = name.toLowerCase();
      return (
        <View key={i} style={styles.item}>
          <SocialIcon
            // placeholder={(<SocialIcon type={data.issuer}/>)}
            type={name}
            size={50}
            style={{alignSelf: 'flex-start'}}
            onPress={() => {
              goToOrgPage(data.issuer, data.secret, data.account, data.type);
            }}
            badge={0}
          />

          <Text>{data.issuer}</Text>
        </View>
      );
    } else {
      return (
        <View key={i} style={styles.item}>
          <Avatar
            placeholder={getLogoSource(name)}
            size={50}
            style={{alignSelf: 'flex-start'}}
            onPress={() => {
              goToOrgPage(data.issuer, data.secret, data.account, data.type);
            }}
            badge={0}
          />
          <Text>{data.issuer}</Text>
        </View>
      );
    }
  };

  const _renderPlaceholder = (i) => <View style={styles.item} key={i} />;

  // const onGetTagetRef = (ref) => {
  //   this.setState({qrFabRef: ref});
  // };

  const filteredTotpList = totpList.filter((item) =>
    item.issuer.toUpperCase().includes(searchQuery.toUpperCase()),
  );
  let privateList = filteredTotpList.filter((item) => item.type === 'private');
  let publicList = filteredTotpList.filter((item) => item.type !== 'private');

  return (
    <Container>
      <Header searchBar rounded  >
        <Left>
          <Button
            onPress={() => {
              props.navigation.toggleDrawer();
            }}>
            <Icon name={'menu'} color={'white'}/>
          </Button>
        </Left>
        <Item>
          <Icon name="ios-search" />
          <Input placeholder="Search" onChangeText={onSearch} />
          {/*<Icon name="desktop" />*/}
        </Item>

      </Header>
      <ScrollView>
        <View
          style={{
            display: 'flex',
            justifyContent: 'space-evenly',
            flexDirection: 'column',
          }}>


          <Card padder>
            <CardItem header bordered>
              <H1 style={{fontFamily:"Rajdhani-Bold"}}> Trasa</H1>
            </CardItem>
            <CardItem cardBody bordered>
              {privateList.length ? null : (
                  <View style={{alignItems: 'center', padding: 50}}>
                    <Text>Looks like you are not enrolled to TRASA.</Text>
                  </View>
              )}
              <Grid
                  style={styles.list}
                  renderItem={_renderItem}
                  renderPlaceholder={_renderPlaceholder}
                  data={privateList}
                  itemsPerRow={4}
              />
            </CardItem>
          </Card>

          <Card transparent>
            <CardItem />
          </Card>



          <Card padder>
            <CardItem header bordered>
              <H1 style={{fontFamily:"Rajdhani-Bold"}}> Personal Services</H1>
            </CardItem>
            <CardItem cardBody bordered>
              {publicList.length ? null : (
                <Text style={{padding: 50}}>
                  You can use this app for 2FA in your personal accounts
                </Text>
              )}
              <Grid
                style={styles.list}
                renderItem={_renderItem}
                renderPlaceholder={_renderPlaceholder}
                data={publicList}
                itemsPerRow={4}
              />
            </CardItem>
          </Card>


              <Card transparent>
                <CardItem>
                  <Body>
                    <Text>Press the + icon below to enroll/add 2FA</Text>
                  </Body>
                </CardItem>
              </Card>



          <Card transparent>
            <CardItem />
            <CardItem />
            <CardItem />
            <CardItem />
          </Card>



        </View>
      </ScrollView>
      <Fab
        active={fabActive}
        direction="up"
        style={{backgroundColor: '#0e0343'}}
        position="bottomRight"
        onPress={(parent, target) => {
          setFabActive(!fabActive);
        }}>
        <Icon name="add" />
        <Button
          onPress={() => {
            goToScanPage();
          }}
          style={{backgroundColor: '#0e0343'}}
          //  ref={this.onGetTagetRef}
        >
          <Icon name="qr-code" />
        </Button>
        <Button
          onPress={() => {
            goToCodeInputPage();
          }}
          style={{backgroundColor: '#0e0343'}}>
          <Icon name="code" />
        </Button>
      </Fab>

      <Fab
        active={fabActive}
        direction="up"
        style={{backgroundColor: '#0e0343'}}
        position="bottomLeft"
        onPress={() => {
          props.navigation.navigate('Notifications');
        }}>
        <Icon name="notifications" />
      </Fab>
      {/*<RNTooltips text={"Long Press Description"} visible={this.state.fabActive} target={this.state.tootipTarget} parent={this.state.tooltipParent} />*/}
    </Container>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    flexDirection: 'column',
  },
  modal: {
    justifyContent: 'center',
    alignItems: 'center',
  },
  containerItem: {
    flex: 1,
    flexDirection: 'column',
    //     backgroundColor: '#fcfcfc',
    padding: 10,
    //      borderBottomColor: '#ddd',
    borderBottomWidth: 1,
  },
  itemContent: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
  },
  title: {
    //     color: '#3c80f7',
    fontSize: 40,
  },
  user: {
    flex: 1,
    fontSize: 14,
    //      color: '#666',
    marginTop: 5,
    marginRight: 5,
  },
  time: {
    //        color: '#666',
    fontSize: 14,
    marginTop: 5,
  },
  icon: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
  },
  noData: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    paddingBottom: 100,
    //      backgroundColor: '#eeeeec'
  },
  listView: {
    //      backgroundColor: '#eeeeec'
  },
  item: {
    flex: 1,
    //height: 160,
    margin: 1,
  },
  list: {
    flex: 1,
  },
  navBar: {
    //       backgroundColor:'#1582dc',
  },
});
