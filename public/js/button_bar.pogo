key = require 'keymaster'

ButtonBar (el) =
  R = React.DOM

  button (type, shape)=
    React.createClass {
        onClick (e) =
          self.props.handleClick(type)

        render () =
          R.div.apply(null, [{className = "button #(type)", onClick = self.onClick}].concat(shape))
      }

  pointer(type, shape) =
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
          if (self.props.buttonClicked == type)
            R.div.apply(null, [{className = "pointer #(type)", style = {top = self.state.top, left = self.state.left}}].concat(shape))
          else
            null
      }

  line = [
    R.div {className = "point top1left1"}
    R.div {className = "point top1left2"}
    R.div {className = "point top1left3"}
  ]
  buttonLine = button('line', line)
  pointerLine = pointer('line', line)

  square = [
    R.div {className = "point top1left1"}
    R.div {className = "point top1left2"}
    R.div {className = "point top2left1"}
    R.div {className = "point top2left2"}
  ]
  buttonSquare = button('square', square)
  pointerSquare = pointer('square', square)

  glider = [
    R.div {className = "point top1left3"}
    R.div {className = "point top2left3"}
    R.div {className = "point top3left3"}
    R.div {className = "point top3left2"}
    R.div {className = "point top2left1"}
  ]
  buttonGlider = button('glider', glider)
  pointerGlider = pointer('glider', glider)

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
