React = require 'react'

GameIsPaused = React.createClass {
  render() =
    if (self.props.show)
      React.DOM.div(
        { className = 'start_menu' }
        React.DOM.p (null) "The game is paused."
        React.DOM.p (null) "Waiting for another player to resume..."
      )
    else
      null
}

module.exports = GameIsPaused
