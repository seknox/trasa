import {
    Body,
    Button,
    Container,
    Content,
    H1,
    H2,
    H3,
    Header,
    Icon,
    Left,
    Right,
    Thumbnail,
    Title
} from "native-base";
import React from "react"
import {Image, Linking, Text} from 'react-native'

export default (props)=>{
    const url="https://www.seknox.com"
    return(
        <Container>
            <Header  >
                <Left>
                    <Button transparent>
                        <Icon name={"arrow-back"} onPress={()=>{props.navigation.navigate("Home")}}/>
                    </Button>
                </Left>
                <Body>
                <Title>Seknox Cybersecurity</Title>
                </Body>
                <Right />
            </Header>
            <Content>
            <H1>Seknox Cybersecurity</H1>
            <H3>Security Services that Empower modern organization</H3>
                <Image source={require("../../Icons/seknox_cover.png")} />
                <H2>Empowering Cybersecurity</H2>

                <H3 style={{fontColor:'red'}}>seknox</H3>
                <Text style={{fontColor:'green'}}>(c-nox)</Text>
                <Text>started as a vision to empower daily IT operations with cyber security solutions.

                    Technology is enabler of modern business and we at seknox ensure our clients to reach full potential without worrying
                    about security threats. We started our juorney by introducing Identity and Access Management solutions (a frontier for
                    information security), notably adding Trusted Access and Session Analytics service (TRASA), which protects unauthorized
                    access to critical systems and provides total visibility to authentication events. With Trasa, we are disrupting global Identity
                    and Access Management services which is USD $18 Billion market.

                    We also offer security vulenrability management, cryptography (data encryption) as a consultancy service.

                    We are continuosly looking for ways to innovate our services</Text>
            <Button onPress={()=>{
                    Linking.openURL(url);
                    }}>
                <Text>Check out website</Text>
            </Button>

                <Text>© 2018 by Seknox Cybersecurity. All rights reserved.</Text>
            </Content>

        </Container>
    )
}