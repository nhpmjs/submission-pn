import { Link } from 'react-router-dom'

export const ErrorMes = ({ message }: { message: string }) => {
  return (
    <div className="flex h-screen items-center justify-center">
      <div className="space-y-4 text-center">
        <div role="status">{message}</div>
        <div>
          <Link
            to="/"
            className="focus-visible:ring-ring bg-primary text-primary-foreground hover:bg-primary/90 inline-flex h-9 w-full cursor-pointer items-center justify-center gap-2 rounded-md px-4 py-2 text-sm font-medium whitespace-nowrap shadow transition-colors focus-visible:ring-1 focus-visible:outline-none disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none"
            type="submit"
          >
            Start a game
          </Link>
        </div>
      </div>
    </div>
  )
}
