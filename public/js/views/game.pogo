Cookies                 = require 'cookies-js'
React                   = require 'react'
WebSocket               = require 'ReconnectingWebSocket'
when                    = require '../when'.when
is                      = require '../when'.is
otherwise               = require '../when'.otherwise
WaitingForAnotherPlayer = require '../waiting_for_another_player'
ButtonBar               = require '../button_bar'
Grid                    = require '../grid'
HelpPopup               = require '../help_popup'
key                     = require 'keymaster'

D = React.DOM

Game = React.createClass {

  getInitialState() =
    {
      waitingForAnotherPlayer  = true
      showHelpPopup            = !Cookies.get("knows-how-to-play")
      withDontShowThisCheckbox = true
    }

  onWsMessage(event) =
    msg = JSON.parse(event.data)
    ack = null

    when (msg.Handshake) [
      is 'wait'
        self.setState {waitingForAnotherPlayer = true}
        ack := {"acknowledged" = "wait"}

      is 'ready'
        self.setState(
          player                  = msg.Player
          cols                    = msg.Cols
          rows                    = msg.Rows
          winSpots                = msg.WinSpots
          waitingForAnotherPlayer = false
        )
        ack := {"acknowledged" = "ready"}

      is 'finish'
        when (msg.Result) [
          is 'won'
            alert "You won"

          is 'lost'
            alert "You lost"

          is 'draw'
            alert "Draw"
        ]
        ack := {"acknowledged" = "finish"}

      otherwise
        if (msg :: Array)
          ack := {"acknowledged" = "game"}
          new cells = self.refs.grid.newCellsToSend()

          self.setState(generation = msg)

          if (new cells.length)
            new cells.for each @(cell)
              cell.State = 1
              cell.Player = self.state.player

            ack.cells = new cells
        else
          console.log("Bad ws response:", msg)
    ]

    if (ack)
      if (self.state.showHelpPopup)
        self.deferredAck = ack
      else
        self.ws.send(JSON.stringify(ack))

  helpPopupWantsToHide() =
    if (self.deferredAck)
      self.ws.send(JSON.stringify(self.deferredAck))
      self.deferredAck = null

    self.setState { showHelpPopup = false }

  onHelpButtonClicked() =
    self.setState { showHelpPopup            = true }
    self.setState { withDontShowThisCheckbox = false }

  componentWillMount() =
    self.ws = @new WebSocket "ws://#(window.location.host)/games/play/#(self.props.params.gameId)"
    self.ws.onmessage = self.onWsMessage
    key('esc', self.helpPopupWantsToHide)

  componentWillUnmount() =
    self.ws.close()

  render() =
    D.div(
      null
      HelpPopup {
        show                     = self.state.showHelpPopup
        wantsToHide              = self.helpPopupWantsToHide
        withDontShowThisCheckbox = self.state.withDontShowThisCheckbox
      }
      WaitingForAnotherPlayer { show = self.state.waitingForAnotherPlayer }
      ButtonBar {
        player              = self.state.player
        show                = !self.state.waitingForAnotherPlayer
        onHelpButtonClicked = self.onHelpButtonClicked
      }
      Grid {
        ref        = "grid"
        show       = !self.state.waitingForAnotherPlayer
        generation = self.state.generation
        player     = self.state.player
        cols       = self.state.cols
        rows       = self.state.rows
        winSpots   = self.state.winSpots
      }
    )
}

module.exports = Game