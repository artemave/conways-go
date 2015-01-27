Cookies           = require 'cookies-js'
React             = require 'react'
WebSocket         = require 'ReconnectingWebSocket'
when              = require '../when'.when
is                = require '../when'.is
otherwise         = require '../when'.otherwise
GameIsPaused      = require '../game_is_paused'
ShareInstructions = require '../share_instructions'
ButtonBar         = require '../button_bar'
Grid              = require '../grid'
HelpPopup         = require '../help_popup'
key               = require 'keymaster'
RR                = require 'react-router'

knowsHowToPlay = Cookies.get("knows-how-to-play")
D              = React.DOM

Game = React.createClass {
  mixins = [RR.Navigation]

  getInitialState() =
    {
      showShareInstructions    = true
      showGameIsPaused         = false
      showGame                 = false
      showHelpPopup            = false
      withDontShowThisCheckbox = false
    }

  onWsMessage(event) =
    msg = JSON.parse(event.data)
    ack = null

    when (msg.Handshake) [
      is 'wait'
        self.setState {
          showShareInstructions = true
          showGameIsPaused      = false
          showGame              = false
        }
        ack := {"acknowledged" = "wait"}

      is 'ready'
        self.setState(
          player                = msg.Player
          cols                  = msg.Cols
          rows                  = msg.Rows
          winSpots              = msg.WinSpots
          freeCellsCount        = msg.FreeCellsCount
          showShareInstructions = false
          showGameIsPaused      = false
          showGame              = true
        )

        if (!knowsHowToPlay)
          self.setState {
            showHelpPopup            = true
            withDontShowThisCheckbox = true
          }
          self.ws.send(JSON.stringify { command = "pause" })

        ack := {"acknowledged" = "ready"}

      is 'pause'
        self.setState {
          showShareInstructions = false
          showGameIsPaused      = true
          showGame              = false
        }

        if (msg.Player == msg.PausedByPlayer)
          self.setState {
            showHelpPopup            = true
            withDontShowThisCheckbox = !knowsHowToPlay
          }

        ack := {"acknowledged" = "pause"}

      is 'resume'
        self.setState {
          player                = msg.Player
          cols                  = msg.Cols
          rows                  = msg.Rows
          winSpots              = msg.WinSpots
          showShareInstructions = false
          showGameIsPaused      = false
          showGame              = true
        }
        ack := {"acknowledged" = "resume"}

      is 'finish'
        self.ws.send(JSON.stringify({"acknowledged" = "finish"}))
        self.ws.close(1000)

        m = when (msg.Result) [
          is 'won'
            "You won"

          is 'lost'
            "You lost"

          is 'draw'
            "Draw"
        ]

        alert(m)
        self.transitionTo "start_menu"

      is 'game_taken'
        alert "This game has already got enough players :("
        self.transitionTo "start_menu"

      is 'game_not_found'
        alert "This game does not exist :("
        self.transitionTo "start_menu"

      is 'game_data'
        ack := {"acknowledged" = "game"}

        self.setState(
          generation     = msg.Generation
          freeCellsCount = msg.FreeCellsCount
        )

      otherwise
        console.log("Bad ws response:", msg)
    ]

    if (ack)
      self.ws.send(JSON.stringify(ack))

  helpPopupWantsToHide() =
    self.ws.send(JSON.stringify { command = "resume" })
    self.setState { showHelpPopup = false }

  onHelpButtonClicked() =
    self.ws.send(JSON.stringify { command = "pause" })
    self.setState {
      showHelpPopup            = true
      withDontShowThisCheckbox = false
    }

  placeShape(e) =
    new cells = e.detail.cells

    new cells.for each @(cell)
      cell.State = 1
      cell.Player = self.state.player

    self.ws.send(JSON.stringify { NewCells = (new cells) })

    self.setState {
      freeCellsCount = self.state.freeCellsCount - new cells.length
    }

  componentWillMount() =
    self.ws = @new WebSocket "ws://#(window.location.host)/games/play/#(self.props.params.gameId)"
    self.ws.onmessage = self.onWsMessage

    key('esc', self.helpPopupWantsToHide)
    document.addEventListener('shape-placed', self.placeShape)

  componentWillUnmount() =
    self.ws.close(1000)
    key.unbind('esc')
    document.removeEventListener('shape-placed', self.placeShape)

  render() =
    D.div(
      null
      HelpPopup {
        show                     = self.state.showHelpPopup
        wantsToHide              = self.helpPopupWantsToHide
        withDontShowThisCheckbox = self.state.withDontShowThisCheckbox
      }
      ShareInstructions {
        show = self.state.showShareInstructions
      }
      GameIsPaused {
        show = self.state.showGameIsPaused
      }
      ButtonBar {
        player              = self.state.player
        show                = self.state.showGame
        freeCellsCount      = self.state.freeCellsCount
        onHelpButtonClicked = self.onHelpButtonClicked
      }
      Grid {
        ref        = "grid"
        show       = self.state.showGame
        generation = self.state.generation
        player     = self.state.player
        cols       = self.state.cols
        rows       = self.state.rows
        winSpots   = self.state.winSpots
      }
    )
}

module.exports = Game
