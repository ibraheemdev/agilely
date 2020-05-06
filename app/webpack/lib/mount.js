import React from "react";
import ReactDOM from "react-dom";

const mount = (Component, nodeId) => {
  document.addEventListener("turbolinks:load", () => {
    const node = document.getElementById(nodeId);
    const props = JSON.parse(node.getAttribute("data-props"));
    ReactDOM.render(<Component {...props} />, node);
  });
};

export default mount;
