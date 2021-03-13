import React from "react";
import axios from "axios";
import { Container, Content, Notification, FlexboxGrid } from "rsuite";

import { AppHeader, AppFooter } from "../../components";
import AppMenu from "./components/AppMenu";
import Home from "./Home";

const { useState, useEffect } = React;

export default function App() {
  // 注销登陆
  const __unLogin = () => {
    // localStorage.removeItem("user");
    window.location.href = "/login";
  };

  const [activeKey, setactiveKey] = useState("1");

  const handleSelect = (eventKey) => {
    // console.log(eventKey);
    setactiveKey(eventKey);
  };

  const AppContents = {
    1: <Home />,
    2: <Home />,
    3: <Home />,
  };

  return (
    <Container style={{ height: "100%" }}>
      <AppHeader user={{ name: "Bin" }} />

      <FlexboxGrid style={{ height: "100%" }}>
        <FlexboxGrid.Item
          style={{
            height: "100%",
            background: "#f7f7fa",
            paddingTop: 15,
            borderRight: "1px solid rgb(229, 229, 234)",
          }}
          colspan={3}
        >
          <AppMenu
            width={260}
            expanded={true}
            activeKey={activeKey}
            onSelect={handleSelect}
          />
        </FlexboxGrid.Item>

        <FlexboxGrid.Item style={{ flex: 1, height: "100%" }}>
          <Content
            style={{
              margin: "auto",
              maxWidth: "1200px",
              minWidth: "1200px",
            }}
          >
            {AppContents[activeKey]}
          </Content>
        </FlexboxGrid.Item>
      </FlexboxGrid>

      <AppFooter />
    </Container>
  );
}
