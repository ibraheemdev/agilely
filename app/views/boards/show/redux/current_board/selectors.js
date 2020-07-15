const getCurrentBoard = (state) => {
  return state.currentBoard;
};

const getListIds = (state) => {
  return state.currentBoard.list_ids
}

export { getCurrentBoard, getListIds };
