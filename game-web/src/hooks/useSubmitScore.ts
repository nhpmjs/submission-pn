import { useMutation, useQueryClient } from '@tanstack/react-query'

export default function useSubmitScore(gamePlayId: string) {
  const client = useQueryClient()
  return useMutation<unknown, unknown, { score: number }>({
    async mutationFn(data) {
      const r = await fetch(import.meta.env.VITE_APP_API_URL + '/game/' + gamePlayId + '/score', {
        method: 'POST',
        body: JSON.stringify(data),
        headers: {
          'content-type': 'application/json',
        },
      })

      if (!r.ok) {
        throw new Error(await r.text())
      }
    },
    onSuccess() {
      client.invalidateQueries({ queryKey: ['game', gamePlayId] })
    },
  })
}
