window.is (expected) (action) =
  action if (actual) matched =
    if (expected == actual)
      action

window.otherwise (action) =
  action if (actual) matched =
    action

window.when (actual, cases) =
  for each @(action if matched) in (cases)
    action = action if (actual) matched
    if (action)
      return (action ())

module.exports = null
