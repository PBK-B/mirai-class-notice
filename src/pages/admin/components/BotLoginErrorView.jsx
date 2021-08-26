import React from 'react';

import useAxios from 'axios-hooks';
import axios from 'axios';

import { Notification, Message, Input, InputGroup } from 'rsuite';
import QRCode from 'qrcode';

const { useState } = React;

const BotLoginErrorView = (props) => {
	const { callbackData, onCallBack } = props;
	const { error, text, url } = callbackData;
	const [qrCodeUrl, setqrCodeUrl] = useState('');

	// 提交滑动认证 ticket
	const [ticketConfig, setticketConfig] = useState({
		value: '',
		netLoding: false,
	});
	const APITicketRequest = (value) => {
		if (!value) {
			Notification.error({
				title: 'Ticket 不得为空！',
			});
			return;
		}

		// 请求中，阻断重复请求
		if (ticketConfig?.netLoding) return;
		setticketConfig({
			...ticketConfig,
			netLoding: true,
		}); // 设置请求中状态

		// 设置请求参数
		const params = new URLSearchParams();
		params.append('ticket', value);

		axios
			.post('/api/bot/login/ticket', params, {})
			.then((res) => {
				setticketConfig({
					...ticketConfig,
					netLoding: false,
				});

				const { data } = res;
				const { code } = data;

				if (code < 1) {
					Notification.error({
						title: '失败，请稍后重试！',
						description: data?.msg || null,
					});
				} else {
					const callback = data?.data;

					// 返回回调
					onCallBack && onCallBack(callback);

					console.log('请求回调', callback);
				}
			})
			.catch((error) => {
				Notification.error({
					title: '失败，请稍后重试！',
					description: '失败，' + error || '失败，请稍后重试！',
				});
				setticketConfig({
					...ticketConfig,
					netLoding: false,
				});
			});
	};

	switch (error) {
		case 10010:
			// 图形验证码获取失败
			return (
				<>
					{text && (
						<div className="bot_login_err_msg_info">
							<Message type="warning" description={text} />
						</div>
					)}
				</>
			);
			break;
		case 10011:
			// 需要输入图形验证码
			return (
				<>
					{text && (
						<div className="bot_login_err_msg_info">
							<Message type="warning" description={text} />
						</div>
					)}
				</>
			);
			break;
		case 10020:
			// 不安全设备错误
			return (
				<>
					{text && (
						<div className="bot_login_err_msg_info">
							<Message type="warning" description={text} />
						</div>
					)}
					<a href={url} target="_blank">
						{url}
					</a>
				</>
			);
			break;
		case 10030:
			// 需要获取短信验证码登陆
			return (
				<>
					{text && (
						<div className="bot_login_err_msg_info">
							<Message type="warning" description={text} />
						</div>
					)}
				</>
			);
			break;
		case 10040:
			// 需要短信验证码获取太频繁
			return (
				<>
					{text && (
						<div className="bot_login_err_msg_info">
							<Message type="warning" description={text} />
						</div>
					)}
				</>
			);
			break;
		case 10050:
			// 需要扫描二维码登陆或获取短信验证码登陆
			QRCode.toDataURL(url)
				.then((qrurl) => {
					setqrCodeUrl(qrurl);
				})
				.catch((err) => {});
			return (
				<>
					{text && (
						<div className="bot_login_err_msg_info">
							<Message type="warning" description={text} />
						</div>
					)}
					<p>
						推荐使用 QQ
						扫描上面二维码认证或者直接打开二次认证页面，二次认证成功后关闭该页面重新点击登陆即可:{' '}
						<a href={url} target="_blank">
							点我打开二次认证页面
						</a>
					</p>
					<img style={{ width: 135, height: 135, marginTop: 15 }} src={qrCodeUrl} alt="qr" />
				</>
			);
			break;
		case 10060:
			// 需要滑动认证
			return (
				<>
					{text && (
						<div className="bot_login_err_msg_info">
							<Message type="warning" description={text} />
						</div>
					)}
					<p>
						1，请先查看获取滑动二维码的 Ticket 代码教程:{' '}
						<a href="https://github.com/Mrs4s/go-cqhttp/blob/master/docs/slider.md" target="_blank">
							https://github.com/Mrs4s/go-cqhttp/blob/master/docs/slider.md
						</a>
					</p>
					<p>
						2，滑动认证页面地址:{' '}
						<a href={url} target="_blank">
							点我打开认证页面
						</a>
					</p>
					<InputGroup style={{ marginTop: 15 }}>
						<Input
							onChange={(value) => {
								setticketConfig({
									...ticketConfig,
									value: value,
								});
							}}
						/>
						<InputGroup.Button
							loading={ticketConfig?.netLoding}
							onClick={() => APITicketRequest(ticketConfig?.value)}
						>
							提交 Ticket 登陆
						</InputGroup.Button>
					</InputGroup>
				</>
			);
			break;
		case 10070:
			// 其他错误
			return (
				<>
					{text && (
						<div className="bot_login_err_msg_info">
							<Message type="warning" description={text} />
						</div>
					)}
				</>
			);
			break;

		default:
			return <></>;
			break;
	}
};

export default BotLoginErrorView;
