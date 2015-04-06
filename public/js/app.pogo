when      = require './when'.when
is        = require './when'.is
otherwise = require './when'.otherwise
React     = require 'react'
RR        = require 'react-router'

RouteHandler = RR.RouteHandler
Route        = React.createFactory(RR.Route)
DefaultRoute = React.createFactory(RR.DefaultRoute)

App = React.createClass {
  render() =
    React.createElement(RouteHandler)
}

routes = Route (
  { handler = App, name = "start_menu", path = "/" }
  DefaultRoute { handler = require './views/start_menu' }
  Route {
    name = 'game'
    handler = require './views/game'
    path = '/games/:gameId'
  }
)

RR.run(routes, RR.HistoryLocation) @(Handler)
  React.render(React.createElement(Handler), document.getElementsByTagName 'main'.0)
