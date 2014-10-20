React = require 'react'

WaitingForAnotherPlayer = React.createClass {
  render() =
    if (self.props.show) @{ React.DOM.div(null, "Waiting for another player to join...") } else @{ null }
}

module.exports = WaitingForAnotherPlayer
