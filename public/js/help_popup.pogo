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
    show    = React.PropTypes.bool
    onClose = React.PropTypes.func
  }

  render() =
    if (self.props.show)
      D.div(
        { className = 'helpPopup' }
        D.div { className = 'icon-cancel', onClick = self.props.onClose }
        D.div { className = 'popupText', dangerouslySetInnerHTML = { __html = helpText } }
        D.div { className = 'hr' }
        D.div(
          { className = 'doNotAutoShowCheckbox' }
          D.label(
            { htmlFor = 'doNotAutoShow' }
            DoNotAutoShowCheckbox { id = 'doNotAutoShow' }
            "Got it. Don't show this again"
          )
        )
        D.div { className = 'icon-play-circled', onClick = self.props.onClose }
      )
    else
      null
}

module.exports = HelpPopup
