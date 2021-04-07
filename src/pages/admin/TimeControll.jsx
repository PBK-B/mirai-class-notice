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

const { useState, useEffect } = React;
const { Column, HeaderCell, Cell, Pagination } = Table;

export default function TimeControll() {
  const [{ data, loading, error }, refetch] = useAxios({
    url: "/api/time/list?count=999",
  });

  // 创建时间点部分区域
  const [showCreateTime, setshowCreateTime] = useState(false);
  const [createTimeLoading, setcreateTimeLoading] = useState(false);
  const [createTime, setcreateTime] = useState({
    group: "",
    start: "",
    end: "",
    remark: "",
  });
  const onCreatTime = (value) => {
    setcreateTime(value);
  };
  const APICreatTime = () => {
    const { group, start, end, remark } = createTime;

    if (!group || !start || !end || !remark) {
      Notification.error({
        title: "时间信息需填写完整",
      });
      return;
    }

    setcreateTimeLoading(true);

    const params = new URLSearchParams();
    params.append("group", group);
    params.append("start", start);
    params.append("end", end);
    params.append("remark", remark);

    axios
      .post("/api/time/create", params, {})
      .then((res) => {
        setcreateTimeLoading(false);

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
          setshowCreateTime(false);
          setcreateTime({
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
        setcreateTimeLoading(false);
      });
  };

  // 修改时间点部分区域
  const [showUpdateTime, setshowUpdateTime] = useState(false);
  const [updateTimeData, setupdateTimeData] = useState();
  const updateTime = (time) => {
    if (!time) {
      Notification.error({
        title: "修改时间异常！",
      });
      return;
    }

    setshowUpdateTime(true);
    setupdateTimeData(time);
  };
  const onUpdateTime = (value) => {
    const newData = { ...updateTimeData, ...value };
    setupdateTimeData(newData);
  };
  const [updateTimeLoading, setupdateTimeLoading] = useState(false);
  const APIUpdateTime = (time) => {
    const { id, group, start, end, remark } = time;

    if (!group || !start || !end || !remark) {
      Notification.error({
        title: "时间信息需填写完整",
      });
      return;
    }

    setupdateTimeLoading(true);

    const params = new URLSearchParams();
    params.append("id", id);
    params.append("group", group);
    params.append("start", start);
    params.append("end", end);
    params.append("remark", remark);

    axios
      .post("/api/time/update", params, {})
      .then((res) => {
        setupdateTimeLoading(false);

        const { data } = res;
        const { code } = data;

        if (code < 1) {
          Notification.error({
            title: data?.msg || "修改失败，请稍后重试！",
          });
        } else {
          const user = data?.data;

          Notification.success({
            title: `修改时间 ${user.remark} 成功！`,
          });

          // 创建用户成功，关闭弹窗，刷新列表数据，清空编辑框数据
          setshowUpdateTime(false);
          setupdateTimeData(null);
          refetch();
        }
      })
      .catch((error) => {
        Notification.error({
          title: "修改失败，" + error || "修改失败，请稍后重试！",
        });
        setupdateTimeLoading(false);
      });
  };

  // 表格数据处理部分代码
  const [times, settimes] = useState([]);
  const [timeList, settimeList] = useState([]);
  useEffect(() => {
    if (data) {
      let { data: data_array } = data;
      data_array = data_array?.map((item, index) => {
        return {
          ...item,
          status: item.status == 1 ? "启用" : "禁用",
          time: new Date().toUTCString(),
        };
      });
      settimeList(data_array);
      onChangePage(data_array, 1);
    }
  }, [data]);

  // 数据分页处理
  const [activePage, setactivePage] = useState(1);
  const lengthMenu = [{ label: <p>10</p>, value: 10 }];
  const [displayLength, setdisplayLength] = useState(lengthMenu[0].value);
  const onChangePage = (list, value) => {
    setactivePage(value);
    const index = displayLength * (value - 1);
    const newArray = list?.slice(index, index + displayLength);
    settimes(newArray);
  };

  // 重新注册上课时间系统定时任务
  const APIrerunTimeTasks = () => {
    axios
      .get("/api/time/reruntasks")
      .then(function (response) {
        const { data } = response;
        const { code } = data;

        if (code < 1) {
          Notification.error({
            title: data?.msg || "定时任务重置失败，请稍后重试！",
          });
        } else {
          Notification.success({
            title: "全部上课时间定时任务重置成功！",
          });
        }
      })
      .catch(function (error) {
        Notification.error({
          title: error
            ? "定时任务重置失败，" + error
            : "定时任务重置失败，请稍后重试！",
        });
      });
  };

  // 手动触发上课通知
  const APInoticeTime = (id) => {
    axios
      .get("/api/time/test01?p=" + id)
      .then(function (response) {
        const { data } = response;
        const { code } = data;

        if (code < 1) {
          Notification.error({
            title: data?.msg || "手动触发上课通知失败，请稍后重试！",
          });
        } else {
          Notification.success({
            title: "手动触发上课通知成功！",
          });
        }
      })
      .catch(function (error) {
        Notification.error({
          title: error
            ? "手动触发上课通知失败，" + error
            : "手动触发上课通知失败，请稍后重试！",
        });
      });
  };

  if (loading) return <Loader backdrop content="loading..." vertical />;

  return (
    <div style={{ marginTop: 25, marginBottom: 25 }}>
      <p>上课时间管理页面</p>
      <br />

      <FlexboxGrid style={{ marginBottom: 15 }} justify="end">
        <FlexboxGrid.Item>
          <Button
            appearance="primary"
            style={{ marginRight: 15, color: "#FFF" }}
            onClick={APIrerunTimeTasks}
          >
            重置定时
          </Button>
          <Button appearance="ghost" onClick={() => setshowCreateTime(true)}>
            添加时间
          </Button>
        </FlexboxGrid.Item>
      </FlexboxGrid>

      <Panel bordered bodyFill>
        <Table
          autoHeight
          data={times}
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
                function editTimeAction(e) {
                  // alert(`id:${rowData.id}`);
                  updateTime({
                    id: rowData.id,
                    group: rowData.group,
                    start: rowData.start,
                    end: rowData.end,
                    remark: rowData.remark,
                  });
                  e.stopPropagation();
                }

                function noticeTime(e) {
                  APInoticeTime(rowData.id);
                  e.stopPropagation();
                }

                return (
                  <span>
                    <a onClick={editTimeAction}> 编辑 </a> |
                    <a onClick={noticeTime}> 通知 </a>
                  </span>
                );
              }}
            </Cell>
          </Column>
        </Table>
        <Table.Pagination
          lengthMenu={lengthMenu}
          displayLength={10}
          activePage={activePage}
          total={timeList?.length || 0}
          onChangeLength={(value) => {
            console.log("改变每页长度", value);
            // setdisplayLength(value);
            // const newArray = coursesList.splice(value);
            // setcourses(newArray);
          }}
          onChangePage={(value) => {
            // console.log("改变的页面", value);
            onChangePage(timeList, value);
          }}
        ></Table.Pagination>
      </Panel>

      <Modal
        show={showCreateTime}
        onHide={() => {
          setshowCreateTime(false);
        }}
        backdrop="static"
      >
        <Modal.Header>
          <Modal.Title>添加上课时间</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <Form fluid onChange={onCreatTime} formValue={createTime}>
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
            onClick={APICreatTime}
            style={{ color: "#FFF" }}
            appearance="primary"
            loading={createTimeLoading}
          >
            添加时间
          </Button>
        </Modal.Footer>
      </Modal>

      <Modal
        show={showUpdateTime}
        onHide={() => {
          setshowUpdateTime(false);
          setupdateTimeData(null);
        }}
        backdrop="static"
      >
        <Modal.Header>
          <Modal.Title>修改上课时间</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <Form fluid onChange={onUpdateTime} formValue={updateTimeData}>
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
            onClick={() => APIUpdateTime(updateTimeData)}
            style={{ color: "#FFF" }}
            appearance="primary"
            loading={updateTimeLoading}
          >
            修改时间
          </Button>
        </Modal.Footer>
      </Modal>
    </div>
  );
}
