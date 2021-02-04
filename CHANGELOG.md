# TO DO
	- Idealy, the components should be set up with only one Websocket connection in the man App.js.
	  This Websocket connection will then pass the broadcast data to each one of the compnents that needs
	  it as a prop. That way we have less overall trafic, less room for error, and dont need to 
	  purge as many idle clients. 

	- Fix the order settleing logic so an order will obsorb another if the ask is below the bid. 
	  Then the ammount purchased will need to be calculated to a reasonalbe float value.

	- Some design so the front end isnt so plain.  




## Changes
