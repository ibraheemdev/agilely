import * as types from "./types";
import axios from "axios";
import authenticityToken from "../lib/authenticity_token";

const setLists = (lists) => ({ type: types.SET_LISTS, payload: [...lists] });

const addList = (board_slug, list_title) => {
  return (dispatch) => {
    event.preventDefault();
    axios
      .post(`/boards/${board_slug}/lists`, {
        authenticity_token: authenticityToken(),
        list: { title: list_title },
      })
      .then((res) => {
        dispatch(addListSuccess(res.data.list));
      });
  };
};

const addListSuccess = (list) => ({
  type: types.ADD_LIST_SUCCESS,
  payload: {
    ...list,
  },
});

const deleteList = (listId) => {
  return (dispatch) => {
    event.preventDefault();
    axios
      .delete(`/lists/${listId}`, {
        data: { authenticity_token: authenticityToken() },
      })
      .then((res) => {
        dispatch(deleteListSuccess(listId));
      });
  };
};

const deleteListSuccess = (listId) => ({
  type: types.DELETE_LIST_SUCCESS,
  payload: listId,
});

export { setLists, addList, deleteList };
