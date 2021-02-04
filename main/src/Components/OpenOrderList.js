import React from 'react';
import '../App.css';
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
        fetch("/api/orders/open")
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
    // useEffect(() => {
    //         fuction fetchdatat() {
    //             try {
    //                 const response = await fetch("/api/orders/open");
    //                 if (!response.ok) { throw Error(response.statusText); }
    //                 const json = await response.json();
    //                 this.setState({ items: json });
    //                 console.log(json);
    //             }
    //             catch (error) { console.log(error); }
    //         }
    //     }
    //     )
    // async fetchdata() {
    //     try {
    //         const response = await fetch("/api/orders/open");
    //         if (!response.ok) { throw Error(response.statusText); }
    //         const json = await response.json();
    //         this.setState({ items: json });
    //         console.log(json);
    //     }
    //     catch (error) { console.log(error); }

    // }


    render() {
        const { error, isLoaded } = this.state;
        const hStyle = { textAlign: 'center', };
        const tableStyle = { padding: '50px' };
        if (error) {
            return <div>Error: {error.message}</div>;
        } else if (!isLoaded) {
            return <div>Loading...</div>;
        } else {
            return (
                <div style={tableStyle}>
                    <h1 style={hStyle} >Open Orders</h1>
                    <Table striped bordered hover>
                        <thead>
                            <tr>
                                <th>ID</th>
                                <th>Pair</th>
                                <th>Amount</th>
                                <th>Price</th>
                                <th>Type</th>
                                <th>Partially Settled</th>
                            </tr>
                        </thead>
                        {this.state.items.map(item => (
                            <tbody key={item.ID}>
                                <tr>
                                    <th>{item.ID}</th>
                                    <th>{item.Trading_Pair.Ticker}</th>
                                    <th>{item.Current_Amount}</th>
                                    <th>{item.Price}</th>
                                    <th>{item.Order_Type}</th>
                                    <th>{item.Partial_Settled.toString()}</th>
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
