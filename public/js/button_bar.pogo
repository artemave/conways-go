React   = require 'react'
key     = require 'keymaster'
shapeOf = require './shape_for_cell'

R = React.DOM

button(type, shortcut) =
  b = React.createClass {
    propTypes = {
      player         = React.PropTypes.number
      freeCellsCount = React.PropTypes.number
      handleClick    = React.PropTypes.func
    }

    disabled = false

    onClick (e) =
      if (!self.disabled)
        self.props.handleClick(type)

    componentDidMount() =
      key(shortcut, self.onClick)

    componentWillUnmount() =
      key.unbind(shortcut)

    render() =
      shape = shapeOf(type)

      if (self.props.player == 2)
        shape := shape.flipAcrossYeqX()

      points = shape.points().map @(point, i)
        className = if (self.props.freeCellsCount > i) @{ 'point' } else @{ 'point disabled' }
        R.div {className = className, style = {left = (point.0 * 7) - 3, top = (point.1 * 7) - 3}}

      if (self.props.freeCellsCount < points.length)
        self.disabled = true
      else
        self.disabled = false

      cName = "button shape #(type) player#(self.props.player)"
      if (self.disabled)
        cName := cName + ' disabled'

      R.div(
        {className = cName, onClick = self.onClick}
        R.div.apply(null, [{className = "null-coordinate"}].concat(points))
      )
  }
  React.createFactory(b)

pointer(type) =
  p = React.createClass {
    getInitialState() =
      { top = 0, left = 0 }

    onMouseMove(e) =
      self.setState { top = e.clientY, left = e.clientX }

    componentDidMount() =
      document.addEventListener('mousemove', self.onMouseMove)

    componentWillUnmount() =
      document.removeEventListener('mousemove', self.onMouseMove)

    render ()=
      shape = shapeOf(type)

      if (self.props.player == 2)
        shape := shape.flipAcrossYeqX()

      points = shape.points().map @(point)
        R.div {className = 'point', style = {left = (point.0 * 7) + 15, top = (point.1 * 7) + 15}}

      if (self.props.buttonClicked == type)
        R.div.apply(null, [{className = "pointer #(type)", style = {top = self.state.top, left = self.state.left}}].concat(points))
      else
        null
  }
  React.createFactory(p)

CellCounter = React.createFactory(
  React.createClass {
    propTypes = {
      show           = React.PropTypes.bool
      player         = React.PropTypes.number
      freeCellsCount = React.PropTypes.number
    }

    render() =
      R.div(
        {className = "cellCounter player#(self.props.player)"}
        "cells left: "
        R.span({className = "counter"}, self.props.freeCellsCount)
      )
  }
)

buttonDot = button('dot', 'p')
pointerDot = pointer('dot')

buttonLine = button('line', 'l')
pointerLine = pointer('line')

buttonSquare = button('square', 's')
pointerSquare = pointer('square')

buttonGlider = button('glider', 'g')
pointerGlider = pointer('glider')

ButtonBar = React.createClass {
  propTypes = {
    player         = React.PropTypes.number
    freeCellsCount = React.PropTypes.number
  }

  getInitialState ()=
    { buttonClicked = 'none', show = false }

  componentDidMount () =
    key('esc', self.cancelPlaceShape)
    document.addEventListener('shape-placed', self.cancelPlaceShape)

  componentWillUnmount()=
    key.unbind('esc')
    document.removeEventListener('shape-placed', self.cancelPlaceShape)

  cancelPlaceShape ()=
    self.setState {buttonClicked = 'none'}

    e = @new CustomEvent "no-shape-wants-to-be-placed"
    document.dispatchEvent(e)

  handleClick (type) =
    self.cancelPlaceShape()
    self.setState {buttonClicked = type}

    e = @new CustomEvent "about-to-place-shape" {detail = {shape = type}}
    document.dispatchEvent(e)

  render () =
    if (self.props.show)
      R.div (
        {className = 'ButtonBar'}
        buttonDot {handleClick = self.handleClick, freeCellsCount = self.props.freeCellsCount, player = self.props.player}
        pointerDot {buttonClicked = self.state.buttonClicked, player = self.props.player}
        buttonLine {handleClick = self.handleClick, freeCellsCount = self.props.freeCellsCount, player = self.props.player}
        pointerLine {buttonClicked = self.state.buttonClicked, player = self.props.player}
        buttonSquare {handleClick = self.handleClick, freeCellsCount = self.props.freeCellsCount, player = self.props.player}
        pointerSquare {buttonClicked = self.state.buttonClicked, player = self.props.player}
        buttonGlider {handleClick = self.handleClick, freeCellsCount = self.props.freeCellsCount, player = self.props.player}
        pointerGlider {buttonClicked = self.state.buttonClicked, player = self.props.player}
        CellCounter {
          player                = self.props.player
          freeCellsCount        = self.props.freeCellsCount
        }
        R.div { className = "button player#(self.props.player) icon-help help", onClick = self.props.onHelpButtonClicked }
      )
    else
      null
}

module.exports = ButtonBar
