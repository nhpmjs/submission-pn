import { useQuery } from '@tanstack/react-query'

export default function useGetScore(gamePlayId: string) {
  return useQuery<Score[], unknown, ScoreBoard>({
    queryKey: ['game', gamePlayId, 'score'],
    async queryFn() {
      const r = await fetch(import.meta.env.VITE_APP_API_URL + '/game/' + gamePlayId + '/score')
      return r.json()
    },
    select(scores) {
      const totals = calculateTotal(scores)
      const frameScores = scoresBoard(scores)
      const frameScoresRaw = framePins(scores)
      return {
        totals,
        frameScores,
        frameScoresRaw,
      }
    },
  })
}

/**
 * Calculate total scores. Returns: `{ [playerId]: total }`.
 */
function calculateTotal(scores: Score[]) {
  const framesByPlayer = scores.reduce(
    (acc, { playerId, score, frame }) => {
      if (!acc[playerId]) {
        acc[playerId] = { frames: {} }
      }
      if (!acc[playerId].frames[frame]) {
        acc[playerId].frames[frame] = []
      }
      acc[playerId].frames[frame].push(parseInt(score, 10))
      return acc
    },
    {} as Record<string, PlayerScore>,
  )

  return Object.entries(framesByPlayer).reduce(
    (results, [playerId, { frames }]) => {
      const frameKeys = Object.keys(frames)
        .map(Number)
        .sort((a, b) => a - b)

      const totalScore = frameKeys.reduce((total, frame, index) => {
        const rolls = frames[frame]
        let frameScore = rolls.reduce((sum, pins) => sum + pins, 0)

        // Strike
        if (rolls.length === 1 && rolls[0] === 10) {
          const nextFrame = frames[frameKeys[index + 1]] || []
          const afterNextFrame = frames[frameKeys[index + 2]] || []

          frameScore +=
            (nextFrame.length === 1 ? nextFrame[0] + (afterNextFrame[0] || 0) : nextFrame[0] + (nextFrame[1] || 0)) || 0
        }
        // Spare
        else if (rolls.length === 2 && frameScore === 10) {
          frameScore += frames[frameKeys[index + 1]]?.[0] || 0
        }

        return total + frameScore
      }, 0)

      results[playerId] = totalScore
      return results
    },
    {} as Record<string, number>,
  )
}

/**
 * Convert raw scores for display (X, /, 1, 2, ...)
 */
function scoresBoard(scores: Score[]) {
  return scores.reduce(
    (all, { playerId, score, frame, roll }) => {
      if (!all[playerId]) {
        all[playerId] = Array.from({ length: 10 }, () => [])
      }

      const numericScore = parseInt(score, 10)
      const playerFrames = all[playerId]
      const frameIndex = frame - 1

      // strike convert
      if (numericScore === 10 && roll === 1) {
        playerFrames[frameIndex].push('X')
      } else {
        playerFrames[frameIndex].push(`${numericScore}`)
      }

      // spare convert
      if (
        playerFrames[frameIndex].length === 2 && // only 2 rolls
        playerFrames[frameIndex][0] !== 'X' && // not a strike
        playerFrames[frameIndex].reduce((a, b) => a + (b === '/' ? 10 : Number(b)), 0) === 10 // total frame score is 10
      ) {
        playerFrames[frameIndex][1] = '/'
      }

      return all
    },
    {} as Record<string, string[][]>,
  )
}

/**
 * Converts raw scores to frames[pins] score. 1-based index.
 */
function framePins(scores: Score[]) {
  return scores.reduce(
    (framesByPlayer, { playerId, score, frame }) => {
      if (!framesByPlayer[playerId]) {
        framesByPlayer[playerId] = { frames: {} }
      }
      if (!framesByPlayer[playerId].frames[frame]) {
        framesByPlayer[playerId].frames[frame] = []
      }
      framesByPlayer[playerId].frames[frame].push(parseInt(score, 10))

      return framesByPlayer
    },
    {} as Record<string, PlayerScore>,
  )
}

interface ScoreBoard {
  totals: Record<string, number>
  frameScores: Record<string, string[][]>
  frameScoresRaw: Record<string, PlayerScore>
}

interface Score {
  id: string
  createdAt: Date
  score: string
  frame: number
  roll: number
  playerId: string
  playername: string
}

interface PlayerScore {
  frames: Record<number, number[]>
}
