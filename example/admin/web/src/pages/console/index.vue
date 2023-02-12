<template>
    <div>
        <q-splitter v-model="splitterModel" style="height: calc(100vh - 50px)">
            <template v-slot:before>
                <q-tabs v-model="tab" vertical inline-label outside-arrows @update:model-value="tablesColumnsFunc"
                        class="bg-blue-10 text-white">
                    <q-tab v-for="(tab, tabIndex) in tables" :key="tabIndex" :name="tab.name"
                           :label="tab.comment"></q-tab>
                </q-tabs>
            </template>
            <template v-slot:after>
                <q-tab-panels v-model="tab" animated swipeable vertical transition-prev="jump-up"
                              transition-next="jump-up">
                    <q-tab-panel :name="tab">
                        <q-list bordered class='rounded-borders'>
                            <!--                            查询条件配置-->
                            <q-expansion-item expand-separator default-opened label='查询条件配置'>
                                <q-card>
                                    <q-card-section>
                                        <div v-if="conf.search.length === 0">
                                            <div class="text-center">
                                                <q-btn icon="sym_o_add" flat label="新增查询配置"
                                                       class="bg-primary text-white"
                                                       @click='conf.search.splice(1, 0, JSON.parse(JSON.stringify(inputDefaultConf)))'></q-btn>
                                            </div>
                                        </div>
                                        <div v-for='(search, searchIndex) in conf.search' :key='searchIndex'>
                                            <div>#{{ searchIndex }}条件</div>
                                            <div class='row items-center q-gutter-md q-mb-md'>
                                                <div class='col'>
                                                    <q-input dense outlined v-model='search.label'
                                                             label='名称'></q-input>
                                                </div>
                                                <div class='col'>
                                                    <q-input dense outlined v-model='search.field'
                                                             label='字段'></q-input>
                                                </div>
                                                <div class='col'>
                                                    <q-select dense outlined v-model='search.type' label='类型'
                                                              @update:model-value="switchInputTypeFunc(search)"
                                                              :options='inputTypeList'></q-select>
                                                </div>
                                                <div class='col-1'>
                                                    <div class='row'>
                                                        <q-btn flat icon='add' class='no-padding' color='primary'
                                                               size='md'
                                                               @click='conf.search.splice(searchIndex + 1, 0, JSON.parse(JSON.stringify(inputDefaultConf)))'></q-btn>
                                                        <q-btn flat icon='close' class='no-padding' color='red'
                                                               size='md'
                                                               @click='conf.search.splice(searchIndex, 1)'></q-btn>
                                                    </div>
                                                </div>
                                            </div>
                                            <div class='q-mx-lg'>
                                                <edit-options :type="search.type"
                                                              :data="search.data"></edit-options>
                                            </div>
                                        </div>
                                    </q-card-section>
                                </q-card>
                            </q-expansion-item>

                            <!--                            工具栏配置-->
                            <q-expansion-item expand-separator default-opened label='工具栏配置'>
                                <q-card>
                                    <q-card-section>
                                        <div v-if="conf.tools.length === 0">
                                            <div class="text-center">
                                                <q-btn icon="sym_o_add" flat label="新增工具按钮"
                                                       class="bg-primary text-white"
                                                       @click='conf.tools.splice(1, 0, JSON.parse(JSON.stringify(editBtnDefaultConf)))'></q-btn>
                                            </div>
                                        </div>
                                        <div v-for='(tool, toolIndex) in conf.tools' :key='toolIndex'>
                                            <div class='q-mb-sm'>#{{ toolIndex }}按钮配置</div>
                                            <div class='row items-center q-gutter-md q-mb-sm'>
                                                <div class='col'>
                                                    <q-input dense outlined v-model='tool.label' label='名称'></q-input>
                                                </div>
                                                <div class='col'>
                                                    <q-select dense outlined v-model='tool.color' label='颜色'
                                                              :options='quasarColorsList'></q-select>
                                                </div>
                                                <div class='col'>
                                                    <q-input dense outlined v-model='tool.url' label='请求地址'></q-input>
                                                </div>
                                                <div class='col'>
                                                    <q-select dense outlined v-model='tool.type' label='类型'
                                                              @update:model-value="switchEditBtnTypeFunc(tool)"
                                                              :options='editBtnToolsTypeList'></q-select>
                                                </div>
                                                <div class='col-1'>
                                                    <div class='row'>
                                                        <q-btn flat icon='add' class='no-padding' color='primary'
                                                               size='md'
                                                               @click='conf.tools.splice(toolIndex + 1, 0, JSON.parse(JSON.stringify(editBtnDefaultConf)))'></q-btn>
                                                        <q-btn flat icon='close' class='no-padding' color='red'
                                                               size='md'
                                                               @click='conf.tools.splice(toolIndex, 1)'></q-btn>
                                                    </div>
                                                </div>
                                            </div>
                                            <div class='row items-center q-gutter-md q-mb-sm'>
                                                <div class='col'>
                                                    <q-input dense outlined v-model='tool.tips'
                                                             label='弹窗提示语句'></q-input>
                                                </div>
                                                <div class='col'>
                                                    <q-input dense outlined v-model='tool.isShow'
                                                             label='显示条件 = 表格数据:scope | 状态管理: store'></q-input>
                                                </div>
                                                <div class='col-1'></div>
                                            </div>
                                            <div class='q-mx-lg' v-if='tool.type !== editBtnTypeObject.checkboxDelete'>
                                                <div class='q-mb-sm'>#{{ toolIndex }}按钮配置=ITEMS</div>
                                                <div v-if="tool.items.length === 0">
                                                    <div class="text-center">
                                                        <q-btn icon="sym_o_add" flat label="新增工具表单配置"
                                                               class="bg-primary text-white"
                                                               @click='tool.items.splice(1, 0, JSON.parse(JSON.stringify(inputDefaultConf)))'></q-btn>
                                                    </div>
                                                </div>
                                                <div v-for='(children, childrenIndex) in tool.items'
                                                     :key='childrenIndex'>
                                                    <div class='row items-center q-gutter-md q-mb-md'>
                                                        <div class='col'>
                                                            <q-input dense outlined v-model='children.label'
                                                                     label='名称'></q-input>
                                                        </div>
                                                        <div class='col'>
                                                            <q-input dense outlined v-model='children.field'
                                                                     label='字段'></q-input>
                                                        </div>
                                                        <div class='col'>
                                                            <q-select dense outlined v-model='children.type' label='类型'
                                                                      @update:model-value="switchInputTypeFunc(children)"
                                                                      :options='inputTypeList'></q-select>
                                                        </div>
                                                        <div class='col-1'>
                                                            <div class='row'>
                                                                <q-btn flat icon='add' class='no-padding'
                                                                       color='primary' size='md'
                                                                       @click='tool.items.splice(childrenIndex + 1, 0, JSON.parse(JSON.stringify(inputDefaultConf)))'></q-btn>
                                                                <q-btn flat icon='close' class='no-padding' color='red'
                                                                       size='md'
                                                                       @click='tool.items.splice(childrenIndex, 1)'></q-btn>
                                                            </div>
                                                        </div>
                                                    </div>
                                                    <div class='q-mx-lg'>
                                                        <edit-options :type="children.type"
                                                                      :data="children.data"></edit-options>
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                    </q-card-section>
                                </q-card>
                            </q-expansion-item>

                            <!--                            数据表格配置-->
                            <q-expansion-item expand-separator default-opened label='数据表格配置'>
                                <q-card>
                                    <q-card-section>
                                        <div>
                                            <div v-for='(column, columnIndex) in conf.columns' :key='columnIndex'>
                                                <div class='row items-center q-gutter-md q-mb-md'>
                                                    <div class='col'>
                                                        <q-input dense outlined v-model='column.label'
                                                                 label='名称'></q-input>
                                                    </div>
                                                    <div class='col'>
                                                        <q-input dense outlined v-model='column.field'
                                                                 label='字段'></q-input>
                                                    </div>
                                                    <div class='col'>
                                                        <q-select dense outlined v-model='column.type' label='类型'
                                                                  @update:model-value="switchInputTypeFunc(column)"
                                                                  :options='inputTypeList'></q-select>
                                                    </div>
                                                    <div class='col-1'>
                                                        <div class='row'>
                                                            <q-btn flat icon='add' class='no-padding' color='primary'
                                                                   size='md'
                                                                   @click='conf.columns.splice(columnIndex + 1, 0, JSON.parse(JSON.stringify(inputDefaultConf)))'></q-btn>
                                                            <q-btn flat icon='close' class='no-padding' color='red'
                                                                   size='md'
                                                                   @click='conf.columns.splice(columnIndex, 1)'></q-btn>
                                                        </div>
                                                    </div>
                                                </div>
                                                <div class='q-mx-lg'>
                                                    <edit-options :type="column.type"
                                                                  :data="column.data"></edit-options>
                                                </div>
                                            </div>

                                            <div v-if="conf.columnOptions.length === 0">
                                                <div class="text-center">
                                                    <q-btn icon="sym_o_add" flat label="新增数据表格操作按钮"
                                                           class="bg-primary text-white"
                                                           @click='conf.columnOptions.splice(1, 0, JSON.parse(JSON.stringify(editBtnDefaultConf)))'></q-btn>
                                                </div>
                                            </div>
                                            <div v-for='(columnOption, columnOptionIndex) in conf.columnOptions'
                                                 :key='columnOptionIndex'>
                                                <div class='q-mb-sm'>#{{ columnOptionIndex }} - 按钮配置</div>
                                                <div class='row items-center q-gutter-md q-mb-sm'>
                                                    <div class='col'>
                                                        <q-input dense outlined v-model='columnOption.label'
                                                                 label='名称'></q-input>
                                                    </div>
                                                    <div class='col'>
                                                        <q-select dense outlined v-model='columnOption.color' label='颜色'
                                                                  :options='quasarColorsList'></q-select>
                                                    </div>
                                                    <div class='col'>
                                                        <q-input dense outlined v-model='columnOption.url'
                                                                 label='请求地址'></q-input>
                                                    </div>
                                                    <div class='col'>
                                                        <q-select dense outlined v-model='columnOption.type' label='类型'
                                                                  @update:model-value="switchEditBtnTypeFunc(columnOption)"
                                                                  :options='editBtnOptionsTypeList'></q-select>
                                                    </div>
                                                    <div class='col-1'>
                                                        <div class='row'>
                                                            <q-btn flat icon='add' class='no-padding' color='primary'
                                                                   size='md'
                                                                   @click='conf.columnOptions.splice(columnOptionIndex + 1, 0, JSON.parse(JSON.stringify(editBtnDefaultConf)))'></q-btn>
                                                            <q-btn flat icon='close' class='no-padding' color='red'
                                                                   size='md'
                                                                   @click='conf.columnOptions.splice(columnOptionIndex, 1)'></q-btn>
                                                        </div>
                                                    </div>
                                                </div>
                                                <div class='row items-center q-gutter-md q-mb-sm'>
                                                    <div class='col'>
                                                        <q-input dense outlined v-model='columnOption.tips'
                                                                 label='弹窗提示语句'></q-input>
                                                    </div>
                                                    <div class='col'>
                                                        <q-input dense outlined v-model='columnOption.isShow'
                                                                 label='显示条件 = 表格数据:scope | 状态管理: store'></q-input>
                                                    </div>
                                                    <div class='col-1'></div>
                                                </div>
                                                <div class='q-pl-md'
                                                     v-if='columnOption.type !== editBtnTypeObject.delete'>
                                                    <div class='q-mb-sm'>按钮弹窗表单配置</div>
                                                    <div v-if="columnOption.items.length === 0">
                                                        <div class="text-center">
                                                            <q-btn icon="sym_o_add" flat label="新增数据表格操作按钮"
                                                                   class="bg-primary text-white"
                                                                   @click='columnOption.items.splice(1, 0, JSON.parse(JSON.stringify(inputDefaultConf)))'></q-btn>
                                                        </div>
                                                    </div>
                                                    <div v-for='(children, childrenIndex) in columnOption.items'
                                                         :key='childrenIndex'>
                                                        <div class='row items-center q-gutter-md q-mb-md'>
                                                            <div class='col'>
                                                                <q-input dense outlined v-model='children.label'
                                                                         label='名称'></q-input>
                                                            </div>
                                                            <div class='col'>
                                                                <q-input dense outlined v-model='children.field'
                                                                         label='字段'></q-input>
                                                            </div>
                                                            <div class='col'>
                                                                <q-select dense outlined v-model='children.type'
                                                                          label='类型'
                                                                          @update:model-value="switchInputTypeFunc(children)"
                                                                          :options='inputTypeList'></q-select>
                                                            </div>
                                                            <div class='col-1'>
                                                                <div class='row'>
                                                                    <q-btn flat icon='add' class='no-padding'
                                                                           color='primary' size='md'
                                                                           @click='columnOption.items.splice(childrenIndex + 1, 0, JSON.parse(JSON.stringify(inputDefaultConf)))'></q-btn>
                                                                    <q-btn flat icon='close' class='no-padding'
                                                                           color='red' size='md'
                                                                           @click='columnOption.items.splice(childrenIndex, 1)'></q-btn>
                                                                </div>
                                                            </div>
                                                        </div>
                                                        <div class='q-mx-lg'>
                                                            <edit-options :type="children.type"
                                                                          :data="children.data"></edit-options>
                                                        </div>
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                    </q-card-section>
                                </q-card>
                            </q-expansion-item>
                        </q-list>

                        <div class='row q-mt-md justify-end'>
                            <div>
                                <q-btn label='生成代码' color='positive' @click='basicModel = true'></q-btn>
                            </div>
                        </div>
                    </q-tab-panel>
                </q-tab-panels>
            </template>
        </q-splitter>

        <q-dialog v-model="basicModel">
            <q-card class="full-width">
                <q-card-section>
                    <div class="text-h6">生成数据表格</div>
                </q-card-section>

                <q-card-section class="q-pt-none">
                    <div class="q-gutter-sm">
                        <q-input dense v-model="conf.basic.indexURL" autofocus label="请求地址"/>
                        <q-input dense v-model="conf.basic.updateURL" label="更新地址"/>
                        <q-input dense v-model="conf.basic.updateKey" label="主键Key"/>
                    </div>

                </q-card-section>

                <q-card-actions align="right" class="text-primary">
                    <q-btn flat label="取消" v-close-popup/>
                    <q-btn flat label="显示代码" v-close-popup @click='showVueTemplateFunc'/>
                </q-card-actions>
            </q-card>
        </q-dialog>

        <q-dialog v-model='showVueTemplate.model' :maximized='showVueTemplate.maximizedToggle' persistent
                  transition-show='slide-up'
                  transition-hide='slide-down'>
            <q-card class="full-width">
                <q-bar class='bg-primary text-white'>
                    <q-space/>
                    <q-btn dense flat icon='minimize' @click='showVueTemplate.maximizedToggle = false'
                           :disable='!showVueTemplate.maximizedToggle'>
                        <q-tooltip v-if='showVueTemplate.maximizedToggle' class='bg-white text-primary'>Minimize
                        </q-tooltip>
                    </q-btn>
                    <q-btn dense flat icon='crop_square' @click='showVueTemplate.maximizedToggle = true'
                           :disable='showVueTemplate.maximizedToggle'>
                        <q-tooltip v-if='!showVueTemplate.maximizedToggle' class='bg-white text-primary'>Maximize
                        </q-tooltip>
                    </q-btn>
                    <q-btn dense flat icon='close' v-close-popup>
                        <q-tooltip class='bg-white text-primary'>Close</q-tooltip>
                    </q-btn>
                </q-bar>
                <q-card-section>
                    <q-input v-model='showVueTemplate.value' type='textarea' outlined autogrow></q-input>
                </q-card-section>
            </q-card>
        </q-dialog>
    </div>
</template>

<script lang="ts">
import {TablesIndexAPI, TablesColumnsAPI} from 'src/api';
import {onMounted, reactive, toRefs} from 'vue'
import {
    inputTypeObject,
    inputDefaultConf,
    inputDefaultOptionsConf,
    quasarColorsObject,
    editBtnTypeObject,
    editBtnDefaultConf
} from 'src/utils/define';
import EditOptions from 'components/editOptions.vue';
import {generateTableFunc} from 'src/utils/template';

export default {
    name: 'ConsoleIndex',
    components: {EditOptions},
    setup() {
        const state = reactive({
            basicModel: false,
            splitterModel: 20,
            tab: '',
            tables: [] as any,
            columns: [] as any,
            inputTypeObject: inputTypeObject,
            inputTypeList: Object.values(inputTypeObject),
            editBtnTypeObject: editBtnTypeObject,
            editBtnToolsTypeList: [
                editBtnTypeObject.create,
                editBtnTypeObject.checkboxDelete
            ],
            editBtnOptionsTypeList: [
                editBtnTypeObject.update,
                editBtnTypeObject.view,
                editBtnTypeObject.delete
            ],
            quasarColorsObject: quasarColorsObject,
            quasarColorsList: Object.values(quasarColorsObject),
            showVueTemplate: {
                model: false,
                maximizedToggle: true,
                value: ''
            },
            conf: {
                basic: {
                    indexURL: '',
                    updateURL: '',
                    updateKey: 'id'
                } as any,
                search: [] as any,
                tools: [] as any,
                columns: [] as any,
                columnOptions: [] as any,
            }
        })

        // 初始化数据
        onMounted(() => {
            TablesIndexAPI().then((res) => {
                state.tables = res
                state.tab = state.tables[0].name
                tablesColumnsFunc()
            })
        })


        // 改变input类型方法
        const switchInputTypeFunc = (conf: any) => {
            switch (conf.type) {
                case inputTypeObject.select:
                case inputTypeObject.selectNumber:
                case inputTypeObject.toggle:
                case inputTypeObject.checkbox:
                case inputTypeObject.editToggle:
                case inputTypeObject.slider:
                case inputTypeObject.range:
                    conf.data = [] as any
                    conf.data.push(JSON.parse(JSON.stringify(inputDefaultOptionsConf)))
                    break;
            }
        }

        // 改变editBtn 类型方法
        const switchEditBtnTypeFunc = (conf: any) => {
            switch (conf.type) {
                case editBtnTypeObject.delete:
                case editBtnTypeObject.checkboxDelete:
                    conf.items = []
            }
        }

        // 获取表属性
        const tablesColumnsFunc = () => {
            TablesColumnsAPI({name: state.tab}).then((res) => {
                state.columns = res
                state.conf.search = []
                state.conf.tools = []

                state.conf.columns = []
                for (let i = 0; i < state.columns.length; i++) {
                    const item = state.columns[i];
                    const inputDefaultConfTmp = JSON.parse(JSON.stringify(inputDefaultConf));
                    state.conf.columns.push(Object.assign(inputDefaultConfTmp, {
                        label: item.comment,
                        field: item.field,
                        type: inputTypeObject.text
                    }));
                }
            })
        }

        // 显示vue代码
        const showVueTemplateFunc = () => {
            state.showVueTemplate.value = generateTableFunc(state.conf)
            state.showVueTemplate.model = true
        }

        return {
            inputDefaultConf,
            editBtnDefaultConf,
            tablesColumnsFunc,
            switchInputTypeFunc,
            switchEditBtnTypeFunc,
            showVueTemplateFunc,
            ...toRefs(state)
        }
    }
}
</script>

<style scoped>

</style>