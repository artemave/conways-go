Grid = require './grid'
require './when'

start up ()=
  game id = window.location.pathname.split "/".pop()

  if (! game id)
    return

  d3.select '#new_game'.remove()

  ws = @new Web socket "ws://#(window.location.host)/games/play/#(game id)"

  grid = nil

  window.onresize () =
    grid.resize()

  ws.onmessage (event) =
    msg = JSON.parse(event.data)

    when (msg.Handshake) [
      is 'wait'
        if (grid)
          grid.hide()

        d3.select '#waiting_for_players'.transition().style 'opacity' 1.each 'end'
          ws.send(JSON.stringify {"acknowledged" = "wait"})

      is 'ready'
        if (!grid)
          grid := @new Grid(msg.Player, msg.Cols, msg.Rows)

        d3.select '#waiting_for_players'.transition().style 'opacity' 0.each 'end'
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
          grid.render next (msg)

          ack = {"acknowledged" = "game"}

          grid.has selection to send @(selection)
            selection.for each @(cell)
              cell.State = 1
              cell.Player = grid.player

            ack.cells = selection

          ws.send(JSON.stringify(ack))
        else
          console.log("Bad ws response:", msg)
    ]

start up()
