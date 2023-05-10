import React from 'react';
import './App.css';
import OpenOrderList from './Components/OpenOrderList'
import ClosedOrderList from "./Components/ClosedOrderList"
import NewOrder from './Components/NewOrder'
import TradingPairs from './Components/TradingPairs'
import WSorders from './Components/WSorders'

function App() {
  return (
    <>
      <TradingPairs />
      <NewOrder />
      <div style={{marginBottom: "200px"}} class="container">
        <div class="row">
          <div class="col">
            <OpenOrderList />
          </div>
          <div class="col">
            <ClosedOrderList />
          </div>
        </div>
      </div>
    </>
  )
}

export default App;
