import React from 'react';
import '../App.css';
import Table from 'react-bootstrap/Table'


class WSorders extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            error: null,
            isLoaded: false,
            items: []
        };
    }

    componentDidMount() {
        let socket = new WebSocket("ws://localhost:8000/ws");
        console.log('Attempting to connect to websocket');
        socket.onopen = () => {
            console.log("Client Connected");
            socket.send("Hello from client")
        }
        socket.onclose = (event) => {
            console.log("Socket Closed Connection: ", event);
        }
        socket.onerror = (error) => {
            console.log("Socker Error: ", error);
        }
        socket.onmessage = (msg) => {
            console.log(typeof msg)
            console.log(msg)
            console.log(typeof msg.data)
            console.log(msg.data)
            // var stringdata = String.fromCharCode.apply(String, msg.data);
            // console.log(stringdata)
        }
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



    render() {
        // let socket = new WebSocket("ws://localhost:8000/ws")
        // console.log('Attempting to connect to websocket')
        // socket.onopen = () => {
        //     console.log("Client Connected");
        //     socket.send("Hello from client")
        // }
        // socket.onclose = (event) => {
        //     console.log("Socket Closed Connection: ", event);
        // }
        // socket.onerror = (error) => {
        //     console.log("Socker Error: ", error);
        // }

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
                    <h1 style={hStyle} >Websocket Orders</h1>
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

export default WSorders;
