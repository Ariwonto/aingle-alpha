import * as React from 'react';
import { inject, observer } from "mobx-react";
import NodeStore from "app/stores/NodeStore";
import Navbar from "react-bootstrap/Navbar";
import Nav from "react-bootstrap/Nav";
import { Dashboard } from "app/components/Dashboard";
import Badge from "react-bootstrap/Badge";
import { RouterStore } from 'mobx-react-router';
import { Explorer } from "app/components/Explorer";
import { NavExplorerSearchbar } from "app/components/NavExplorerSearchbar";
import { Redirect, Route, Switch } from 'react-router-dom';
import { LinkContainer } from 'react-router-bootstrap';
import { ExplorerTransactionQueryResult } from "app/components/ExplorerTransactionQueryResult";
import { ExplorerBundleQueryResult } from "app/components/ExplorerBundleQueryResult";
import { ExplorerAddressQueryResult } from "app/components/ExplorerAddressResult";
import { ExplorerTagQueryResult } from "app/components/ExplorerTagResult";
import { Explorer404 } from "app/components/Explorer404";
import { Misc } from "app/components/Misc";
import { Neighbors } from "app/components/Neighbors";
import { Visualizer } from "app/components/Visualizer";
import { Explorer420 } from "app/components/Explorer420";
import { Helmet } from 'react-helmet'
import * as style from '../../assets/main.css';

interface Props {
    history: any;
    routerStore?: RouterStore;
    nodeStore?: NodeStore;
}

@inject("nodeStore")
@inject("routerStore")
@observer
export class Root extends React.Component<Props, any> {
    renderDevTool() {
        if (process.env.NODE_ENV !== 'production') {
            const DevTools = require('mobx-react-devtools').default;
            return <DevTools />;
        }
    }

    componentDidMount(): void {
        this.props.nodeStore.connect();
    }

    render() {
        return (
            <>
                <Navbar expand="lg" bg="light" variant="light" className={`mb-4 ${style.hornetNavbar}`}>
                    <Navbar.Brand>
                        <img
                            src="/assets/favicon.svg"
                            width="40"
                            className="d-inline-block"
                            alt="Hornet"
                        />
                    </Navbar.Brand>
                    <Navbar.Toggle aria-controls="main-navbar-nav" />
                    <Navbar.Collapse id="main-navbar-nav">
                        <Nav className="mr-auto">
                            <LinkContainer to="/dashboard" >
                                <Nav.Link>Dashboard</Nav.Link>
                            </LinkContainer>
                            <LinkContainer to="/neighbors">
                                <Nav.Link>Neighbors</Nav.Link>
                            </LinkContainer>
                            <LinkContainer to="/explorer">
                                <Nav.Link>Main Explorer</Nav.Link>
                            </LinkContainer>
                            <LinkContainer to="/visualizer">
                                <Nav.Link>Visualizer</Nav.Link>
                            </LinkContainer>
                            <LinkContainer to="/debug">
                                <Nav.Link>Misc</Nav.Link>
                            </LinkContainer>
                        </Nav>
                        <NavExplorerSearchbar />
                        {!this.props.nodeStore.websocketConnected &&
                            <Navbar.Text>
                                <Badge variant="danger">WS not connected!</Badge>
                            </Navbar.Text>
                        }
                        <a href="https://aingle.ai">
                            <svg xmlns="http://www.w3.org/2000/svg"
                                className="navbar-brand"
                                width="32"
                                height="32"
                                viewBox="0 0 16 16"
                                focusable="false">
                                <path fill="currentColor"
                                    fillRule="evenodd"
                                    d="M8 0C3.58 0 0 3.58 0 8s3.58 8 8 8 8-3.58 8-8-3.58-8-8-8zM7 3h2v2H7V3zm3 10H6v-1h1V7H6V6h3v6h1v1z" />
                            </svg>
                        </a>
                    </Navbar.Collapse>
                </Navbar>
                <Helmet defer={false}>
                    <title>{this.props.nodeStore.documentTitle}</title>
                </Helmet>
                <Switch>
                    <Route exact path="/dashboard" component={Dashboard} />
                    <Route exact path="/debug" component={Misc} />
                    <Route exact path="/neighbors" component={Neighbors} />
                    <Route exact path="/explorer/tx/:hash" component={ExplorerTransactionQueryResult} />
                    <Route exact path="/explorer/bundle/:hash" component={ExplorerBundleQueryResult} />
                    <Route exact path="/explorer/addr/:hash" component={ExplorerAddressQueryResult} />
                    <Route exact path="/explorer/tag/:hash" component={ExplorerTagQueryResult} />
                    <Route exact path="/explorer/404/:search" component={Explorer404} />
                    <Route exact path="/explorer/420" component={Explorer420} />
                    <Route exact path="/explorer" component={Explorer} />
                    <Route exact path="/visualizer" component={Visualizer} />
                    <Redirect to="/dashboard" />
                </Switch>
                {this.props.children}
                {this.renderDevTool()}
            </>
        );
    }
}
