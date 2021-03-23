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

export default function TimeControll() {
  const [{ data, loading, error }, refetch] = useAxios({
    url: "/api/time/list",
  });

  const [showCreateUser, setshowCreateUser] = useState(false);
  const [createUserLoading, setcreateUserLoading] = useState(false);
  const [createUser, setcreateUser] = useState({
    group: "",
    start: "",
    end: "",
    remark: "",
  });

  const onCreatUser = (value) => {
    setcreateUser(value);
  };
  const APICreatUser = () => {
    const { group, start, end, remark } = createUser;

    if (!group || !start || !end || !remark) {
      Notification.error({
        title: "时间信息需填写完整",
      });
      return;
    }

    setcreateUserLoading(true);

    const params = new URLSearchParams();
    params.append("action", "createUser");
    params.append("group", group);
    params.append("start", start);
    params.append("end", end);
    params.append("remark", remark);

    axios
      .post("/api/time/create", params, {})
      .then((res) => {
        setcreateUserLoading(false);

        const { data } = res;
        const { code } = data;

        if (code < 1) {
          Notification.error({
            title: data?.msg || "添加失败，请稍后重试！",
          });
        } else {
          const user = data?.data;

          Notification.success({
            title: `添加时间 ${user.remark} 成功！`,
          });

          // 创建用户成功，关闭弹窗，刷新列表数据，清空编辑框数据
          setshowCreateUser(false);
          setcreateUser({
            group: group,
            start: "",
            end: "",
            remark: "",
          });
          refetch();
        }
      })
      .catch((error) => {
        Notification.error({
          title: "添加失败，" + error || "添加失败，请稍后重试！",
        });
        setcreateUserLoading(false);
      });
  };

  let users = [];

  if (data) {
    let { data: data_array } = data;
    data_array = data_array?.map((item, index) => {
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
      <p>上课时间管理页面</p>
      <br />

      <FlexboxGrid style={{ marginBottom: 15 }} justify="end">
        <FlexboxGrid.Item>
          <Button appearance="ghost" onClick={() => setshowCreateUser(true)}>
            添加时间
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
            <HeaderCell>分类</HeaderCell>
            <Cell dataKey="group" />
          </Column>

          <Column width={200}>
            <HeaderCell>上课时间</HeaderCell>
            <Cell dataKey="start" />
          </Column>

          <Column width={200}>
            <HeaderCell>下课时间</HeaderCell>
            <Cell dataKey="end" />
          </Column>

          <Column width={260}>
            <HeaderCell>备注</HeaderCell>
            <Cell dataKey="remark" />
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
                    <a onClick={handleAction}> 编辑 </a>
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
          <Modal.Title>添加上课时间</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <Form fluid onChange={onCreatUser} formValue={createUser}>
            <FormGroup>
              <ControlLabel>分类：</ControlLabel>
              <FormControl name="group" />
            </FormGroup>

            <FormGroup>
              <ControlLabel>上课时间：</ControlLabel>
              <FormControl name="start" />
            </FormGroup>

            <FormGroup>
              <ControlLabel>下课时间：</ControlLabel>
              <FormControl name="end" />
            </FormGroup>

            <FormGroup>
              <ControlLabel>备注：</ControlLabel>
              <FormControl name="remark" />
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
            添加时间
          </Button>
        </Modal.Footer>
      </Modal>
    </div>
  );
}
