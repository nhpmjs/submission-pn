import { useMutation } from '@tanstack/react-query'

export interface NewGamePayload {
  players: string[]
}

export default function useNewGame() {
  return useMutation<{ gameId: string }, unknown, NewGamePayload>({
    async mutationFn(data) {
      const r = await fetch(import.meta.env.VITE_APP_API_URL + '/game/new', {
        method: 'POST',
        body: JSON.stringify(data),
        headers: {
          'content-type': 'application/json',
        },
      })

      if (r.ok) {
        const d = await r.json()
        return {
          gameId: d.ID,
        }
      }
      throw new Error(await r.text())
    },
  })
}
