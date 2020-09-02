import {Body, Button, Container, Content, H1, H3, Header, Icon, Left, Right, Title} from "native-base";
import React from "react";

export default (props)=>{
    return(
        <Container>
            <Header  >
                <Left>
                    <Button transparent>
                        <Icon name={"arrow-back"} onPress={()=>{props.navigation.navigate("Home")}}/>
                    </Button>
                </Left>
                <Body>
                <Title>Help</Title>
                </Body>
                <Right />
            </Header>
            <Content>
                <H1>Help</H1>
            </Content>


        </Container>
    )
}

