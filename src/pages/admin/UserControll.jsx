import React, { useEffect } from 'react';

import useAxios from 'axios-hooks';
import axios from 'axios';

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
	Row,
} from 'rsuite';

import { UserStore } from './stores';
import { observer } from 'mobx-react';

const { useState } = React;
const { Column, HeaderCell, Cell, Pagination } = Table;

function UserControll() {
	const [{ data, loading, error }, refetch] = useAxios({
		url: '/api/user/list',
	});

	// 创建用户操作
	const [showCreateUser, setshowCreateUser] = useState(false);
	const [createUserLoading, setcreateUserLoading] = useState(false);
	const [createUser, setcreateUser] = useState({
		name: '',
		passwd: '',
	});
	const onCreateUser = (value) => {
		setcreateUser(value);
	};
	const APICreateUser = () => {
		const { name, passwd } = createUser;

		if (!name || !passwd) {
			Notification.error({
				title: '账号密码不得为空！',
			});
			return;
		}

		setcreateUserLoading(true);

		const params = new URLSearchParams();
		params.append('action', 'createUser');
		params.append('name', name);
		params.append('password', passwd);

		axios
			.post('/api/user/create', params, {})
			.then((res) => {
				setcreateUserLoading(false);

				const { data } = res;
				const { code } = data;

				if (code < 1) {
					Notification.error({
						title: data?.msg || '创建失败，请稍后重试！',
					});
				} else {
					const user = data?.data;

					Notification.success({
						title: `创建用户 ${user.name} 成功！`,
					});

					// 创建用户成功，关闭弹窗，刷新列表数据，清空编辑框数据
					setshowCreateUser(false);
					setcreateUser({ name: '', passwd: '' });
					refetch();
				}
			})
			.catch((error) => {
				Notification.error({
					title: '创建失败，' + error || '创建失败，请稍后重试！',
				});
				setcreateUserLoading(false);
			});
	};

	// 修改用户密码操作
	const [showUpdateUser, setshowUpdateUser] = useState(false);
	const [updateUserLoading, setupdateUserLoading] = useState(false);
	const [updateUser, setupdateUser] = useState({
		id: 0,
		name: '',
		passwd: '',
	});
	const onUpdateUser = (value) => {
		setupdateUser({ ...updateUser, ...value });
	};
	const APIUpdateUser = (id, passwd) => {
		if (!id || !passwd) {
			Notification.error({
				title: '新密码不得为空！',
			});
			return;
		}

		setupdateUserLoading(true);

		const params = new URLSearchParams();
		params.append('id', id);
		params.append('password', passwd);

		axios
			.post('/api/user/update', params, {})
			.then((res) => {
				setupdateUserLoading(false);

				const { data } = res;
				const { code } = data;

				if (code < 1) {
					Notification.error({
						title: data?.msg || '修改失败，请稍后重试！',
					});
				} else {
					const user = data?.data;

					Notification.success({
						title: `修改用户 ${user.name} 成功！`,
					});

					// 修改用户成功，关闭弹窗，刷新列表数据，清空编辑框数据
					setshowUpdateUser(false);
					setcreateUser({ id: 0, name: '', passwd: '' });
					refetch();
				}
			})
			.catch((error) => {
				Notification.error({
					title: '修改失败，' + error || '修改失败，请稍后重试！',
				});
				setupdateUserLoading(false);
			});
	};

	// 用户列表数据处理
	let users = [];

	if (data) {
		let { data: data_array } = data;
		data_array = data_array?.map((item, index) => {
			return {
				...item,
				status: item.status == 1 ? '启用' : '禁用',
				time: new Date().toUTCString(),
			};
		});
		users = data_array;
	}

	// 禁用用户
	const APIDisableUser = (id, isDisable) => {
		const msgStrTitle = isDisable ? '禁用' : '启用';
		const params = new URLSearchParams();
		params.append('id', id);

		axios
			.post('/api/user/upstatus', params, {})
			.then((res) => {
				const { data } = res;
				const { code } = data;

				if (code < 1) {
					Notification.error({
						title: data?.msg || msgStrTitle + '失败，请稍后重试！',
					});
				} else {
					const user = data?.data;

					Notification.success({
						title: `${msgStrTitle}用户 ${user.name} 成功！`,
					});

					// 禁用用户成功，关闭弹窗，刷新列表数据，清空编辑框数据
					refetch();
				}
			})
			.catch((error) => {
				Notification.error({
					title: msgStrTitle + '失败，' + error || msgStrTitle + '失败，请稍后重试！',
				});
			});
	};

	if (loading) return <Loader backdrop content="loading..." vertical />;

	return (
		<div style={{ marginTop: 25, marginBottom: 25 }}>
			<p>后台用户管理页面</p>
			<br />

			<FlexboxGrid style={{ marginBottom: 15 }} justify="end">
				<FlexboxGrid.Item>
					<Button
						appearance="ghost"
						onClick={() => setshowCreateUser(true)}
						disabled={UserStore?.me?.id !== 1}
					>
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

						<Cell style={{ padding: 0, display: 'flex', alignItems: 'center' }}>
							{(rowData) => {
								function editAction(e) {
									setupdateUser({
										id: rowData.id,
										name: rowData.name,
										passwd: '',
									});
									setshowUpdateUser(true);

									// 结束事件分发
									e.stopPropagation();
								}

								function disableAction(e) {
									// 结束事件分发
									e.stopPropagation();
									APIDisableUser(rowData.id, rowData.status === '启用');
								}

								return (
									<Row style={{ overflow: 'hidden' }}>
										<Button
											appearance="link"
											onClick={editAction}
											disabled={
												rowData.id === UserStore?.me?.id || UserStore?.me?.id === 1
													? false
													: true
											}
										>
											{' '}
											编辑{' '}
										</Button>{' '}
										|
										<Button
											appearance="link"
											onClick={disableAction}
											disabled={
												rowData.id === UserStore?.me?.id || UserStore?.me?.id === 1
													? false
													: true
											}
										>
											{rowData.status === '启用' ? ' 禁用 ' : ' 启用 '}
										</Button>
									</Row>
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
					<Form fluid onChange={onCreateUser} formValue={createUser}>
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
						onClick={APICreateUser}
						style={{ color: '#FFF' }}
						appearance="primary"
						loading={createUserLoading}
					>
						创建账号
					</Button>
				</Modal.Footer>
			</Modal>

			<Modal
				show={showUpdateUser}
				onHide={() => {
					setshowUpdateUser(false);
				}}
				backdrop="static"
			>
				<Modal.Header>
					<Modal.Title>修改用户密码</Modal.Title>
				</Modal.Header>
				<Modal.Body>
					<Form fluid onChange={onUpdateUser} formValue={updateUser}>
						<FormGroup>
							<ControlLabel>账号：</ControlLabel>
							<FormControl name="name" disabled />
						</FormGroup>

						<FormGroup>
							<ControlLabel>密码：</ControlLabel>
							<FormControl name="passwd" />
						</FormGroup>
					</Form>
				</Modal.Body>
				<Modal.Footer>
					<Button
						onClick={() => APIUpdateUser(updateUser.id, updateUser.passwd)}
						style={{ color: '#FFF' }}
						appearance="primary"
						loading={updateUserLoading}
					>
						修改密码
					</Button>
				</Modal.Footer>
			</Modal>
		</div>
	);
}

export default observer(UserControll);
