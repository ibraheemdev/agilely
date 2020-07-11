import React from "react";

const CardDetail = () => {
  return (
    <div className="flex items-center fixed z-20 left-0 top-0 w-full h-full overflow-auto bg-black bg-opacity-50">
      <div className="bg-white p-4 mx-4 border w-full">
        <div className="flex md:flex-row-reverse flex-wrap">
          <button className="top-0 right-0">
          <svg fill="currentColor" viewBox="0 0 20 20"><path fillRule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clipRule="evenodd"></path></svg>
          </button>
          <div className="w-full md:w-3/4 bg-gray-500 p-4 text-center text-gray-200">
            1
          </div>
          <div className="w-full md:w-1/4 bg-gray-400 p-4 text-center text-gray-700">
            2
          </div>
        </div>
      </div>
    </div>
  );
};

export default CardDetail;
