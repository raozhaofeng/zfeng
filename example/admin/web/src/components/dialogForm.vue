<template>
    <q-dialog v-model='isShow' persistent>
        <q-card class='full-width'>
            <q-card-section v-if='title !== ""'>
                <div class='text-h6'>{{ title }}</div>
            </q-card-section>

            <q-card-section class='q-pt-none'>
                <div class='q-gutter-sm'>
                    <div v-for='(input, inputIndex) in items' :key='inputIndex'>
                        <div v-if="input.type === inputTypeObject.json">
                            <div>{{input.label}}</div>
                            <div v-for="(rows, rowsIndex) in input.data" :key="rowsIndex" class="row q-gutter-sm">
                                <div v-for="(col, colIndex) in rows" :key="colIndex" class="col q-mb-sm">
                                    <dynamic-input :type='col.type' v-model='params[input.field][col.field]'
                                                   :readonly="type === editBtnTypeObject.view"
                                                   :data='col.data'
                                                   :label='col.label'></dynamic-input>
                                </div>
                            </div>
                        </div>
                        <dynamic-input :type='input.type' v-model='params[input.field]' v-else
                                       :readonly="type === editBtnTypeObject.view"
                                       :data='dynamicData?.hasOwnProperty(input.field) ? dynamicData[input.field] : input.data'
                                       :label='input.label'></dynamic-input>
                    </div>
                </div>
            </q-card-section>

            <q-card-actions align='right' class='text-primary'>
                <q-btn flat label='取消' v-close-popup/>
                <q-btn flat label='保存' @click="submitFunc"/>
            </q-card-actions>
        </q-card>
    </q-dialog>
</template>

<script lang='ts'>
import {api} from 'boot/axios';
import {positiveNotify} from 'src/utils';
import DynamicInput from 'src/components/input.vue';
import {editBtnTypeObject, inputTypeObject} from 'src/utils/define';
import {reactive, toRefs, watch} from 'vue';

export default {
    name: 'DialogForm',
    props: {
        // 标题
        title: {
            type: String
        },
        type: {
            type: String,
        },
        // 提交路由
        url: {
            type: String
        },
        // input 条数
        items: {
            type: Array
        },
        // 动态替换data
        dynamicData: {
            type: Object
        },
        // 字段值
        values: {
            type: Object,
            default: () => {
                return {};
            }
        },
        // 扩展值
        extendValues: {
            type: Object,
            default: () => {
                return {};
            }
        },
        afterFunc: {
            type: Function,
        }
    },
    components: {DynamicInput},
    setup(props: any) {
        const state = reactive({
            isShow: false,
            params: {} as any
        });

        watch(props, (newProps) => {
            state.params = JSON.parse(JSON.stringify(Object.assign(newProps.extendValues), state.params));
            console.log(newProps)
            newProps.items?.forEach((item: any) => {
                //  如果是JSON 格式
                if (item.type === inputTypeObject.json) {
                    state.params[item.field] = {} as any
                    console.log(item.field)
                    item.data.forEach((itemRows: any) => {
                        itemRows.forEach((itemCol: any) => {
                            if (newProps.values?.hasOwnProperty(item.field) && newProps.values[item.field].hasOwnProperty(itemCol.field)) {
                                state.params[item.field][itemCol.field] = newProps.values[item.field][itemCol.field]
                            }
                        })
                    })
                } else {
                    // 默认值
                    if (newProps.values?.hasOwnProperty(item.field)) {
                        state.params[item.field] = newProps.values[item.field];
                    }
                }
            });
        });

        // 提交数据
        const submitFunc = () => {
            if (props.type === editBtnTypeObject.update || props.type === editBtnTypeObject.create) {
                api.post(props.url, state.params).then(() => {
                    positiveNotify('保存成功')
                    if (props.afterFunc !== undefined) {
                        props.afterFunc(state.params)
                    }
                    state.isShow = false
                })
            }
        };

        return {
            editBtnTypeObject,
            inputTypeObject,
            submitFunc,
            ...toRefs(state)
        };
    }
};
</script>

<style scoped>

</style>
