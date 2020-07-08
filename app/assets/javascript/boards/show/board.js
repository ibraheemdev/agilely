import React, { useState } from "react";
import Sidebar from "./sidebar/sidebar";
import Header from "./header/header";
import NewList from "./list/new_list";
import List from "./list/index";
import { DragDropContext, Droppable } from "react-beautiful-dnd";
import { authenticityToken, url, midString } from "./lib";
import axios from "axios";

const Board = (props) => {
  const [isOpen, toggleSidebar] = useState(false);
  const [board, updateBoard] = useState(props.board);
  const [newListTitle, setNewListTitle] = useState("");
  const [view, toggleView] = useState("BOARD")

  const can_edit = props.role === "admin" || props.role === "editor"

  const handleNewCard = (listId, card) => {
    const newBoard = { ...board };
    newBoard.lists.find((list) => list._id.$oid === listId).cards.push(card);
  };

  const handleNewList = () => {
    event.preventDefault();
    axios
      .post(`${url}/boards/${board.slug}/lists`, {
        authenticity_token: authenticityToken(),
        list: { title: newListTitle },
      })
      .then((res) => {
        let newBoard = { ...board };
        res.data.list.cards = [];
        newBoard.lists.push(res.data.list);
        updateBoard(newBoard);
      })
      .catch((err) => console.log(err))
      .then(() => {
        setNewListTitle("");
        document.getElementById("newListTitle").focus();
      });
  };

  const handleUpdateTitle = (title) => {
    event.preventDefault();
    if (title !== board.title) {
      axios.patch(`${url}/boards/${board.slug}`, {
        authenticity_token: authenticityToken(),
        board: { title: title },
      });
    }
  };

  const handleDeleteList = (targetList) => {
    event.preventDefault();
    axios
      .delete(`${url}/boards/${board.slug}/lists/${targetList._id.$oid}`, {
        data: { authenticity_token: authenticityToken() },
      })
      .then(() => {
        updateBoard((oldBoard) => {
          const newLists = oldBoard.lists.slice();
          newLists.splice(oldBoard.lists.indexOf(targetList), 1);
          return { ...oldBoard, lists: newLists };
        });
      })
      .catch((err) => console.log(err));
  };

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
      console.log(board.lists.find((l) => l ))
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
      {props.current_user && (
        <Sidebar
          toggleSidebar={toggleSidebar}
          isOpen={isOpen}
          boards={props.boards_titles}
        />
      )}
      <div className="h-screen flex body-scrollbar">
        <div className="flex-1 min-w-0 flex flex-col bg-white">
          <Header
            toggleSidebar={toggleSidebar}
            title={board.title}
            handleUpdateTitle={handleUpdateTitle}
            can_edit={can_edit}
            lists={board.lists}
            current_user={props.current_user}
            users={props.users}
            toggleView={toggleView}
            view={view}
          />
          <div className="flex-1 overflow-auto">
            <main className={`${view === "BOARD" && "inline-flex"} p-3 h-full`}>
              <DragDropContext onDragEnd={onDragEnd}>
                <Droppable
                  droppableId="board"
                  direction={view === "BOARD" ? "horizontal" : "vertical"}
                  type="list"
                >
                  {(provided) => (
                    <div
                      className={`${view === "BOARD" && "inline-flex"}`}
                      {...provided.droppableProps}
                      ref={provided.innerRef}
                    >
                      {board.lists.map((list, index) => (
                        <List
                          index={index}
                          list={list}
                          handleDeletion={handleDeleteList}
                          key={list._id.$oid}
                          board_slug={board.slug}
                          handleNewCard={handleNewCard}
                          can_edit={can_edit}
                          view={view}
                        />
                      ))}
                      {provided.placeholder}
                    </div>
                  )}
                </Droppable>
              </DragDropContext>
              {props.current_user && (
                <NewList
                  title={newListTitle}
                  setTitle={setNewListTitle}
                  handleSubmit={handleNewList}
                />
              )}
            </main>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Board;
