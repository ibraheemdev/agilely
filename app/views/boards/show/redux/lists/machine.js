import { machineTypes } from "@redux/machine";
import * as types from "./types";

console.log(types);
const machine = Object.freeze({
  [types.SET_LISTS]: {
    [machineTypes.IDLE]: machineTypes.SUCCESSFUL,
  },
});

export default machine;
