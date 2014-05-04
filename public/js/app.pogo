d3   = require 'd3'
$    = require 'jquery'
Grid = require './grid'

is (expected) (action) =
  action if (actual) matched =
    if (expected == actual)
      action

otherwise (action) =
  action if (actual) matched =
    action

when (actual, cases) =
  for each @(action if matched) in (cases)
    action = action if (actual) matched
    if (action)
      return (action ())

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
    is 'ready'
      $'#waiting_for_players'.fade out
        grid.show()

    is 'wait'
      grid.hide()
      $'#waiting_for_players'.fade in()

    otherwise
      console.log("Bad ws response:", event.data)
  ]