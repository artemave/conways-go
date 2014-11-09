React    = require 'react'
Cookies  = require 'cookies-js'
helpText = require './help-text'

D = React.DOM

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
    show                     = React.PropTypes.bool
    wantsToHide              = React.PropTypes.func
    withDontShowThisCheckbox = React.PropTypes.bool
  }

  getDefaultProps() =
    { withDontShowThisCheckbox = false }

  render() =
    if (self.props.show)
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
    else
      null
}

module.exports = HelpPopup
