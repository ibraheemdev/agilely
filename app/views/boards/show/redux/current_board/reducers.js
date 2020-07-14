import * as types from "./types";

const reducer = (state = {}, action) => {
  switch (action.type) {
    case types.SET_CURRENT_BOARD:
      return {
        ...action.payload,
      };
    default:
      return state;
  }
};

export default reducer;
