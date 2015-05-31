React  = require 'react'
D3Grid = require './d3grid'

Grid = React.createClass {
  propTypes = {
    player     = React.PropTypes.number
    cols       = React.PropTypes.number
    rows       = React.PropTypes.number
    winSpots   = React.PropTypes.arrayOf(React.PropTypes.object)
    generation = React.PropTypes.arrayOf(React.PropTypes.object)
  }

  componentDidMount() =
    self.grid = @new D3Grid(self.getDOMNode(), self.props)

  componentWillUnmount() =
    self.grid.unbindResize()

  shouldComponentUpdate(nextProps, nextState) =
    if (nextProps.generation && self.props.generation != nextProps.generation)
      self.grid.renderNext(nextProps.generation)

    false

  render() =
    React.DOM.div {className = 'D3Grid'}
}

module.exports = Grid
