React = require 'react'

Button bar (el) =
  self.el = el
  R = React.DOM

  ButtonLine = React.createClass {
      render () =
        R.div(
          {className = "button line"}
          R.div {className = "point top1left1"}
          R.div {className = "point top1left2"}
          R.div {className = "point top1left3"}
        )
    }

  ButtonSquare = React.createClass {
      render () =
        R.div(
          {className = "button square"}
          R.div {className = "point top1left1"}
          R.div {className = "point top1left2"}
          R.div {className = "point top2left1"}
          R.div {className = "point top2left2"}
        )
    }

  ButtonGlider = React.createClass {
      render () =
        R.div(
          {className = "button glider"}
          R.div {className = "point top1left3"}
          R.div {className = "point top2left3"}
          R.div {className = "point top3left3"}
          R.div {className = "point top4left2"}
          R.div {className = "point top3left1"}
        )
    }

  bb = React.createClass {
      render () = 
        R.div (
          {}
          ButtonLine()
          ButtonSquare()
          ButtonGlider()
        )
    }

  React.render (bb(), el) component

  self.hide() =
    self.el.style.display = 'none'

  self.show() =
    self.el.style.display = 'block'

  self

module.exports = Button bar
