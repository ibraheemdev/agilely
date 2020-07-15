const getCards = (state) => {
  return state.cards;
};

const getCard = (state, id) => {
  return state.cards[id];
};
export { getCards, getCard };
