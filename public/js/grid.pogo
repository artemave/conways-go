React  = require 'react'
D3Grid = require './d3grid'

Grid = React.createClass {
  propTypes = {
    show       = React.PropTypes.bool
    player     = React.PropTypes.number
    cols       = React.PropTypes.number
    rows       = React.PropTypes.number
    winSpots   = React.PropTypes.arrayOf(React.PropTypes.object)
    generation = React.PropTypes.arrayOf(React.PropTypes.object)
  }

  shouldComponentUpdate(nextProps, nextState) =
    if (nextProps.show && self.grid && self.props.generation)

      if (self.props.generation != nextProps.generation)
        self.grid.renderNext(self.props.generation)

      false
    else
      true

  componentDidUpdate() =
    if (self.props.show)
      if (!self.grid)
        self.grid = @new D3Grid(self.getDOMNode(), self.props)
    else
      if (self.grid)
        self.grid.unbindResize()
        self.grid = nil

  render() =
    if (self.props.show)
      React.DOM.div {className = 'D3Grid'}
    else
      null
}

module.exports = Grid
