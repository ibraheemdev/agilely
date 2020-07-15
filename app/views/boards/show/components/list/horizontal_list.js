import React from "react";
import Card from "../card";
import { Droppable, Draggable } from "react-beautiful-dnd";
import TextareaAutosize from "react-textarea-autosize";
import { listActions } from "@redux/lists";
import { useDispatch } from "react-redux";

const HorizontalList = (props) => {
  const dispatch = useDispatch();
  return (
    <Draggable
      draggableId={props.id}
      index={props.index}
      isDragDisabled={!props.canEdit}
    >
      {(provided) => (
        <div
          className="flex flex-col list-scrollbar w-72 rounded-md mr-2"
          {...provided.draggableProps}
          ref={provided.innerRef}
        >
          <div className="relative -mb-2" {...provided.dragHandleProps}>
            <div className="flex items-center justify-between rounded-t-md bg-lightgray pt-3 pb-1">
              <input
                id={`list-${props.id}-title`}
                value={props.list.title}
                className={`w-full text-sm font-medium text-gray-700 bg-lightgray mx-3 py-1 px-1 focus:bg-white rounded-md focus:cursor-auto hover:cursor-pointer ${
                  props.canEdit && "hover:bg-gray-500 hover:bg-opacity-25"
                }`}
                onChange={(event) => props.setTitle(event.target.value)}
                onKeyPress={(e) => {
                  if (e && e.charCode == 13) {
                    document.activeElement.blur();
                  }
                }}
                onBlur={() => props.handleUpdateTitle()}
                disabled={!props.canEdit}
              />
              <span className="mr-3">{props.list.position}</span>
              {props.canEdit && (
                <button
                  onClick={() => {
                    dispatch(listActions.deleteList(props.id));
                  }}
                  className="text-gray-700 focus:outline-none hover:bg-gray-400 p-1 mr-2 rounded-md"
                  disabled={!props.canEdit}
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
              )}
            </div>
          </div>
          <div
            id={`list${props.id}top`}
            className={`min-h-0 overflow-y-auto bg-lightgray ${
              (!props.canEdit || props.newCardIsOpen) && "rounded-b-md pb-2"
            }`}
          >
            <div className="py-1 px-3">
              <Droppable droppableId={props.id} type="card">
                {(provided) => (
                  <div
                    {...provided.droppableProps}
                    ref={provided.innerRef}
                    className="min-h-1"
                  >
                    {props.list.card_ids.map((id, index) => (
                      <Card key={id} id={id} index={index} />
                    ))}
                    {provided.placeholder}
                  </div>
                )}
              </Droppable>
              {props.newCardIsOpen && (
                <form
                  ref={props.newCardRef}
                  onSubmit={() => props.handleNewCard(props.newCardTitle)}
                >
                  <div className="mt-2">
                    <div
                      href="#"
                      className="block py-2 px-3 bg-white rounded-md shadow"
                    >
                      <div className="flex justify-between">
                        <TextareaAutosize
                          id={`new-card-title-${props.id}`}
                          maxRows={16}
                          onChange={(e) =>
                            props.updateNewCardTitle(e.target.value)
                          }
                          value={props.newCardTitle}
                          className="body-scrollbar w-full min-h-16 focus:outline-none text-sm font-normal leading-snug text-gray-900"
                          placeholder="Enter a title for this card..."
                        />
                      </div>
                    </div>
                  </div>
                  <div className={`flex items-center mt-3 rounded-sm`}>
                    <button
                      type="submit"
                      className="flex px-3 py-1.5 text-sm font-medium text-white bg-indigo-600 rounded-md hover:bg-indigo-500"
                    >
                      <span className="text-center">Add Card</span>
                    </button>
                    <button
                      onClick={() => props.toggleNewCard(false)}
                      type="button"
                      className="ml-2 text-gray-700 focus:outline-none hover:bg-gray-400 p-1 mr-3 rounded-md"
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
              )}
            </div>
          </div>
          {!props.newCardIsOpen && props.canEdit && (
            <div className="py-2 px-2 bg-lightgray rounded-b-md">
              <button
                onClick={async () => {
                  await props.toggleNewCard(true);
                  let el = document.getElementById(
                    `new-card-title-${props.id}`
                  );
                  await el.scrollIntoView({ behavior: "smooth" });
                  await el.focus({ preventScroll: true });
                }}
                className="focus:outline-none hover:bg-gray-400 hover:bg-opacity-50 flex items-center w-full py-1 px-1 rounded-md"
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
                <h3
                  className={`flex-shrink-0 text-sm text-gray-600 capitalize`}
                >
                  Add a card
                </h3>
              </button>
            </div>
          )}
        </div>
      )}
    </Draggable>
  );
};
export default HorizontalList;
