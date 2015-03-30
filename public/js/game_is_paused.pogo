React = require 'react'

GameIsPaused = React.createClass {
  render() =
    React.DOM.div(
      { className = 'start_menu' }
      React.DOM.p (null) "The game is paused."
      React.DOM.p (null) "Waiting for another player to resume..."
    )
}

module.exports = GameIsPaused
