import * as React from 'react';
import Container from "react-bootstrap/Container";
import NodeStore from "app/stores/NodeStore";
import {inject, observer} from "mobx-react";
import ExplorerStore from "app/stores/ExplorerStore";
import Col from "react-bootstrap/Col";
import Row from "react-bootstrap/Row";

interface Props {
    nodeStore?: NodeStore;
    explorerStore?: ExplorerStore;
    match?: {
        params: {
            search: string,
        }
    }
}

@inject("nodeStore")
@inject("explorerStore")
@observer
export class Explorer404 extends React.Component<Props, any> {

    render() {
        let {search} = this.props.match.params;
        return (
            <Container fluid>
                <h3>Tangle Explorer 404</h3>
                <p>
                    The search via '{search}' did not yield any results.
                </p>
                <Row>
                    <Col>
                        <img src={"/assets/crying.png"} width={250}/>
                    </Col>
                </Row>
            </Container>
        );
    }
}