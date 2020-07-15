import React, { useState } from "react";
import { Draggable } from "react-beautiful-dnd";
import { useSelector } from "react-redux";
import { cardSelectors } from "@redux/cards";
import { metadataSelectors } from "@redux/metadata";
import CardDetail from "./card_detail";

const Card = ({ id, index }) => {
  const [modalIsOpen, openModal] = useState(false);
  const { title, position } = useSelector((state) =>
    cardSelectors.getCard(state, id)
  );
  const canEdit = useSelector((state) => metadataSelectors.canEdit(state));
  return (
    <div className="cursor-pointer">
      <Draggable draggableId={id} index={index} isDragDisabled={!canEdit}>
        {(provided) => (
          <div
            ref={provided.innerRef}
            {...provided.draggableProps}
            {...provided.dragHandleProps}
            className="mt-2"
            onClick={() => openModal(true)}
          >
            <div className="py-2 px-3 bg-white rounded-md shadow flex flex-wrap justify-between items-baseline">
              <span
                className="text-sm font-normal leading-snug text-gray-900 break-all"
                id={`card-${id}-title`}
              >
                {title}
              </span>
              <span className="text-xs text-gray-600">{position}</span>
            </div>
          </div>
        )}
      </Draggable>
      {modalIsOpen && <CardDetail />}
    </div>
  );
};

export default Card;
