React   = require 'react'
Cookies = require 'cookies-js'

D = React.DOM

DoNotAutoShowCheckbox = React.createClass {
  doNotAutoShow() =
    if (self.refs.checkbox.getDOMNode().checked)
      Cookies.set "knows-how-to-play" "true"
    else
      Cookies.expire "knows-how-to-play"

  render() =
    self.transferPropsTo(
      D.input { type = 'checkbox', onChange = self.doNotAutoShow, ref = 'checkbox' }
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
        D.div { className = 'popupCloseButton', onClick = self.props.onClose }
        D.div { className = 'popupText' } "help text"
        D.label(
          { htmlFor = 'doNotAutoShow' }
          DoNotAutoShowCheckbox { id = 'doNotAutoShow' }
          "Don't show this again"
        )
      )
    else
      null
}

module.exports = HelpPopup
