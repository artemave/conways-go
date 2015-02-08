d3    = require 'd3'
Hover = require './hover'
_     = require 'lodash'
require 'd3-tip'

D3Grid (el, opts) =
  self = this
  grid = []
  svg  = null

  el width () =
    window.getComputedStyle(el).get property value "width".replace "px" ''

  el height () =
    window.getComputedStyle(el).get property value "height".replace "px" ''

  setGridHeight() =
    el.set attribute 'style' "height:#(self.y(opts.rows))px"
    svg.attr("width", el width())
    svg.attr("height", el height())

  scale xy () =
    self.x = d3.scale.linear().domain([0, opts.cols]).rangeRound([0, el width()])
    self.y = d3.scale.linear().domain([0, el height() / self.x(1)]).rangeRound([0, el height()])

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

  self.render next (generation) =
    closest to (d) live cell is at least (number of) cells away =
      !_(generation).find @(gc)
        Math.abs(gc.Row - d.Row) < number of && Math.abs(gc.Col - d.Col) < number of && gc.Player == opts.player

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

    if (self.cellUnderCursor)
      self.hover.maybeDrawShape(self.cellUnderCursor)

  resize() =
    scale xy()

    svg.select 'rect 'all.attr 'width' @{ self.x(0.8) }.
    attr 'height' @{ self.y(0.8) }.
    attr 'x' @(d) @{ self.x(d.Col) + self.x(0.1) }.
    attr 'y' @(d) @{ self.y(d.Row) + self.y(0.1) }

    setGridHeight()

  win spot at (cell) =
    p = {Row = cell.Row, Col = cell.Col}
    s = _(opts.winSpots).find @(spot)
      _(p).isEqual(spot.Point)

  calculate initial class (d) =
    s = win spot at (d)
    "dead#(if (s) @{" winSpot winSpot#(s.Player)"} else @{''})"

  mark hovered cells new() =
    s = svg.selectAll('rect.hover')
    s.classed('new', true).classed('hover', false)

    e = @new CustomEvent "shape-placed" {detail = {cells = (s.data())}}
    document.dispatchEvent(e)

  self.unbindResize() =
    window.removeEventListener('resize', resize)

  init() =
    svg := d3.select(el).append "svg"

    scale xy()

    for (ey = 0, ey < opts.rows, ey:=ey+1)
      for (ex = 0, ex < opts.cols, ex:=ex+1)
        grid.push {
          Row = ey
          Col = ex
        }

    setGridHeight()

    self.hover = @new Hover(self, opts.player)

    svg.select 'rect' all.data(grid).enter().append 'rect'.
    on 'mousemove' @(d)
      self.cellUnderCursor = d
      self.hover.maybeDrawShape(d)
    .
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
        if (opts.player == win spot.Player) @{ "Enemy flag" } else @{ "Your flag" }

    svg.call(tip)
    svg.selectAll 'rect.winSpot'.on 'mousemove' (tip.show).on 'mouseout' (tip.hide)

    window.addEventListener('resize', resize)

  init()
  self

module.exports = D3Grid
