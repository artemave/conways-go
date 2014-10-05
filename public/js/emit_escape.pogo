emitEscape() =
  e = document.createEvent "Events"
  e.initEvent("keydown", true, true)
  e.keyCode = 27
  e.which = 27
  e.keyIdentifier = "U+001B"
  document.dispatchEvent(e)

module.exports = emitEscape
