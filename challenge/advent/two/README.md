# Day Two  

Of the advent of code 2022 challenges. 

# Context 

The elves are deciding access rights via Rock Paper Scissors. An elf has obtained a premonition of the moves to be thrown by their opponent. They've made a response to each move while trying to avoid suspicion of winning too many games. Help the elf know how well they will do with the strategy in place.

Rock.Beats(Scissors)  // AX ~ CZ
Paper.Beats(Rock)     // BY ~ AX
Scissors.Beats(Paper) // CZ ~ BY

Rock.Worth(1)
Paper.Worth(2)
Scissors.Worth(3)
Win.Worth(6)
Tie.Worth(3)
Lose.Worth(0)

# Goal

Determine how many points the strategy results in.
Consider a Tournament to be a pair of players and a sequence of rounds
A player will have a Score, or even a Round history

The file is a series of Rounds
A Round has an outcome, the decision of winnder and loser, i.e. the score increase for each player.
Tournament.Track(Round.Outcome())

We can improve the strategy by winning more, the best ones to win are the high worth games and the best to lose are the low worth games. 50/50 on middle worth games would ensure that the same number of win/losses nets a overall win. And the floor could be cleared as well with a win on every round. Consider making it easy to compute these answers in addition to the primary goal by choosing a solution with ample coverage.
Remember that tie games net both players a win in points.

We can cheaply track both player's scores rather than one player's score, and we could modify a Round before getting the outcome and track another piece of metadata when we do so.

Round {
  ABC Rps
  XYZ Rps
}

RPS int
const (
  Rock      //
  Paper     //
  Scissors  //
)
