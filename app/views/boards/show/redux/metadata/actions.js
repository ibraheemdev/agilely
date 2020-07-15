import * as types from "./types";

const setMetadata = (metadata) => ({
  type: types.SET_METADATA,
  payload: { ...metadata },
});

export { setMetadata };
