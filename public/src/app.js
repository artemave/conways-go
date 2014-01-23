var d3   = require("d3")
var Grid = require("./grid")

var svg = d3.select("body").append("svg")
  .attr("width", window.innerWidth)
  .attr("height", window.innerHeight);

var grid = new Grid(svg, window);

var ws = new WebSocket("ws://" + window.location.host + "/go-ws");

ws.onmessage = function(event) {
  grid.renderNextGeneration(JSON.parse(event.data));
  grid.hasSelectionToSend(function(selection) {
      ws.send(JSON.stringify(selection));
  })
}

window.onresize = function() {
  grid.resize();
}
