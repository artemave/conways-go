when       = require './when'.when
is         = require './when'.is
otherwise  = require './when'.otherwise
React      = require 'react'
RR         = require 'react-router'

App = React.createClass {
  render() =
    self.props.activeRouteHandler()
}

routes = RR.Routes (
  { location = 'history' }
  RR.Route (
    {
      name = 'app'
      handler = App
      path = ''
    }
    RR.Route {
      name = 'game'
      handler = require './views/game'
      path = '/games/:gameId'
    }
    RR.Route {
      name = 'start_menu'
      handler = require './views/start_menu'
      path = '/'
    }
  )
)

React.renderComponent(routes, document.getElementsByTagName 'main'.0)
