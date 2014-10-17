React = require 'react'

R = React.DOM

StartMenu = React.createClass {
  onSubmit(e) =
    e.preventDefault()

  render() =
    R.form(
      { onSubmit = self.onSubmit }
      R.label { htmlFor = "gameSize" } "Select field size"
      R.select(
        { id = "gameSize" }
        R.option { value = "small" } "Small"
        R.option { value = "medium" } "Medium"
        R.option { value = "large" } "Large"
      )
      R.input { type = "submit", value = "Create" }
    )

}

module.exports = StartMenu
