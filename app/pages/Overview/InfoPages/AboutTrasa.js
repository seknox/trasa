import {Body, Button, Container, Content, H1, H3, Header, Icon, Left, Right, Title} from "native-base";
import React from "react";
import {Linking, Text} from "react-native";

export default (props)=>{
    const url="https://www.seknox.com/trasa/"
    return(
        <Container>
            <Header  >
                <Left>
                    <Button transparent>
                        <Icon name={"arrow-back"} onPress={()=>{props.navigation.navigate("Home")}}/>
                    </Button>
                </Left>
                <Body>
                <Title>TRASA</Title>
                </Body>
                <Right />
            </Header>
<Content>
            <H1>TRASA</H1>
            <H3>Trusted Access and Session Analytics</H3>
    <Text>

    </Text>
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

