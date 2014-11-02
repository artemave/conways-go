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

  newCellsToSend() =
    if (self.grid)
      self.grid.newCellsToSend()
    else
      []

  shouldComponentUpdate(nextProps, nextState) =
    if (self.grid)
      if (self.props.generation)
        self.grid.renderNext(self.props.generation)

        if (nextProps.show == self.props.show)
          return (false)
    else
      if (self.props.player && self.props.cols && self.props.rows && self.props.winSpots)
        self.grid = @new D3Grid(self.getDOMNode(), self.props)

    true

  render() =
    display = if (self.props.show) @{ 'block' } else @{ 'none' }
    React.DOM.div {className = 'D3Grid', style = {display = display}}
}

module.exports = Grid
