import * as types from "./types";

const reducer = (state = {}, action) => {
  switch (action.type) {
    case types.SET_LISTS:
      return {
        ...action.payload.reduce(
          (obj, item) => Object.assign(obj, { [item._id]: item }),
          {}
        ),
      };
    case types.ADD_LIST_SUCCESS:
      return {
        ...state,
        lists: [...state.lists, { ...action.payload, cards: [] }],
      };
    case types.DELETE_LIST_SUCCESS:
      const newLists = state.lists.slice();
      const targetIndex = state.lists.findIndex(
        (l) => l._id.$oid === action.payload
      );
      console.log(action.payload);
      newLists.splice(targetIndex, 1);
      return {
        ...state,
        lists: [...newLists],
      };
    default:
      return state;
  }
};

export default reducer;
