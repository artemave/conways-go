React = require 'react'

ShareInstructions = React.createClass {
  render() =
    React.DOM.div(
      { className = 'start_menu' }
      React.DOM.p (null) "Send the link to this page to your opponent."
      React.DOM.p (null) "Once they join, the game will start."
    )
}

module.exports = ShareInstructions
