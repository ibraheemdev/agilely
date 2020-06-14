import React from "react";
import Logo from "../../../../assets/images/exlogo.png";

const Sidebar = props => {
  return (
    <div
      className={`transform transition-transform fixed z-10 inset-y-0 left-0 w-64 px-8 py-4 bg-gray-100 border-r overflow-auto ${
        props.isOpen
          ? "translate-x-0 ease-out duration-150"
          : "-translate-x-full ease-in duration-150"
      }`}
    >
      <div className="-mx-3 pl-3 pr-1 flex items-center justify-between">
        <span>
          <img src={Logo} alt="" className="h-10 w-10" />
        </span>
        <button
          onClick={() => props.toggleSidebar(false)}
          className="text-gray-700"
        >
          <svg fill="none" viewBox="0 0 24 24" className="h-6 w-6">
            <path
              d="M6 18L18 6M6 6L18 18"
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              stroke="currentColor"
            />
          </svg>
        </button>
      </div>
      <nav className="mt-8">
        <h3 className="text-xs font-semibold text-gray-600 uppercase tracking-wide">
          Issues
        </h3>
        <div className="mt-2 -mx-3">
          <a
            href="#"
            className="flex justify-between items-center px-3 py-2 bg-gray-200 rounded-lg"
          >
            <span className="text-sm font-medium text-gray-900">All</span>
            <span className="text-xs font-semibold text-gray-700">36</span>
          </a>
          <a
            href="#"
            className="flex justify-between items-center px-3 py-2 bg-gray-200 rounded-lg"
          >
            <span className="text-sm font-medium text-gray-700">
              Assigned to me
            </span>
            <span className="text-xs font-semibold text-gray-700">2</span>
          </a>
          <a
            href="#"
            className="flex justify-between items-center px-3 py-2 bg-gray-200 rounded-lg"
          >
            <span className="text-sm font-medium text-gray-700">
              Created by me
            </span>
            <span className="text-xs font-semibold text-gray-700">1</span>
          </a>
          <a
            href="#"
            className="flex justify-between items-center px-3 py-2 bg-gray-200 rounded-lg"
          >
            <span className="text-sm font-medium text-gray-700">Archived</span>
          </a>
        </div>

        <h3 className="mt-8 text-xs font-semibold text-gray-600 uppercase tracking-wide">
          Tags
        </h3>
        <div className="mt-2 -mx-3">
          <a
            href="#"
            className="flex justify-between items-center px-3 py-2 bg-gray-200 rounded-lg"
          >
            <span className="text-sm font-medium text-gray-700">Bug</span>
          </a>
          <a
            href="#"
            className="flex justify-between items-center px-3 py-2 bg-gray-200 rounded-lg"
          >
            <span className="text-sm font-medium text-gray-700">
              Feature Request
            </span>
          </a>
          <a
            href="#"
            className="flex justify-between items-center px-3 py-2 bg-gray-200 rounded-lg"
          >
            <span className="text-sm font-medium text-gray-700">Marketing</span>
          </a>
          <a
            href="#"
            className="flex justify-between items-center px-3 py-2 bg-gray-200 rounded-lg"
          >
            <span className="text-sm font-medium text-gray-700">v2.0</span>
          </a>
          <a
            href="#"
            className="flex justify-between items-center px-3 py-2 bg-gray-200 rounded-lg"
          >
            <span className="text-sm font-medium text-gray-700">
              Enhancement
            </span>
          </a>
          <a
            href="#"
            className="flex justify-between items-center px-3 py-2 bg-gray-200 rounded-lg"
          >
            <span className="text-sm font-medium text-gray-700">Design</span>
          </a>
        </div>
        <button className="mt-2 -ml-1 flex items-center text-sm font-medium text-gray-600">
          <svg
            className="h-5 w-5 text-gray-500"
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
          <span className="ml-1">New Project</span>
        </button>
      </nav>
    </div>
  );
};

export default Sidebar;
