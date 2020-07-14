import * as types from "./types";

const setLists = (lists) => ({ type: types.SET_LISTS, payload: [...lists] });

export { setLists };
