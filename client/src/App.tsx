import React, { useState, useEffect } from "react";

const App = () => {
  let phrase = "Help us decide ";

  type status = "has-poll" | "has-focus" | "idle";
  const [status, setStatus] = useState<status>("idle");

  const pollArray = [
    "where to go this evening",
    "what movie to watch",
    "what to eat",
    "whatever",
    "what to wear to the party",
    "what time to meet up",
  ];

  interface Poll {
    title: string;
  }
  const [poll, setPoll] = useState<Poll>({
    title:
      phrase + pollArray[Math.floor(Math.random() * pollArray.length)],
  });

  useEffect(() => {
    const timeout = setInterval(() => {
      while (status === "idle") {
        let nxt = Math.floor(Math.random() * pollArray.length);
        if (pollArray[nxt] !== poll.title) {
          setPoll({...poll, title: phrase + pollArray[nxt]});
          break;
        }
      }
    }, 2000);
    return () => clearInterval(timeout);
  }, [poll, pollArray]);
  
  return (
    <div className="font-nunito">
      <div className="mx-6">
        <input
          type="text"
          value={poll.title}
          onFocus={() => {
            setStatus("has-poll");
            setPoll({...poll, title: phrase});
          }}
          onBlur={() => poll.title === phrase && setStatus("idle")}
          onChange={(event) => {
            setPoll({...poll, title: event.target.value});
            poll.title === "phrase" ? setStatus("idle") : setStatus("has-poll");
          }}
          className="shadow-inner w-full py-3 mx-auto mt-64 text-center text-2xl placeholder-gray-500 text-gray-600 font-semibold tracking-tight border border-gray-500 rounded-lg focus:outline-none focus:border-blue-400"
        />
      </div>
    </div>
  );
};

export default App;
