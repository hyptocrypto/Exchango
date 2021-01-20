import React from 'react';
import './App.css';
import Card from 'react-bootstrap/Card'

class NewOrder extends React.Component {
    constructor(props) {
        super(props);
        this.state = { value: '' };

        this.handleTickerChange = this.handleTickerChange.bind(this);
        this.handleAmountChange = this.handleAmountChange.bind(this);
        this.handleOrderTypeChange = this.handleOrderTypeChange.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
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

    handleTickerChange(event) {
        this.setState({ ticker: event.target.value });
    }
    handleAmountChange(event) {
        this.setState({ amount: event.target.value });
    }
    handleOrderTypeChange(event) {
        this.setState({ order_type: event.target.value });
    }

    handleSubmit(event) {
        const postdata = {
            method: 'POST',
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(
                {
                    "Trading_Pair": this.state.Trading_Pair.value,
                    "Ammount": this.state.Amout.value,
                    "Order_Type": this.state.Order_Type.value,
                }

            )
        }
        console.log(postdata);
        event.preventDefault();
    }

    render() {
        const CardStyle = { padding: '50px' }
        const elmStyle = { padding: '20px' }
        const hStyle = { textAlign: 'center', };
        const { items, error, isLoaded } = this.state;
        if (error) {
            return <div>Error: {error.message}</div>;
        } else if (!isLoaded) {
            return <div>Loading...</div>;
        } else {
            return (
                <div style={CardStyle}>
                    <Card>
                        <Card.Body>
                            <h3 style={hStyle}>New Order</h3>
                            <form onSubmit={this.handleSubmit}>
                                <label style={elmStyle}>
                                    Trading Pair:
                            <select name='Trading_Pair' value={this.state.Ticker} onChange={this.handleTickerChange}>
                                        {items.map(item => (
                                            <option key={item.ID} value={item.Ticker}>{item.Ticker}</option>
                                        ))}
                                    </select>
                                </label>
                                <label style={elmStyle}>
                                    Amount:
                                    <input name='Amount' onChange value={this.state.Amout} type='text' pattern='[0-9]*' placeholder='Integer Only' />
                                </label>
                                <label style={elmStyle}>
                                    Order Type:
                                    <select name='Order_Type' value={this.state.Order_Type} onChange={this.handleChange}>
                                        <option vlaue='Buy'>Buy</option>
                                        <option vlaue='Sell'>Sell</option>
                                    </select>
                                </label>
                                <button onClick={this.handleSubmit}> Submit </button>
                            </form>
                        </Card.Body>
                    </Card>
                </div>
            );
        }
    }
}

export default NewOrder;
