import * as types from "./types";

const reducer = (state = {}, action) => {
  switch (action.type) {
    case types.SET_CARDS:
      return {
        ...action.payload.reduce(
          (obj, item) => Object.assign(obj, { [item._id]: item }),
          {}
        ),
      };
    default:
      return state;
  }
};

export default reducer;
