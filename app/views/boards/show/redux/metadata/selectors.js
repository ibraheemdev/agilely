const getRole = (state) => {
  return state.metadata.role;
};

const canEdit = (state) => {
  return state.metadata.role === "admin" || state.metadata.role === "editor"
}

export { getRole, canEdit };
