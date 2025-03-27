import clx from 'classnames'
import { memo } from 'react'
import { Link, useParams } from 'react-router-dom'

import { ErrorMes } from '../components/ErrorMes'
import { LoadingSpinner } from '../components/LoadingSpinner'
import { MAX_POINTS, TOTAL_FRAME } from '../constants'
import useGame from '../hooks/useGame'
import useGetScore from '../hooks/useGetScore'
import useSubmitScore from '../hooks/useSubmitScore'

export default function GamePlay() {
  const params = useParams()
  const gameId = params.gameId as string
  const gameQuery = useGame(gameId)
  const scoreQuery = useGetScore(gameId)
  const scoreMut = useSubmitScore(gameId)
  const handleScoreSubmit = (score: number) => {
    scoreMut.mutate({ score })
  }

  const loading = gameQuery.isLoading || scoreQuery.isLoading

  if (loading) {
    return <LoadingSpinner />
  }

  if (!gameQuery.isSuccess || !scoreQuery.isSuccess) {
    if (gameQuery.error && gameQuery.error.status === 404) {
      return <ErrorMes message={gameQuery.error.message} />
    }
    return <h3>Something went wrong</h3>
  }

  const game = gameQuery.data
  const { currentRoll, currentUser, currentFrame } = game
  const { totals, frameScores, frameScoresRaw } = scoreQuery.data
  const users = game.participants || []
  const isDone = game.status === 'done'
  const curUser = users[currentUser]

  const rolls = frameScoresRaw[curUser.playerId]?.frames[currentFrame] || []
  const remainingPins = gameQuery.isFetching ? null : calculateRemainingPins(rolls, currentFrame)

  const winnerId = isDone ? Object.entries(totals).sort(([, scoreA], [, scoreB]) => scoreB - scoreA)[0]?.[0] : null
  return (
    <div className="bg-background relative flex min-h-svh flex-col">
      <div className="bg-background">
        <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
          <div className="w-full md:max-w-6xl">
            <div className="flex flex-col gap-6">
              <div className="bg-card text-card-foreground rounded-xl border shadow">
                <div className="flex flex-col space-y-1.5 p-6">
                  <span className="rounded-sm bg-gray-600 px-2.5 py-0.5 text-center text-xl font-semibold tracking-tight text-white">
                    Score board
                  </span>
                </div>
                <div className="mb-6 overflow-auto p-6 pt-0">
                  <>
                    <table className="w-full border-collapse rounded-md border p-2 text-black md:table-fixed">
                      <thead>
                        <tr className="border bg-gray-100">
                          <th className="w-48 border border-gray-400 bg-gray-200 p-2" rowSpan={2}>
                            Player
                          </th>
                          <th className="border border-gray-400 bg-gray-200 p-2" colSpan={10}>
                            Frames
                          </th>
                          <th className="border border-gray-400 bg-gray-200 p-2" rowSpan={2}>
                            Total
                          </th>
                        </tr>
                        <tr className="border bg-gray-100">
                          {Array.from({ length: TOTAL_FRAME }, (_, i) => (
                            <th
                              key={i}
                              className={`border border-gray-400 bg-gray-200 p-3 whitespace-nowrap text-black ${i + 1 === currentFrame ? 'bg-green-300' : ''}`}
                            >
                              {i + 1}
                            </th>
                          ))}
                        </tr>
                      </thead>
                      <tbody>
                        {users.map((user, userIndex) => (
                          <tr key={userIndex}>
                            <td
                              className={clx(`truncate border p-2 font-bold whitespace-nowrap`, {
                                'bg-green-300': !isDone && userIndex === currentUser,
                              })}
                            >
                              {user.name}
                            </td>
                            {Array.from({ length: TOTAL_FRAME }, (_, i) => {
                              const frame = i + 1
                              const frameScore = frameScores[user.playerId]?.[i] || []
                              const active = frame === currentFrame && userIndex === currentUser

                              return (
                                <td key={frame} className="border p-3 whitespace-nowrap text-black">
                                  <div className="flex justify-between">
                                    <Roll active={active && currentRoll === 1} score={frameScore[0]} />
                                    <Roll active={active && currentRoll === 2} score={frameScore[1]} />
                                    {frame === 10 && (
                                      <Roll active={active && currentRoll === 3} score={frameScore[2]} />
                                    )}
                                  </div>
                                </td>
                              )
                            })}
                            <td className="border p-3 whitespace-nowrap text-black">{totals[user.playerId]}</td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </>
                </div>
                {!isDone && (
                  <div>
                    <div className="flex flex-col space-y-1.5 p-6 pt-0">
                      <span className="rounded-sm bg-gray-600 px-2.5 py-0.5 text-center text-xl font-semibold tracking-tight text-white">
                        ⬇️ Throw pins below ⬇️
                      </span>

                      <div className="flex justify-between gap-2 space-x-2 p-6 text-center">
                        <Pins remainingPins={remainingPins} handleScoreSubmit={handleScoreSubmit} />
                      </div>
                    </div>
                  </div>
                )}

                {isDone && (
                  <div className="p-6 pt-0">
                    <h4 className="text-2xl font-bold">Winner: {users.find((u) => u.playerId === winnerId)?.name}</h4>
                    <div className="mx-auto w-90">
                      <Link
                        to="/"
                        className="focus-visible:ring-ring bg-primary text-primary-foreground hover:bg-primary/90 inline-flex h-9 w-full cursor-pointer items-center justify-center gap-2 rounded-md px-4 py-2 text-sm font-medium whitespace-nowrap shadow transition-colors focus-visible:ring-1 focus-visible:outline-none disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none"
                        type="submit"
                      >
                        New game
                      </Link>
                    </div>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

function calculateRemainingPins(rolls: number[], currentFrame: number) {
  const maxPins = MAX_POINTS
  const knockedPins = rolls.reduce((sum, roll) => sum + roll, 0)
  let remainingPins = maxPins - knockedPins
  if (currentFrame === TOTAL_FRAME) {
    if (rolls.length * MAX_POINTS === knockedPins) {
      return maxPins
    }
    remainingPins = rolls.length <= 2 ? maxPins * rolls.length - knockedPins : maxPins
  }
  return remainingPins
}

const Pins = memo(function Pins({
  remainingPins,
  handleScoreSubmit,
}: {
  remainingPins: number | null
  handleScoreSubmit(v: number): void
}) {
  return Array.from({ length: MAX_POINTS + 1 }, (_, i) => i).map((pins) => {
    return (
      <button
        key={pins}
        disabled={remainingPins ? pins > remainingPins : false}
        onClick={() => handleScoreSubmit(pins)}
        className="focus-visible:ring-ring hover:text-accent-foreground flex w-[4rem] cursor-pointer flex-col items-center gap-2 rounded-md border border-gray-200 p-1 text-sm font-normal whitespace-nowrap hover:bg-green-300 focus-visible:ring-1 focus-visible:outline-none disabled:pointer-events-none disabled:bg-gray-400 disabled:opacity-50 md:h-[5rem]"
      >
        <img src="/pin.webp" className="w-[2.5rem]" />
        <span>{pins}</span>
      </button>
    )
  })
})

const Roll: React.FC<{ active: boolean; score?: string }> = ({ active, score }) => {
  return (
    <span
      className={clx({
        'w-6 rounded-md bg-green-300 text-center': active,
      })}
    >
      {score ?? '﹒'}
    </span>
  )
}
