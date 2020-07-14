import * as types from "./types";

const setCurrentBoard = (board) => ({ type: types.SET_CURRENT_BOARD, payload: { ...board } });

export { setCurrentBoard };
