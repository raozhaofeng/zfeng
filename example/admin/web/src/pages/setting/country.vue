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
    name: 'SettingCountry',
    components: {DynamicInput, DialogForm, TableBodyCell},
    setup() {
        const visibleStorage = '_table_visible_008'
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
                basic: {indexURL: '/country/index', updateURL: '/country/update', updateKey: 'id'},
                search: [{label: '管理', field: 'admin_name', type: 'text'}, {label: '名称', field: 'name', type: 'text'}, {
                    label: '别名',
                    field: 'alias',
                    type: 'text'
                }, {label: 'ISO', field: 'iso1', type: 'text'}, {label: '区号', field: 'code', type: 'text'}, {
                    label: '时间',
                    field: 'created_at',
                    type: 'dateRangePicker'
                }],
                tools: [{
                    label: '新增',
                    url: '/country/create',
                    color: 'primary',
                    type: 'create',
                    tips: '',
                    isShow: '',
                    items: [{
                        label: '图标',
                        field: 'icon',
                        type: 'image'
                    }, {label: '名称', field: 'name', type: 'text'}, {
                        label: '别名',
                        field: 'alias',
                        type: 'text'
                    }, {label: 'ISO', field: 'iso1', type: 'text'}, {
                        label: '区号',
                        field: 'code',
                        type: 'text'
                    }, {label: '数据', field: 'data', type: 'textarea'}]
                }, {
                    label: '删除',
                    url: '/country/delete',
                    color: 'negative',
                    type: 'checkboxDelete',
                    tips: '',
                    isShow: '',
                    items: []
                }],
                visibleColumns: localStorage.getItem(visibleStorage) ?
                    JSON.parse(<string>localStorage.getItem(visibleStorage)) :
                    ['id', 'admin_name', 'name', 'alias', 'iso1', 'icon', 'code', 'sort', 'status', 'created_at', 'options'],
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
                    sortable: true
                }, {
                    label: '名称',
                    field: 'name',
                    name: 'name',
                    type: 'editText',
                    align: 'left',
                    sortable: true
                }, {
                    label: '别名',
                    field: 'alias',
                    name: 'alias',
                    type: 'editText',
                    align: 'left',
                    sortable: true
                }, {
                    label: 'ISO',
                    field: 'iso1',
                    name: 'ISO',
                    type: 'editText',
                    align: 'left',
                    sortable: true
                },{
                    label: '图标',
                    field: 'icon',
                    name: 'icon',
                    type: 'image',
                    align: 'left',
                    sortable: true
                }, {
                    label: '区号',
                    field: 'code',
                    name: 'code',
                    type: 'editText',
                    align: 'left',
                    sortable: true
                }, {
                    label: '排序',
                    field: 'sort',
                    name: 'sort',
                    type: 'editNumber',
                    align: 'left',
                    sortable: true
                }, {
                    label: '状态',
                    field: 'status',
                    name: 'status',
                    type: 'editToggle',
                    align: 'left',
                    sortable: true,
                    data: [{label: '激活', value: '10'}, {label: '禁用', value: '-1'}]
                }, {
                    label: '创建时间',
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
                        url: '/country/update',
                        color: 'primary',
                        type: 'update',
                        tips: '',
                        isShow: '',
                        items: [{
                            label: '图标',
                            field: 'icon',
                            type: 'image'
                        }, {label: '名称', field: 'name', type: 'text'}, {
                            label: '别名',
                            field: 'alias',
                            type: 'text'
                        }, {label: 'ISO', field: 'iso1', type: 'text'}, {label: '区号', field: 'code', type: 'text'}, {
                            label: '数据',
                            field: 'data',
                            type: 'textarea'
                        }, {
                            label: '状态',
                            field: 'status',
                            type: 'selectNumber',
                            data: [{label: '激活', value: '10'}, {label: '禁用', value: '-1'}]
                        }]
                    }, {
                        label: '删除',
                        url: '/country/delete',
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