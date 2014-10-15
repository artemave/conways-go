require './when'
Grid = require './grid'
Button bar = require './button_bar'

start up ()=
  game id = window.location.pathname.split "/".pop()

  if (! game id)
    return

  d3.select '#new_game'.remove()

  ws = @new Web socket "ws://#(window.location.host)/games/play/#(game id)"

  grid = nil
  button bar = nil

  window.onresize () =
    if (grid)
      grid.resize()

  ws.onmessage (event) =
    msg = JSON.parse(event.data)

    when (msg.Handshake) [
      is 'wait'
        if (grid)
          grid.hide()

        if (button bar)
          button bar.hide()

        d3.select '#waiting_for_players'.transition().style 'opacity' 1.each 'end'
          ws.send(JSON.stringify {"acknowledged" = "wait"})

      is 'ready'
        if (!grid)
          grid := @new Grid(msg.Player, msg.Cols, msg.Rows, msg.WinSpots)

        if (!button bar)
          button bar := @new Button bar(msg.Player)
          button bar.render(document. get element by id "button-bar")

        d3.select '#waiting_for_players'.transition().style 'opacity' 0.each 'end'
          button bar.show()
          grid.show()
          ws.send(JSON.stringify {"acknowledged" = "ready"})

      is 'finish'
        when (msg.Result) [
          is 'won'
            alert "You won"

          is 'lost'
            alert "You lost"

          is 'draw'
            alert "Draw"
        ]
        ws.send(JSON.stringify {"acknowledged" = "finish"})

      otherwise
        if (msg :: Array)
          ack = {"acknowledged" = "game"}
          new cells = grid.new cells to send()

          grid.render next (msg)

          if (new cells.length)
            new cells.for each @(cell)
              cell.State = 1
              cell.Player = grid.player

            ack.cells = new cells

          ws.send(JSON.stringify(ack))
        else
          console.log("Bad ws response:", msg)
    ]

/* start up() */

R = React
RR = ReactRouter

App = React.createClass {
  render() =
    self.props.activeRouteHandler()
}

routes = RR.Routes (
  { location = 'history' }
  RR.Route (
    {
      name = 'app'
      handler = App
      path = ''
    }
    RR.Route {
      name = 'game'
      handler = require './views/game'
      path = '/games/:gameId'
    }
    RR.Route {
      name = 'start_menu'
      handler = require './views/start_menu'
      path = '/'
    }
  )
)

React.renderComponent(routes, document.getElementsByTagName 'main'.0)
