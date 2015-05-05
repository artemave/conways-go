Cookies   = require 'cookies-js'
React     = require 'react'
request   = require 'superagent'
HelpPopup = React.createFactory(require '../help_popup')
key       = require 'keymaster'
RR        = require 'react-router'

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

    gameId = request.post '/games'.type 'form'.send {gameSize = gameSize}.end(^)!.text
    Cookies.set("knows-how-to-play", "true")
    self.context.router.transitionTo("/games/#(gameId)")

  practiceWall(e) =
    gameId = request.post '/practice'.end(^)!.text
    Cookies.set("knows-how-to-play", "true")
    self.context.router.transitionTo "/games/#(gameId)"

  showHelpPopup() =
    self.setState { showHelpPopup = true }

  showLeaderboards() =
    self.context.router.transitionTo '/leaderboards'

  hideHelpPopup() =
    self.setState { showHelpPopup = false }

  render() =
    D.div(
      { className = 'start_menu' }
      D.div(
        { className = 'start_menu_button', onClick = self.newGame }
        D.span { className = 'button_label' } 'new'
        D.select(
          { ref = "gameSize", onClick = @(e) @{ e.stopPropagation() } }
          D.option { value = "small" } "SMALL"
          D.option { value = "medium" } "MEDIUM"
          D.option { value = "large" } "LARGE"
        )
        D.span { className = 'button_label' } 'game'
      )
      D.div(
        { className = 'start_menu_button', onClick = self.practiceWall }
        D.span { className = 'button_label' } 'practice wall'
      )
      D.div(
        { className = 'start_menu_button', onClick = self.showHelpPopup }
        D.span { className = 'button_label' } 'how to play'
      )
      D.div(
        { className = 'start_menu_button', onClick = self.showLeaderboards }
        D.span { className = 'button_label' } 'leaderboards'
      )
      HelpPopup {
        show        = self.state.showHelpPopup
        wantsToHide = self.hideHelpPopup
      }
    )

}

module.exports = StartMenu
