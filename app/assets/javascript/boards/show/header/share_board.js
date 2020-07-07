import React, { useRef, useState } from "react";
import { useOutsideAlerter } from "../lib";
import Participation from "./participation";

const ShareBoardModal = (props) => {
  const wrapperRef = useRef(null);
  useOutsideAlerter(wrapperRef, () => props.toggleSelf(false));
  const [emails, setEmails] = useState([])
  const [currentEmail, setCurrentEmail] = useState("")
  const [error, setError] = useState("")

  const handleKeyDown = evt => {
    if (["Enter", "Tab", ","].includes(evt.key)) {
      evt.preventDefault();

      var value = this.state.value.trim();

      if (value && this.isValid(value)) {
        this.setState({
          items: [...this.state.items, this.state.value],
          value: ""
        });
      }
    }
  };

  const handleDeleteEmail = email => {
    setEmails(oldEmails => {
      newEmails = oldEmails.filter(e => e !== email)
    })
  }

  const handleChange = (event) => {
    setCurrentEmail(event.target.value)
  }

  const handlePaste = event => {
    event.preventDefault();

    var paste = event.clipboardData.getData("text");
    var emails = paste.match(/[\w\d\.-]+@[\w\d\.-]+\.[\w\d\.-]+/g);

    if (emails) {
      var toBeAdded = emails.filter(email => !isInList(email));
     setEmails({emails: [...emails, ...toBeAdded]})
    }
  };

  const isValid = (email) => {
    let error = null;

    if (isInList(email)) {
      error = `${email} has already been added.`;
    }

    if (isEmail(email)) {
      error = `${email} is not a valid email address.`;
    }

    if (error) {
      setError(error)
      return false;
    }

    return true;
  }

  const isInList = (email) => {
    return this.state.items.includes(email);
  }

  const isEmail = (email) => {
    return /[\w\d\.-]+@[\w\d\.-]+\.[\w\d\.-]+/.test(email);
  }

  return (
    <div className="flex items-center fixed z-20 left-0 top-0 w-full h-full overflow-auto bg-black bg-opacity-50">
      <div ref={wrapperRef} className="bg-white p-6 mx-auto max-w-2xl rounded-lg">
        <div>
          <h1 className="text-xl">Share with people or teams</h1>
          <p className="text-sm text-gray-500 tracking-wide">
            Invite users by their email address, and they will recieve an email
            inviting them to your board
          </p>
          {emails.map(e => (
            <div className="tag-item" key={e}>
            {e}
            <button
              type="button"
              className="p-2"
              onClick={() => handleDeleteEmail(e)}
            >
              &times;
            </button>
          </div>
          ))}
          <input
            type="text"
            placeholder="Add people or teams"
            value={currentEmail}
            onKeyDown={() => handleKeyDown}
            onChange={() => handleChange}
            onPaste={() => handlePaste}
            className="mt-2 w-full px-3 py-2 border border-indigo-400 rounded-md placeholder-gray-500 focus:outline-none focus:border-indigo-600 transition duration-150 ease-in-out sm:text-sm sm:leading-5"
          />
          {error && <p className="text-red-600">{error}</p>}
          {props.participations.map((p) => (
            <Participation participation={p} key={p.id}/>
          ))}
        </div>

      </div>
    </div>
  );
};

export default ShareBoardModal;
