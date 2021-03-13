import React from "react";
import { Sidenav, Nav, Icon } from "rsuite";

export default function AppMenu(props) {
  return (
    <Sidenav {...props}>
      <Sidenav.Body>
        <Nav>
          <Nav.Item eventKey="1" icon={<Icon icon="dashboard" />}>
            系统状态
          </Nav.Item>
          <Nav.Item eventKey="2" icon={<Icon icon="peoples" />}>
            用户管理
          </Nav.Item>
          <Nav.Item eventKey="3" icon={<Icon icon="gear" />}>
            上课时间
          </Nav.Item>
          <Nav.Item eventKey="4" icon={<Icon icon="gear" />}>
            课程管理
          </Nav.Item>
        </Nav>
      </Sidenav.Body>
    </Sidenav>
  );
}
