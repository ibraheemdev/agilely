import * as types from "./types";

const setParticipants = (participants) => ({
  type: types.SET_PARTICIPANTS,
  payload: [...participants],
});

export { setParticipants };
