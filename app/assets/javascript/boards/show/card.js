import React, { useState } from "react";
import { Draggable } from "react-beautiful-dnd";
import CardDetail from "./card_detail"

const Card = (props) => {
  const [modalIsOpen, openModal] = useState(false);
  return (
    <div className="cursor-pointer">
      <Draggable draggableId={props.card.id.toString()} index={props.index} isDragDisabled={!props.editor}>
        {(provided) => (
          <div
            ref={provided.innerRef}
            {...provided.draggableProps}
            {...provided.dragHandleProps}
            className="mt-2"
            onClick={() => openModal(true)}
          >
            <div className="block py-2 px-3 bg-white rounded-md shadow flex flex-wrap justify-between items-baseline">
              <span className="text-sm font-normal leading-snug text-gray-900 break-all">
                {props.card.title}
              </span>
              <span className="text-xs text-gray-600">
                {props.card.position}
              </span>
            </div>
          </div>
        )}
      </Draggable>
      {modalIsOpen && <CardDetail />}
    </div>
  );
};

export default Card;
