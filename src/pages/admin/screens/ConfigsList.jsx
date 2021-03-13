import React from "react";
import {
  FlexboxGrid,
  Button,
  Table,
  Panel,
  Loader,
  Tag,
  Drawer,
  SelectPicker,
  Row,
} from "rsuite";
import useAxios from "axios-hooks";
import MonacoEditor from "react-monaco-editor";
import axios from "axios";

const { useState, useEffect } = React;
const { Column, HeaderCell, Cell, Pagination } = Table;

const ConfigInfoDrawer = (props) => {
  const { show, onClose, configData } = props;
  const { id, name, remarks, status, json } = configData;

  const selectLanguages = [
    {
      label: "Android Java - okhttp",
      value: "android-java-okhttp",
      role: "Master",
    },
    {
      label: "JavaScript - jQuery",
      value: "javascript-jquery",
      role: "Master",
    },
  ];
  const selectModes = [
    {
      label: "id",
      value: "id",
      role: "Master",
    },
    {
      label: "name",
      value: "name",
      role: "Master",
    },
  ];
  const [selectLanguage, setselectLanguage] = useState(
    selectLanguages[0].value
  );
  const [selectMode, setselectMode] = useState(selectModes[0].value);

  const [code, setcode] = useState({
    text: "",
    language: "javasrcipt",
  });
  const [jsonStr, setjsonStr] = useState("");
  const options = {
    scrollBeyondLastLine: false,
    automaticLayout: true,
    readOnly: true, // 编辑器只读
    showFoldingControls: "always",
    formatOnPaste: true,
    formatOnType: true,
    folding: true,
    wordWrap: true,
  };

  function formatJSON(val = "") {
    try {
      const res = JSON.parse(val);
      return JSON.stringify(res, null, 2);
    } catch {
      return val;
    }
  }

  const loadCode = (language = "android-java-okhttp", mode = "id") => {
    let codeObj = code;

    const domain = window.location.protocol + "//" + window.location.host;
    const parameter = mode === "id" ? "id=" + id : "name=" + name;

    switch (language) {
      case "android-java-okhttp":
        codeObj = {
          language: "java",
          text: `
OkHttpClient client = new OkHttpClient().newBuilder().build();
Request request = new Request.Builder()
  .url("${domain}/api/appConfig.php?${parameter}&action=useConfig")
  .method("GET", null)
  .build();
Response response = client.newCall(request).execute();
`,
        };
        break;
      case "javascript-jquery":
        codeObj = {
          language: "javascript",
          text: `
var settings = {
  "url": "${domain}/api/appConfig.php?${parameter}&action=useConfig",
  "method": "GET",
  "timeout": 0,
};

$.ajax(settings).done(function (response) {
  console.log(response);
});
`,
        };
        break;
    }

    setcode(codeObj);
  };

  useEffect(() => {
    setjsonStr(json);
    loadCode();
  }, [configData]);

  return (
    <Drawer show={show} onHide={onClose}>
      <Drawer.Header>
        <Drawer.Title>配置文件详情</Drawer.Title>
      </Drawer.Header>
      <Drawer.Body>
        <div>
          配置 ID：<b>{id}</b>
        </div>
        <div style={{ marginTop: 10 }}>
          配置名称：<b>{name}</b>
        </div>
        <div style={{ marginTop: 10 }}>
          状态：
          {status === "0" ? (
            <Tag color="red">禁用</Tag>
          ) : (
            <Tag color="cyan">启用</Tag>
          )}
        </div>

        {remarks && (
          <div style={{ marginTop: 10 }}>
            备注：
            <span style={{ color: "#0006" }}>{remarks}</span>
          </div>
        )}

        <Row style={{ width: "100%", marginTop: 25, paddingLeft: 5 }}>
          <SelectPicker
            data={selectModes}
            defaultValue={selectMode}
            onChange={(value, event) => {
              setselectMode(value);
              loadCode(selectLanguage, value);
            }}
            style={{ width: 110, marginRight: 15 }}
            searchable={false}
            cleanable={false}
          />

          <SelectPicker
            data={selectLanguages}
            defaultValue={selectLanguage}
            onChange={(value, event) => {
              setselectLanguage(value);
              loadCode(value, selectMode);
            }}
            style={{ width: 220 }}
            searchable={false}
            cleanable={false}
          />
        </Row>

        {code && (
          <div
            style={{
              marginTop: 25,
              borderRadius: 5,
              background: "#1e1e1e",
              padding: "20px 0",
              minHeight: 550,
            }}
          >
            <MonacoEditor
              width="100%"
              height="550"
              language={code.language}
              theme="vs-dark"
              value={code.text}
              options={options}
              editorDidMount={(editor, monaco) => {
                // console.log("editorDidMount", editor);
                // editor.focus();
              }}
            />
          </div>
        )}

        <div
          style={{
            marginTop: 25,
            borderRadius: 5,
            background: "#1e1e1e",
            padding: "20px 0",
            minHeight: 200,
          }}
        >
          <MonacoEditor
            width="100%"
            height="200"
            language="json"
            theme="vs-dark"
            value={formatJSON(jsonStr)}
            options={options}
            onChange={() => {}}
            editorDidMount={(editor, monaco) => {
              // 格式化代码
              editor.trigger("anyString", "editor.action.formatDocument");
              editor.getAction("editor.action.formatDocument").run();
            }}
          />
        </div>
      </Drawer.Body>
      <Drawer.Footer>
        <Button
          style={{ marginBottom: 15 }}
          onClick={onClose}
          appearance="subtle"
        >
          关闭
        </Button>
      </Drawer.Footer>
    </Drawer>
  );
};

export default function ConfigsList(props) {
  const { toPage } = props;

  // 获取用户
  let user = localStorage.getItem("user") || null;
  try {
    user = user ? JSON.parse(user) : null;
  } catch (e) {
    user = null;
  }

  const [showDrawer, setshowDrawer] = useState(false);
  const [configData, setconfigData] = useState({});

  const onShowDrawer = (configItem) => {
    setconfigData(configItem);
    setshowDrawer(true);
  };

  // 获取配置列表请求
  const [{ data, loading, error }, refetch] = useAxios({
    url: "/api/appConfig.php?action=queryConfigList",
    headers: {
      token: user?.token,
    },
  });

  let configs = [];
  if (data) {
    const { data: data_array } = data;
    configs = data_array;
  }

  const APIUpdateConfigStatus = (id) => {
    if (!id) {
      return;
    }

    const params = new URLSearchParams();
    params.append("action", "updateConfigStatus");
    params.append("id", id);

    axios
      .post("/api/appConfig.php", params, {
        headers: {
          token: user?.token,
        },
      })
      .then((res) => {
        // 更新状态成功，刷新配置列表
        refetch();
      })
      .catch((error) => {});
  };

  useEffect(() => {
    refetch();
  }, [toPage]);

  if (loading) return <Loader backdrop content="loading..." vertical />;

  return (
    <div>
      <FlexboxGrid style={{ marginBottom: 15 }} justify="end">
        <FlexboxGrid.Item>
          <Button appearance="ghost" onClick={() => toPage("CreateConfig")}>
            创建配置
          </Button>
        </FlexboxGrid.Item>
      </FlexboxGrid>

      <Panel bordered bodyFill>
        <Table
          autoHeight
          data={configs}
          onRowClick={(data) => {
            // console.log(data);
            onShowDrawer(data);
          }}
        >
          <Column width={70} align="center" fixed>
            <HeaderCell>ID</HeaderCell>
            <Cell dataKey="id" />
          </Column>

          <Column width={200} fixed>
            <HeaderCell>配置名称</HeaderCell>
            <Cell dataKey="name" />
          </Column>

          <Column width={200}>
            <HeaderCell>备注</HeaderCell>
            <Cell dataKey="remarks" />
          </Column>

          <Column width={260}>
            <HeaderCell>时间</HeaderCell>
            <Cell dataKey="time" />
          </Column>

          <Column width={90}>
            <HeaderCell>状态</HeaderCell>
            <Cell dataKey="status">
              {(rowData) =>
                rowData.status === "0" ? (
                  <Tag color="red">禁用</Tag>
                ) : (
                  <Tag color="cyan">启用</Tag>
                )
              }
            </Cell>
          </Column>

          <Column width={160} fixed="right">
            <HeaderCell>操作</HeaderCell>

            <Cell
              onClick={(event) => {
                // 禁止操作栏点击事件冒泡
                event.stopPropagation();
              }}
            >
              {(rowData) => {
                const { id, status } = rowData;
                return (
                  <span>
                    <Button
                      style={{ margin: 0, marginRight: 10, padding: 0 }}
                      appearance="link"
                      onClick={() => {
                        toPage("CreateConfig", {
                          configData: rowData,
                        });
                      }}
                    >
                      编辑
                    </Button>
                    |
                    <Button
                      style={{ margin: 0, marginLeft: 10, padding: 0 }}
                      appearance="link"
                      onClick={() => APIUpdateConfigStatus(id)}
                    >
                      {status === "0" ? " 启用 " : " 禁用 "}
                    </Button>
                  </span>
                );
              }}
            </Cell>
          </Column>
        </Table>

        {configData && (
          <ConfigInfoDrawer
            show={showDrawer}
            configData={configData}
            onClose={() => setshowDrawer(false)}
          />
        )}
      </Panel>
    </div>
  );
}
