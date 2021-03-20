import React from "react";

import { CourseList, CreateCourse } from "./screens";

const { useState, useRef, useEffect } = React;

export default function CourseControll() {
  const [screenKey, setscreenKey] = useState("CourseList");
  const screenPropsRef = useRef({
    toPage: (name, params = undefined) => {
      if (params) {
        screenPropsRef.current = {
          toPage: screenPropsRef.current.toPage,
          params,
        };
      } else {
        screenPropsRef.current = {
          toPage: screenPropsRef.current.toPage,
        };
      }
      setscreenKey(name);
    },
  });

  let screens = {
    CourseList: <CourseList {...screenPropsRef.current} />,
    CreateCourse: <CreateCourse {...screenPropsRef.current} />,
  };

  return (
    <div style={{ height: "100%", marginTop: 25 }}>
      {screenKey && screens[screenKey]}
    </div>
  );
}