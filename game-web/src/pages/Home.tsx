import { useForm, useFieldArray, Controller } from 'react-hook-form'
import { useNavigate } from 'react-router-dom'
import classNames from 'classnames'

import useNewGame from '../hooks/useNewGame'

interface Player {
  name: string
}

interface NewGameFormData {
  players: Player[]
}
export default function Home() {
  const navigate = useNavigate()
  const newGameMut = useNewGame()

  const { handleSubmit, control, formState } = useForm<NewGameFormData>({
    defaultValues: {
      players: [{ name: '' }],
    },
    mode: 'all',
  })

  const { fields, append, remove } = useFieldArray({
    control,
    name: 'players',
    rules: {
      required: true,
      minLength: 2,
    },
  })

  const onSubmit = (data: NewGameFormData) => {
    newGameMut.mutate(
      {
        players: data.players.map((f) => f.name),
      },
      {
        onSuccess(d) {
          navigate('/game/' + d.gameId)
        },
      },
    )
  }

  return (
    <div className="bg-background relative flex min-h-svh flex-col">
      <div className="themes-wrapper bg-background">
        <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
          <div className="w-full max-w-sm">
            <div className="flex flex-col gap-6">
              <div className="bg-card text-card-foreground rounded-xl border shadow">
                <div className="flex flex-col space-y-1.5 p-6">
                  <div className="text-2xl font-semibold tracking-tight">New Game</div>
                  <div className="text-muted-foreground text-sm">Enter at least 2 players and start the game.</div>
                </div>
                <div className="p-6 pt-0">
                  <form onSubmit={handleSubmit(onSubmit)}>
                    <div className="flex flex-col gap-6">
                      {fields.map((field, index) => {
                        return (
                          <Controller
                            key={field.id}
                            name={`players.${index}.name`}
                            control={control}
                            rules={{
                              required: true,
                              minLength: 1,
                              validate(f) {
                                return f.trim() !== ''
                              },
                            }}
                            render={({ field, fieldState }) => {
                              const isLast = index === fields.length - 1
                              return (
                                <>
                                  <div className="flex gap-2">
                                    <input
                                      {...field}
                                      type="text"
                                      className="border-input placeholder:text-muted-foreground focus-visible:ring-ring flex h-9 w-full rounded-md border bg-transparent px-3 py-1 text-base shadow-sm transition-colors focus-visible:ring-1 focus-visible:outline-none disabled:cursor-not-allowed disabled:opacity-50 md:text-sm"
                                      placeholder="Player name (*)"
                                    />
                                    <button
                                      title="Add"
                                      type="button"
                                      onClick={() => (isLast ? append({ name: '' }) : remove(index))}
                                      disabled={fieldState.invalid}
                                      className={classNames(
                                        '"focus-visible:ring-ring bg-primary text-primary-foreground hover:bg-primary/90 [&_svg]:shrink-0" inline-flex h-9 w-9 cursor-pointer items-center justify-center rounded-md text-sm font-medium whitespace-nowrap shadow transition-colors focus-visible:ring-1 focus-visible:outline-none disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4',
                                        {
                                          'bg-red-400': !isLast,
                                        },
                                      )}
                                    >
                                      {isLast ? <AddIcon /> : <RemoveIcon />}
                                    </button>
                                  </div>
                                </>
                              )
                            }}
                          />
                        )
                      })}

                      <button
                        disabled={!formState.isValid}
                        className="focus-visible:ring-ring bg-primary text-primary-foreground hover:bg-primary/90 inline-flex h-9 w-full cursor-pointer items-center justify-center gap-2 rounded-md px-4 py-2 text-sm font-medium whitespace-nowrap shadow transition-colors focus-visible:ring-1 focus-visible:outline-none disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0"
                        type="submit"
                      >
                        Start game
                      </button>
                    </div>
                  </form>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

const RemoveIcon = () => (
  <svg
    stroke="currentColor"
    fill="currentColor"
    strokeWidth="0"
    viewBox="0 0 512 512"
    height="1em"
    width="1em"
    xmlns="http://www.w3.org/2000/svg"
  >
    <path
      fill="none"
      strokeMiterlimit="10"
      strokeWidth="32"
      d="M448 256c0-106-86-192-192-192S64 150 64 256s86 192 192 192 192-86 192-192z"
    ></path>
    <path
      fill="none"
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth="32"
      d="M320 320 192 192m0 128 128-128"
    ></path>
  </svg>
)

const AddIcon = () => (
  <svg
    stroke="currentColor"
    fill="currentColor"
    strokeWidth="0"
    viewBox="0 0 512 512"
    height="1em"
    width="1em"
    xmlns="http://www.w3.org/2000/svg"
  >
    <path d="M256 48C141.31 48 48 141.31 48 256s93.31 208 208 208 208-93.31 208-208S370.69 48 256 48zm80 224h-64v64a16 16 0 0 1-32 0v-64h-64a16 16 0 0 1 0-32h64v-64a16 16 0 0 1 32 0v64h64a16 16 0 0 1 0 32z"></path>
  </svg>
)
