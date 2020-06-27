import React, { useRef, useState } from "react";
import Fuse from "fuse.js";
import { useOutsideAlerter } from "../lib" 

const Search = (props) => {
  const [query, updateQuery] = useState("");
  const [isOpen, toggleResults] = useState(false);

  const fuseOptions = {
    shouldSort: true,
    minMatchCharLength: 2,
    threshold: 0.4,
    keys: ["title"],
  };

  const fuse = new Fuse(props.lists.map((x) => x.cards).flat(), fuseOptions);
  const results = fuse.search(query);

  const wrapperRef = useRef(null);
  useOutsideAlerter(wrapperRef, () => {toggleResults(false); updateQuery("");});

  return (
    <div className="flex-shrink-1 relative w-64" ref={wrapperRef}>
      <span className="absolute inset-y-0 left-0 pl-3 flex items-center">
        <svg className="h-6 w-6 text-gray-600" viewBox="0 0 24 24" fill="none">
          <path
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
          />
        </svg>
      </span>
      <input
        className="block w-64 rounded-md border border-gray-400 pl-10 pr-4 py-2 text-sm text-gray-900 placeholder-gray-600 flex-grow-0"
        placeholder="Search"
        value={query}
        onChange={(event) => updateQuery(event.target.value)}
        onClick={() => toggleResults(true)}
      />
      {results.length > 0 && isOpen && (
        <div className="absolute z-30 w-144 bg-lightgray rounded-md mt-2 flex flex-col pb-3 pl-4 shadow-md max-h-2xl overflow-y-auto">
          <h3 className="text-xs font-semibold text-gray-600 uppercase tracking-wide mt-3">
            Your Cards
          </h3>
          {results.map((el, i) => (
            <button key={i} className="cursor-pointer w-64 mt-2 text-left focus:outline-none" onClick={() => {
              toggleResults(false)
              let title = document.getElementById(`card-${el.item.id}-title`)
              title.classList.add("highlight")
              setTimeout(() => {
                title.classList.remove("highlight")
              }, 2000)
            }}>
              <div className="py-2 px-3 bg-white shadow-base rounded-md shadow flex flex-wrap justify-between items-baseline">
                <span className="text-sm font-normal leading-snug text-gray-900 break-all">
                  {el.item.title}
                </span>
              </div>
              <span className="whitespace-no-wrap text-gray-700 font-semibold text-sm mt-1 ml-1">
                in {props.lists.find((l) => l.id === el.item.list_id).title}
              </span>
            </button>
          ))}
        </div>
      )}
    </div>
  );
};

export default Search;
