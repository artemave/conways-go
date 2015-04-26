// because phantomjs does not implement CustomEvent yet
window.CustomEvent = window.CustomEvent @or @(name, params)
  params := params || { bubbles = true, cancelable = true, detail = undefined }
  e = document.createEvent 'CustomEvent'
  e.initCustomEvent(name, params.bubbles, params.cancelable, params.detail)
  e
