import * as types from "./types";
import stateMachine from "./machine";

const reducer = (state = types.IDLE, action) => {
  if (stateMachine[action.type] != undefined) {
    return stateMachine[action.type][state];
  } else {
    return state;
  }
};

export default reducer;
