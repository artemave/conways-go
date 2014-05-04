d3   = require 'd3'
$    = require 'jquery'
Grid = require './grid'
_ = require './when'

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

  _.when (msg.handshake) [
    _.is 'ready'
      $'#waiting_for_players'.fade out
        grid.show()

    _.is 'wait'
      grid.hide()
      $'#waiting_for_players'.fade in()

    _.otherwise
      console.log("Bad ws response:", event.data)
  ]
