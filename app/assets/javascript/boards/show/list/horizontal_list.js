import React, { useState } from "react";
import Card from "../card";
import { Droppable, Draggable } from "react-beautiful-dnd";
import AutosizeInput from "react-input-autosize";

const HorizontalList = (props) => {
  const [cardsAreVisible, toggleVisibility] = useState(false);
  return (
    <Draggable
      draggableId={props.list._id.$oid}
      index={props.index}
      isDragDisabled={!props.can_edit}
    >
      {(provided) => (
        <div
          {...provided.draggableProps}
          ref={provided.innerRef}
          className="mb-3 mx-2"
        >
          <div className="relative -mb-2" {...provided.dragHandleProps}>
            <div className="flex items-center">
              <button
                onClick={() => toggleVisibility(!cardsAreVisible)}
                className="p-2 focus:outline-none hover:bg-gray-500 hover:bg-opacity-25 rounded-md"
              >
                <svg
                  width="7"
                  height="12"
                  transform={cardsAreVisible ? "rotate(90)" : undefined}
                  viewBox="0 0 7 12"
                  fill="none"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    d="M1 1L6 6L1 11"
                    stroke="#4A5568"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  />
                </svg>
              </button>
              <AutosizeInput
                id={`list-${props.list._id.$oid}-title`}
                value={props.title}
                inputClassName={`max-w-xs text-md font-medium text-gray-700 ml-2 mr-1 py-1 px-2 focus:bg-white rounded-md focus:cursor-auto hover:cursor-pointer ${
                  props.can_edit && "hover:bg-gray-500 hover:bg-opacity-25"
                }`}
                onChange={(event) => props.setTitle(event.target.value)}
                onKeyPress={(e) => {
                  if (e && e.charCode == 13) {
                    document.activeElement.blur();
                  }
                }}
                onBlur={() => props.handleUpdateTitle()}
                disabled={!props.can_edit}
              />
              {props.can_edit && (
                <div className="flex items-center">
                  <button
                    className="hover:bg-gray-500 hover:bg-opacity-25 p-2 rounded-md"
                    onClick={async () => {
                      await toggleVisibility(true);
                      await props.toggleNewCard(true);
                      let el = document.getElementById(
                        `new-card-title-${props.list._id.$oid}`
                      );
                      await el.scrollIntoView({ behavior: "smooth" });
                      await el.focus({ preventScroll: true });
                    }}
                  >
                    <svg
                      width="12"
                      height="12"
                      className="text-gray-700"
                      viewBox="0 0 12 12"
                      fill="none"
                      xmlns="http://www.w3.org/2000/svg"
                    >
                      <g clipPath="url(#clip0)">
                        <path
                          d="M11.3333 6.00002H0.666664H11.3333ZM6 0.666687V11.3334V0.666687Z"
                          stroke="currentColor"
                          strokeWidth="2"
                          strokeLinecap="round"
                          strokeLinejoin="round"
                        />
                      </g>
                      <defs>
                        <clipPath id="clip0">
                          <rect width="12" height="12" fill="white" />
                        </clipPath>
                      </defs>
                    </svg>
                  </button>
                  <button
                    onClick={() => props.handleDeletion(props.list)}
                    className="hover:bg-gray-500 hover:bg-opacity-25 p-2 rounded-md"
                    disabled={!props.can_edit}
                  >
                    <svg
                      className="h-4 w-4 text-gray-700"
                      fill="currentColor"
                      viewBox="0 0 20 20"
                    >
                      <path
                        d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z"
                        clipRule="evenodd"
                        fillRule="evenodd"
                      ></path>
                    </svg>
                  </button>
                </div>
              )}
            </div>
          </div>
          <div
            id={`list${props.list._id.$oid}top`}
            className={`min-h-0 overflow-y-auto ${
              (!props.can_edit || props.newCardIsOpen) && "rounded-b-md pb-2"
            }`}
          >
            <div className="py-1 px-3">
              <Droppable droppableId={props.list._id.$oid} type="card">
                {(provided) => (
                  <div
                    {...provided.droppableProps}
                    ref={provided.innerRef}
                    className="min-h-1"
                  >
                    {cardsAreVisible &&
                      props.list.cards.map((card, index) => (
                        <Card
                          key={card._id.$oid}
                          card={card}
                          index={index}
                          can_edit={props.can_edit}
                        />
                      ))}
                    {provided.placeholder}
                  </div>
                )}
              </Droppable>
              {props.newCardIsOpen && (
                <form
                  ref={props.newCardRef}
                  onSubmit={() =>
                    props.handleNewCard(props.newCardTitle, props.list._id.$oid)
                  }
                >
                  <div className="mt-2 py-2 px-3 bg-white rounded-md shadow flex flex-wrap justify-between items-baseline">
                    <input
                      type="text"
                      value={props.newCardTitle}
                      onChange={(e) => props.updateNewCardTitle(e.target.value)}
                      id={`new-card-title-${props.list._id.$oid}`}
                      className="text-sm font-normal leading-snug text-gray-900 outline-none"
                    />
                  </div>
                  <button
                    onClick={() => props.toggleNewCard(false)}
                    type="button"
                    className="ml-2 text-gray-700 focus:outline-none hover:bg-gray-400 p-1 mr-3 rounded-md"
                  >
                    <svg
                      className="h-4 w-4 text-gray-700"
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
                </form>
              )}
            </div>
          </div>
        </div>
      )}
    </Draggable>
  );
};
export default HorizontalList;
