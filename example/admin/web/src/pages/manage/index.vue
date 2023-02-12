<template>
    <div class='q-pa-md'>
        <!--        请求参数-->
        <div class='row q-gutter-sm items-center'>
            <div v-for='(search, searchIndex) in conf.search' :key='searchIndex'>
                <dynamic-input v-model="params[search.field]" :label='search.label' :type='search.type'
                               :data='search.data'></dynamic-input>
            </div>
            <div v-if='conf.search.length > 0'>
                <q-btn icon='search' color="primary"
                       @click='requestTableFunc({ pagination: params.pagination })'></q-btn>
            </div>
        </div>

        <!--        数据表格-->
        <q-card flat bordered class="q-mt-md">
            <!--            数据表格工具栏-->
            <q-card-section>
                <div class='row'>
                    <div v-for='(tool, toolIndex) in conf.tools' :key='toolIndex' class='q-mr-sm q-mb-sm'>
                        <q-btn :label='tool.label' :color="tool.color"
                               v-if='editBtnEvalIsShowFunc(tool.isShow, {scope: {}}) && isAuth(tool.url)'
                               @click='showDialogFormFunc(tool)'></q-btn>
                    </div>
                </div>
            </q-card-section>
            <q-card-section>
                <q-table flat :rows='rows' :columns='conf.columns' :row-key='conf.basic.updateKey' selection='multiple'
                         @request='requestTableFunc' v-model:selected='checkboxList'
                         :visible-columns="conf.visibleColumns"
                         v-model:pagination='params.pagination'>
                    <template v-slot:top>
                        <q-space/>
                        <q-btn color="secondary" class="q-mr-sm" icon-right="archive" label="下载CSV文件" no-caps
                               @click="exportFileFunc"/>
                        <q-select v-model="conf.visibleColumns" multiple outlined dense options-dense
                                  @update:model-value="updateVisibleColumnsFunc"
                                  :display-value="$q.lang.table.columns" emit-value map-options :options="conf.columns"
                                  option-value="name" options-cover style="min-width: 150px"/>
                    </template>
                    <template v-slot:body-cell='scope'>
                        <q-td :style="scope.col.type !== 'options' ? {maxWidth: '160px'} : {}">
                            <table-body-cell :data="scope" :update-key="conf.basic.updateKey" @refresh="refreshFunc"
                                             :extend-values="extendValues" :dynamic-data="dynamicData"
                                             :update-url="conf.basic.updateURL"></table-body-cell>
                        </q-td>
                    </template>
                </q-table>
            </q-card-section>
        </q-card>
    </div>

    <dialog-form ref='dialogFormRef' :title='dialogForm.label' :items='dialogForm.items' :type="dialogForm.type"
                 :dynamic-data='dialogForm.dynamicData' :url="dialogForm.url" :after-func='dialogForm.afterFunc'
                 :values='dialogForm.values'></dialog-form>
</template>

<script lang="ts">
import {api} from 'boot/axios';
import {onMounted, reactive, ref, toRefs} from 'vue'
import DialogForm from 'components/dialogForm.vue';
import DynamicInput from 'components/input.vue'
import TableBodyCell from 'components/tableBodyCell.vue';
import {editBtnEvalIsShowFunc, editBtnTypeObject} from 'src/utils/define';
import {confirmBoxDialog} from 'src/utils';
import {exportCSVFile} from 'src/utils/export';
import {isAuth} from 'src/utils/auth';

export default {
    name: 'ManageIndex',
    components: {DynamicInput, DialogForm, TableBodyCell},
    setup() {
        const visibleStorage = '_table_visible_001'
        const dialogFormRef = ref(null) as any;
        const state = reactive({
            params: {
                pagination: {sortBy: 'id', descending: true, page: 0, rowsPerPage: 10}
            } as any,
            dialogForm: {} as any,
            rows: [] as any,
            dynamicData: {} as any,
            extendValues: {} as any,
            checkboxList: [] as any,
            conf: {
                basic: {indexURL: '/manage/index', updateURL: '/manage/update', updateKey: 'id'},
                search: [{label: '上级用户', field: 'parent_name', type: 'text'}, {
                    label: '域名',
                    field: 'domain',
                    type: 'text'
                }, {
                    label: '用户名',
                    field: 'username',
                    type: 'text'
                }, {
                    label: '邮箱',
                    field: 'email',
                    type: 'text'
                }, {label: '昵称', field: 'nickname', type: 'text'}, {
                    label: '状态',
                    field: 'status',
                    type: 'select',
                    data: [{label: '激活', value: 10}, {label: '禁用', value: -1}]
                }, {label: '时间', field: 'updated_at', type: 'dateRangePicker'}],
                tools: [{
                    label: '新增',
                    url: '/manage/create',
                    color: 'primary',
                    type: 'create',
                    tips: '',
                    isShow: '',
                    items: [{label: '用户名', field: 'username', type: 'text'}, {
                        label: '角色',
                        field: 'roles',
                        type: 'checkbox',
                        data: []
                    }, {label: '密码', field: 'password', type: 'password'}, {
                        label: '安全密码',
                        field: 'security_key',
                        type: 'password'
                    }, {label: '域名(多个用 "," 分割)', field: 'domain', type: 'text'}, {
                        label: '过期时间',
                        field: 'expired_at',
                        type: 'datePicker'
                    }, {
                        label: 'Token设置', field: 'data', type: 'json',
                        data: [[{label: '加密Key', field: 'key', type: 'text'}, {
                            label: '登陆唯一',
                            field: 'only',
                            type: 'select',
                            data: [{label: '唯一登陆', value: true}, {label: '无限登陆', value: false}]
                        }, {label: '过期时间(s)', field: 'expire', type: 'number'}], [{
                            label: '白名单(用逗号分割)',
                            field: 'whitelist',
                            type: 'textarea'
                        }], [{label: '黑名单(用逗号分割)', field: 'blacklist', type: 'textarea'}]]
                    }]
                }, {
                    label: '删除',
                    url: '/manage/delete',
                    color: 'negative',
                    type: 'checkboxDelete',
                    tips: '',
                    isShow: '',
                    items: []
                }],
                visibleColumns: localStorage.getItem(visibleStorage) ?
                    JSON.parse(<string>localStorage.getItem(visibleStorage)) :
                    ['id', 'avatar', 'parent_name', 'domain', 'username', 'invite_code', 'role', 'email', 'nickname', 'money', 'status', 'expired_at', 'updated_at', 'options'],
                columns: [{
                    label: '主键',
                    field: 'id',
                    name: 'id',
                    type: 'text',
                    align: 'left',
                    sortable: true
                }, {
                    label: '头像',
                    field: 'avatar',
                    name: 'avatar',
                    type: 'image',
                    align: 'left',
                    sortable: true
                }, {
                    label: '上级用户',
                    field: 'parent_name',
                    name: 'parent_name',
                    type: 'text',
                    align: 'left',
                    sortable: false
                }, {
                    label: '域名',
                    field: 'domain',
                    name: 'domain',
                    type: 'editTextarea',
                    align: 'left',
                    sortable: true,
                }, {
                    label: '用户名',
                    field: 'username',
                    name: 'username',
                    type: 'text',
                    align: 'left',
                    sortable: true
                }, {
                    label: '邀请码',
                    field: 'invite_code',
                    name: 'invite_code',
                    type: 'text',
                    align: 'left',
                    sortable: false
                }, {
                    label: '角色',
                    field: 'role',
                    name: 'role',
                    type: 'text',
                    align: 'left',
                    sortable: true
                }, {
                    label: '邮箱',
                    field: 'email',
                    name: 'email',
                    type: 'editText',
                    align: 'left',
                    sortable: true
                }, {
                    label: '昵称',
                    field: 'nickname',
                    name: 'nickname',
                    type: 'editText',
                    align: 'left',
                    sortable: true
                }, {
                    label: '金额',
                    field: 'money',
                    name: 'money',
                    type: 'text',
                    align: 'left',
                    sortable: true
                }, {
                    label: '状态',
                    field: 'status',
                    name: 'status',
                    type: 'editToggle',
                    align: 'left',
                    sortable: true,
                    data: [{label: '激活', value: 10}, {label: '禁用', value: -1}]
                }, {
                    label: '过期时间',
                    field: 'expired_at',
                    name: 'expired_at',
                    type: 'datePicker',
                    align: 'left',
                    sortable: true
                }, {
                    label: '更新时间',
                    field: 'updated_at',
                    name: 'updated_at',
                    type: 'datePicker',
                    align: 'left',
                    sortable: true
                }, {
                    label: '操作',
                    field: '',
                    name: 'options',
                    type: 'options',
                    align: 'left',
                    items: [{
                        label: '编辑',
                        url: '/manage/update',
                        color: 'secondary',
                        type: 'update',
                        tips: '',
                        isShow: '',
                        items: [{label: '头像', field: 'avatar', type: 'image'}, {
                            label: '角色',
                            field: 'roles',
                            type: 'checkbox',
                            data: []
                        }, {
                            label: '状态',
                            field: 'status',
                            type: 'select',
                            data: [{label: '激活', value: 10}, {label: '禁用', value: -1}]
                        }, {
                            label: '域名(多个用 "," 分割)',
                            field: 'domain',
                            type: 'text'
                        }, {label: '邮箱', field: 'email', type: 'text'}, {
                            label: '昵称',
                            field: 'nickname',
                            type: 'text'
                        }, {label: '密码', field: 'password', type: 'text'}, {
                            label: '安全密码',
                            field: 'security_key',
                            type: 'text'
                        }, {label: '过期时间', field: 'expired_at', type: 'datePicker'}, {
                            label: 'Token设置', field: 'data', type: 'json',
                            data: [[{label: '加密Key', field: 'key', type: 'text'}, {
                                label: '登陆唯一',
                                field: 'only',
                                type: 'select',
                                data: [{label: '唯一登陆', value: true}, {label: '无限登陆', value: false}]
                            }, {label: '过期时间(s)', field: 'expire', type: 'number'}], [{
                                label: '白名单(用逗号分割)',
                                field: 'whitelist',
                                type: 'textarea'
                            }], [{label: '黑名单(用逗号分割)', field: 'blacklist', type: 'textarea'}]]
                        }]
                    }, {
                        label: '删除',
                        url: '/manage/delete',
                        color: 'negative',
                        type: 'delete',
                        tips: '',
                        isShow: '',
                        items: []
                    }]
                }]
            }
        })

        onMounted(() => {
            // 请求管理员角色列表
            api.post('/role/roles', {}, {showLoading: false} as any).then((res: any) => {
                state.dynamicData['roles'] = res
            })

            requestTableFunc({pagination: state.params.pagination});
        })

        // 请求数据表格
        const requestTableFunc = (props: { pagination: any }) => {
            state.params.pagination = props.pagination;
            api.post(state.conf.basic.indexURL, state.params).then((res: any) => {
                state.rows = res.items;
                state.params.pagination.rowsNumber = res.count;
            });
        };

        // 刷新当前页面
        const refreshFunc = () => {
            requestTableFunc({pagination: state.params.pagination});
        }

        // 显示模态框方法
        const showDialogFormFunc = (conf: any) => {
            switch (conf.type) {
                case editBtnTypeObject.checkboxDelete:
                    const params = {} as any
                    params[state.conf.basic.updateKey] = []
                    for (let i = 0; i < state.checkboxList.length; i++) {
                        params[state.conf.basic.updateKey].push(state.checkboxList[i][state.conf.basic.updateKey])
                    }
                    confirmDeleteFunc(conf.url, '是否确认删除【' + params[state.conf.basic.updateKey].toString() + '】', params)
                    break;
                case editBtnTypeObject.create:
                    state.dialogForm = conf
                    state.dialogForm.afterFunc = refreshFunc
                    state.dialogForm.dynamicData = state.dynamicData
                    dialogFormRef.value.isShow = true
                    break;
            }
        }

        // 确认删除数据
        const confirmDeleteFunc = (url: string, msg: string, params: any) => {
            confirmBoxDialog('删除数据', msg, () => {
                api.post(url, params).then(() => {
                    refreshFunc()
                })
            })
        }

        // 更新字段显示内容
        const updateVisibleColumnsFunc = () => {
            localStorage.setItem(visibleStorage, JSON.stringify(state.conf.visibleColumns))
        }

        // 下载csv文件
        const exportFileFunc = () => {
            exportCSVFile(state.conf.columns, state.conf.visibleColumns, state.rows)
        }

        return {
            isAuth,
            exportFileFunc,
            updateVisibleColumnsFunc,
            dialogFormRef,
            refreshFunc,
            showDialogFormFunc,
            editBtnEvalIsShowFunc,
            requestTableFunc,
            ...toRefs(state)
        }
    }
}
</script>

<style scoped>

</style>
