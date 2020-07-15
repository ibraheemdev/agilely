import * as types from "./types";

const reducer = (state = { role: "guest" }, action) => {
  switch (action.type) {
    case types.SET_METADATA:
      return {...action.payload}
    default:
      return state;
  }
};

export default reducer;
