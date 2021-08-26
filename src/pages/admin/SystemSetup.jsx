import React from 'react';
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
	DatePicker,
	SelectPicker,
	Notification,
	Modal,
} from 'rsuite';

import useAxios from 'axios-hooks';
import axios from 'axios';
import BotLoginErrorView from './components/BotLoginErrorView';

const { useState, useEffect } = React;

export default function SystemSetup() {
	const [systemInfo, setsystemInfo] = useState();
	const [{ data: systemInfoData, loading: systemLoading, error: systemError }, systemRefch] = useAxios({
		url: '/api/system/info',
	});

	const [botInfo, setbotInfo] = useState();
	const [{ data: botInfoData, loading: botLoading, error: botError }, botRefetch] = useAxios({
		url: '/api/bot/info',
	});

	useEffect(() => {
		if (botInfoData?.data) {
			setbotInfo(botInfoData?.data);
			setbotLoginConfig({
				...botLoginConfig,
				account: botInfoData?.data?.account,
			});
		}

		if (systemInfoData?.data) {
			const info = systemInfoData?.data;
			const newInfo = {
				school_time: new Date(info?.school_time * 1000), // 开学时间设置
				few_weeks: info?.few_weeks, // 这学期共几周
				notice_minute: info?.notice_minute, // 提前多少分钟通知
			};
			setsystemInfo(newInfo);
		}
	}, [botInfoData, systemInfoData]);

	const [schoolTime, setschoolTime] = useState();
	const [schoolTimeLoad, setschoolTimeLoad] = useState(false);
	const APIUpdateSchoolTime = (time) => {
		if (!time) {
			Notification.error({
				title: '参数异常！',
			});
			return;
		}

		setschoolTimeLoad(true);

		const params = new URLSearchParams();
		params.append('time', time);

		axios
			.post('/api/system/schooltime', params, {})
			.then((res) => {
				const { data } = res;
				const { code } = data;
				setschoolTimeLoad(false);

				if (code < 1) {
					Notification.error({
						title: '修改失败，请稍后重试！',
					});
				} else {
					const user = data?.data;

					Notification.success({
						title: `修改开学时间数据成功！`,
					});

					// 修改数据成功，刷新页面数据
					systemRefch();
				}
			})
			.catch((error) => {
				Notification.error({
					title: '修改失败，' + error || '修改失败，请稍后重试！',
				});
				setschoolTimeLoad(false);
			});
	};

	const inputFewWeeksRef = React.createRef();
	const [fewWeeks, setfewWeeks] = useState(0);
	const [fewWeeksLoad, setfewWeeksLoad] = useState(false);
	const APIUpdateFewWeeks = (weeks) => {
		if (!weeks) {
			Notification.error({
				title: '参数异常！',
			});
			return;
		}

		setfewWeeksLoad(true);

		const params = new URLSearchParams();
		params.append('weeks', weeks);

		axios
			.post('/api/system/fewweeks', params, {})
			.then((res) => {
				const { data } = res;
				const { code } = data;
				setfewWeeksLoad(false);

				if (code < 1) {
					Notification.error({
						title: '修改失败，请稍后重试！',
					});
				} else {
					const user = data?.data;

					Notification.success({
						title: `修改这学期共几周数据成功！`,
					});

					// 修改数据成功，刷新页面数据
					systemRefch();
				}
			})
			.catch((error) => {
				Notification.error({
					title: '修改失败，' + error || '修改失败，请稍后重试！',
				});
				setfewWeeksLoad(false);
			});
	};

	const inputNoticeMinuteRef = React.createRef();
	const [noticeMinute, setnoticeMinute] = useState(0);
	const [noticeMinuteLoad, setnoticeMinuteLoad] = useState(false);
	const APIUpdateNoticeMinute = (minute) => {
		if (!minute) {
			Notification.error({
				title: '参数异常！',
			});
			return;
		}

		setnoticeMinuteLoad(true);

		const params = new URLSearchParams();
		params.append('minute', minute);

		axios
			.post('/api/system/noticeminute', params, {})
			.then((res) => {
				const { data } = res;
				const { code } = data;
				setnoticeMinuteLoad(false);

				if (code < 1) {
					Notification.error({
						title: '修改失败，请稍后重试！',
					});
				} else {
					const user = data?.data;

					Notification.success({
						title: `修改提前多少分钟通知时间成功！`,
					});

					// 修改数据成功，刷新页面数据
					systemRefch();
				}
			})
			.catch((error) => {
				Notification.error({
					title: '修改失败，' + error || '修改失败，请稍后重试！',
				});
				setnoticeMinuteLoad(false);
			});
	};

	const [groupCode, setgroupCode] = useState(0);
	const [groupCodeLoad, setgroupCodeLoad] = useState(false);
	const APIUpdateGroupCode = (code) => {
		if (!code) {
			Notification.error({
				title: '参数异常！',
			});
			return;
		}

		setgroupCodeLoad(true);

		const params = new URLSearchParams();
		params.append('group_code', code);

		axios
			.post('/api/bot/update/groupcode', params, {})
			.then((res) => {
				const { data } = res;
				const { code } = data;
				setgroupCodeLoad(false);

				if (code < 1) {
					Notification.error({
						title: '修改失败，请稍后重试！',
					});
				} else {
					const user = data?.data;

					Notification.success({
						title: `修改通知的QQ群号成功！`,
					});

					// 修改数据成功，刷新页面数据
					systemRefch();
				}
			})
			.catch((error) => {
				Notification.error({
					title: '修改失败，' + error || '修改失败，请稍后重试！',
				});
				setgroupCodeLoad(false);
			});
	};

	const [reLoginLoad, setreLoginLoad] = useState(false);
	const APIreLogin = () => {
		if (botInfo?.on_line) {
			Notification.error({
				title: 'Bot 账号已在线',
			});
			return;
		}

		setreLoginLoad(true);

		const params = new URLSearchParams();

		axios
			.post('/api/bot/relogin', params, {})
			.then((res) => {
				const { data } = res;
				const { code } = data;
				setreLoginLoad(false);

				if (code < 1) {
					Notification.error({
						title: 'Bot 账号重新登陆失败，请稍后重试！',
					});
				} else {
					const user = data?.data;

					Notification.success({
						title: `Bot 账号重新登陆成功！`,
						onClose: () => {
							// Bot 账号重新登陆成功，刷新页面数据
							botRefetch();
						},
					});

					// Bot 账号重新登陆成功，刷新页面数据
					// botRefetch();
				}
			})
			.catch((error) => {
				Notification.error({
					title: '重新登陆失败，' + error || '重新登陆失败，请稍后重试！',
				});
				setreLoginLoad(false);
			});
	};

	// 登陆账号
	const [botLoginConfig, setbotLoginConfig] = useState({
		account: undefined, // 账号
		password: '',
		showModal: false, // 显示弹窗
		netLoding: false, // 标识请求中状态
		callback: {}, // 登陆回调数据
	});
	const APIBotLogin = () => {
		const { account, password } = botLoginConfig;

		if (!account || !password) {
			Notification.error({
				title: '账号密码不得为空！',
			});
			return;
		}

		// 请求中，不能重复发起请求
		if (botLoginConfig.netLoding) {
			return;
		}

		setbotLoginConfig({
			...botLoginConfig,
			netLoding: true,
		});

		const params = new URLSearchParams();
		params.append('account', account);
		params.append('password', password);

		axios
			.post('/api/bot/login', params, {})
			.then((res) => {
				setbotLoginConfig({
					...botLoginConfig,
					netLoding: false,
				});

				const { data } = res;
				const { code } = data;

				if (code < 1) {
					Notification.error({
						title: data?.msg || '失败，请稍后重试！',
					});
				} else {
					const callback = data?.data;

					console.log('回调', callback);

					if (callback?.error) {
						// 需要二次认证
						setbotLoginConfig({
							...botLoginConfig,
							showModal: true,
							callback: callback,
						});
					} else {
						setbotLoginConfig({
							...botLoginConfig,
							showModal: false,
						});
						botRefetch();
					}
				}
			})
			.catch((error) => {
				Notification.error({
					title: '失败，' + error || '失败，请稍后重试！',
				});
				setbotLoginConfig({
					...botLoginConfig,
					netLoding: false,
				});
			});
	};

	return systemInfo ? (
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
					<Button style={{ color: '#FFF' }} appearance="primary" disabled={true}>
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
							<DatePicker
								style={{ width: 320 }}
								placeholder="选择日期"
								defaultValue={systemInfo?.school_time}
								showWeekNumbers
								onChange={(date) => {
									const time = parseInt(date.getTime() / 1000);
									// console.log("时间", time);
									setschoolTime(time);
								}}
							/>
						</Col>
						<Col xs={1}></Col>
						<Button
							style={{ color: '#FFF' }}
							appearance="primary"
							loading={schoolTimeLoad}
							onClick={() => APIUpdateSchoolTime(schoolTime)}
						>
							保存修改
						</Button>
					</Row>

					<Row style={{ marginTop: 20 }}>
						<Col justify="center" xs={4}>
							<p className="item-title">这学期共几周：</p>
						</Col>
						<Col xs={6}>
							<InputGroup>
								<InputGroup.Button
									onClick={() => {
										inputFewWeeksRef.current.handleMinus();
									}}
								>
									-
								</InputGroup.Button>
								<InputNumber
									className={'custom-input-number'}
									defaultValue={systemInfo?.few_weeks}
									ref={inputFewWeeksRef}
									max={99}
									min={1}
									onChange={(value) => {
										setfewWeeks(value);
									}}
								/>
								<InputGroup.Button
									onClick={() => {
										inputFewWeeksRef.current.handlePlus();
									}}
								>
									+
								</InputGroup.Button>
							</InputGroup>
						</Col>
						<Col xs={1}></Col>
						<Button
							style={{ color: '#FFF' }}
							loading={fewWeeksLoad}
							onClick={() => APIUpdateFewWeeks(fewWeeks)}
							appearance="primary"
						>
							保存修改
						</Button>
					</Row>

					<Row style={{ marginTop: 70 }}>
						<Col justify="center" xs={4}>
							<p className="item-title">提前多少分钟通知：</p>
						</Col>
						<Col xs={6}>
							<InputGroup>
								<InputGroup.Button
									onClick={() => {
										inputNoticeMinuteRef.current.handleMinus();
									}}
								>
									-
								</InputGroup.Button>
								<InputNumber
									className={'custom-input-number'}
									ref={inputNoticeMinuteRef}
									defaultValue={systemInfo?.notice_minute}
									max={999}
									min={1}
									onChange={(value) => {
										setnoticeMinute(value);
									}}
								/>
								<InputGroup.Button
									onClick={() => {
										inputNoticeMinuteRef.current.handlePlus();
									}}
								>
									+
								</InputGroup.Button>
							</InputGroup>
						</Col>
						<Col xs={1}></Col>
						<Button
							style={{ color: '#FFF' }}
							appearance="primary"
							loading={noticeMinuteLoad}
							onClick={() => APIUpdateNoticeMinute(noticeMinute)}
						>
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
									<Input
										placeholder="请输入机器人 QQ 账号"
										value={botLoginConfig?.account}
										onChange={(value) => {
											setbotLoginConfig({
												...botLoginConfig,
												account: value,
											});
										}}
									/>
								</Col>
							</FlexboxGrid>
							<FlexboxGrid style={{ marginBottom: 20 }} justify="end">
								<Col xs={4}>
									<p className="item-title">QQ 密码：</p>
								</Col>
								<Col xs={10}>
									<Input
										placeholder="请输入密码"
										type="password"
										value={botLoginConfig?.password}
										onChange={(value) => {
											setbotLoginConfig({
												...botLoginConfig,
												password: value,
											});
										}}
									/>
								</Col>
							</FlexboxGrid>
							<FlexboxGrid justify="end">
								{botInfo?.account && !botInfo?.on_line ? (
									<Button
										appearance="primary"
										loading={reLoginLoad}
										style={{ marginRight: 20 }}
										onClick={() => APIreLogin()}
									>
										重新登陆
									</Button>
								) : null}

								<Button
									appearance="primary"
									loading={botLoginConfig?.netLoding}
									onClick={() => {
										APIBotLogin();
									}}
								>
									登陆账号
								</Button>
							</FlexboxGrid>
							{botInfo?.group_list && botInfo?.group_list?.length > 1 ? (
								<>
									<FlexboxGrid style={{ marginTop: 70, marginBottom: 20 }} justify="end">
										<Col xs={6}>
											<p className="item-title">通知 QQ 群号：</p>
										</Col>
										<Col xs={10}>
											<SelectPicker
												data={botInfo?.group_list?.map((item, index) => {
													return {
														label: `${item.code} (${item.name})`,
														value: item.code,
														role: 'Master',
													};
												})}
												defaultValue={botInfo?.group_code || botInfo?.group_list[0]?.code}
												onChange={(value) => {
													// console.log("选中群号", value);
													setgroupCode(value);
												}}
												style={{ width: 224 }}
												searchable={false}
											/>
										</Col>
									</FlexboxGrid>

									<FlexboxGrid justify="end">
										<Button
											loading={groupCodeLoad}
											appearance="primary"
											onClick={() => {
												APIUpdateGroupCode(groupCode);
											}}
										>
											保存修改
										</Button>
									</FlexboxGrid>
								</>
							) : null}
						</Col>

						<Col xs={12} style={{ paddingLeft: 80 }}>
							{botInfo?.account ? (
								<>
									<Avatar circle size="lg" style={{ background: '#0001' }} src={botInfo?.avatar} />
									<p className="qq-name">{botInfo?.nick_name}</p>
									<Tag color={botInfo?.on_line ? 'green' : 'red'}>
										{botInfo?.on_line ? '在线' : '离线'}
									</Tag>
								</>
							) : null}
						</Col>
					</FlexboxGrid>
				</div>
			</div>

			{/* 立即登陆机器人账号弹窗 */}
			<Modal
				show={botLoginConfig.showModal}
				onHide={() => {
					setbotLoginConfig({
						...botLoginConfig,
						showModal: false,
					});
				}}
				backdrop="static"
			>
				<Modal.Header>
					<Modal.Title>本次登陆需要认证</Modal.Title>
				</Modal.Header>
				<Modal.Body>
					{botLoginConfig?.callback?.error ? (
						<BotLoginErrorView
							callbackData={botLoginConfig?.callback}
							onCallBack={(data) => {
								// 需要二次认证
								setbotLoginConfig({
									...botLoginConfig,
									callback: data,
								});
							}}
						/>
					) : null}
				</Modal.Body>
				<Modal.Footer>
					<Button
						onClick={APIBotLogin}
						style={{ color: '#FFF' }}
						appearance="primary"
						disabled={botLoginConfig?.callback?.error ? true : false}
						loading={botLoginConfig.netLoding}
					>
						立即登陆
					</Button>
				</Modal.Footer>
			</Modal>
		</div>
	) : null;
}
