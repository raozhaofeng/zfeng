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
    name: 'UserIndex',
    components: {DynamicInput, DialogForm, TableBodyCell},
    setup() {
        const visibleStorage = '_table_visible_021'
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
                basic: {indexURL: '/user/index', updateURL: '/user/update', updateKey: 'id'},
                search: [{label: '管理', field: 'admin_name', type: 'text'}, {
                    label: '邀请者',
                    field: 'parent_name',
                    type: 'text'
                }, {label: '国家', field: 'country_name', type: 'text'}, {
                    label: '用户',
                    field: 'username',
                    type: 'text'
                }, {label: '邮箱', field: 'email', type: 'text'}, {
                    label: '手机',
                    field: 'telephone',
                    type: 'text'
                }, {label: '昵称', field: 'nickname', type: 'text'}, {
                    label: '性别',
                    field: 'sex',
                    type: 'selectNumber',
                    data: [{label: '未知', value: '-1'}, {label: '男', value: '1'}, {label: '女', value: '2'}]
                }, {
                    label: '类型',
                    field: 'type',
                    type: 'selectNumber',
                    data: [{label: '虚拟用户', value: '-1'}, {label: '普通用户', value: '10'}]
                }, {
                    label: '状态',
                    field: 'status',
                    type: 'selectNumber',
                    data: [{label: '激活', value: '10'}, {label: '冻结', value: '1'}, {label: '禁用', value: '-1'}]
                }, {label: '生日范围', field: 'birthday', type: 'dateRangePicker'}, {
                    label: '时间',
                    field: 'updated_at',
                    type: 'dateRangePicker'
                }],
                tools: [{
                    label: '新增',
                    url: '/user/create',
                    color: 'primary',
                    type: 'create',
                    tips: '',
                    isShow: '',
                    items: [{label: '用户名', field: 'username', type: 'text'}, {
                        label: '密码',
                        field: 'password',
                        type: 'password'
                    }, {
                        label: '安全密钥',
                        field: 'security_key',
                        type: 'password',
                    },{
                        label: '类型',
                        field: 'type',
                        type: 'selectNumber',
                        data: [{label: '普通用户', value: '10'}, {label: '虚拟用户', value: '-1'}]
                    }]
                }, {
                    label: '加减款',
                    url: '/user/amount',
                    color: 'secondary',
                    type: 'create',
                    tips: '',
                    isShow: '',
                    items: [{
                        label: '类型',
                        field: 'type',
                        type: 'selectNumber',
                        data: [{label: '加款', value: '1'}, {label: '减款', value: '2'}]
                    }, {label: '用户名', field: 'username', type: 'text'}, {label: '金额', field: 'money', type: 'number'}]
                }, {
                    label: '删除',
                    url: '/user/delete',
                    color: 'negative',
                    type: 'checkboxDelete',
                    tips: '',
                    isShow: '',
                    items: []
                }],
                visibleColumns: localStorage.getItem(visibleStorage) ?
                    JSON.parse(<string>localStorage.getItem(visibleStorage)) :
                    ['id', 'admin_name', 'parent_name', 'country_name', 'username', 'invite_code', 'nickname', 'email', 'telephone', 'avatar', 'sex', 'birthday', 'money', 'freeze_money', 'type', 'status', 'ip4', 'updated_at', 'options'],
                columns: [{
                    label: 'ID',
                    field: 'id',
                    name: 'id',
                    type: 'text',
                    align: 'left',
                    sortable: true
                }, {
                    label: '管理',
                    field: 'admin_name',
                    name: 'admin_name',
                    type: 'text',
                    align: 'left',
                    sortable: false
                }, {
                    label: '邀请者',
                    field: 'parent_name',
                    name: 'parent_name',
                    type: 'text',
                    align: 'left',
                    sortable: false
                }, {
                    label: '国家',
                    field: 'country_name',
                    name: 'country_name',
                    type: 'text',
                    align: 'left',
                    sortable: false
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
                    label: '昵称',
                    field: 'nickname',
                    name: 'nickname',
                    type: 'editText',
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
                    label: '手机号码',
                    field: 'telephone',
                    name: 'telephone',
                    type: 'editText',
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
                    label: '性别',
                    field: 'sex',
                    name: 'sex',
                    type: 'selectNumber',
                    align: 'left',
                    sortable: true,
                    data: [{label: '未知', value: '-1'}, {label: '男', value: '1'}, {label: '女', value: '2'}]
                }, {
                    label: '生日',
                    field: 'birthday',
                    name: 'birthday',
                    type: 'datePicker',
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
                    label: '冻结金额',
                    field: 'freeze_money',
                    name: 'freeze_money',
                    type: 'text',
                    align: 'left',
                    sortable: true
                }, {
                    label: '类型',
                    field: 'type',
                    name: 'type',
                    type: 'selectNumber',
                    align: 'left',
                    sortable: true,
                    data: [{label: '虚拟用户', value: '-1'}, {label: '普通用户', value: '10'}]
                }, {
                    label: '状态',
                    field: 'status',
                    name: 'status',
                    type: 'selectNumber',
                    align: 'left',
                    sortable: true,
                    data: [{label: '禁用', value: '-1'}, {label: '冻结', value: '1'}, {label: '激活', value: '10'}]
                },{
                    label: 'IP4',
                    field: 'ip4',
                    name: 'ip4',
                    type: 'ip4',
                    align: 'left',
                    sortable: true
                }, {
                    label: '时间',
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
                        label: '更新',
                        url: '/user/update',
                        color: 'primary',
                        type: 'update',
                        tips: '',
                        isShow: '',
                        items: [{label: '头像', field: 'avatar', type: 'image'}, {
                            label: '邀请者ID',
                            field: 'parent_id',
                            type: 'number'
                        }, {label: '国家ID', field: 'country_id', type: 'number'}, {
                            label: '邮箱',
                            field: 'email',
                            type: 'text'
                        }, {label: '昵称', field: 'nickname', type: 'text'}, {
                            label: '手机号',
                            field: 'telephone',
                            type: 'text'
                        }, {
                            label: '性别',
                            field: 'sex',
                            type: 'selectNumber',
                            data: [{label: '未知', value: '-1'}, {label: '男', value: '1'}, {label: '女', value: '2'}]
                        }, {label: '生日', field: 'birthday', type: 'datePicker'}, {
                            label: '密码',
                            field: 'password',
                            type: 'text'
                        }, {label: '安全密钥', field: 'security_key', type: 'text'}, {
                            label: '状态',
                            field: 'status',
                            type: 'selectNumber',
                            data: [{label: '激活', value: '10'}, {label: '冻结', value: '1'}, {label: '禁用', value: '-1'}]
                        }]
                    }, {
                        label: '删除',
                        url: '/user/delete',
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