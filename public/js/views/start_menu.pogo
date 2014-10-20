React   = require 'react'
request = require 'superagent'
RR      = require 'react-router'

R = React.DOM

StartMenu = React.createClass {
  mixins = [RR.Navigation]

  onSubmit(e) =
    e.preventDefault()
    gameSize = self.refs.gameSize.getDOMNode().value

    request.post '/games'.type 'form'.send {gameSize = gameSize}.end @(error) @(res)
      gameId = res.text
      self.transitionTo("/games/#(gameId)")

  render() =
    R.form(
      { onSubmit = self.onSubmit }
      R.label { htmlFor = "gameSize" } "Select field size"
      R.select(
        { id = "gameSize", ref = "gameSize" }
        R.option { value = "small" } "Small"
        R.option { value = "medium" } "Medium"
        R.option { value = "large" } "Large"
      )
      R.input { type = "submit", value = "Create" }
    )

}

module.exports = StartMenu
