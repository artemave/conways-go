React = require 'react'

WaitingForAnotherPlayer = React.createClass {
  render() =
    if (self.props.show)
      if (self.props.showShareInstructions)
        React.DOM.div(
          { className = 'start_menu' }
          React.DOM.p (null) "Send a url to this page to your opponent."
          React.DOM.p (null) "Once they join the game will start."
        )
      else
        React.DOM.div(
          { className = 'start_menu' }
          React.DOM.p (null) "The game is paused."
          React.DOM.p (null) "Waiting for another player to resume..."
        )
    else
      null
}

module.exports = WaitingForAnotherPlayer
