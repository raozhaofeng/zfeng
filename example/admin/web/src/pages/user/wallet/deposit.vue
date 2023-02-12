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
    name: 'UserWalletDeposit',
    components: {DynamicInput, DialogForm, TableBodyCell},
    setup() {
        const visibleStorage = '_table_visible_017'
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
                basic: {indexURL: '/wallet/deposit/index', updateURL: '/wallet/deposit/update', updateKey: 'id'},
                search: [{label: '管理', field: 'admin_name', type: 'text'}, {
                    label: '用户',
                    field: 'username',
                    type: 'text'
                }, {label: '订单号', field: 'order_sn', type: 'text'}, {
                    label: '姓名｜Token',
                    field: 'payment_real_name',
                    type: 'text'
                }, {label: '卡号｜地址', field: 'payment_card_number', type: 'text'}, {
                    label: '状态',
                    field: 'status',
                    type: 'selectNumber',
                    data: [{label: '审核', value: '10'}, {label: '拒绝', value: '-1'}, {label: '通过', value: '20'}]
                }, {label: '数据', field: 'data', type: 'text'}, {
                    label: '时间',
                    field: 'created_at',
                    type: 'dateRangePicker'
                }],
                tools: [{
                    label: '新增',
                    url: '/wallet/deposit/create',
                    color: 'primary',
                    type: 'create',
                    tips: '',
                    isShow: '',
                    items: [{label: '用户名', field: 'username', type: 'text'}, {
                        label: '金额',
                        field: 'money',
                        type: 'number'
                    }, {label: '凭证', field: 'proof', type: 'image'}]
                }, {
                    label: '删除',
                    url: '/wallet/deposit/delete',
                    color: 'negative',
                    type: 'checkboxDelete',
                    tips: '',
                    isShow: '',
                    items: []
                }],
                visibleColumns: localStorage.getItem(visibleStorage) ?
                    JSON.parse(<string>localStorage.getItem(visibleStorage)) :
                    ['order_sn', 'admin_name', 'username', 'payment_real_name', 'payment_card_number', 'balance', 'money', 'status', 'proof', 'data', 'created_at', 'options'],
                columns: [{
                    label: '订单号',
                    field: 'order_sn',
                    name: 'order_sn',
                    type: 'text',
                    align: 'left',
                    sortable: true
                }, {
                    label: '管理',
                    field: 'admin_name',
                    name: 'admin_name',
                    type: 'text',
                    align: 'left',
                    sortable: true
                }, {
                    label: '用户',
                    field: 'username',
                    name: 'username',
                    type: 'text',
                    align: 'left',
                    sortable: true
                }, {
                    label: '姓名｜Token',
                    field: 'payment_real_name',
                    name: 'payment_real_name',
                    type: 'text',
                    align: 'left',
                    sortable: true
                }, {
                    label: '卡号｜地址',
                    field: 'payment_card_number',
                    name: 'payment_card_number',
                    type: 'text',
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
                    label: '余额',
                    field: 'balance',
                    name: 'balance',
                    type: 'text',
                    align: 'left',
                    sortable: true
                }, {
                    label: '状态',
                    field: 'status',
                    name: 'status',
                    type: 'selectNumber',
                    align: 'left',
                    sortable: true,
                    data: [{label: '拒绝', value: '-1'}, {label: '待处理', value: '10'}, {label: '通过', value: '20'}]
                }, {
                    label: '凭证',
                    field: 'proof',
                    name: 'proof',
                    type: 'image',
                    align: 'left',
                    sortable: true
                }, {
                    label: '数据',
                    field: 'data',
                    name: 'data',
                    type: 'editText',
                    align: 'left',
                    sortable: true
                }, {
                    label: '时间',
                    field: 'created_at',
                    name: 'created_at',
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
                        url: '/wallet/deposit/update',
                        color: 'primary',
                        type: 'update',
                        tips: '',
                        isShow: '',
                        items: [{label: '数据', field: 'data', type: 'text'}, {
                            label: '创建时间',
                            field: 'created_at',
                            type: 'datePicker'
                        }]
                    }, {
                        label: '删除',
                        url: '/wallet/deposit/delete',
                        color: 'negative',
                        type: 'delete',
                        tips: '',
                        isShow: '',
                        items: []
                    }, {
                        label: '审核',
                        url: '/wallet/deposit/status',
                        color: 'secondary',
                        type: 'update',
                        tips: '',
                        isShow: 'scope.row.status === 10',
                        items: [{
                            label: '状态',
                            field: 'status',
                            type: 'selectNumber',
                            data: [{label: '通过', value: '20'}, {label: '待审核', value: '10'}, {label: '拒绝', value: '-1'}]
                        }, {label: '失败原因', field: 'data', type: 'textarea'}]
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