React   = require 'react'
request = require 'superagent'
_ = require 'lodash'

D = React.DOM

Leaderboards = React.createClass {

  getInitialState() =
    { scores = {} }

  componentWillMount() =
    scores = request.get "/scores".end(^)!.body

    if (scores)
      self.setState { scores = scores }
    else
      window.location.href = "#(window.location.protocol)//#(window.location.hostname)/fetch_leaderboards"

  render() =
    boardsDivs = _.map(self.state.scores) @(board)
      playerScore = if (board.playerScore)
        "Your rank: #(board.playerScore.formattedScoreRank), score: #(board.playerScore.formattedScore)"
      else
        "You: ---"

      scoresRows = _(board.items || []).map @(score)
        D.tr(
          null
          D.td(null, score.formattedScoreRank)
          D.td(null, score.formattedScore)
          D.td(null, score.player.displayName)
          D.td(null, D.img { src = score.player.avatarImageUrl, width = 40 })
        )
      .tap @(rows) @{ rows.unshift(null) }.value()

      scoresTable = D.table(
        null
        D.thead(
          null
          D.tr(
            null
            D.th(null, 'Rank')
            D.th(null, 'Score')
            D.th({ colSpan = 2 }, 'Player')
          )
        )
        D.tbody.apply(D, scoresRows)
      )

      D.div(
        { className = "board" }
        D.h2(null, board.name)
        D.div({ className = "playerScore" }, playerScore)
        D.div({ className = "scoreTableContainer" }, scoresTable)
      )

    boardsDivs.unshift { className = "boardsContainer" }

    D.div(
      { className = "boardsPage" }
      D.h1(null, "Boards")
      D.div.apply(D, boardsDivs)
    )
}

module.exports = Leaderboards
