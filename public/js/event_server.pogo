EventEmitter = require 'eventemitter2'.EventEmitter2

if (!window.eventServer)
  window.eventServer = @new EventEmitter { newListener = false }

module.exports = window.eventServer
