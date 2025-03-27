import { useQuery } from '@tanstack/react-query'

export default function useGame(gamePlayId: string) {
  return useQuery<unknown, APIError, GamePlay>({
    queryKey: ['game', gamePlayId],
    async queryFn() {
      const r = await fetch(import.meta.env.VITE_APP_API_URL + '/game/' + gamePlayId)
      if (!r.ok) {
        const e = await r.json()
        const error = new APIError(e.message)

        error.status = r.status
        throw error
      }
      return r.json()
    },
  })
}

export class APIError extends Error {
  status?: number
}

export interface GamePlay {
  id: string
  createdAt: Date
  status?: Omit<string, 'done'> | 'done'
  name: string
  currentFrame: number
  currentRoll: number
  currentUser: number
  participants: Participant[]
}

export interface Participant {
  playerId: string
  name: string
}
