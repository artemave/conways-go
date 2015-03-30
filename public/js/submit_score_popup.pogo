React = require 'react'
D     = React.DOM

SubmitScorePopup = React.createClass {
  render() =
    D.div(
      { className = 'popupOverlay' }
      D.div(
        { className = 'submitScorePopup' }
        D.div { className = 'icon-cancel', onClick = self.props.wantsToHide }
        D.div(
          { className = "popupText" }
          D.h3 (null) "You won!"
          D.p (null) "Would you like to submit your score?"
          D.div(
            { className = "popup-button-container" }
            D.div { className = "popup-button" } "Yes"
            D.div { className = "popup-button", onClick = self.props.wantsToHide } "No"
          )
        )
      )
    )
}

module.exports = SubmitScorePopup
