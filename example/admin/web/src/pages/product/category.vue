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
    name: 'ProductCategory',
    components: {DynamicInput, DialogForm, TableBodyCell},
    setup() {
        const visibleStorage = '_table_visible_005'
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
                basic: {indexURL: '/product/category/index', updateURL: '/product/category/update', updateKey: 'id'},
                search: [{label: '父级分类', field: 'parent_name', type: 'text'}, {
                    label: '管理员',
                    field: 'admin_name',
                    type: 'text'
                }, {label: '标题(语言字典)', field: 'name', type: 'text'}, {
                    label: '状态',
                    field: 'status',
                    type: 'selectNumber',
                    data: [{label: '禁用', value: '-1'}, {label: '激活', value: '10'}]
                }, {
                    label: '类型',
                    field: 'type',
                    type: 'selectNumber',
                    data: [{label: '商品', value: '1'}, {label: '区块链', value: '10'}]
                },{
                    label: '推荐',
                    field: 'recommend',
                    type: 'selectNumber',
                    data: [{label: '关闭', value: '-1'}, {label: '推荐', value: '10'}]
                }, {label: '时间', field: 'updated_at', type: 'dateRangePicker'}],
                tools: [{
                    label: '新增',
                    url: '/product/category/create',
                    color: 'primary',
                    type: 'create',
                    tips: '',
                    isShow: '',
                    items: [{label: '封面', field: 'image', type: 'image'}, {
                        label: '父级ID',
                        field: 'parent_id',
                        type: 'number'
                    }, {
                        label: '类型',
                        field: 'type',
                        type: 'selectNumber',
                        data: [{label: '商品', value: '1'}, {label: '区块链', value: '10'}]
                    }, {label: '标题(语言字典)', field: 'name', type: 'text'}, {label: '数据', field: 'data', type: 'textarea'}]
                }, {
                    label: '删除',
                    url: '/product/category/delete',
                    color: 'negative',
                    type: 'checkboxDelete',
                    tips: '',
                    isShow: '',
                    items: []
                }],
                visibleColumns: localStorage.getItem(visibleStorage) ?
                    JSON.parse(<string>localStorage.getItem(visibleStorage)) :
                    ['id', 'parent_name', 'admin_name', 'type', 'name', 'image', 'sort', 'status', 'recommend', 'updated_at', 'options'],
                columns: [{
                    label: '主键',
                    field: 'id',
                    name: 'id',
                    type: 'text',
                    align: 'left',
                    sortable: true
                }, {
                    label: '父级分类',
                    field: 'parent_name',
                    name: 'parent_name',
                    type: 'text',
                    align: 'left',
                    sortable: true
                }, {
                    label: '管理员',
                    field: 'admin_name',
                    name: 'admin_name',
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
                    data: [{label: '商品', value: '1'}, {label: '区块链', value: '10'}]
                }, {
                    label: '标题(语言字典)',
                    field: 'name',
                    name: 'name',
                    type: 'editText',
                    align: 'left',
                    sortable: true
                }, {
                    label: '封面',
                    field: 'image',
                    name: 'image',
                    type: 'image',
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
                    label: '推荐',
                    field: 'recommend',
                    name: 'recommend',
                    type: 'editToggle',
                    align: 'left',
                    sortable: true,
                    data: [{label: '推荐', value: '10'}, {label: '关闭', value: '-1'}]
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
                        label: '编辑',
                        url: '/product/category/update',
                        color: 'primary',
                        type: 'update',
                        tips: '',
                        isShow: '',
                        items: [{label: '封面', field: 'image', type: 'image'}, {
                            label: '父级ID',
                            field: 'parent_id',
                            type: 'number'
                        }, {
                            label: '类型',
                            field: 'type',
                            type: 'selectNumber',
                            data: [{label: '商品', value: '1'}, {label: '区块链', value: '10'}]
                        },{label: '标题(语言字典)', field: 'name', type: 'text'}, {
                            label: '数据',
                            field: 'data',
                            type: 'textarea'
                        }]
                    }, {
                        label: '删除',
                        url: '/product/category/delete',
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
