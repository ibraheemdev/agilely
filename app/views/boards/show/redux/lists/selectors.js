const getLists = (state) => {
  return state.lists;
};

const getList = (state, id) => {
  return state.lists[id]
};

export { getLists, getList };
