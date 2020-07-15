import * as types from "./types";

const reducer = (
  state = { role: "guest", board_direction: "vertical" },
  action
) => {
  switch (action.type) {
    case types.SET_ROLE:
      return { ...state, role: action.payload };
    default:
      return state;
  }
};

export default reducer;
