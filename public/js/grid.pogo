Grid (svg, window, columns: 150, rows: 100) =
  self = this
  self.svg = svg
  self.columns = columns
  self.rows = rows
  self.selection = []
  self.selection is in progress = false
  self.window = window
  self.grid = []
  self.player = nil

  (cell) is being drawn =
    self.selection is in progress && cell.get 'class' attribute == 'new'

  cantors pairing (a, b) =
    0.5 * (a + b) * (a + b + 1) + b

  scale xy () =
    self.x = d3.scale.linear().domain([0, self.columns - 1]).rangeRound([0, self.window.innerWidth])
    self.y = d3.scale.linear().domain([0, self.window.innerHeight / self.x(1)]).rangeRound([0, self.window.innerHeight])

  add (cell) to selection =
    if (self.selection is in progress && !_(this.class list).contains 'fog')
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
    calculate live class (d) =
      if ((this) is being drawn)
        'new'
      else
        "player#(d.Player) live"

    calculate dead class (d) =
      my closest live cell is at least (number of) cells away =
        !_(generation).find @(gc)
          Math.abs(gc.Row - d.Row) < number of && Math.abs(gc.Col - d.Col) < number of && gc.Player == self.player

      if ((this) is being drawn)
        'new'
      else
        'dead' + if (my closest live cell is at least 5 cells away) @{' fog'} else @{''}

    rect = self.svg.select 'rect' all.data (generation) @(d)
      cantors pairing (d.Row, d.Col)

    rect.attr('class', calculate live class)
    rect.exit().attr('class', calculate dead class)

  self.resize () =
    self.svg.attr("width", self.window.innerWidth)

    scale xy()

    self.svg.select 'rect 'all.attr 'width' @{ self.x(1) }.
    attr 'height' @{ self.y(1) }.
    attr 'x' @(d) @{ self.x(d.Col) }.
    attr 'y' @(d) @{ self.y(d.Row) }

    d3.select '#viewport'.style({'height' = "#(self.y(self.rows))px"})

  scale xy()

  for (ey = 0, ey < self.rows, ey:=ey+1)
    for (ex = 0, ex < self.columns, ex:=ex+1)
      self.grid.push {
          Row = ey
          Col = ex
        }

  self.svg.select 'rect' all.data(self.grid).enter().append 'rect'.
  on 'mousedown' @{ self.selection is in progress = true }.
  on 'mousemove' (add to selection).
  on 'mouseup' @{ self.selection is in progress = false }.
  attr 'width' @{ self.x(1) }.
  attr 'height' @{ self.y(1) }.
  attr 'class' 'dead'.
  attr 'rx' @(d) @{ self.x(0.1)}.
  attr 'ry' @(d) @{ self.y(0.1)}.
  attr 'x' @(d) @{ self.x(d.Col) }.
  attr 'y' @(d) @{ self.y(d.Row) }

  d3.select '#viewport'.style({'height' = "#(self.y(self.rows))px"})

  self

module.exports = Grid
