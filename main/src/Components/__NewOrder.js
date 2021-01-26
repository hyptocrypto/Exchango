import React, { Component } from 'react'
import '../App.css';
import Card from 'react-bootstrap/Card'
import Select from 'react-select'


export default class NewOrder extends Component {
    constructor(props) {
        super(props);
        this.state = {
            selectedOption: '',
            ticker: '',
            amount: '',
            order_type: '',
            trading_pairs: []
        }
    }

    componentDidMount() {
        fetch("/api/all")
            .then(res => res.json())
            .then(
                (result) => {
                    this.setState({
                        trading_pairs: result
                    });
                },
                // Note: it's important to handle errors here
                // instead of a catch() block so that we don't swallow
                // exceptions from actual bugs in components.
                (error) => {
                    this.setState({
                        error
                    });
                }
            )
    }
    handelChange(selectedOption) {
        this.setState({ selectedOption });
    }
    render() {
        let options = this.state.trading_pairs.map(function (pair) {
            return pair.Ticker;
        })
        return (
            <div>

            </div>
        )
    }
}

