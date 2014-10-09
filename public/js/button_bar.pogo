key = require 'keymaster'
shapeOf = require './shape_for_cell'

ButtonBar(player) =
  R = React.DOM

  button(type) =
    React.createClass {
        disabled = false

        onClick (e) =
          if (!self.disabled)
            self.props.handleClick(type)

        render() =
          shape = shapeOf(type)

          if (player == 2)
            shape := shape.flipAcrossYeqX()

          points = shape.points().map @(point, i)
            className = if (self.props.freeCellsCount > i) @{ 'point' } else @{ 'point disabled' }
            R.div {className = className, style = {left = (point.0 * 7) - 3, top = (point.1 * 7) - 3}}

          if (self.props.freeCellsCount < points.length)
            self.disabled = true
          else
            self.disabled = false

          cName = "button #(type) player#(player)"
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

          if (player == 2)
            shape := shape.flipAcrossYeqX()

          points = shape.points().map @(point)
            R.div {className = 'point', style = {left = (point.0 * 7) + 15, top = (point.1 * 7) + 15}}

          if (self.props.buttonClicked == type)
            R.div.apply(null, [{className = "pointer #(type)", style = {top = self.state.top, left = self.state.left}}].concat(points))
          else
            null
      }

  CellCounter = React.createClass {
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
        1000

    shapePlacedHandler(e) =
      newCount = self.state.freeCellsCount - e.detail.shapeCellCount
      self.setState {freeCellsCount = newCount}
      self.props.publishFreeCellsCount(newCount)
      self.replenishCellCount()

    render() =
      R.div(
        {className = "cellCounter player#(player)"}
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

  self.reactComponent = React.createClass {
      getInitialState ()=
        { buttonClicked = 'none', freeCellsCount = 5 }

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
        R.div (
          null
          buttonDot {handleClick = self.handleClick, freeCellsCount = self.state.freeCellsCount}
          pointerDot {buttonClicked = self.state.buttonClicked}
          buttonLine {handleClick = self.handleClick, freeCellsCount = self.state.freeCellsCount}
          pointerLine {buttonClicked = self.state.buttonClicked}
          buttonSquare {handleClick = self.handleClick, freeCellsCount = self.state.freeCellsCount}
          pointerSquare {buttonClicked = self.state.buttonClicked}
          buttonGlider {handleClick = self.handleClick, freeCellsCount = self.state.freeCellsCount}
          pointerGlider {buttonClicked = self.state.buttonClicked}
          CellCounter {publishFreeCellsCount = self.publishFreeCellsCount}
        )
    }

  self.render(el) =
    self.el = el
    React.render (self.reactComponent(), el) component

  self.hide() =
    self.el.style.display = 'none'

  self.show() =
    self.el.style.display = 'block'

  self

module.exports = ButtonBar
