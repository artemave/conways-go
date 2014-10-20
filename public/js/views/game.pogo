React                   = require 'react'
WebSocket               = require 'ReconnectingWebSocket'
when                    = require '../when'.when
is                      = require '../when'.is
otherwise               = require '../when'.otherwise
WaitingForAnotherPlayer = require '../waiting_for_another_player'
ButtonBar               = require '../button_bar'
Grid                    = require '../grid'

D = React.DOM

Game = React.createClass {

  getInitialState() =
    { waitingForAnotherPlayer = true }

  componentWillMount() =
    self.ws = @new WebSocket "ws://#(window.location.host)/games/play/#(self.props.params.gameId)"
    shmelf = self

    self.ws.onmessage (event) =
      msg = JSON.parse(event.data)

      when (msg.Handshake) [
        is 'wait'
          shmelf.setState {waitingForAnotherPlayer = true}
          shmelf.ws.send(JSON.stringify {"acknowledged" = "wait"})

        is 'ready'
          shmelf.setState(
            player                  = msg.Player
            cols                    = msg.Cols
            rows                    = msg.Rows
            winSpots                = msg.WinSpots
            waitingForAnotherPlayer = false
          )
          shmelf.ws.send(JSON.stringify {"acknowledged" = "ready"})

        is 'finish'
          when (msg.Result) [
            is 'won'
              alert "You won"

            is 'lost'
              alert "You lost"

            is 'draw'
              alert "Draw"
          ]
          shmelf.ws.send(JSON.stringify {"acknowledged" = "finish"})

        otherwise
          if (msg :: Array)
            ack = {"acknowledged" = "game"}
            new cells = shmelf.refs.grid.newCellsToSend()

            shmelf.setState(generation = msg)

            if (new cells.length)
              new cells.for each @(cell)
                cell.State = 1
                cell.Player = shmelf.state.player

              ack.cells = new cells

            shmelf.ws.send(JSON.stringify(ack))
          else
            console.log("Bad ws response:", msg)
      ]

  componentWillUnmount() =
    self.ws.close()

  render() =
    D.div(
      null
      WaitingForAnotherPlayer { show = self.state.waitingForAnotherPlayer }
      ButtonBar { player = self.state.player, show = !self.state.waitingForAnotherPlayer }
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
