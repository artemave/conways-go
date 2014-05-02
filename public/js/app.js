var d3   = require("d3-browserify")
var $    = require('jquery')
var Grid = require("./grid")

var gameId = window.location.pathname.split("/").pop()
var ws = new WebSocket("ws://" + window.location.host + '/games/play/' + gameId);

ws.onmessage = function(event) {
  var msg = JSON.parse(event.data)

  switch (msg.handshake) {
  case 'ready':
    $('#waiting_for_players').fadeOut(function() {
      var svg = d3.select("body").append("svg")
      .attr("width", window.innerWidth)
      .attr("height", window.innerHeight);

      var grid = new Grid(svg, window);

      window.onresize = function() {
        grid.resize();
      }
    })

    // grid.renderNextGeneration(JSON.parse(event.data.game));
    // grid.hasSelectionToSend(function(selection) {
    //     ws.send(JSON.stringify(selection));
    // })
    break;

  case 'wait':
    console.log("Waiting for players to join...");
    break;
  case 'game_taken':
    break;

  default:
    console.log("Bad ws response:", event.data);
  }
}
