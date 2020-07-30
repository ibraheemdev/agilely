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
  const [title, updateTitle] = useState(list.title);
  const view = useSelector((state) =>
    metadataSelectors.getBoardDirection(state)
  );
  const canEdit = useSelector((state) => metadataSelectors.canEdit(state));

  const newCardRef = useRef(null);
  useOutsideAlerter(newCardRef, () => toggleNewCard(false));

  const listProps = {
    ...props,
    list: list,
    title: title,
    updateTitle: updateTitle,
    newCardIsOpen: newCardIsOpen,
    toggleNewCard: toggleNewCard,
    canEdit: canEdit,
    newCardRef: newCardRef,
    updateNewCardTitle: updateNewCardTitle,
    newCardTitle: newCardTitle,
  };

  if (view === "vertical") {
    return <VerticalList {...listProps} />;
  } else if (view === "horizontal") {
    return <HorizontalList {...listProps} />;
  }
};

export default List;
