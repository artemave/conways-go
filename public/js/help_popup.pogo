React    = require 'react'
key      = require 'keymaster'
Cookies  = require 'cookies-js'
helpText = require './help-text'
_        = require 'lodash'

D = React.DOM

DoNotAutoShowCheckbox = React.createFactory(
  React.createClass {
    onChange() =
      if (self.refs.checkbox.getDOMNode().checked)
        Cookies.set "knows-how-to-play" "true"
      else
        Cookies.expire "knows-how-to-play"

    render() =
      D.input(_.assign(
          { type = 'checkbox', onChange = self.onChange, ref = 'checkbox' }
          self.props
        )
      )
  }
)

HelpPopup = React.createClass {
  propTypes = {
    wantsToHide              = React.PropTypes.func
    withDontShowThisCheckbox = React.PropTypes.bool
  }

  componentWillMount() =
    key('esc', 'help_popup', self.props.wantsToHide)

  componentWillUnmount() =
    key.unbind('esc', 'help_popup')

  getDefaultProps() =
    { withDontShowThisCheckbox = false }

  render() =
    key.setScope 'help_popup'

    D.div(
      { className = 'helpPopup' }
      D.div { className = 'icon-cancel', onClick = self.props.wantsToHide }
      D.div { className = 'popupText', dangerouslySetInnerHTML = { __html = helpText } }
      D.div { className = 'hr' }
      D.div(
        {
          className = 'doNotAutoShowCheckbox'
          style = { display = if (self.props.withDontShowThisCheckbox) @{'block'} else @{'none'} }
        }
        D.label(
          { htmlFor = 'doNotAutoShow' }
          DoNotAutoShowCheckbox { id = 'doNotAutoShow' }
          "Got it. Don't show this again"
        )
      )
      D.div { className = 'icon-play-circled', onClick = self.props.wantsToHide }
    )
}

module.exports = HelpPopup
