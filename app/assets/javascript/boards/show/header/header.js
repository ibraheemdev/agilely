import React, { useState } from "react";
import GuestAlert from "./guest_alert";
// import ShareBoardModal from "./share_board";
import AutosizeInput from "react-input-autosize";
// import Search from "./search";

const Header = (props) => {
  const [title, setTitle] = useState(props.title);
  const [shareBoardModalIsOpen, toggleShareBoardModal] = useState(true);

  return (
    <div className="sm:border-b-2 sm:border-gray-200">
      {!props.current_user && (
        <div className="border-b-2 border-gray-200">
          <GuestAlert />
        </div>
      )}
      {/* s */}
      <header>
        <div className="px-6">
          <div className="flex justify-between items-center py-3 border-b border-gray-200">
            <div className="flex-1 min-w-0 flex">
              {props.current_user && (
                <button
                  onClick={() => props.toggleSidebar(true)}
                  className="text-gray-600 mr-3"
                >
                  <svg className="h-6 w-6" viewBox="0 0 24 24" fill="none">
                    <path
                      d="M4 6h16M4 12h16M4 18h7"
                      stroke="currentColor"
                      strokeWidth="2"
                      strokeLinecap="round"
                    />
                  </svg>
                </button>
              )}
              {/* <Search lists={props.lists} /> */}
            </div>
            <div className="ml-6 flex-shrink-0 flex items-center">
              {/* <button>
              <svg
                fill="currentColor"
                className="h-6 w-6 text-gray-600"
                viewBox="0 0 20 20"
              >
                <path d="M10 2a6 6 0 00-6 6v3.586l-.707.707A1 1 0 004 14h12a1 1 0 00.707-1.707L16 11.586V8a6 6 0 00-6-6zM10 18a3 3 0 01-3-3h6a3 3 0 01-3 3z"></path>
              </svg>
            </button> */}
              <button className="ml-6">
                <img
                  className="h-8 w-8 rounded-full object-cover"
                  src="https://upload.wikimedia.org/wikipedia/commons/7/7c/Profile_avatar_placeholder_large.png"
                  alt=""
                />
              </button>
            </div>
          </div>
          <div className="flex items-center justify-between py-2">
            <div className="sm:flex sm:items-center -ml-3">
              <AutosizeInput
                id="boardtitle"
                value={title}
                inputClassName={`box-content text-2xl font-semibold text-gray-900 leading-tight py-1 px-3 bg-white focus:bg-white rounded-md ${
                  props.can_edit &&
                  "hover:cursor-pointer hover:bg-gray-500 hover:bg-opacity-25"
                } focus:cursor-auto`}
                onChange={(event) => setTitle(event.target.value)}
                onKeyPress={(e) => {
                  if (e && e.charCode == 13) {
                    document.activeElement.blur();
                  }
                }}
                onBlur={() => props.handleUpdateTitle(title)}
                disabled={!props.can_edit}
              />
              {/* <div className="mt-1 flex items-center sm:mt-0 sm:ml-6">
              <span className="rounded-full border-2 border-white">
                <img
                  className="h-6 w-6 rounded-full object-cover"
                  src="https://upload.wikimedia.org/wikipedia/commons/7/7c/Profile_avatar_placeholder_large.png"
                  alt=""
                />
              </span>
              <span className="-ml-2 rounded-full border-2 border-white">
                <img
                  className="h-6 w-6 rounded-full object-cover"
                  src="https://upload.wikimedia.org/wikipedia/commons/7/7c/Profile_avatar_placeholder_large.png"
                  alt=""
                />
              </span>
            </div> */}
            </div>
            <div className="flex">
              <span className="hidden sm:inline-flex p-1 border bg-gray-200 rounded-md">
                <button className={`px-2 py-1 rounded focus:outline-none ${props.view === "LIST" && "bg-white shadow"}`} onClick={() => props.toggleView("LIST")}>
                  <svg
                    className="h-6 w-6 text-gray-600"
                    viewBox="0 0 24 24"
                    fill="none"
                  >
                    <path
                      d="M4 6h16M4 10h16M4 14h16M4 18h16"
                      stroke="currentColor"
                      strokeWidth="2"
                      strokeLinecap="round"
                    />
                  </svg>
                </button>
                <button className={`px-2 py-1 rounded focus:outline-none ${props.view === "BOARD" && "bg-white shadow"}`} onClick={() => props.toggleView("BOARD")}>
                  <svg
                    className="h-6 w-6 text-gray-600"
                    viewBox="0 0 24 24"
                    fill="#718096"
                  >
                    <path d="M16 12c0-1.656 1.344-3 3-3s3 1.344 3 3-1.344 3-3 3-3-1.344-3-3zm1 0c0-1.104.896-2 2-2s2 .896 2 2-.896 2-2 2-2-.896-2-2zm-8 0c0-1.656 1.344-3 3-3s3 1.344 3 3-1.344 3-3 3-3-1.344-3-3zm1 0c0-1.104.896-2 2-2s2 .896 2 2-.896 2-2 2-2-.896-2-2zm-8 0c0-1.656 1.344-3 3-3s3 1.344 3 3-1.344 3-3 3-3-1.344-3-3zm1 0c0-1.104.896-2 2-2s2 .896 2 2-.896 2-2 2-2-.896-2-2z" />
                  </svg>
                </button>
              </span>
              <button onClick={() => toggleShareBoardModal(true)}className="flex-shrink-0 ml-5 flex items-center pl-2 pr-4 py-2 text-sm font-medium text-white bg-indigo-600 rounded-md hover:bg-indigo-500">
                <svg className="h-6 w-6" viewBox="0 0 24 24" fill="none">
                  <path
                    d="M12 7v10m5-5H7"
                    stroke="currentColor"
                    strokeWidth="2"
                    strokeLinecap="round"
                  />
                </svg>
                <span className="ml-1">Share Board</span>
              </button>
            </div>
          </div>
        </div>
        <div className="flex px-4 p-1 border-t border-b bg-gray-200 sm:hidden">
          <button className={`inline-flex items-center justify-center w-1/2 px-2 py-1 rounded focus:outline-none ${props.view === "LIST" && "bg-white shadow"}`} onClick={() => props.toggleView("LIST")}>
            <svg
              className="h-6 w-6 text-gray-600"
              viewBox="0 0 24 24"
              fill="none"
            >
              <path
                d="M4 6h16M4 10h16M4 14h16M4 18h16"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
              />
            </svg>
            <span className="ml-2 text-sm font-medium text-gray-600 leading-none">
              List
            </span>
          </button>
          <button className={`inline-flex items-center justify-center w-1/2 px-2 py-1 rounded focus:outline-none ${props.view === "BOARD" && "bg-white shadow"}`} onClick={() => props.toggleView("BOARD")}>
            <svg
              className="h-6 w-6 text-gray-600"
              viewBox="0 0 24 24"
              fill="#718096"
            >
              <path d="M16 12c0-1.656 1.344-3 3-3s3 1.344 3 3-1.344 3-3 3-3-1.344-3-3zm1 0c0-1.104.896-2 2-2s2 .896 2 2-.896 2-2 2-2-.896-2-2zm-8 0c0-1.656 1.344-3 3-3s3 1.344 3 3-1.344 3-3 3-3-1.344-3-3zm1 0c0-1.104.896-2 2-2s2 .896 2 2-.896 2-2 2-2-.896-2-2zm-8 0c0-1.656 1.344-3 3-3s3 1.344 3 3-1.344 3-3 3-3-1.344-3-3zm1 0c0-1.104.896-2 2-2s2 .896 2 2-.896 2-2 2-2-.896-2-2z" />
            </svg>
            <span className="ml-2 text-sm font-medium text-gray-900 leading-none">
              Board
            </span>
          </button>
        </div>
      </header>
    </div>
  );
};

export default Header;