import {
  BoardState,
  PlayerState,
  ResponseBody,
  TileState,
  checkForWinner,
  cn,
  initialBoardState,
  isDraw,
} from '~/utils'
import { For, createMemo, createSignal } from 'solid-js'

export default function Home() {
  const [boardState, setBoardState] = createSignal(initialBoardState)
  const [winner, setWinner] = createSignal('')
  const gameIsDraw = createMemo(() => isDraw(boardState()))

  const handleClick = async (tileIndex: number) => {
    if (winner()) return
    const newBoardState = new Map(boardState())
    newBoardState.set(tileIndex, TileState.X)
    setBoardState(newBoardState)

    if (checkForWinner(newBoardState)) {
      setWinner(PlayerState.X)
    } else if (!isDraw(newBoardState)) {
      const newBoardState2 = await getAIMove()
      if (checkForWinner(newBoardState2)) {
        setWinner(PlayerState.O)
      }
    }
  }

  const getAIMove = async (): Promise<BoardState> => {
    const body: Record<number, TileState> = {}

    Array.from(boardState()).forEach(([key, value]) => {
      body[key] = value
    })

    const res = await fetch('http://localhost:8080/move', {
      method: 'post',
      body: JSON.stringify(body),
      headers: {
        'Content-Type': 'application/json',
      },
    })

    const { tileIndex }: ResponseBody = await res.json()

    const newBoardState = new Map(boardState())
    newBoardState.set(tileIndex, TileState.O)
    setBoardState(newBoardState)
    return newBoardState
  }

  const restartGame = () => {
    setBoardState(initialBoardState)
    setWinner('')
  }

  return (
    <div>
      <div class='grid gap-3 grid-cols-3 grid-rows-3 w-fit m-auto mt-20'>
        <For each={Array.from(boardState().entries())}>
          {([tileIndex, value]) => (
            <div
              class={cn(
                'w-20 h-20 bg-white flex justify-center items-center',
                value === TileState.Empty && !winner()
                  ? 'cursor-pointer'
                  : value === TileState.X
                  ? 'text-red-500'
                  : 'text-blue-500'
              )}
              onclick={() => handleClick(tileIndex)}
            >
              <div class='text-center text-6xl font-bold'>{value}</div>
            </div>
          )}
        </For>
      </div>
      <div class='text-center mt-2'>
        {winner() && <h2>The winner is {winner()}!</h2>}
        {gameIsDraw() && <h2>Draw</h2>}
        <button onclick={restartGame}>Restart</button>
      </div>
    </div>
  )
}
