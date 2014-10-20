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

  componentDidUpdate() =
    if (self.props.show)
      if (self.grid)
        self.grid.renderNext(self.props.generation)
      else
        if (self.props.player && self.props.cols && self.props.rows && self.props.winSpots)
          self.grid = @new D3Grid(self.getDOMNode(), self.props)
    else
      self.grid = null

  render() =
    if (self.props.show) @{ React.DOM.div {className = 'D3Grid'} } else @{ null }
}

module.exports = Grid
