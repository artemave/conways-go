Grid (player, columns, rows, winSpots) =
  self = this
  self.columns = columns
  self.rows = rows
  self.selection = []
  self.selection is in progress = false
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

  cantors pairing (a, b) =
    0.5 * (a + b) * (a + b + 1) + b

  set viewport height() =
    viewport.set attribute 'style' "height:#(self.y(self.rows))px"

  scale xy () =
    self.x = d3.scale.linear().domain([0, self.columns]).rangeRound([0, viewport width()])
    self.y = d3.scale.linear().domain([0, viewport height() / self.x(1)]).rangeRound([0, viewport height()])

  (cell) is being drawn =
    self.selection is in progress && cell.get 'class' attribute == 'new'

  self.has selection to send (callback) =
    if (!self.selection is in progress && self.selection.length > 0)
      callback(self.selection)
      self.selection = []

  self.show () =
    svg.style 'visibility' 'visible'

  self.hide () =
    svg.style 'visibility' 'hidden'

  self.render next (generation) =
    calculate live class (d) =
      cell = this
      c = if ((cell) is being drawn)
        'new'
      else
        "player#(d.Player) live"

      if (win spot at (d))
        c := c+" winSpot#(d.Player)"

      c

    calculate dead class (d) =
      cell = this

      my closest live cell is at least (number of) cells away =
        !_(generation).find @(gc)
          Math.abs(gc.Row - d.Row) < number of && Math.abs(gc.Col - d.Col) < number of && gc.Player == self.player

      c = if ((cell) is being drawn)
        'new'
      else
        'dead' + if (my closest live cell is at least 5 cells away) @{' fog'} else @{''}

      w = win spot at (d)

      if (w)
        c := c+" winSpot#(w.Player)"

      c

    rect = svg.select 'rect' all.data (generation) @(d)
      cantors pairing (d.Row, d.Col)

    rect.attr('class', calculate live class)
    rect.exit().attr('class', calculate dead class)

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
    "dead#(if (s) @{" winSpot#(s.Player)"} else @{''})"

  set viewport height()
  svg.attr("width", viewport width())
  svg.attr("height", viewport height())

  svg.select 'rect' all.data(self.grid).enter().append 'rect'.
  attr 'width' @{ self.x(0.8) }.
  attr 'height' @{ self.y(0.8) }.
  attr 'class' (calculate initial class).
  attr 'rx' @(d) @{ self.x(0.2)}.
  attr 'ry' @(d) @{ self.y(0.2)}.
  attr 'x' @(d) @{ self.x(d.Col) + self.x(0.1) }.
  attr 'y' @(d) @{ self.y(d.Row) + self.y(0.1) }

  self

module.exports = Grid
