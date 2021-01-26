import React, { Component } from 'react'
import '../App.css';
import Table from 'react-bootstrap/Table'

class TradingPairs extends Component {
    constructor(props) {
        super(props);
        this.state = {
            error: null,
            isLoaded: false,
            items: []
        };
    }

    componentDidMount() {
        fetch("/api/all")
            .then(res => res.json())
            .then(
                (result) => {
                    this.setState({
                        isLoaded: true,
                        items: result
                    });
                },
                // Note: it's important to handle errors here
                // instead of a catch() block so that we don't swallow
                // exceptions from actual bugs in components.
                (error) => {
                    this.setState({
                        isLoaded: true,
                        error
                    });
                }
            )
    }
    render() {
        const { error, isLoaded, items } = this.state;
        const hStyle = { textAlign: 'center', };
        const tableStyle = { padding: '50px' };
        if (error) {
            return <div>Error: {error.message}</div>;
        } else if (!isLoaded) {
            return <div>Loading...</div>;
        } else {
            return (
                <div style={tableStyle}>
                    <h1 style={hStyle} >Trading Pairs</h1>
                    <Table striped bordered hover>
                        <thead>
                            <tr>
                                <th>Pair</th>
                                <th>Price</th>
                                <th>Daily_Volume</th>
                                <th>Daily_High</th>
                                <th>Daily_Low</th>
                                <th>Daily_Change</th>
                            </tr>
                        </thead>
                        {items.map(item => (
                            <tbody key={item.ID}>
                                <tr>
                                    <th>{item.Ticker}</th>
                                    <th>{item.Price}</th>
                                    <th>{item.Daily_Volume}</th>
                                    <th>{item.Daily_High}</th>
                                    <th>{item.Daily_Low}</th>
                                    <th>{item.Precent_Change}</th>
                                </tr>
                            </tbody>
                        ))}
                    </Table>
                </div>
            );
        }
    }
}
export default TradingPairs
