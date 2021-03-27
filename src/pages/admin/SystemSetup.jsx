import React from "react";
import {
  FlexboxGrid,
  Row,
  Col,
  Input,
  Divider,
  Button,
  Avatar,
  Tag,
  InputNumber,
  InputGroup,
} from "rsuite";
import useAxios from "axios-hooks";

const { useState, useEffect } = React;

export default function SystemSetup() {
  const [botInfo, setbotInfo] = useState({});
  const [{ data: botInfoData, loading, error }, refetch] = useAxios({
    url: "/api/bot/info",
  });

  useEffect(() => {
    if (botInfoData?.data) {
      setbotInfo(botInfoData?.data);
    }
  }, [botInfoData]);

  return (
    <div className="page-system" style={{ marginTop: 25, marginBottom: 25 }}>
      <div style={{ marginBottom: 30 }}>通知系统设置页面</div>
      <div className="list-view">
        <Divider>网站设定</Divider>
        <Row className="item-view">
          <Col justify="center" xs={4}>
            <p className="item-title">网站标题：</p>
          </Col>
          <Col xs={6}>
            <Input placeholder="Default Input" />
          </Col>
          <Col xs={1}></Col>
          <Button
            style={{ color: "#FFF" }}
            appearance="primary"
            disabled={true}
          >
            保存修改
          </Button>
        </Row>

        <Divider>上课设定</Divider>
        <div className="item-view">
          <Row>
            <Col justify="center" xs={4}>
              <p className="item-title">开学时间设置：</p>
            </Col>
            <Col xs={6}>
              <Input placeholder="Default Input" />
            </Col>
            <Col xs={1}></Col>
            <Button style={{ color: "#FFF" }} appearance="primary">
              保存修改
            </Button>
          </Row>
          <Row style={{ marginTop: 20 }}>
            <Col justify="center" xs={4}>
              <p className="item-title">这学期共几周：</p>
            </Col>
            <Col xs={6}>
              <InputGroup>
                <InputGroup.Button onClick={() => {}}>-</InputGroup.Button>
                <InputNumber
                  className={"custom-input-number"}
                  defaultValue={1}
                  max={999}
                  min={1}
                  onChange={(value) => {}}
                />
                <InputGroup.Button onClick={() => {}}>+</InputGroup.Button>
              </InputGroup>
            </Col>
            <Col xs={1}></Col>
            <Button style={{ color: "#FFF" }} appearance="primary">
              保存修改
            </Button>
          </Row>

          <Row style={{ marginTop: 70 }}>
            <Col justify="center" xs={4}>
              <p className="item-title">提前多少分钟通知：</p>
            </Col>
            <Col xs={6}>
              <InputGroup>
                <InputGroup.Button onClick={() => {}}>-</InputGroup.Button>
                <InputNumber
                  className={"custom-input-number"}
                  defaultValue={1}
                  max={999}
                  min={1}
                  onChange={(value) => {}}
                />
                <InputGroup.Button onClick={() => {}}>+</InputGroup.Button>
              </InputGroup>
            </Col>
            <Col xs={1}></Col>
            <Button style={{ color: "#FFF" }} appearance="primary">
              保存修改
            </Button>
          </Row>
        </div>

        <Divider>机器人设定</Divider>

        <div className="item-view">
          <FlexboxGrid>
            <Col xs={12} justify="end">
              <FlexboxGrid style={{ marginBottom: 20 }} justify="end">
                <Col xs={4}>
                  <p className="item-title">QQ 账号：</p>
                </Col>
                <Col xs={10}>
                  <Input placeholder="Default Input" />
                </Col>
              </FlexboxGrid>
              <FlexboxGrid style={{ marginBottom: 20 }} justify="end">
                <Col xs={4}>
                  <p className="item-title">QQ 密码：</p>
                </Col>
                <Col xs={10}>
                  <Input placeholder="Default Input" />
                </Col>
              </FlexboxGrid>
              <FlexboxGrid justify="end">
                <Button appearance="primary">登陆账号</Button>
              </FlexboxGrid>

              <FlexboxGrid
                style={{ marginTop: 70, marginBottom: 20 }}
                justify="end"
              >
                <Col xs={6}>
                  <p className="item-title">通知 QQ 群号：</p>
                </Col>
                <Col xs={10}>
                  <Input placeholder="Default Input" />
                </Col>
              </FlexboxGrid>
              <FlexboxGrid justify="end">
                <Button appearance="primary">保存修改</Button>
              </FlexboxGrid>
            </Col>

            <Col xs={12} style={{ paddingLeft: 80 }}>
              {botInfo?.account ? (
                <>
                  <Avatar
                    circle
                    size="lg"
                    style={{ background: "#0001" }}
                    src={botInfo?.avatar}
                  />
                  <p className="qq-name">{botInfo?.nick_name}</p>
                  <Tag color={botInfo?.on_line ? "green" : "red"}>
                    {botInfo?.on_line ? "在线" : "离线"}
                  </Tag>
                </>
              ) : null}
            </Col>
          </FlexboxGrid>
        </div>
      </div>
    </div>
  );
}
