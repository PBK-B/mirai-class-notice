import React from 'react';
import useAxios from 'axios-hooks';
import axios from 'axios';

import { Table, Panel, Loader, FlexboxGrid, Button, Alert, Tag, Pagination, Modal, Notification } from 'rsuite';

const { useState, useEffect, useRef } = React;
const { Column, HeaderCell, Cell } = Table;

export default function CourseList(props) {
    const { toPage, params = {} } = props;

    const [activePage, setactivePage] = useState(1);
    const [courses, setcourses] = useState([]);

    const [{ data, loading, error }, refetch] = useAxios({
        url: '/api/course/list?count=999',
    });

    // 禁用课程
    const APIDisableCourse = (id, isDisable) => {
        const msgStrTitle = isDisable ? '禁用' : '启用';

        if (!id) {
            Notification.error({
                title: msgStrTitle + '失败，课程数据异常！',
            });
            return;
        }

        const params = new URLSearchParams();
        params.append('id', id);

        axios
            .post('/api/course/upstatus', params, {})
            .then((res) => {
                const { data } = res;
                const { code } = data;

                if (code < 1) {
                    Notification.error({
                        title: data?.msg || msgStrTitle + '失败，请稍后重试！',
                    });
                } else {
                    const course = data?.data;

                    Notification.success({
                        title: `${msgStrTitle}课程通知 ${course.title} 成功！`,
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

    // 删除课程
    const [showDeleteModal, setshowDeleteModal] = useState(false);
    const deleteCourseData = useRef();
    const APIDeleteCourse = (id) => {
        const msgStrTitle = '删除';

        if (!id) {
            Notification.error({
                title: msgStrTitle + '失败，课程数据异常！',
            });
            return;
        }

        const params = new URLSearchParams();
        params.append('id', id);

        axios
            .post('/api/course/delete', params, {})
            .then((res) => {
                const { data } = res;
                const { code } = data;

                if (code < 1) {
                    Notification.error({
                        title: data?.msg || msgStrTitle + '失败，请稍后重试！',
                    });
                } else {
                    const course = data?.data;

                    Notification.success({
                        title: `${msgStrTitle}课程通知 ${course.title} 成功！`,
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

    const [coursesList, setcoursesList] = useState([]);
    useEffect(() => {
        if (data) {
            let { data: data_array } = data;
            const week_int = ['零', '一', '二', '三', '四', '五', '六', '七', '八', '九', '十'];
            data_array = data_array?.map((item, index) => {
                return {
                    ...item,
                    time: new Date().toUTCString(),
                    time_str: item.time.start,
                    week_time: '星期' + week_int[item.week_time],
                };
            });
            // setcourses(data_array);
            setcoursesList(data_array);
            onChangePage(data_array, 1);
        }
    }, [data]);

    const lengthMenu = [{ label: <p>10</p>, value: 10 }];
    const [displayLength, setdisplayLength] = useState(lengthMenu[0].value);

    const onChangePage = (coursesList, value) => {
        setactivePage(value);
        const index = displayLength * (value - 1);
        const newArray = coursesList?.slice(index, index + displayLength);
        setcourses(newArray);
    };

    if (loading) return <Loader backdrop content="loading..." vertical />;

    return (
        <div style={{ marginTop: 25, marginBottom: 25 }}>
            <p>全部课程管理页面</p>
            <br />

            <FlexboxGrid style={{ marginBottom: 15 }} justify="end">
                <FlexboxGrid.Item style={{ marginRight: 15 }}>
                    <Button
                        appearance="ghost"
                        onClick={() => {
                            Alert.warning('导入课表功能正在开发…');
                        }}
                    >
                        导入课表
                    </Button>
                </FlexboxGrid.Item>
                <FlexboxGrid.Item>
                    <Button appearance="ghost" onClick={() => toPage('CreateCourse')}>
                        添加课程
                    </Button>
                </FlexboxGrid.Item>
            </FlexboxGrid>

            <Panel bordered bodyFill>
                <Table
                    autoHeight
                    data={courses}
                    onRowClick={(data) => {
                        console.log(data);
                    }}
                >
                    <Column width={70} align="center" fixed>
                        <HeaderCell>ID</HeaderCell>
                        <Cell dataKey="id" />
                    </Column>

                    <Column width={200} fixed>
                        <HeaderCell>课程名称</HeaderCell>
                        <Cell dataKey="title" />
                    </Column>

                    <Column width={80}>
                        <HeaderCell>星期</HeaderCell>
                        <Cell dataKey="week_time" />
                    </Column>

                    <Column width={160}>
                        <HeaderCell>上课时间</HeaderCell>
                        <Cell dataKey="time_str" />
                    </Column>

                    <Column width={120}>
                        <HeaderCell>授课老师</HeaderCell>
                        <Cell dataKey="teacher" />
                    </Column>

                    <Column width={260}>
                        <HeaderCell>备注</HeaderCell>
                        <Cell dataKey="remark" />
                    </Column>

                    <Column width={90}>
                        <HeaderCell>状态</HeaderCell>
                        <Cell dataKey="status">
                            {(rowData) =>
                                rowData.status === 0 ? <Tag color="red">禁用</Tag> : <Tag color="cyan">启用</Tag>
                            }
                        </Cell>
                    </Column>

                    <Column width={160} fixed="right">
                        <HeaderCell>操作</HeaderCell>

                        <Cell>
                            {(rowData) => {
                                function handleAction(e) {
                                    alert(`id:${rowData.id}`);
                                    e.stopPropagation();
                                }

                                function goToEdit(e) {
                                    toPage('CreateCourse', { id: rowData.id });
                                    e.stopPropagation();
                                }

                                function onDisable(e) {
                                    APIDisableCourse(rowData.id, rowData.status != 0 ? true : false);
                                    e.stopPropagation();
                                }

                                function onDelete(e) {
                                    deleteCourseData.current = { ...rowData };
                                    setshowDeleteModal(true);
                                    e.stopPropagation();
                                }

                                return (
                                    <span>
                                        <a onClick={goToEdit}> 编辑 </a> |
                                        <a onClick={onDisable}>{rowData.status != 0 ? ' 禁用 ' : ' 启用 '}</a> |
                                        <a onClick={onDelete}> 删除 </a>
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
                    total={coursesList?.length || 0}
                    onChangeLength={(value) => {
                        // console.log("改变每页长度", value);
                        // setdisplayLength(value);
                        // const newArray = coursesList.splice(value);
                        // setcourses(newArray);
                    }}
                    onChangePage={(value) => {
                        // console.log("改变的页面", value);
                        onChangePage(coursesList, value);
                    }}
                ></Table.Pagination>
            </Panel>

            <Modal
                size="xs"
                show={showDeleteModal}
                onHide={() => {
                    setshowDeleteModal(false);
                }}
            >
                <Modal.Header>
                    <Modal.Title>友情提示</Modal.Title>
                </Modal.Header>
                <Modal.Body>
                    <p>你确定要删除该课程通知吗？删除后不可恢复！！！</p>
                </Modal.Body>
                <Modal.Footer>
                    <Button
                        onClick={() => {
                            APIDeleteCourse(deleteCourseData.current?.id);
                            setshowDeleteModal(false);
                        }}
                        color="red"
                        style={{ color: '#FFF' }}
                        appearance="primary"
                    >
                        确定删除
                    </Button>
                    <Button
                        onClick={() => {
                            setshowDeleteModal(false);
                        }}
                        appearance="subtle"
                    >
                        取消
                    </Button>
                </Modal.Footer>
            </Modal>
        </div>
    );
}
