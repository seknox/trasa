import React from 'react';
import AsyncStorage from '@react-native-community/async-storage';
import {Image, ScrollView, StyleSheet, TouchableOpacity, View} from 'react-native';
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
import {SocialIcon} from 'react-native-elements';
import {iconName} from '../../utils/icons'


export function Home(props) {
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


    const onSearch = (q) => {
        setSearchQuery(q);
        //   setState({totpList:this.totpList.filter(item=>item.issuer.toUpperCase().includes(q.toUpperCase()))})
    };
    const trimIssuerName= (name) =>{
        const maxChars = 11
        if(name?.length>maxChars){
            return name.slice(0,maxChars-1)+"..."
        }else {
            return name
        }
    }




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

                                            <Text >{trimIssuerName("seknox-okta.okta.com")}</Text>
                                            {/*<Text >{trimIssuerName("item.issuer")}</Text>*/}
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
                                        <View key={i} style={styles.issuer}>

                                            <SocialIcon
                                                // placeholder={(<SocialIcon type={data.issuer}/>)}
                                                type={iconName(item.issuer)}
                                                light
                                                raised
                                                iconSize={30}
                                                style={styles.icons}
                                                onPress={() => {
                                                    goToOrgPage(item.issuer, item.secret, item.account, item.type);
                                                }}
                                            />

                                            <Text>{trimIssuerName(item.issuer)}</Text>

                                        </View>


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
            <BottomFabs navigation={props.navigation}/>
            {/*<RNTooltips text={"Long Press Description"} visible={this.state.fabActive} target={this.state.tootipTarget} parent={this.state.tooltipParent} />*/}
        </Container>
    );
}


const BottomFabs = (props)=> {
    const [fabActive, setFabActive] = React.useState(false);
    const goToScanPage = () => {

        //{rand:Date.now()} is little hack to call useEffect in TotpScan page
        props.navigation.navigate('ScanPage',{rand:Date.now()});
    };

    const goToCodeInputPage = () => {
        props.navigation.navigate('CodeInput');
    };


    return(
        <View>
            <Fab
                active={fabActive}
                direction="up"
                style={{backgroundColor: '#0e0343'}}
                position="bottomRight"
                onPress={(parent, target) => {
                    setFabActive(!fabActive);
                }}
            >
                <Icon name="add" />
                <Button
                    onPress={() => {
                        setFabActive(false);
                        goToScanPage();
                    }}
                    style={{backgroundColor: '#0e0343'}}
                    //  ref={this.onGetTagetRef}
                >
                    <Icon name="qr-code" />
                </Button>
                <Button
                    onPress={() => {
                        setFabActive(false);
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
        </View>
    )

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


