import * as types from "./types";

const setCards = (cards) => ({ type: types.SET_CARDS, payload: [ ...cards ] });

export { setCards };
