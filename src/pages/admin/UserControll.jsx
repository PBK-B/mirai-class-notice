import React from "react";

import useAxios from "axios-hooks";
import axios from "axios";

import {
  Table,
  Panel,
  Loader,
  FlexboxGrid,
  Button,
  Modal,
  Notification,
  Form,
  FormControl,
  FormGroup,
  ControlLabel,
} from "rsuite";

const { useState } = React;
const { Column, HeaderCell, Cell, Pagination } = Table;

export default function UserControll() {
  const [{ data, loading, error }, refetch] = useAxios({
    url: "/api/user/list",
  });

  const [showCreateUser, setshowCreateUser] = useState(false);
  const [createUserLoading, setcreateUserLoading] = useState(false);
  const [createUser, setcreateUser] = useState({
    name: "",
    passwd: "",
  });

  const onCreatUser = (value) => {
    setcreateUser(value);
  };
  const APICreatUser = () => {
    const { name, passwd } = createUser;

    if (!name || !passwd) {
      Notification.error({
        title: "账号密码不得为空！",
      });
      return;
    }

    setcreateUserLoading(true);

    const params = new URLSearchParams();
    params.append("action", "createUser");
    params.append("name", name);
    params.append("password", passwd);

    axios
      .post("/api/user/create", params, {})
      .then((res) => {
        setcreateUserLoading(false);

        const { data } = res;
        const { code } = data;

        if (code < 1) {
          Notification.error({
            title: data?.msg || "创建失败，请稍后重试！",
          });
        } else {
          const user = data?.data;

          Notification.success({
            title: `创建用户 ${user.name} 成功！`,
          });

          // 创建用户成功，关闭弹窗，刷新列表数据，清空编辑框数据
          setshowCreateUser(false);
          setcreateUser({ name: "", passwd: "" });
          refetch();
        }
      })
      .catch((error) => {
        Notification.error({
          title: "创建失败，" + error || "创建失败，请稍后重试！",
        });
        setcreateUserLoading(false);
      });
  };

  let users = [];

  if (data) {
    let { data: data_array } = data;
    data_array = data_array.map((item, index) => {
      return {
        ...item,
        status: item.status == 1 ? "启用" : "禁用",
        time: new Date().toUTCString(),
      };
    });
    users = data_array;
  }

  if (loading) return <Loader backdrop content="loading..." vertical />;

  return (
    <div style={{ marginTop: 25, marginBottom: 25 }}>
      <p>Hello this is admin UserControl Page!!!</p>
      <br />

      <FlexboxGrid style={{ marginBottom: 15 }} justify="end">
        <FlexboxGrid.Item>
          <Button appearance="ghost" onClick={() => setshowCreateUser(true)}>
            创建用户
          </Button>
        </FlexboxGrid.Item>
      </FlexboxGrid>

      <Panel bordered bodyFill>
        <Table
          autoHeight
          data={users}
          onRowClick={(data) => {
            console.log(data);
          }}
        >
          <Column width={70} align="center" fixed>
            <HeaderCell>ID</HeaderCell>
            <Cell dataKey="id" />
          </Column>

          <Column width={200} fixed>
            <HeaderCell>账号</HeaderCell>
            <Cell dataKey="name" />
          </Column>

          <Column width={200}>
            <HeaderCell>状态</HeaderCell>
            <Cell dataKey="status" />
          </Column>

          <Column width={260}>
            <HeaderCell>登陆时间</HeaderCell>
            <Cell dataKey="time" />
          </Column>

          <Column width={160} fixed="right">
            <HeaderCell>操作</HeaderCell>

            <Cell>
              {(rowData) => {
                function handleAction() {
                  alert(`id:${rowData.id}`);
                }
                return (
                  <span>
                    <a onClick={handleAction}> 编辑 </a> |
                    <a onClick={handleAction}> 禁用 </a>
                  </span>
                );
              }}
            </Cell>
          </Column>
        </Table>
      </Panel>

      <Modal
        show={showCreateUser}
        onHide={() => {
          setshowCreateUser(false);
        }}
        backdrop="static"
      >
        <Modal.Header>
          <Modal.Title>创建用户</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <Form fluid onChange={onCreatUser} formValue={createUser}>
            <FormGroup>
              <ControlLabel>账号：</ControlLabel>
              <FormControl name="name" />
            </FormGroup>

            <FormGroup>
              <ControlLabel>密码：</ControlLabel>
              <FormControl name="passwd" />
            </FormGroup>
          </Form>
        </Modal.Body>
        <Modal.Footer>
          <Button
            onClick={APICreatUser}
            style={{ color: "#FFF" }}
            appearance="primary"
            loading={createUserLoading}
          >
            创建账号
          </Button>
        </Modal.Footer>
      </Modal>
    </div>
  );
}
