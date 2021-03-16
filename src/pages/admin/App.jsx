import React from "react";
import axios from "axios";
import { Container, Content, Notification, FlexboxGrid } from "rsuite";

import { AppHeader, AppFooter } from "../../components";
import AppMenu from "./components/AppMenu";
import Home from "./Home";
import UserControll from "./UserControll";
import TimeControll from "./TimeControll";

const { useState, useEffect } = React;

export default function App() {
  // 注销登陆
  const __unLogin = () => {
    document.cookie = "u_token=; expires=Thu, 01 Jan 1970 00:00:00 GMT";
    window.location.href = "/login";
  };

  const [activeKey, setactiveKey] = useState("1");
  const [user, setuser] = useState({});

  const apiGetMe = () => {
    axios
      .get("/api/user/me")
      .then((res) => {
        const { data } = res;
        const { code } = data;
        if (code < 1) {
          Notification.error({
            title: data?.msg || "登陆失效，请重新登陆！",
            onClose: () => {
              __unLogin();
            },
          });
        } else {
          const _user = data?.data;
          setuser(_user);
        }
      })
      .catch((err) => {
        Notification.error({
          title: "请重新登陆！" + err,
          onClose: () => {
            __unLogin();
          },
        });
      });
  };

  useEffect(() => {
    apiGetMe();
  }, []);

  const handleSelect = (eventKey) => {
    // console.log(eventKey);
    setactiveKey(eventKey);
  };

  const AppContents = {
    1: <Home />,
    2: <UserControll />,
    3: <TimeControll />,
  };

  return (
    <Container style={{ height: "100%" }}>
      <AppHeader user={user} />

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
