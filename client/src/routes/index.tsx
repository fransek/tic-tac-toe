import { For, createSignal } from 'solid-js'
import {
  PlayerState,
  TileState,
  checkForWinner,
  cn,
  initialBoardState,
} from '~/utils'

export default function Home() {
  const [boardState, setBoardState] = createSignal(initialBoardState)
  const [isPlayerXTurn, setIsPlayerXTurn] = createSignal(true)
  const [winner, setWinner] = createSignal('')

  const handleClick = (tileIndex: number, player: PlayerState) => {
    if (winner()) return
    const newBoardState = new Map(boardState())
    newBoardState.set(tileIndex, isPlayerXTurn() ? TileState.X : TileState.O)
    setBoardState(newBoardState)
    setIsPlayerXTurn(!isPlayerXTurn())
    if (checkForWinner(boardState())) setWinner(player)
  }

  const restartGame = () => {
    setBoardState(initialBoardState)
    setWinner('')
    setIsPlayerXTurn(true)
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
              onclick={() =>
                handleClick(
                  tileIndex,
                  isPlayerXTurn() ? PlayerState.X : PlayerState.O
                )
              }
            >
              <div class='text-center text-6xl font-bold'>{value}</div>
            </div>
          )}
        </For>
      </div>
      {winner() && (
        <div class='text-center mt-2'>
          <h2>The winner is {winner()}!</h2>
          <button onclick={restartGame}>Restart</button>
        </div>
      )}
    </div>
  )
}
