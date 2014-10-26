React   = require 'react'
Cookies = require 'cookies-js'

D = React.DOM

helpText = "
Objective: capture enemy flag or eliminate all enemy cells.

You do that by placing cells on the battlefield.

Each cell clears a small area of fog around it, thus allowing you to place more cells.

Cells interact with each other according to the rules of Conway's game of life. So be careful how you place them, as, for example, they might suddenly die from over population.

Cbol adds one extra rule to those of game of life: a cell belongs to a player. When cells of different players collide, they produce neutral cells. When a cell of a player meets a neutral cell, the result (if any) is that player’s cells.

Cells are placed in shapes. Of which there are four: point, line, square and glider.

Each shape costs the amount of cells it consists from (e.g. the price to place a square is four cells). This is taken out from your pool of cells. That pool is being replenished over time.

Shapes:

- Point. Useless on its own (as it won’t live through to the next generation), it is good to disrupt enemy ranks when it is dropped right in the middle of it. Think of it as a grenade. And a cheap one too: costs only 1 cell.
- Line. Clear the fog with these. Costs 3 cells.
- Square. A better alternative to Line. It does not pulsate and so the edge of the fog does not pulsate making it easier to place next shape. Costs 4 cells.
- Glider. Unlike all others, glider is moving towards the enemy. Costs 5 cells.
"

DoNotAutoShowCheckbox = React.createClass {
  onChange() =
    if (self.refs.checkbox.getDOMNode().checked)
      Cookies.set "knows-how-to-play" "true"
    else
      Cookies.expire "knows-how-to-play"

  render() =
    self.transferPropsTo(
      D.input { type = 'checkbox', onChange = self.onChange, ref = 'checkbox' }
    )
}

HelpPopup = React.createClass {
  propTypes = {
    show    = React.PropTypes.bool
    onClose = React.PropTypes.func
  }

  render() =
    if (self.props.show)
      D.div(
        { className = 'helpPopup' }
        D.div { className = 'icon-cancel', onClick = self.props.onClose }
        D.div { className = 'popupText' } (helpText)
        D.div(
          { className = 'doNotAutoShowCheckbox' }
          D.label(
            { htmlFor = 'doNotAutoShow' }
            DoNotAutoShowCheckbox { id = 'doNotAutoShow' }
            "Don't show this again"
          )
        )
      )
    else
      null
}

module.exports = HelpPopup
