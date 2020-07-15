import React, { useLayoutEffect } from "react";
import Board from "./components/board";
import configureStore from "./redux/store";
import { Provider, useDispatch, useSelector } from "react-redux";
import { currentBoardActions } from "@redux/current_board";
import { participantActions } from "@redux/participants";
import { boardTitleActions } from "@redux/board_titles";
import { listActions } from "@redux/lists";
import { cardActions } from "@redux/cards";
import { metadataActions } from "./redux/metadata";

const Index = (props) => {
  const reduxStore = configureStore();
  return (
    <Provider store={reduxStore}>
      <App {...props} />
    </Provider>
  );
};

export default Index;

const App = (props) => {
  const dispatch = useDispatch();
  console.log(props);
  useLayoutEffect(() => {
    dispatch(currentBoardActions.setCurrentBoard(props.board));
    dispatch(listActions.setLists(props.lists));
    dispatch(cardActions.setCards(props.cards));
    dispatch(participantActions.setParticipants(props.participants));
    dispatch(boardTitleActions.setBoardTitles(props.boards_titles));
    dispatch(metadataActions.setRole(props.role));
  }, []);
  return <Board />;
};
