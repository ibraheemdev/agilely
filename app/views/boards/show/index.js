import React from "react";
import Board from "./components/board";
import configureStore from "./redux/store";
import { Provider } from "react-redux";

const reduxStore = configureStore();

const Index = (props) => (
  <Provider store={reduxStore}>
    <Board {...props}/>
  </Provider>
);

export default Index;
