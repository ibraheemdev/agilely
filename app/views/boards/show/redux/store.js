import { createStore, applyMiddleware, combineReducers } from "redux";
import * as reducers from "./reducers";
import thunk from "redux-thunk";

const configureStore = () => {
  const rootReducer = combineReducers(reducers);
  const store = createStore(rootReducer, applyMiddleware(thunk));
  return store;
};

export default configureStore;
