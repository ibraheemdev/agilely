const getRole = (state) => {
  return state.metadata.role;
};

const getBoardDirection = (state) => {
  return state.metadata.board_direction;
};

const canEdit = (state) => {
  return state.metadata.role === "admin" || state.metadata.role === "editor"
}

export { getRole, canEdit, getBoardDirection };
