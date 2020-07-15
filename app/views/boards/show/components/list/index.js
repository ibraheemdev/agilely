import React, { useState, useRef } from "react";
import { useOutsideAlerter } from "../lib";
import VerticalList from "./vertical_list";
import HorizontalList from "./horizontal_list";
import { useSelector } from "react-redux";
import { listSelectors } from "@redux/lists";
import { metadataSelectors } from "@redux/metadata";

const List = (props) => {
  const [newCardIsOpen, toggleNewCard] = useState(false);
  const [newCardTitle, updateNewCardTitle] = useState("");
  const list = useSelector((state) => listSelectors.getList(state, props.id));
  const view = useSelector((state) =>
    metadataSelectors.getBoardDirection(state)
  );
  const canEdit = useSelector((state) => metadataSelectors.canEdit(state));

  const newCardRef = useRef(null);
  useOutsideAlerter(newCardRef, () => toggleNewCard(false));

  if (view === "vertical") {
    return (
      <VerticalList
        {...props}
        newCardIsOpen={newCardIsOpen}
        list={list}
        canEdit={canEdit}
        toggleNewCard={toggleNewCard}
        newCardRef={newCardRef}
        updateNewCardTitle={updateNewCardTitle}
        newCardTitle={newCardTitle}
      />
    );
  } else if (view === "horizontal") {
    return (
      <HorizontalList
        {...props}
        newCardIsOpen={newCardIsOpen}
        list={list}
        canEdit={canEdit}
        toggleNewCard={toggleNewCard}
        newCardRef={newCardRef}
        updateNewCardTitle={updateNewCardTitle}
        newCardTitle={newCardTitle}
      />
    );
  }
};

export default List;
