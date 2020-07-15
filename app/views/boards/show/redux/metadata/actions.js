import * as types from "./types";

const setRole = (metadata) => ({
  type: types.SET_ROLE,
  payload: metadata,
});

export { setRole };
