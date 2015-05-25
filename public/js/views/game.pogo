require '../wrap_websocket'
Cookies           = require 'cookies-js'
React             = require 'react'
WebSocket         = require 'ReconnectingWebSocket'
when              = require '../when'.when
is                = require '../when'.is
otherwise         = require '../when'.otherwise
RR                = require 'react-router'
GameIsPaused      = React.createFactory(require '../game_is_paused')
ShareInstructions = React.createFactory(require '../share_instructions')
ButtonBar         = React.createFactory(require '../button_bar')
Grid              = React.createFactory(require '../grid')
HelpPopup         = React.createFactory(require '../help_popup')
SubmitScorePopup  = React.createFactory(require '../submit_score_popup')

D = React.DOM

Game = React.createClass {
  mixins = [RR.Navigation]

  getInitialState() =
    {
      showShareInstructions    = true
      showGameIsPaused         = false
      showGame                 = false
      showHelpPopup            = false
      showSubmitScore          = false
      withDontShowThisCheckbox = false
      knowsHowToPlay           = Cookies.get("knows-how-to-play")
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

        if (!self.state.knowsHowToPlay)
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
            withDontShowThisCheckbox = !self.state.knowsHowToPlay
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

        when (msg.Result) [
          is 'won'
            self.setState { showSubmitScore = true }

          is 'lost'
            alert "You lost"
            self.context.router.transitionTo "start_menu"

          is 'draw'
            alert "Draw"
            self.context.router.transitionTo "start_menu"
        ]

      is 'game_taken'
        alert "This game has already got enough players :("
        self.context.router.transitionTo "start_menu"

      is 'game_not_found'
        alert "This game does not exist :("
        self.context.router.transitionTo "start_menu"

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

  showSubmitScoreWantsToHide() =
    self.context.router.transitionTo "start_menu"

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
    WS = window.WrapWebSocket()
    self.ws = @new WS "#(self.props.wsHost)/games/play/#(self.context.router.getCurrentParams().gameId)"
    self.ws.onmessage = self.onWsMessage

    window.eventServer.on('shape-placed', self.placeShape)

  componentWillUnmount() =
    self.ws.close(1000)
    window.eventServer.off('shape-placed', self.placeShape)

  render() =
    D.div(
      null
      HelpPopup {
        show                     = self.state.showHelpPopup
        wantsToHide              = self.helpPopupWantsToHide
        withDontShowThisCheckbox = self.state.withDontShowThisCheckbox
      }
      (if (self.state.showShareInstructions) @{ ShareInstructions() } else @{ null })
      (if (self.state.showGameIsPaused) @{ GameIsPaused() } else @{ null })
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
      (if (self.state.showSubmitScore)
        SubmitScorePopup {
          wantsToHide = self.showSubmitScoreWantsToHide
          gameId      = self.context.router.getCurrentParams().gameId
        }
      else @{ null })
    )
}

module.exports = Game
