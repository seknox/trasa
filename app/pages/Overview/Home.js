import React from 'react';
import AsyncStorage from '@react-native-community/async-storage';
import {Image, ScrollView, StyleSheet, TouchableOpacity, View} from 'react-native';
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
              <H1 style={{fontFamily:"Rajdhani-Bold"}}> TRASA Accounts</H1>
            </CardItem>
            <CardItem cardBody bordered>
              {privateList.length ? null : (
                  <View style={{alignItems: 'center', padding: 50}}>
                    <Text>Looks like you are not enrolled to TRASA.</Text>
                  </View>
              )}
              <View style={styles.list} >
                {privateList.map((item,i)=>(
                    <TouchableOpacity key={i} style={styles.item}  onPress={() => {
                        goToOrgPage(item.issuer, item.secret, item.account, item.type);
                    }}>
                        <View style={styles.issuer}>
                            <Image
                                source={require('../Icons/trasa.png')}
                                resizeMethod={'scale'}
                                width={10}
                                height={10}
                                size={10}
                                style={styles.icons}

                                badge={0}
                            />

                            <Text >{item.issuer}</Text>
                        </View>


                    </TouchableOpacity>
                ))}
              </View>
            </CardItem>
          </Card>

          <Card transparent>
            <CardItem />
          </Card>

          <Card padder>
            <CardItem header bordered>
              <H1 style={{fontFamily:"Rajdhani-Bold"}}> Personal Accounts</H1>
            </CardItem>
            <CardItem cardBody bordered>
              {publicList.length ? null : (
                <Text style={{padding: 50}}>
                  You can use this app for 2FA in your personal accounts
                </Text>
              )}
              <View style={styles.list} >
                {publicList.map((item,i)=>(
                    <View key={i} style={styles.item}>

                        {
                            ['angellist','codepen','envelope','etsy','facebook','flickr','foursquare','github-alt','github','gitlab','instagram','linkedin','medium','pinterest','quora','reddit-alien','soundcloud','stack-overflow','steam','stumbleupon','tumblr','twitch','twitter','google','google-plus-official','vimeo','vk','weibo','wordpress','youtube'].
                            includes(item.issuer.toLowerCase())?
                                (
                                    <View style={styles.issuer}>
                                        <SocialIcon
                                            // placeholder={(<SocialIcon type={data.issuer}/>)}
                                            type={item.issuer.toLowerCase()}
                                            light
                                            raised
                                            iconSize={30}
                                            style={styles.icons}
                                            onPress={() => {
                                                goToOrgPage(item.issuer, item.secret, item.account, item.type);
                                            }}
                                            badge={0}
                                        />
                                        <Text >{item.issuer}</Text>
                                    </View>
                                ):
                                (
                                    <View key={i} style={styles.issuer}>
                                        <Avatar
                                            placeholder={getLogoSource('default')}
                                            size={50}
                                            style={styles.icons}
                                            onPress={() => {
                                                goToOrgPage(item.issuer, item.secret, item.account, item.type);
                                            }}
                                            badge={0}
                                        />
                                        <Text>{item.issuer}</Text>
                                    </View>
                                )


                        }





                    </View>
                ))}
              </View>
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

    list: {
      marginVertical: 10,
        // marginHorizontal: 5,
         flex: 1,
        flexDirection:'row',

        flexWrap: 'wrap',

        justifyContent:'flex-start',
        alignContent: 'flex-start',

    },
  item: {
    // flex: 1,
    flexDirection:'column',
      width: 100,
      height: 100,
    minHeight: 90,
      minWidth: 80,
      maxHeight: 100,


      alignItems: 'flex-start',
      justifyContent: 'space-evenly',
      alignContent: 'flex-start',


  },

    icons: {

        // alignItems: 'center',
        // alignSelf: 'flex-start',
        height: 50,
        width: 50,
    },
    issuer: {
      alignItems: 'center'
    },

});


