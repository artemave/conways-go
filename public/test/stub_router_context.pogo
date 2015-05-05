React = require 'react'
_ = require 'lodash'

StubRouterContext (component, props, stubs) =
  routerStub() = @{}

  _.assign(routerStub, {
    makePath ()                = @{}
    makeHref ()                = @{}
    transitionTo ()            = @{}
    replaceWith ()             = @{}
    goBack ()                  = @{}
    getCurrentPath ()          = @{}
    getCurrentRoutes ()        = @{}
    getCurrentPathname ()      = @{}
    getCurrentParams ()        = @{}
    getCurrentQuery ()         = @{}
    isActive ()                = @{}
    getRouteAtDepth()          = @{}
    setRouteComponentAtDepth() = @{}
  }, stubs)

  React.createClass(
    childContextTypes = {
      router     = React.PropTypes.func
      routeDepth = React.PropTypes.number
    }

    getChildContext() =
      {
        router     = routerStub
        routeDepth = 0
      }

    render() =
      component(props)
  )

module.exports = StubRouterContext
