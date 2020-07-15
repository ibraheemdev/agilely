import React, { useState, useEffect } from "react";
import Sidebar from "./sidebar/sidebar";
import Header from "./header/header";
import NewList from "./list/new_list";
import List from "./list/index";
import { useSelector } from "react-redux";
import { DragDropContext, Droppable } from "react-beautiful-dnd";
import { authenticityToken, url, midString } from "./lib";
import axios from "axios";
import { currentBoardSelectors } from "@redux/current_board";
import { metadataSelectors } from "@redux/metadata";

const Board = () => {
  const [isOpen, toggleSidebar] = useState(false);
  const [newListTitle, setNewListTitle] = useState("");

  const view = useSelector((state) => metadataSelectors.getBoardDirection(state))
  const listIds = useSelector((state) =>
    currentBoardSelectors.getListIds(state)
  );
  const canEdit = useSelector((state) => metadataSelectors.canEdit(state));

  const onDragEnd = (result) => {
    const { destination, source, draggableId, type } = result;

    if (!destination) return;
    if (
      destination.droppableId === source.droppableId &&
      destination.index === source.index
    ) {
      return;
    }

    const getMidstring = (cards, destination) => {
      const above_card = cards[destination - 1];
      const below_card = cards[destination + 1];
      var below_position = Object.is(below_card, undefined)
        ? (below_position = "")
        : (below_position = below_card.position);
      var above_position = Object.is(above_card, undefined)
        ? (above_position = "")
        : (above_position = above_card.position);
      return midString(above_position, below_position);
    };

    if (type === "list") {
      console.log(result);
      const newLists = Array.from(board.lists);
      const targetList = board.lists.find(
        (list) => list._id.$oid === draggableId
      );
      newLists.splice(source.index, 1);
      newLists.splice(destination.index, 0, targetList);
      const midstring = getMidstring(newLists, destination.index);
      newLists[destination.index].position = midstring;
      updateBoard((oldBoard) => {
        const newBoard = { ...oldBoard };
        newBoard.lists = newLists;
        return newBoard;
      });

      axios.patch(`${url}/boards/${board.slug}/lists/${targetList._id.$oid}`, {
        authenticity_token: authenticityToken(),
        list: {
          position: midstring,
        },
      });
    } else {
      const startList = board.lists.find(
        (list) => list._id.$oid === source.droppableId
      );
      const endList = board.lists.find(
        (list) => list._id.$oid === destination.droppableId
      );
      console.log(board.lists.find((l) => l));
      const targetCard = startList.cards.find(
        (card) => card._id.$oid === draggableId
      );

      if (startList === endList) {
        const newCards = Array.from(startList.cards);
        newCards.splice(source.index, 1);
        newCards.splice(destination.index, 0, targetCard);
        const midstring = getMidstring(newCards, destination.index);
        newCards[destination.index].position = midstring;

        const newBoard = { ...board };
        newBoard.lists[
          board.lists.findIndex((l) => l === startList)
        ].cards = newCards;
        updateBoard({ ...newBoard });

        axios.patch(`${url}/cards/${targetCard._id.$oid}`, {
          authenticity_token: authenticityToken(),
          card: {
            position: midstring,
          },
        });
      } else {
        const startCards = Array.from(startList.cards);
        startCards.splice(source.index, 1);
        const endCards = Array.from(endList.cards);
        endCards.splice(destination.index, 0, targetCard);
        const midstring = getMidstring(endCards, destination.index);
        endCards[destination.index].position = midstring;

        const newBoard = { ...board };
        newBoard.lists[
          board.lists.findIndex((l) => l === startList)
        ].cards = startCards;
        newBoard.lists[
          board.lists.findIndex((l) => l === endList)
        ].cards = endCards;
        updateBoard({ ...newBoard });

        axios.patch(`${url}/cards/${targetCard._id.$oid}`, {
          authenticity_token: authenticityToken(),
          card: {
            position: midstring,
            list_id: endList._id.$oid,
          },
        });
      }
    }
  };

  return (
    <div>
      {true && <Sidebar toggleSidebar={toggleSidebar} isOpen={isOpen} />}
      <div className="h-screen flex body-scrollbar">
        <div className="flex-1 min-w-0 flex flex-col bg-white">
          <Header toggleSidebar={toggleSidebar} />
          <div className="flex-1 overflow-auto">
            <main className={`${view === "BOARD" && "inline-flex"} p-3 h-full`}>
              <DragDropContext onDragEnd={onDragEnd}>
                <Droppable
                  droppableId="board"
                  direction={view}
                  type="list"
                >
                  {(provided) => (
                    <div
                      className={`${view === "vertical" && "inline-flex"}`}
                      {...provided.droppableProps}
                      ref={provided.innerRef}
                    >
                      {listIds.map((id, index) => (
                        <List id={id} key={id} index={index} />
                      ))}
                      {provided.placeholder}
                    </div>
                  )}
                </Droppable>
              </DragDropContext>
              {/* {props.current_user && (
                <NewList
                  title={newListTitle}
                  setTitle={setNewListTitle}
                  handleSubmit={() => {
                    dispatch(
                      listActions.addList(props.board.slug, newListTitle)
                    );
                    setNewListTitle("");
                    document.getElementById("newListTitle").focus();
                  }}
                />
              )} */}
            </main>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Board;
