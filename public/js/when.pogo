module.exports.is (expected) (action) =
  action if (actual) matched =
    if (expected == actual)
      action

module.exports.otherwise (action) =
  action if (actual) matched =
    action

module.exports.when (actual, cases) =
  for each @(action if matched) in (cases)
    action = action if (actual) matched
    if (action)
      return (action ())
