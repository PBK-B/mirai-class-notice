/*
 * @Author: Bin
 * @Date: 2021-09-01
 * @FilePath: /class_notice/src/pages/admin/stores/UserStore.js
 */
import { makeAutoObservable } from 'mobx';

class UserStore {
	me = {};

	constructor() {
		makeAutoObservable(this);
	}

	setMe(me) {
		this.me = me;
	}
}

export default new UserStore();
