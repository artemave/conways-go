Hover = require './hover'
emitEscape = require './emit_escape'

Grid (player, columns, rows, winSpots) =
  self = this
  self.columns = columns
  self.rows = rows
  self.grid = []
  self.player = player
  self.winSpots = winSpots

  viewport = window.document.get element by id 'viewport'

  svg = d3.select "#viewport".append "svg".
  style "visibility" "hidden"

  viewport width () =
    window.getComputedStyle(viewport).get property value "width".replace "px" ''

  viewport height () =
    window.getComputedStyle(viewport).get property value "height".replace "px" ''

  set viewport height() =
    viewport.set attribute 'style' "height:#(self.y(self.rows))px"

  scale xy () =
    self.x = d3.scale.linear().domain([0, self.columns]).rangeRound([0, viewport width()])
    self.y = d3.scale.linear().domain([0, viewport height() / self.x(1)]).rangeRound([0, viewport height()])

  self.new cells to send() =
    s = svg.selectAll 'rect.new'.data()

  self.add class (class) to (data) =
    svg.selectAll 'rect'.data(data) @(d) @{ "#(d.Row)_#(d.Col)" }.
    classed(class, true).
    exit().classed(class, false)

  self.clearClass(class) =
    svg.selectAll "rect.#(class)".classed(class, false)

  self.any of (cells) classed with any of (classes) =
    !svg.selectAll(["rect.#(c)", where: c <- classes].join()).filter @(d)
      // check if cell from selection is withing cells from args
      _(cells).any @(c)
        c.Row == d.Row && c.Col == d.Col
    .empty()

  self.show () =
    svg.style 'visibility' 'visible'

  self.hide () =
    svg.style 'visibility' 'hidden'

  self.render next (generation) =
    closest to (d) live cell is at least (number of) cells away =
      !_(generation).find @(gc)
        Math.abs(gc.Row - d.Row) < number of && Math.abs(gc.Col - d.Col) < number of && gc.Player == self.player

    maybeFog(d) =
      closest to (d) live cell is at least 5 cells away

    rect = svg.select 'rect' all.data (generation) @(d)
      "#(d.Row)_#(d.Col)"

    rect.
    classed('new', false).
    classed('dead', false).
    classed('live', true).
    classed('player1', @(d) @{ d.Player == 1 }, true).
    classed('player2', @(d) @{ d.Player == 2 }, true)

    rect.exit().
    classed('new', false).
    classed('live', false).
    classed('dead', true).
    classed('fog', maybeFog)

  self.resize () =
    scale xy()

    svg.select 'rect 'all.attr 'width' @{ self.x(0.8) }.
    attr 'height' @{ self.y(0.8) }.
    attr 'x' @(d) @{ self.x(d.Col) + self.x(0.1) }.
    attr 'y' @(d) @{ self.y(d.Row) + self.y(0.1) }

    set viewport height()
    svg.attr("width", viewport width())
    svg.attr("height", viewport height())

  scale xy()

  for (ey = 0, ey < self.rows, ey:=ey+1)
    for (ex = 0, ex < self.columns, ex:=ex+1)
      self.grid.push {
          Row = ey
          Col = ex
        }

  win spot at (cell) =
    p = {Row = cell.Row, Col = cell.Col}
    s = _(self.winSpots).find @(spot)
      _(p).isEqual(spot.Point)

  calculate initial class (d) =
    s = win spot at (d)
    "dead#(if (s) @{" winSpot winSpot#(s.Player)"} else @{''})"

  mark hovered cells new() =
    s = svg.selectAll('rect.hover')
    cellCount = s.size()

    s.classed('new', true).classed('hover', false)

    e = @new CustomEvent "shape-placed" {detail = {shapeCellCount = cellCount}}
    document.dispatchEvent(e)

    emitEscape()

  set viewport height()
  svg.attr("width", viewport width())
  svg.attr("height", viewport height())

  hover = @new Hover(self)

  svg.select 'rect' all.data(self.grid).enter().append 'rect'.
  on 'mousemove' @(d) @{ hover.maybeDrawShape(d) }.
  on 'click' (mark hovered cells new).
  attr 'width' @{ self.x(0.8) }.
  attr 'height' @{ self.y(0.8) }.
  attr 'class' (calculate initial class).
  attr 'rx' @{ self.x(0.2)}.
  attr 'ry' @{ self.y(0.2)}.
  attr 'x' @(d) @{ self.x(d.Col) + self.x(0.1) }.
  attr 'y' @(d) @{ self.y(d.Row) + self.y(0.1) }

  tip = d3.tip().attr 'class' 'd3-tip'.offset([self.y(-0.6),0]).html @(d)
    win spot = win spot at (d)
    if (win spot)
      if (self.player == win spot.Player) @{ "Enemy flag" } else @{ "Your flag" }

  svg.call(tip)
  svg.selectAll 'rect.winSpot'.on 'mousemove' (tip.show).on 'mouseout' (tip.hide)

  self

module.exports = Grid
