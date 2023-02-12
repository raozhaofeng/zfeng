<template>
    <q-layout view='hHh lpR fFf'>
        <q-header reveal bordered class='bg-primary text-white'>
            <q-toolbar>
                <q-btn dense flat round icon='menu' @click='toggleLeftDrawer'/>
                <div @click="$router.push('/')" class="row justify-start items-center cursor-pointer">
                    <q-avatar>
                        <q-img :src='logoImg' no-spinner width='30px' height='30px'/>
                    </q-avatar>
                    <div class="text-h6">BeeGo</div>
                </div>
                <q-toolbar-title></q-toolbar-title>
                <q-btn dense flat :label="'邀请码' + userInfo.invite_code"></q-btn>
                <q-btn dense flat :label='"剩余:" + remainingDays() + "天"'></q-btn>
                <q-btn-dropdown flat>
                    <template v-slot:label>
                        <q-avatar rounded size='30px'>
                            <q-img :src='imageSrc(userInfo.avatar)' width='30px' height='30px' no-spinner/>
                        </q-avatar>
                        <div class='q-ml-xs'>{{ userInfo.username }}</div>
                    </template>
                    <q-list class='bg-primary text-white text-body2'>
                        <q-item v-close-popup clickable @click="updateDialogFormFunc('info')">
                            <q-item-section>
                                <div class='row q-gutter-sm items-center'>
                                    <q-icon name='sym_o_contact_mail' size='sm'></q-icon>
                                    <div>用户信息</div>
                                </div>
                            </q-item-section>
                        </q-item>
                        <q-item v-close-popup clickable @click="updateDialogFormFunc('password')">
                            <q-item-section>
                                <div class='row q-gutter-sm items-center'>
                                    <q-icon name='sym_o_lock' size='sm'></q-icon>
                                    <div>登陆密码</div>
                                </div>
                            </q-item-section>
                        </q-item>
                        <q-item v-close-popup clickable @click="updateDialogFormFunc('security')">
                            <q-item-section>
                                <div class='row q-gutter-sm items-center'>
                                    <q-icon name='sym_o_security' size='sm'></q-icon>
                                    <div>安全密码</div>
                                </div>
                            </q-item-section>
                        </q-item>
                        <q-separator></q-separator>
                        <q-item v-close-popup clickable @click="adminLogout">
                            <q-item-section>
                                <div class='row q-gutter-sm items-center'>
                                    <q-icon name='sym_o_logout' size='sm'></q-icon>
                                    <div>退出登陆</div>
                                </div>
                            </q-item-section>
                        </q-item>
                    </q-list>
                </q-btn-dropdown>
                <q-btn dense flat round icon='sym_o_fullscreen' @click='$q.fullscreen.toggle()'></q-btn>
                <q-btn dense flat icon="sym_o_groups" :label="'在线客服'"
                       @click='toggleRightDrawer'>
                    <q-badge color="red" floating v-if="userInfo.unread_nums > 0">{{ userInfo.unread_nums }}</q-badge>
                </q-btn>
            </q-toolbar>
        </q-header>

        <q-drawer show-if-above v-model='leftDrawerOpen' side='left' bordered @mouseover='miniState = false'
                  @mouseout='miniState = true' :mini='miniState' mini-to-overlay :width='220'>
            <menu-list :data='menuListData'></menu-list>
            <div style='height: 100px'></div>
        </q-drawer>

        <q-drawer v-model='rightDrawerOpen' side='right' bordered>
            <conversation :drawer="rightDrawerOpen"></conversation>
        </q-drawer>

        <q-page-container>
            <router-view/>
        </q-page-container>
    </q-layout>

    <dialog-form ref='dialogFormRef' :title='dialogForm.title' :items='dialogForm.items' :type="'update'"
                 :dynamic-data='dialogForm.dynamicData' :url="dialogForm.url" :after-func='dialogForm.afterFunc'
                 :values='dialogForm.values'></dialog-form>
</template>

<script lang="ts">
import {imageSrc} from 'src/utils';
import store from 'src/store';
import MenuList from 'src/components/menu.vue';
import Conversation from 'src/components/conversation.vue';
import {ref, reactive, toRefs, watch} from 'vue';
import DialogForm from 'src/components/dialogForm.vue';
import {AdminInfoAPI} from 'src/api';
import router from 'src/router';

export default {
    name: 'MainLayout',
    components: {MenuList, Conversation, DialogForm},
    setup() {
        const leftDrawerOpen = ref(false);
        const rightDrawerOpen = ref(false);
        const dialogFormRef = ref(null) as any;

        const state = reactive({
            miniState: true,
            menuListData: store.state.user.menuList,
            logoImg: new URL('/icons/favicon.png', import.meta.url).href,
            userInfo: JSON.parse(JSON.stringify(store.state.user.info)),
            updateForm: {
                info: {
                    title: '更新用户信息',
                    url: '/update',
                    values: {} as any,
                    items: [
                        {
                            label: '头像',
                            field: 'avatar',
                            type: 'image'
                        }, {
                            label: '昵称',
                            field: 'nickname',
                            type: 'text'
                        }, {
                            label: '邮箱',
                            field: 'email',
                            type: 'text',
                        },
                    ]
                },
                password: {
                    title: '更新密码',
                    url: '/update/password',
                    items: [
                        {
                            label: '密码类型',
                            field: 'type',
                            type: 'select',
                            data: [
                                {label: '登陆密码', value: 1},
                                {label: '安全密码', value: 2}
                            ]
                        }, {
                            label: '旧密码',
                            field: 'old_password',
                            type: 'password'
                        }, {
                            label: '新密码',
                            field: 'new_password',
                            type: 'password'
                        }
                    ]
                }
            } as any,
            dialogForm: {} as any
        });

        // 更新资料方法
        const updateDialogFormFunc = (opt: any) => {
            switch (opt) {
                case 'info':
                    state.dialogForm = state.updateForm.info;
                    state.dialogForm.values = state.userInfo
                    // 设置执行完后回调方法
                    state.dialogForm.afterFunc = (params: any) => {
                        state.userInfo = Object.assign(JSON.parse(JSON.stringify(state.userInfo)), params)
                        store.commit('user/updateInfo', state.userInfo)
                    }
                    break;
                case 'password':
                    state.updateForm.password.values = {type: 1}
                    state.dialogForm = state.updateForm.password;
                    break;
                case 'security':
                    state.updateForm.password.values = {type: 2}
                    state.dialogForm = state.updateForm.password;
                    break;
            }
            dialogFormRef.value.isShow = true
        }

        // 计算剩余天数
        const remainingDays = () => {
            if (state.userInfo.expired_at <= 0) {
                return 365;
            }

            const nowTime = Date.parse(new Date) / 1000;
            return Math.floor((state.userInfo.expired_at - nowTime) / 86400);
        };

        //  监控在线人数, 未读消息数量
        watch(() => store.state.user.info, (val) => {
            state.userInfo = val
        })

        //  更新用户信息
        AdminInfoAPI().then((res: any) => {
            store.commit('user/updateInfo', res)
        })

        const toggleRightDrawer = () => {
            const userInfo = JSON.parse(JSON.stringify(store.state.user.info))
            userInfo.unread_nums = 0
            store.commit('user/updateInfo', userInfo)
            rightDrawerOpen.value = !rightDrawerOpen.value;
        }

        // 管理员退出
        const adminLogout = () => {
            localStorage.clear();
            store.commit('user/updateToken', '');
            setTimeout(() => {
                void router.push({name: 'Login'})
            }, 100)
        }

        return {
            ...toRefs(state),
            dialogFormRef,
            imageSrc,
            adminLogout,
            leftDrawerOpen,
            remainingDays,
            updateDialogFormFunc,
            toggleLeftDrawer() {
                leftDrawerOpen.value = !leftDrawerOpen.value;
            },
            rightDrawerOpen,
            toggleRightDrawer,
        };
    }
};
</script>
