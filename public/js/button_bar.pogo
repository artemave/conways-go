key = require 'keymaster'
shapeOf = require './shape_for_cell'

ButtonBar(player) =
  R = React.DOM

  button (type)=
    React.createClass {
        onClick (e) =
          self.props.handleClick(type)

        render () =
          shape = shapeOf(type)

          if (player == 2)
            shape := shape.flipAcrossYeqX()

          points = shape.points().map @(point)
            R.div {className = 'point', style = {left = (point.0 * 7) - 3, top = (point.1 * 7) - 3}}

          R.div(
            {className = "button #(type) player#(player)", onClick = self.onClick}
            R.div.apply(null, [{className = "null-coordinate"}].concat(points))
          )
      }

  pointer(type) =
    React.createClass {
        getInitialState()=
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

  buttonLine = button('line')
  pointerLine = pointer('line')

  buttonSquare = button('square')
  pointerSquare = pointer('square')

  buttonGlider = button('glider')
  pointerGlider = pointer('glider')

  self.reactComponent = React.createClass {
      getInitialState ()=
        { buttonClicked = 'none' }

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

      render () =
        R.div (
          null
          buttonLine {handleClick = self.handleClick}
          pointerLine {buttonClicked = self.state.buttonClicked}
          buttonSquare {handleClick = self.handleClick}
          pointerSquare {buttonClicked = self.state.buttonClicked}
          buttonGlider {handleClick = self.handleClick}
          pointerGlider {buttonClicked = self.state.buttonClicked}
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
