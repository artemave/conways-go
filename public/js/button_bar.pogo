React   = require 'react'
key     = require 'keymaster'
shapeOf = require './shape_for_cell'

R = React.DOM

button(type) =
  React.createClass {
    propTypes = {
      player = React.PropTypes.number
    }

    disabled = false

    onClick (e) =
      if (!self.disabled)
        self.props.handleClick(type)

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

      cName = "button #(type) player#(self.props.player)"
      if (self.disabled)
        cName := cName + ' disabled'

      R.div(
        {className = cName, onClick = self.onClick}
        R.div.apply(null, [{className = "null-coordinate"}].concat(points))
      )
  }

pointer(type) =
  React.createClass {
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

CellCounter = React.createClass {
  propTypes = {
    player = React.PropTypes.number
  }

  getInitialState() =
    { freeCellsCount = 5 }

  componentDidMount() =
    document.addEventListener('shape-placed', self.shapePlacedHandler)

  componentWillUnmount() =
    document.removeEventListener('shape-placed', self.shapePlacedHandler)

  replenishCellCount() =
    if (self.state.freeCellsCount < 5)
      setTimeout
        newCount = self.state.freeCellsCount + 1
        self.setState {freeCellsCount = newCount}
        self.props.publishFreeCellsCount(newCount)
        self.replenishCellCount()
      2000

  shapePlacedHandler(e) =
    newCount = self.state.freeCellsCount - e.detail.shapeCellCount
    self.setState {freeCellsCount = newCount}
    self.props.publishFreeCellsCount(newCount)
    self.replenishCellCount()

  render() =
    R.div(
      {className = "cellCounter player#(self.props.player)"}
      "cells left: "
      R.span({className = "counter"}, self.state.freeCellsCount)
    )
}

buttonDot = button('dot')
pointerDot = pointer('dot')

buttonLine = button('line')
pointerLine = pointer('line')

buttonSquare = button('square')
pointerSquare = pointer('square')

buttonGlider = button('glider')
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

  componentWillUnmount()=
    key.unbind('esc', self.cancelPlaceShape)

  cancelPlaceShape ()=
    self.setState {buttonClicked = 'none'}

    e = @new CustomEvent "no-shape-wants-to-be-placed"
    document.dispatchEvent(e)

  handleClick (type) =
    self.cancelPlaceShape()
    self.setState {buttonClicked = type}

    e = @new CustomEvent "about-to-place-shape" {detail = {shape = type}}
    document.dispatchEvent(e)

  publishFreeCellsCount(count) =
    self.setState {freeCellsCount = count}

  render () =
    if (self.props.show)
      R.div (
        null
        buttonDot {handleClick = self.handleClick, freeCellsCount = self.state.freeCellsCount, player = self.props.player}
        pointerDot {buttonClicked = self.state.buttonClicked, player = self.props.player}
        buttonLine {handleClick = self.handleClick, freeCellsCount = self.state.freeCellsCount, player = self.props.player}
        pointerLine {buttonClicked = self.state.buttonClicked, player = self.props.player}
        buttonSquare {handleClick = self.handleClick, freeCellsCount = self.state.freeCellsCount, player = self.props.player}
        pointerSquare {buttonClicked = self.state.buttonClicked, player = self.props.player}
        buttonGlider {handleClick = self.handleClick, freeCellsCount = self.state.freeCellsCount, player = self.props.player}
        pointerGlider {buttonClicked = self.state.buttonClicked, player = self.props.player}
        CellCounter {publishFreeCellsCount = self.publishFreeCellsCount, player = self.props.player}
      )
    else
      null
}

module.exports = ButtonBar
