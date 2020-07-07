import React from "react";
import md5 from "md5";

const Participation = ({ participation }) => (
  <div className="flex items-center hover:bg-gray-300 rounded-full py-1">
    <img
      className="rounded-full h-10 w-10 mr-2"
      src={`https://www.gravatar.com/avatar/${md5(
        participation.user.email
      )}.jpg`}
    />
    <div className="flex w-full justify-between items-center">
      <div>
        <h2 className="text-md font-medium">{participation.user.name}</h2>
        <p className="text-sm text-gray-600">{participation.user.email}</p>
      </div>
      <div>{participation.role}</div>
    </div>
  </div>
);

export default Participation;
