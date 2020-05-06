// Run this example by adding <%= javascript_pack_tag 'hello_react' %> to the head of your layout file,
// like app/views/layouts/application.html.erb. All it does is render <div>Hello React</div> at the bottom
// of the page.

import React from "react";
import mount from "../lib/mount";

const Hello = (props) => {
  console.log(props);
  return (
    <div>
      <div>Hello {props.name}!</div>
    </div>
  );
};
mount(Hello, "hello_react");
