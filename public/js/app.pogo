Grid = require './grid'
require './when'

game id = window.location.pathname.split "/".pop()
ws = @new Web socket "ws://#(window.location.host)/games/play/#(game id)"

svg = d3.select "body".append "svg".
attr "width" (window.inner width).
attr "height" (window.inner height).
style "visibility" "hidden"

grid = @new Grid(svg, window)

window.onresize () =
  grid.resize()

ws.onmessage (event) =
  msg = JSON.parse(event.data)

  when (msg.Handshake) [
    is 'wait'
      grid.hide()
      $'#waiting_for_players'.fade in()
        ws.send(JSON.stringify {"acknowledged" = "wait"})

    is 'ready'
      grid.player := msg.Player
      $'#waiting_for_players'.fade out
        grid.show()
        ws.send(JSON.stringify {"acknowledged" = "ready"})

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
