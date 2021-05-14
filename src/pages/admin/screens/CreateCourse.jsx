import React from 'react';
import {
	Input,
	Row,
	Col,
	FlexboxGrid,
	Button,
	Notification,
	InputGroup,
	InputNumber,
	InputPicker,
	RadioGroup,
	Radio,
	Modal,
	Icon,
} from 'rsuite';
import useAxios from 'axios-hooks';
import axios from 'axios';

import { PageHead } from '../components';

const { useState, useEffect } = React;

function WeekTimeSelector(props) {
	const { data, style, onDataChange, operation, onOperationChange } = props;
	const [selectorData, setselectorData] = useState([]);

	useEffect(() => {
		let newData;
		switch (operation) {
			case 'all':
				newData = data?.map((item, index) => {
					return {
						index,
						value: item,
						select: true,
					};
				});
				break;
			case 'singular':
				newData = data?.map((item, index) => {
					return {
						index,
						value: item,
						select: !((index + 1) % 2 == 0),
					};
				});
				break;
			case 'even':
				newData = data?.map((item, index) => {
					return {
						index,
						value: item,
						select: (index + 1) % 2 == 0,
					};
				});
				break;
			default:
				newData = data?.map((item, index) => {
					return {
						index,
						value: item,
						select: false,
					};
				});
				break;
		}
		setselectorData(newData);
		onDataChange(newData);
	}, [operation]);

	const onClickSelector = (index) => {
		const newData = [...selectorData];
		const oldItem = newData[index];
		newData[index] = {
			...oldItem,
			select: !oldItem?.select,
		};
		// console.log("数据", index, newData);
		setselectorData(newData);
		onDataChange(newData);
		onOperationChange('null');
	};

	return (
		<Row style={style}>
			{selectorData?.map((item, index) => (
				<Col key={index} xs={1} justify="center" style={{ marginBottom: 20, marginRight: 5 }}>
					{item?.select ? (
						<Button
							style={{ color: '#FFF' }}
							appearance="primary"
							onClick={() => onClickSelector(item?.index)}
						>
							{item?.value}
						</Button>
					) : (
						<Button appearance="ghost" onClick={() => onClickSelector(item?.index)}>
							{item?.value}
						</Button>
					)}
				</Col>
			))}
		</Row>
	);
}

export default function CreateCourse(props) {
	const { toPage, params = {} } = props;
	const { id = undefined } = params;

	const [title, settitle] = useState();
	const [classroom, setclassroom] = useState();
	const [classroomId, setclassroomId] = useState();
	const [teacher, setteacher] = useState();
	const [remarks, setremarks] = useState();
	const [weekTime, setweekTime] = useState(1);
	const [lessonSerial, setlessonSerial] = useState(1);
	const [cycle, setcycle] = useState([]);
	const [timeId, settimeId] = useState();

	const [modalShow, setmodalShow] = useState(false);
	const [createLoading, setcreateLoading] = useState(false);

	const [courseData, setcourseData] = useState();
	useEffect(() => {
		if (id) {
			// console.log('课程ID', id);
			axios
				.get('/api/course/get?id=' + id)
				.then((res) => {
					const { data } = res;
					const { code, data: cData, msg } = data;

					if (code != 1 || !cData) {
						Notification.error({
							title: msg || '课程数据获取失败，请稍后重试！',
						});
						return;
					}

					setcourseData(cData);

					settitle(cData?.title);
					setclassroom(cData?.classroom);
					setclassroomId(cData?.classroom_id);
					setteacher(cData?.teacher);
					setremarks(cData?.remark);
					setweekTime(cData?.week_time);
					setlessonSerial(cData?.lesson_serial);
					setcycle(cData?.cycle);
					settimeId(cData?.time?.id + '');

					// 处理上课周数选择
					// TODO: 目前这里只是简单处理一下上课周数处理，以后需要做更加适配的数据对比处理
					if (cData?.cycle?.length == 20) {
						setwtOperation('all');
						setwtOperationGroup('all');
					} else if (cData?.cycle[1] == 3) {
						setwtOperation('singular');
						setwtOperationGroup('singular');
					} else if (cData?.cycle[1] == 4) {
						setwtOperation('even');
						setwtOperationGroup('even');
					}

					// console.log('课程数据', cData);
				})
				.catch((error) => {
					Notification.error({
						title: '课程数据获取失败，' + error || '课程数据获取失败，请稍后重试！',
					});
				});
		}
	}, [id]);

	function goBack() {
		toPage('CourseList');
	}

	const [{ data: timeListData, loading: timeListLoading, error: timeListError }, timeListRefetch] = useAxios({
		url: '/api/time/list?count=999',
	});
	const [timeList, settimeList] = useState([]);
	useEffect(() => {
		const list =
			timeListData?.data?.map((item, index) => {
				return {
					label: `${item.group} - ${item.remark}`,
					value: `${item.id}`,
					role: item.group,
				};
			}) || [];
		if (list.length > 0) {
			settimeList(list);
			const defTimeItem = list[0];
			settimeId(defTimeItem?.value);
		}
	}, [timeListData]);

	const RequestCreateCourse = (course) => {
		const { title, classroom, classroomId, teacher, remarks, weekTime, lessonSerial, cycle, timeId } = course;

		if (
			!title ||
			!classroom ||
			!classroomId ||
			!teacher ||
			!weekTime ||
			!lessonSerial ||
			!timeId ||
			cycle.length < 1
		) {
			Notification.error({
				title: '课程内容不完整！',
			});
			return;
		}

		setcreateLoading(true);

		const params = new URLSearchParams();
		params.append('title', title);
		params.append('classroom', classroom);
		params.append('classroom_id', classroomId);
		params.append('teacher', teacher);
		params.append('remarks', remarks || '');
		params.append('week_time', weekTime);
		params.append('lesson_serial', lessonSerial);
		params.append('cycle', JSON.stringify(cycle));
		params.append('time_id', timeId);

		axios
			.post('/api/course/create', params)
			.then((res) => {
				setcreateLoading(false);

				const { data: resData } = res;
				const { code, msg, data } = resData;

				if (code < 1) {
					Notification.error({
						title: msg || '添加失败，请稍后重试！',
					});
				} else {
					setmodalShow(true);
					setcreateLoading(false);
				}
			})
			.catch((error) => {
				Notification.error({
					title: '添加失败，' + error || '添加失败，请稍后重试！',
				});
				setcreateLoading(false);
			});
	};

	const RequestUpdateCourse = (course) => {
		const { id, title, classroom, classroomId, teacher, remarks, weekTime, lessonSerial, cycle, timeId } = course;

		if (
			!id ||
			!title ||
			!classroom ||
			!classroomId ||
			!teacher ||
			!weekTime ||
			!lessonSerial ||
			!timeId ||
			cycle.length < 1
		) {
			Notification.error({
				title: '课程数据不完整！',
			});
			return;
		}

		setcreateLoading(true);

		const params = new URLSearchParams();
		params.append('id', id);
		params.append('title', title);
		params.append('classroom', classroom);
		params.append('classroom_id', classroomId);
		params.append('teacher', teacher);
		params.append('remarks', remarks || '');
		params.append('week_time', weekTime);
		params.append('lesson_serial', lessonSerial);
		params.append('cycle', JSON.stringify(cycle));
		params.append('time_id', timeId);

		axios
			.post('/api/course/update', params)
			.then((res) => {
				setcreateLoading(false);

				const { data: resData } = res;
				const { code, msg, data } = resData;

				if (code < 1) {
					Notification.error({
						title: msg || '课程数据修改失败，请稍后重试！',
					});
				} else {
					Notification.success({
						title: '课程数据修改成功！',
					});
					setcreateLoading(false);
				}
			})
			.catch((error) => {
				Notification.error({
					title: '课程数据修改失败，' + error || '课程数据修改失败，请稍后重试！',
				});
				setcreateLoading(false);
			});
	};

	const lessonSerialInpRef = React.createRef();
	const lsHandleMinus = () => {
		lessonSerialInpRef.current.handleMinus();
	};
	const lsHandlePlus = () => {
		lessonSerialInpRef.current.handlePlus();
	};

	const weekTimeInpRef = React.createRef();
	const wtHandleMinus = () => {
		weekTimeInpRef.current.handleMinus();
	};
	const wtHandlePlus = () => {
		weekTimeInpRef.current.handlePlus();
	};

	const [wtOperation, setwtOperation] = useState('null');
	const [wtOperationGroup, setwtOperationGroup] = useState('null');

	const emptyInput = () => {
		settitle('');
		setclassroom('');
		setclassroomId('');
		setteacher('');
		setremarks('');
		setweekTime(1);
		setlessonSerial(1);
		setcycle([]);
		settimeId();
		setwtOperation('null');
		setwtOperationGroup('null');
	};

	return (
		<div>
			<PageHead title={id ? '编辑课程' : '添加课程'} onLeftClick={() => goBack()} />

			<Row>
				<Col xs={6}>
					<p style={{ marginBottom: 15 }}>课程名称</p>
					<Input
						defaultValue={title}
						value={title}
						placeholder="例如：Java 程序设计"
						onChange={(value) => settitle(value)}
					/>
				</Col>
				<Col xs={2} />
				<Col xs={16}>
					<p style={{ marginBottom: 15 }}>课程备注</p>
					<Input
						defaultValue={remarks}
						value={remarks}
						placeholder="例如：请同学们带好《Java 程序设计》书，以及笔和作业本"
						onChange={(value) => setremarks(value)}
					/>
				</Col>
			</Row>

			<Row style={{ marginTop: 40 }}>
				<Col xs={4}>
					<p style={{ marginBottom: 15 }}>教学楼</p>
					<Input
						defaultValue={classroom}
						value={classroom}
						placeholder="例如：知行楼"
						onChange={(value) => setclassroom(value)}
					/>
				</Col>
				<Col xs={1} />
				<Col xs={4}>
					<p style={{ marginBottom: 15 }}>教室</p>
					<Input
						defaultValue={classroomId}
						value={classroomId}
						placeholder="例如：406"
						onChange={(value) => setclassroomId(value)}
					/>
				</Col>

				<Col xs={1} />
				<Col xs={4}>
					<p style={{ marginBottom: 15 }}>教师名称</p>
					<Input
						defaultValue={teacher}
						value={teacher}
						placeholder="例如：刘老师"
						onChange={(value) => setteacher(value)}
					/>
				</Col>
			</Row>

			<Row style={{ marginTop: 40 }}>
				<Col xs={3}>
					<p style={{ marginBottom: 15 }}>星期几</p>
					<InputGroup>
						<InputGroup.Button onClick={wtHandleMinus}>-</InputGroup.Button>
						<InputNumber
							className={'custom-input-number'}
							ref={weekTimeInpRef}
							defaultValue={weekTime}
							value={weekTime}
							max={7}
							min={1}
							onChange={(value) => setweekTime(value)}
						/>
						<InputGroup.Button onClick={wtHandlePlus}>+</InputGroup.Button>
					</InputGroup>
				</Col>
				<Col xs={1} />
				<Col xs={3}>
					<p style={{ marginBottom: 15 }}>当天第几节课</p>
					<InputGroup>
						<InputGroup.Button onClick={lsHandleMinus}>-</InputGroup.Button>
						<InputNumber
							className={'custom-input-number'}
							ref={lessonSerialInpRef}
							defaultValue={lessonSerial}
							value={lessonSerial}
							max={20}
							min={1}
							onChange={(value) => setlessonSerial(value)}
						/>
						<InputGroup.Button onClick={lsHandlePlus}>+</InputGroup.Button>
					</InputGroup>
				</Col>

				<Col xs={1} />
				<Col xs={6}>
					<p style={{ marginBottom: 15 }}>上课时间</p>
					<InputPicker
						data={timeList}
						groupBy="role"
						defaultValue={timeId}
						value={timeId}
						onChange={(value) => {
							settimeId(value);
						}}
					/>
				</Col>
			</Row>

			<Col style={{ marginTop: 40 }}>
				<p style={{ marginBottom: 20 }}>上课周数</p>

				<WeekTimeSelector
					data={[1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20]}
					operation={wtOperation}
					onDataChange={(data) => {
						const dataArray = [];
						data?.map((item, index) => {
							if (!item.select) {
								return;
							}
							dataArray.push(item.value);
						});
						setcycle(dataArray);
						// console.log("返回数据", JSON.stringify(dataArray));
					}}
					onOperationChange={(value) => {
						setwtOperationGroup(value);
					}}
				/>

				<Col>
					<RadioGroup
						name="radioList"
						inline
						appearance="picker"
						value={wtOperationGroup}
						onChange={(value) => {
							if (value == wtOperation) {
								setwtOperation('null');
								setwtOperationGroup('null');
								return;
							}
							setwtOperation(value);
							setwtOperationGroup(value);
						}}
					>
						<Radio value="all">全选</Radio>
						<Radio value="singular">单周</Radio>
						<Radio value="even">双周</Radio>
					</RadioGroup>
				</Col>
			</Col>

			<FlexboxGrid style={{ marginTop: 80, marginBottom: 20 }} justify="end">
				<Button style={{ marginRight: 20 }} appearance="ghost" onClick={goBack}>
					{id ? '取消保存' : '取消添加'}
				</Button>
				<Button
					style={{ color: '#FFF' }}
					appearance="primary"
					loading={createLoading}
					onClick={() => {
						if (id) {
							RequestUpdateCourse({
								id,
								title,
								classroom,
								classroomId,
								teacher,
								remarks,
								weekTime,
								lessonSerial,
								cycle,
								timeId,
							});
						} else {
							RequestCreateCourse({
								title,
								classroom,
								classroomId,
								teacher,
								remarks,
								weekTime,
								lessonSerial,
								cycle,
								timeId,
							});
						}
					}}
				>
					{id ? '保存课程' : '添加课程'}
				</Button>
			</FlexboxGrid>

			<Modal backdrop="static" show={modalShow} onHide={() => setmodalShow(false)} size="xs">
				<Modal.Body style={{ marginTop: 10 }}>
					<Icon
						icon="ok-circle"
						style={{
							color: '#81bf63',
							marginRight: 10,
						}}
					/>
					课程添加成功！
				</Modal.Body>
				<Modal.Footer>
					<Button size="xs" onClick={() => setmodalShow(false)} appearance="default">
						继续添加
					</Button>
					<Button
						size="xs"
						onClick={() => {
							emptyInput();
							setmodalShow(false);
						}}
						appearance="default"
					>
						添加新的
					</Button>
					<Button
						size="xs"
						onClick={() => {
							setmodalShow(false);
							goBack();
						}}
						appearance="default"
					>
						返回列表
					</Button>
				</Modal.Footer>
			</Modal>
		</div>
	);
}
