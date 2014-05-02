var d3 = require('d3');

function Grid(svg, window, cells_in_row) {
  var self = this;
  self.svg                   = svg;
  self.cells_in_row          = cells_in_row || 150;
  self.selection             = [];
  self.selection_in_progress = false;
  self.window = window;
  self.grid = [];

  function cell_is_being_drawn(elem) {
    return self.selection_in_progress && elem.getAttribute('class') == 'new'
  };

  function calculateLiveClass(d) {
    if (cell_is_being_drawn(this)) {
      return 'new';
    }
    return 'player' + d.PlayerId + ' live';
  };

  function calculateDeadClass(d) {
    if (cell_is_being_drawn(this)) {
      return 'new';
    }
    return 'dead';
  };

  function cantors_pairing(a, b) {
    return 0.5 * (a + b) * (a + b + 1) + b;
  };

  function scaleXY() {
    self.x = d3.scale.linear().domain([0, self.cells_in_row-1]).rangeRound([0, self.window.innerWidth]);
    self.y = d3.scale.linear().domain([0, self.window.innerHeight / self.x(1)]).rangeRound([0, self.window.innerHeight]);
  }

  function addCellToSelection(data) {
    if (self.selection_in_progress) {
      this.setAttribute('class', 'new');
      self.selection.push(data);
    }
  }

  self.show = function() {
    self.svg.style('visibility', 'visible')
  }

  self.hide = function() {
    self.svg.style('visibility', 'hidden')
  }

  self.hasSelectionToSend = function(callback) {
    if (!self.selection_in_progress && self.selection.length > 0) {
      callback(self.selection)
      self.selection = [];
    }
  };

  self.renderNextGeneration = function(data) {
    var rect = self.svg.selectAll('rect').data(data, function(d) {
        return cantors_pairing(d.Row, d.Col);
    });

    rect.attr('class', calculateLiveClass);
    rect.exit().attr('class', calculateDeadClass);
  };

  self.resize = function() {
    self.svg.attr("width", self.window.innerWidth).attr("height", self.window.innerHeight);

    self.scaleXY();

    self.svg.selectAll('rect')
      .attr('width', function(d) { return self.x(1) })
      .attr('height', function(d) { return self.y(1) })
      .attr('x', function(d) { return self.x(d.Col) })
      .attr('y', function(d) { return self.y(d.Row) });
  }

  scaleXY();

  // XXX Why swapping two fors fucks things up???
  for (var ey = 0; ey < self.window.innerHeight/self.x(1); ey++) {
    for (var ex = 0; ex < self.cells_in_row; ex++) {
      self.grid.push({Row: ey, Col: ex});
    }
  }

  self.svg.selectAll('rect').data(self.grid)
    .enter().append('rect')
    .on('mousedown', function() { self.selection_in_progress = true })
    .on('mousemove', addCellToSelection)
    .on('mouseup', function() { self.selection_in_progress = false })
    .attr('width', function(d) { return self.x(1) })
    .attr('height', function(d) { return self.y(1) })
    .attr('class', 'dead')
    .attr('x', function(d) { return self.x(d.Col) })
    .attr('y', function(d) { return self.y(d.Row) });
}

module.exports = Grid;
