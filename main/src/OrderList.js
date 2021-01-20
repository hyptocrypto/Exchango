import React from 'react';
import './App.css';
import Table from 'react-bootstrap/Table'


class OrderList extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            error: null,
            isLoaded: false,
            items: []
        };
    }

    componentDidMount() {
        fetch("/api/orders/all")
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
                    <h1 style={hStyle} >Current Orders</h1>
                    <Table striped bordered hover>
                        <thead>
                            <tr>
                                <th>ID</th>
                                <th>Pair</th>
                                <th>Amount</th>
                                <th>Type</th>
                                <th>Active</th>
                            </tr>
                        </thead>
                        {items.map(item => (
                            <tbody>
                                <tr key={item.ID}>
                                    <th>{item.ID}</th>
                                    <th>{item.Trading_Pair.Ticker}</th>
                                    <th>{item.Amount}</th>
                                    <th>{item.Order_Type}</th>
                                    <th>{item.Settled.toString()}</th>
                                </tr>
                            </tbody>
                        ))}
                    </Table>
                </div>
            );
        }
    }
}

export default OrderList;
