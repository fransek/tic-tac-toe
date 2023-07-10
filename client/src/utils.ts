export enum TileState {
  Empty = '',
  X = 'X',
  O = 'O',
}

export type BoardState = Map<number, TileState>

export enum PlayerState {
  X = 'X',
  O = 'O',
}

export const initialBoardState = new Map([
  [0, TileState.Empty],
  [1, TileState.Empty],
  [2, TileState.Empty],
  [3, TileState.Empty],
  [4, TileState.Empty],
  [5, TileState.Empty],
  [6, TileState.Empty],
  [7, TileState.Empty],
  [8, TileState.Empty],
])

export const checkForThreeInARow = (
  boardState: BoardState,
  startingTileIndex: number,
  increment: number
) => {
  const tile1 = boardState.get(startingTileIndex)
  const tile2 = boardState.get(startingTileIndex + increment)
  const tile3 = boardState.get(startingTileIndex + increment * 2)

  return tile1 !== TileState.Empty && tile1 === tile2 && tile2 === tile3
}

export const checkForWinner = (boardState: BoardState) => {
  for (let i = 0; i < 3; i++) {
    if (checkForThreeInARow(boardState, i * 3, 1)) return true
    if (checkForThreeInARow(boardState, i, 3)) return true
  }
  if (checkForThreeInARow(boardState, 0, 4)) return true
  if (checkForThreeInARow(boardState, 2, 2)) return true
  return false
}

export const isDraw = (boardState: BoardState) => {
  for (let i = 0; i < 9; i++) {
    if (boardState.get(i) === TileState.Empty) return false
  }
  return true
}

export const cn = (...classes: string[]) => classes.filter(Boolean).join(' ')

export type RequestBody = {
  0: TileState
  1: TileState
  2: TileState
  3: TileState
  4: TileState
  5: TileState
  6: TileState
  7: TileState
  8: TileState
}

export type ResponseBody = {
  tileIndex: number
}
