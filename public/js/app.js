var d3   = require("d3-browserify")
var $    = require('browserify-zepto')
var Grid = require("./grid")

var ws = new WebSocket("ws://" + window.location.host + window.location.pathname);

ws.onmessage = function(event) {
  switch (event.data.handshake) {
  case 'ready':
    $('#waiting_for_players').fadeOut(function() {
      var svg = d3.select("body").append("svg")
      .attr("width", window.innerWidth)
      .attr("height", window.innerHeight);

      var grid = new Grid(svg, window);
    })

    // grid.renderNextGeneration(JSON.parse(event.data.game));
    // grid.hasSelectionToSend(function(selection) {
    //     ws.send(JSON.stringify(selection));
    // })
    break;

  case 'wait':
    break;
  case 'game_taken':
    break;

  default:
    console.log("Bad ws response:", event.data);
  }
}

window.onresize = function() {
  grid.resize();
}
