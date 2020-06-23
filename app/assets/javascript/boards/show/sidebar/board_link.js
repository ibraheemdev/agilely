import React from "react";
import { url } from "../lib";

const BoardLink = (props) => (
  <a
    href={`${url}/b/${props.board.slug}`}
    className="flex justify-between items-center px-3 py-2 bg-gray-200 rounded-lg"
  >
    <span className="text-sm font-medium text-gray-900">{props.board.title}</span>
    <span className="text-xs font-semibold text-gray-700">36</span>
  </a>
);

export default BoardLink;
