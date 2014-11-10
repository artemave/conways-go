Cookies   = require 'cookies-js'
React     = require 'react'
request   = require 'superagent'
RR        = require 'react-router'
HelpPopup = require '../help_popup'
key       = require 'keymaster'

D = React.DOM

StartMenu = React.createClass {
  mixins = [RR.Navigation]

  getInitialState() =
    { showHelpPopup = false }

  componentWillMount() =
    key('esc', self.hideHelpPopup)

  componentWillUnmount() =
    key.unbind('esc')

  newGame(e) =
    e.preventDefault()
    gameSize = self.refs.gameSize.getDOMNode().value

    request.post '/games'.type 'form'.send {gameSize = gameSize}.end @(error) @(res)
      gameId = res.text
      Cookies.set("knows-how-to-play", "true")
      self.transitionTo("/games/#(gameId)")

  showHelpPopup() =
    self.setState { showHelpPopup = true }

  hideHelpPopup() =
    self.setState { showHelpPopup = false }

  render() =
    D.div(
      { className = 'start_menu' }
      D.div(
        { className = 'start_menu_button', onClick = self.newGame }
        D.span { className = 'button_label' } 'new'
        D.select(
          { ref = "gameSize", onClick = @(e) @{ e.preventDefault(), false } }
          D.option { value = "small" } "Small"
          D.option { value = "medium" } "Medium"
          D.option { value = "large" } "Large"
        )
        D.span { className = 'button_label' } 'game'
      )
      D.div(
        { className = 'start_menu_button', onClick = self.showHelpPopup }
        D.span { className = 'button_label' } 'how to play'
      )
      HelpPopup {
        show        = self.state.showHelpPopup
        wantsToHide = self.hideHelpPopup
      }
    )

}

module.exports = StartMenu
