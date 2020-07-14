import * as types from "./types";

const reducer = (state = [], action) => {
  switch (action.type) {
    case types.SET_BOARD_TITLES:
      return action.payload
    default:
      return state;
  }
};

export default reducer;
