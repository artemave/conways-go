d3 = require 'd3'

Grid (svg, window, columns: 150) =
  self = this
  self.svg = svg
  self.columns = columns
  self.selection = []
  self.selection is in progress = false
  self.window = window
  self.grid = []

  (cell) is being drawn =
    self.selection is in progress && cell.get 'class' attribute == 'new'

  calculate live class (d) =
    if ((this) is being drawn)
      'new'
    else
      "player#(d.player id) live"

  calculate dead class (d) =
    if ((this) is being drawn)
      'new'
    else
      'dead'

  cantors pairing (a, b) =
    0.5 * (a + b) * (a + b + 1) + b

  scale xy () =
    self.x = d3.scale.linear().domain([0, self.columns - 1]).rangeRound([0, self.window.innerWidth])
    self.y = d3.scale.linear().domain([0, self.window.innerHeight / self.x(1)]).rangeRound([0, self.window.innerHeight])

  add (cell) to selection =
    if (self.selection is in progress)
      this.set 'class' attribute 'new'
      self.selection.push(cell)

  self.show () =
    self.svg.style 'visibility' 'visible'

  self.hide () =
    self.svg.style 'visibility' 'hidden'

  self.has selection to send (callback) =
    if (!self.selection is in progress && self.selection.length > 0)
      callback(self.selection)
      self.selection = []

  self.render next (generation) =
    rect = self.svg.select 'rect' all.data (generation) @(d)
      cantors pairing (d.Row, d.Col)

  self.resize () =
    self.svg.
    attr("width", self.window.innerWidth).
    attr("height", self.window.innerHeight)

    scale xy()

    self.svg.select 'rect 'all.attr 'width' @{ self.x(1) }.
    attr 'height' @{ self.y(1) }.
    attr 'x' @{ self.x(d.Col) }.
    attr 'y' @{ self.y(d.Row) }

  scale xy()

  for (ey = 0, ey < self.window.innerHeight/self.x(1), ey:=ey+1)
    for (ex = 0, ex < self.columns, ex:=ex+1)
      self.grid.push {
          Row = ey
          Col = ex
        }

  self.svg.select 'rect' all.data(self.grid).enter().append 'rect'.
  on 'mousedown' @{ self.selection_in_progress = true }.
  on 'mousemove' (add to selection).
  on 'mouseup' @{ self.selection_in_progress = false }.
  attr 'width' @{ self.x(1) }.
  attr 'height' @{ self.y(1) }.
  attr 'class' 'dead'.
  attr 'x' @(d) @{ self.x(d.Col) }.
  attr 'y' @(d) @{ self.y(d.Row) }

  self

module.exports = Grid
