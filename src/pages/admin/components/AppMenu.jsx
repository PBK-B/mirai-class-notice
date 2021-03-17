import React from "react";
import { Sidenav, Nav, Icon } from "rsuite";

export default function AppMenu(props) {
  return (
    <Sidenav {...props}>
      <Sidenav.Body>
        <Nav>
          <Nav.Item eventKey="1" icon={<Icon icon="dashboard" />}>
            <p className="app-menu-text">系统状态</p>
          </Nav.Item>
          <Nav.Item eventKey="2" icon={<Icon icon="peoples" />}>
            <p className="app-menu-text">用户管理</p>
          </Nav.Item>
          <Nav.Item eventKey="3" icon={<Icon icon="realtime" />}>
            <p className="app-menu-text">上课时间</p>
          </Nav.Item>
          <Nav.Item eventKey="4" icon={<Icon icon="calendar" />}>
            <p className="app-menu-text">课程管理</p>
          </Nav.Item>
        </Nav>
      </Sidenav.Body>
    </Sidenav>
  );
}
