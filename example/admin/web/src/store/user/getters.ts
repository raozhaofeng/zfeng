import { GetterTree } from 'vuex';
import { StateInterface } from '../index';
import { UserStateInterface } from './state';

const getters: GetterTree<UserStateInterface, StateInterface> = {
    // 判断是否登陆
    isLogin(state: UserStateInterface) {
        return state.token.length > 0;
    }
};

export default getters;
