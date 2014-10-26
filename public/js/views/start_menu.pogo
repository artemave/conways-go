Cookies = require 'cookies-js'
React   = require 'react'
request = require 'superagent'
RR      = require 'react-router'

D = React.DOM

StartMenu = React.createClass {
  mixins = [RR.Navigation]

  onSubmit(e) =
    e.preventDefault()
    gameSize = self.refs.gameSize.getDOMNode().value

    request.post '/games'.type 'form'.send {gameSize = gameSize}.end @(error) @(res)
      gameId = res.text
      Cookies.set("knows-how-to-play", "true")
      self.transitionTo("/games/#(gameId)")

  render() =
    D.form(
      { onSubmit = self.onSubmit }
      D.label { htmlFor = "gameSize" } "Select field size"
      D.select(
        { id = "gameSize", ref = "gameSize" }
        D.option { value = "small" } "Small"
        D.option { value = "medium" } "Medium"
        D.option { value = "large" } "Large"
      )
      D.input { type = "submit", value = "Create" }
    )

}

module.exports = StartMenu
