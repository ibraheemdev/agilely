import React, { useState } from "react";

const NewList = (props) => {
  const [isOpen, toggleOpen] = useState(false);
  return (
    <div className="flex flex-col w-84 rounded-md">
      <div className="bg-lightgray rounded-md py-3 px-2 ">
        <button
          onClick={async () => {
            await toggleOpen(true);
            document.getElementById("newListTitle").focus();
          }}
          className={`focus:outline-none w-full rounded-md py-1 px-1 hover:bg-gray-400 hover:bg-opacity-50 flex items-center ${
            isOpen ? "hidden" : ""
          }`}
        >
          <svg
            className="h-6 w-6 text-gray-600"
            viewBox="0 0 24 24"
            fill="none"
          >
            <path
              d="M12 7v10m5-5H7"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
            />
          </svg>
          <h3 className="flex-shrink-0 text-sm text-gray-600 capitalize">
            Add another list
          </h3>
        </button>
        <div className={`${isOpen ? "" : "hidden"}`}>
          <form onSubmit={() => props.handleSubmit()}>
            <input
              id="newListTitle"
              className="p-2 border-2 border-indigo-600 rounded-md body-scrollbar w-full focus:outline-none text-sm text-gray-700"
              placeholder="Enter list title.."
              onChange={(e) => props.setTitle(e.target.value)}
              value={props.title}
            />
            <div className="flex items-center mt-2 rounded-sm">
              <button
                type="submit"
                className="focus:outline-none flex px-3 py-1.5 text-sm font-medium text-white bg-indigo-600 rounded-md hover:bg-indigo-500"
              >
                <span className="text-center">Add List</span>
              </button>
              <button
                type="button"
                onClick={() => toggleOpen(false)}
                className="ml-1 text-gray-700 focus:outline-none hover:bg-gray-400 p-1 rounded-md"
              >
                <svg
                  className="h-6 w-6 text-gray-700"
                  fill="currentColor"
                  viewBox="0 0 20 20"
                >
                  <path
                    d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
                    clipRule="evenodd"
                    fillRule="evenodd"
                  ></path>
                </svg>
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default NewList;
