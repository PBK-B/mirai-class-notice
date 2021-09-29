import React, { useState, useEffect } from 'react';
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

const { Column, HeaderCell, Cell, Pagination } = Table;

export default function PluginControl() {
	const [{ data, loading, error }, refetch] = useAxios({
		url: '/api/plugin/list',
	});

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

	if (loading) return <Loader backdrop content="loading..." vertical />;

	return (
		<div style={{ marginTop: 25, marginBottom: 25 }}>
			<p>插件管理页面</p>
			<br />

			<FlexboxGrid style={{ marginBottom: 15 }} justify="end">
				<FlexboxGrid.Item>
					<Button appearance="ghost" onClick={() => {}} disabled={UserStore?.me?.id !== 1}>
						安装插件
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
						<HeaderCell>名称</HeaderCell>
						<Cell dataKey="name" />
					</Column>

					<Column width={200}>
						<HeaderCell>包名</HeaderCell>
						<Cell dataKey="package" />
					</Column>

					<Column width={100}>
						<HeaderCell>作者</HeaderCell>
						<Cell>
							{(rowData) => (
								<a href={rowData?.website || '#'} target="_blank">
									{'@' + rowData?.author}
								</a>
							)}
						</Cell>
					</Column>

					<Column width={70}>
						<HeaderCell>版本</HeaderCell>
						<Cell dataKey="version" />
					</Column>

					<Column width={70}>
						<HeaderCell>状态</HeaderCell>
						<Cell dataKey="status" />
					</Column>

					<Column width={260}>
						<HeaderCell>安装时间</HeaderCell>
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
											配置
										</Button>
										|
										<Button appearance="link" onClick={() => {}}>
											日志
										</Button>
										|
										<Button
											appearance="link"
											onClick={disableAction}
											disabled={UserStore?.me?.id !== 1}
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

			<Modal show={false} onHide={() => {}} backdrop="static">
				<Modal.Header>
					<Modal.Title>安装插件</Modal.Title>
				</Modal.Header>
				<Modal.Body>
					<Form fluid onChange={() => {}} formValue={() => {}}>
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
			</Modal>
		</div>
	);
}
