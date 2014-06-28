d3   = require 'd3'
$    = require 'jquery'
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

  when (msg.handshake) [
    is 'wait'
      grid.hide()
      $'#waiting_for_players'.fade in()
        ws.send(JSON.stringify {"acknowledged" = "wait"})

    is 'ready'
      $'#waiting_for_players'.fade out
        grid.show()
        ws.send(JSON.stringify {"acknowledged" = "ready"})

      /* grid.has selection to send @(selection) */
      /*   ws.send(JSON.stringify(selection)) */
    otherwise
      if (msg :: Array)
        grid.render next (msg)
        ws.send(JSON.stringify {"acknowledged" = "game"})
      else
        console.log("Bad ws response:", msg)
  ]
