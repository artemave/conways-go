module.exports = {
  is (expected) (action) =
    action if (actual) matched =
      if (expected == actual)
        action

  otherwise (action) =
    action if (actual) matched =
      action

  when (actual, cases) =
    for each @(action if matched) in (cases)
      action = action if (actual) matched
      if (action)
        return (action ())
}
