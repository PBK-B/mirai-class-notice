import React from "react";
import { Input, Row, Col, FlexboxGrid, Button, Notification } from "rsuite";
import axios from "axios";
import MonacoEditor from "react-monaco-editor";

import { PageHead } from "../components";

const { useState } = React;

export default function CreateConfig(props) {
  const { toPage, params = {} } = props;
  const { configData = undefined } = params;

  // if (configData) {
  //   console.log("配置", configData);
  // }

  let user = localStorage.getItem("user") || null;
  try {
    user = user ? JSON.parse(user) : null;
  } catch (e) {
    user = null;
  }

  const [name, setname] = useState(configData?.name);
  const [remarks, setremarks] = useState(configData?.remarks);

  const [code, setcode] = useState(
    configData?.json ? formatJSON(configData.json) : "{}"
  );
  const options = {
    scrollBeyondLastLine: false,
    automaticLayout: true,
    readOnly: false, // 编辑器只读
    showFoldingControls: "always",
    formatOnPaste: true,
    formatOnType: true,
    folding: true,
  };

  const editorDidMount = (editor, monaco) => {
    // console.log("editorDidMount", editor);
    // editor.focus();
  };

  const [createLoading, setcreateLoading] = useState(false);

  const APICreatConfig = (name, remarks, code) => {
    if (!name || !code) {
      Notification.error({
        title: "配置名称和内容不得为空！",
      });
      return;
    }

    try {
      JSON.parse(code);
    } catch (error) {
      Notification.error({
        title: "配置内容不是 JSON 数据格式！",
      });
      return;
    }

    setcreateLoading(true);

    const params = new URLSearchParams();
    params.append("action", "createConfig");
    params.append("name", name);
    params.append("remarks", remarks || "");
    params.append("json", code);

    axios
      .post("/api/appConfig.php", params, {
        headers: {
          token: user?.token,
        },
      })
      .then((res) => {
        setcreateLoading(false);

        const { data } = res;
        const { code } = data;

        if (code < 1) {
          Notification.error({
            title: data?.msg || "创建失败，请稍后重试！",
          });
        } else {
          Notification.success({
            title: `创建配置 ${name} 成功！`,
            onClose: () => {
              // 创建配置成功，关闭弹窗，刷新列表数据，清空编辑框数据
              goBack();
            },
          });
          setcreateLoading(false);
        }
      })
      .catch((error) => {
        Notification.error({
          title: "创建失败，" + error || "创建失败，请稍后重试！",
        });
        setcreateLoading(false);
      });
  };

  const APIUpdateConfig = (id, name, remarks, code) => {
    if (!id || !name || !code) {
      Notification.error({
        title: "配置名称和内容不得为空！",
      });
      return;
    }

    if (
      contrastString(name, configData?.name) &&
      contrastString(remarks, configData?.remarks) &&
      contrastString(code, configData?.json)
    ) {
      Notification.error({
        title: "配置未发生任何变动！",
      });
      return;
    }

    try {
      JSON.parse(code);
    } catch (error) {
      Notification.error({
        title: "配置内容不是 JSON 数据格式！",
      });
      return;
    }

    setcreateLoading(true);

    const params = new URLSearchParams();
    params.append("action", "updateConfig");
    params.append("id", id);
    params.append("name", name);
    params.append("remarks", remarks || "");
    params.append("json", code);

    axios
      .post("/api/appConfig.php", params, {
        headers: {
          token: user?.token,
        },
      })
      .then((res) => {
        setcreateLoading(false);

        const { data } = res;
        const { code } = data;

        if (code < 1) {
          Notification.error({
            title: data?.msg || "修改失败，请稍后重试！",
          });
        } else {
          Notification.success({
            title: `修改配置 ${name} 成功！`,
            onClose: () => {
              // 修改配置成功，关闭弹窗，刷新列表数据，清空编辑框数据
              goBack();
            },
          });
          setcreateLoading(false);
        }
      })
      .catch((error) => {
        Notification.error({
          title: "修改失败，" + error || "修改失败，请稍后重试！",
        });
        setcreateLoading(false);
      });
  };

  function goBack() {
    toPage("ConfigsList");
  }

  function formatJSON(val = "") {
    try {
      const res = JSON.parse(val);
      return JSON.stringify(res, null, 2);
    } catch {
      return val;
    }
  }

  function contrastString(str, subStr) {
    if (
      str?.replace(/\t|\n|\v|\r| |\f/g, "") ===
      subStr?.replace(/\t|\n|\v|\r| |\f/g, "")
    ) {
      return true;
    } else {
      return false;
    }
  }

  return (
    <div>
      <PageHead
        title={configData ? "编辑配置" : "创建配置"}
        onLeftClick={() => goBack()}
      />

      <Row>
        <Col xs={6}>
          <p style={{ marginBottom: 15 }}>配置名称</p>
          <Input
            defaultValue={name}
            placeholder="建议使用英文"
            onChange={(value) => setname(value)}
          />
        </Col>
        <Col xs={2} />
        <Col xs={16}>
          <p style={{ marginBottom: 15 }}>配置备注</p>
          <Input
            defaultValue={remarks}
            placeholder="可不填，备注该配置作用"
            onChange={(value) => setremarks(value)}
          />
        </Col>
      </Row>

      <div
        style={{
          marginTop: 20,
          borderRadius: 5,
          background: "#1e1e1e",
          padding: "20px 0",
          minHeight: 600,
        }}
      >
        <MonacoEditor
          width="100%"
          height="600"
          language="json"
          theme="vs-dark"
          defaultValue={code}
          options={options}
          onChange={(value) => setcode(value)}
          editorDidMount={editorDidMount}
        />
      </div>

      <FlexboxGrid style={{ marginTop: 20, marginBottom: 20 }} justify="end">
        <Button style={{ marginRight: 20 }} appearance="ghost" onClick={goBack}>
          {configData ? "取消保存" : "取消创建"}
        </Button>
        <Button
          style={{ color: "#FFF" }}
          appearance="primary"
          loading={createLoading}
          onClick={() =>
            configData
              ? APIUpdateConfig(configData.id, name, remarks, code)
              : APICreatConfig(name, remarks, code)
          }
        >
          {configData ? "保存配置" : "创建配置"}
        </Button>
      </FlexboxGrid>
    </div>
  );
}
