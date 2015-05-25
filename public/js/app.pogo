React     = require 'react'
RR        = require 'react-router'
Game = require './views/game'
EventEmitter = require 'eventemitter2'.EventEmitter2

RouteHandler  = RR.RouteHandler
Route         = React.createFactory(RR.Route)
DefaultRoute  = React.createFactory(RR.DefaultRoute)

window.eventServer = @new EventEmitter()

GameWithProps = React.createClass {
  wsHost = if (window.location.protocol == 'http:') @{ 'ws:' } else @{ 'wss:' } + "//#(window.location.host)"

  render() =
    React.createElement(Game, { wsHost = self.wsHost })
}

App = React.createClass {
  render() =
    React.createElement(RouteHandler)
}

routes = Route (
  { handler = App, name = "start_menu", path = "/" }
  DefaultRoute { handler = require './views/start_menu' }
  Route {
    name    = 'game'
    handler = GameWithProps
    path    = '/games/:gameId'
  }
  Route {
    name    = 'leaderboards'
    handler = require './views/leaderboards'
    path    = '/leaderboards'
  }
)

RR.run(routes, RR.HistoryLocation) @(Handler)
  React.render(React.createElement(Handler), document.getElementsByTagName 'main'.0)
