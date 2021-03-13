import React, { useState } from "react";
import axios from "axios";
import {
  Button,
  Container,
  Content,
  FlexboxGrid,
  Panel,
  Form,
  FormGroup,
  ControlLabel,
  FormControl,
  ButtonToolbar,
  Whisper,
  Tooltip,
  Notification,
} from "rsuite";
import md5 from "js-md5";
import { AppHeader, AppFooter } from "../../components";

export default function App() {
  const [formValue, setformValue] = useState({ name: "", password: "" });

  const _handleSubmit = () => {
    const user = formValue.name;
    const passwd = formValue.password;
    // const passwd_md5 = md5(formValue.password);
    // const passwd_md5 = "cd4e731ce352369c9a273b1a69509995";

    if (!user || !passwd) {
      Notification.error({
        title: "账号或密码未输入！",
      });
      return;
    }

    axios
      .get("/api/login?name=" + user + "&password=" + passwd)
      .then((res) => {
        const { data } = res;
        if (data?.code != 1) {
          const msg = data.msg || "服务器错误！";
          Notification.error({
            title: msg,
          });
        } else {
          _goAdminPage(data?.data);
        }
      })
      .catch((err) => {
        Notification.error({
          title: "服务器错误：" + err,
        });
      });
  };

  function setCookie(cname, cvalue, exdays) {
    var d = new Date();
    d.setTime(d.getTime() + exdays * 24 * 60 * 60 * 1000);
    var expires = "expires=" + d.toGMTString();
    document.cookie = cname + "=" + cvalue + "; " + expires;
  }

  const _goAdminPage = (data) => {
    // console.log("数据", data);
    window.location.href = "/admin";
  };

  const __tooltip = (msg) => <Tooltip>{msg}</Tooltip>;

  return (
    <Container style={{ height: "100%" }}>
      <AppHeader />
      <Content>
        <FlexboxGrid style={{ height: "100%" }} justify="center">
          <FlexboxGrid.Item style={{ margin: "auto" }} colspan={7}>
            <Panel header={<h3>欢迎回来，登陆</h3>} bordered>
              <Form
                fluid
                onChange={(formValue) => {
                  setformValue(formValue);
                }}
              >
                <FormGroup>
                  <ControlLabel>账号</ControlLabel>
                  <FormControl name="name" />
                </FormGroup>
                <FormGroup>
                  <ControlLabel>密码</ControlLabel>
                  <FormControl name="password" type="password" />
                </FormGroup>
                <FormGroup>
                  <ButtonToolbar>
                    <Button
                      onClick={_handleSubmit}
                      style={{ color: "#FFF" }}
                      appearance="primary"
                    >
                      立即登陆
                    </Button>
                    <Whisper
                      placement="top"
                      trigger="hover"
                      speaker={__tooltip(
                        "如需重置密码请登陆服务器或联系超级管理员操作！"
                      )}
                    >
                      <Button appearance="link">重置密码?</Button>
                    </Whisper>
                  </ButtonToolbar>
                </FormGroup>
              </Form>
            </Panel>
          </FlexboxGrid.Item>
        </FlexboxGrid>
      </Content>
      <AppFooter />
    </Container>
  );
}
