import React, { useState } from "react";
import Sidebar from "./sidebar";
import Header from "./header";
import List from "./list";
import { DragDropContext, Droppable } from "react-beautiful-dnd";
import { authenticityToken, url, midString } from "./lib";
import axios from "axios";

const Board = (props) => {
  const [newListIsOpen, toggleNewList] = useState(false);
  const [isOpen, toggleSidebar] = useState(false);
  const [board, updateBoard] = useState(props.board);
  const [newListTitle, setNewListTitle] = useState("");

  const handleNewCard = (listId, card) => {
    const newBoard = { ...board };
    newBoard.lists.find((list) => list.id === listId).cards.push(card);
  };

  const handleNewList = () => {
    event.preventDefault();
    toggleNewList(false);
    axios
      .post(`${url}/boards/${board.slug}/lists`, {
        authenticity_token: authenticityToken(),
        list: { title: newListTitle },
      })
      .then((res) => {
        let newBoard = { ...board };
        newBoard.lists.push(res.data.list);
        updateBoard(newBoard);
      })
      .catch((err) => console.log(err))
      .then(() => setNewListTitle(""));
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
      .delete(`${url}/lists/${targetList.id}`, {
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

    const getPositions = (cards, destination) => {
      const above_card = cards[destination - 1];
      const below_card = cards[destination + 1];
      var below_position = Object.is(below_card, undefined)
        ? (below_position = "")
        : (below_position = below_card.position);
      var above_position = Object.is(above_card, undefined)
        ? (above_position = "")
        : (above_position = above_card.position);
      return {
        above_position: above_position,
        below_position: below_position,
        midstring: midString(above_position, below_position),
      };
    };

    if (type === "list") {
      console.log(result);
      const newLists = Array.from(board.lists);
      const targetList = board.lists.find(
        (list) => list.id === parseInt(draggableId)
      );
      newLists.splice(source.index, 1);
      newLists.splice(destination.index, 0, targetList);
      const { above_position, below_position, midstring } = getPositions(
        newLists,
        destination.index
      );
      newLists[destination.index].position = midstring;

      updateBoard((oldBoard) => {
        const newBoard = { ...oldBoard };
        newBoard.lists = newLists;
        return newBoard;
      });

      axios.patch(`${url}/lists/${targetList.id}/reorder`, {
        authenticity_token: authenticityToken(),
        above: above_position,
        below: below_position,
      });
    } else {
      const startList = board.lists.find(
        (list) => list.id === parseInt(source.droppableId)
      );
      const endList = board.lists.find(
        (list) => list.id === parseInt(destination.droppableId)
      );
      const targetCard = startList.cards.find(
        (card) => card.id === parseInt(draggableId)
      );

      if (startList === endList) {
        const newCards = Array.from(startList.cards);
        newCards.splice(source.index, 1);
        newCards.splice(destination.index, 0, targetCard);
        const { above_position, below_position, midstring } = getPositions(
          newCards,
          destination.index
        );
        newCards[destination.index].position = midstring;

        const newBoard = { ...board };
        newBoard.lists[
          board.lists.findIndex((l) => l === startList)
        ].cards = newCards;
        updateBoard({ ...newBoard });

        axios.patch(`${url}/cards/${targetCard.id}/reorder`, {
          authenticity_token: authenticityToken(),
          above: above_position,
          below: below_position,
        });
      } else {
        const startCards = Array.from(startList.cards);
        startCards.splice(source.index, 1);
        const endCards = Array.from(endList.cards);
        endCards.splice(destination.index, 0, targetCard);
        const { above_position, below_position, midstring } = getPositions(
          endCards,
          destination.index
        );
        endCards[destination.index].position = midstring;

        const newBoard = { ...board };
        newBoard.lists[
          board.lists.findIndex((l) => l === startList)
        ].cards = startCards;
        newBoard.lists[
          board.lists.findIndex((l) => l === endList)
        ].cards = endCards;
        updateBoard({ ...newBoard });

        axios.patch(`${url}/cards/${targetCard.id}/reorder`, {
          authenticity_token: authenticityToken(),
          above: above_position,
          below: below_position,
          new_list: endList.id,
        });
      }
    }
  };

  return (
    <div>
      <div className="h-screen flex body-scrollbar">
        <Sidebar toggleSidebar={toggleSidebar} isOpen={isOpen} />
        <div className="flex-1 min-w-0 flex flex-col bg-white">
          <Header
            toggleSidebar={toggleSidebar}
            title={board.title}
            handleUpdateTitle={handleUpdateTitle}
          />
          <div className="flex-1 overflow-auto">
            <main className="p-3 h-full inline-flex">
              <DragDropContext onDragEnd={onDragEnd}>
                <Droppable
                  droppableId="board"
                  direction="horizontal"
                  type="list"
                >
                  {(provided) => (
                    <div
                      className="inline-flex"
                      {...provided.droppableProps}
                      ref={provided.innerRef}
                    >
                      {board.lists.map((list, index) => (
                        <List
                          index={index}
                          list={list}
                          handleDeletion={handleDeleteList}
                          key={list.id}
                          board_slug={board.slug}
                          handleNewCard={handleNewCard}
                        />
                      ))}
                      {provided.placeholder}
                    </div>
                  )}
                </Droppable>
              </DragDropContext>
              <div className="flex flex-col w-84 rounded-md">
                <div className="bg-lightgray rounded-md py-3 px-2 ">
                  <button
                    onClick={async () => {
                      await toggleNewList(true);
                      document.getElementById("newListTitle").focus();
                    }}
                    className={`focus:outline-none w-full rounded-md py-1 px-1 hover:bg-gray-400 hover:bg-opacity-50 flex items-center ${
                      newListIsOpen ? "hidden" : ""
                    }`}
                  >
                    <svg
                      className="h-6 w-6 text-gray-600"
                      viewBox="0 0 24 24"
                      fill="none"
                    >
                      <path
                        d="M12 7v10m5-5H7"
                        stroke="currentColor"
                        strokeWidth="2"
                        strokeLinecap="round"
                      />
                    </svg>
                    <h3 className="flex-shrink-0 text-sm text-gray-600 capitalize">
                      Add another list
                    </h3>
                  </button>
                  <div className={`${newListIsOpen ? "" : "hidden"}`}>
                    <form onSubmit={() => handleNewList()}>
                      <input
                        id="newListTitle"
                        className="p-2 border-2 border-indigo-600 rounded-md body-scrollbar w-full focus:outline-none text-sm text-gray-700"
                        placeholder="Enter list title.."
                        onChange={(e) => setNewListTitle(e.target.value)}
                        value={newListTitle}
                      />
                      <div className="flex items-center mt-2 rounded-sm">
                        <button
                          type="submit"
                          className="focus:outline-none flex px-3 py-1.5 text-sm font-medium text-white bg-indigo-600 rounded-md hover:bg-indigo-500"
                        >
                          <span className="text-center">Add List</span>
                        </button>
                        <button
                          type="button"
                          onClick={() => toggleNewList(false)}
                          className="ml-1 text-gray-700 focus:outline-none hover:bg-gray-400 p-1 rounded-md"
                        >
                          <svg
                            className="h-6 w-6 text-gray-700"
                            fill="currentColor"
                            viewBox="0 0 20 20"
                          >
                            <path
                              d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
                              clipRule="evenodd"
                              fillRule="evenodd"
                            ></path>
                          </svg>
                        </button>
                      </div>
                    </form>
                  </div>
                </div>
              </div>
            </main>
          </div>
        </div>
      </div>  
    </div>
  );
};

export default Board;
