import React from 'react';
import './App.css';
import OpenOrderList from './Components/OpenOrderList'
import ClosedOrderList from "./Components/ClosedOrderList"
import NewOrder from './Components/NewOrder'
import TradingPairs from './Components/TradingPairs'

function App() {
  return (
    <>
      <TradingPairs />
      <NewOrder />
      <OpenOrderList />
      <ClosedOrderList />
    </>
  )
}

export default App;
