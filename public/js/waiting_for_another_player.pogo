React = require 'react'

WaitingForAnotherPlayer = React.createClass {
  render() =
    if (self.props.show)
      React.DOM.div({ className = 'start_menu' }, "Waiting for another player to join...")
    else
      null
}

module.exports = WaitingForAnotherPlayer
