import { createStore, applyMiddleware } from "redux";
import { listReducer } from "./index";
import thunk from "redux-thunk";

export default function configureStore() {
  const store = createStore(listReducer, applyMiddleware(thunk));
  return store;
}
