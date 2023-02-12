<template>
    <div class="q-ma-md">
        <q-card flat bordered>
            <q-tabs v-model="tab" dense class="text-grey" active-color="primary" indicator-color="primary"
                    align="justify" narrow-indicator>
                <q-tab :name="group.id" :label="group.name" v-for="(group, index) in groups"
                       :key="index"></q-tab>
            </q-tabs>
            <q-separator/>

            <q-tab-panels v-model="tab" animated>
                <q-tab-panel :name="group.id" v-for="(group, index) in groups" :key="index">
                    <div class="q-gutter-md">
                        <div v-for="(children, childrenIndex) in group.children" :key="childrenIndex">
                            <div class="q-my-xs text-bold" v-if="children.type !== inputTypeObject.checkbox">
                                {{ children.name }}
                            </div>
                            <!--                            文章设置-->
                            <div v-if="children.type === inputTypeObject.children">
                                <div class="row justify-end q-mb-md q-gutter-xs">
                                    <div>
                                        <q-btn dense flat outline class="bg-secondary text-white q-pa-xs" size="xs"
                                               @click="addChildrenFunc(children)"
                                               label="添加"></q-btn>
                                    </div>
                                    <div>
                                        <q-btn dense flat outline class="bg-red text-white q-pa-xs" size="xs"
                                               @click="delChildrenFunc(children)"
                                               label="删除"></q-btn>
                                    </div>
                                </div>
                                <div v-for="(childrenValue, childrenValueIndex) in JSON.parse(children.value)"
                                     :key="childrenValueIndex">
                                    <div class="row q-gutter-md items-center">
                                        <div class="col q-mb-md"
                                             v-for="(childrenData, childrenDataIndex) in JSON.parse(children.data)"
                                             :key="childrenDataIndex">
                                            <dynamic-input :type="childrenData.type" :label="childrenData.name"
                                                           v-model="params[children.field].value[childrenValueIndex][childrenData.field]"
                                                           :data="[]"></dynamic-input>
                                        </div>
                                    </div>
                                </div>
                            </div>
                            <div v-else-if="children.type === inputTypeObject.json">
                                <div v-for="(arrValue, arrIndex) in JSON.parse(children.data)" :key="arrIndex">
                                    <div class="row q-gutter-sm">
                                        <div class="col q-mb-sm" v-for="(rowValue, rowIndex) in arrValue" :key="rowIndex">
                                            <dynamic-input :type="rowValue.type" :label="rowValue.name"
                                                           v-model="params[children.field].value[rowValue.field]"
                                                           :data="rowValue.data"></dynamic-input>
                                        </div>
                                    </div>
                                </div>
                            </div>
                            <dynamic-input :type="children.type" :label="children.name" v-else
                                           v-model="params[children.field].value"
                                           :data="children.data === '' ? [] : JSON.parse(children.data)"></dynamic-input>
                        </div>
                    </div>
                    <div class="text-right q-mt-md">
                        <q-btn color="primary" label="保存设置" @click="submitFunc"></q-btn>
                    </div>
                </q-tab-panel>
            </q-tab-panels>
        </q-card>
    </div>
</template>

<script lang="ts">
import {api} from 'boot/axios';
import DynamicInput from '/src/components/input.vue'
import {inputTypeObject} from 'src/utils/define';
import {positiveNotify} from 'src/utils';
import {onMounted, reactive, toRefs} from 'vue'

export default {
    name: 'SettingIndex',
    components: {DynamicInput},
    setup() {
        const state = reactive({
            tab: '',
            groups: [] as any,
            params: {} as any
        })

        onMounted(() => {
            api.post('/setting/index').then((res: any) => {
                state.groups = res.groups
                state.tab = state.groups[0].id
                for (let i = 0; i < state.groups.length; i++) {
                    let item = state.groups[i]
                    // 初始值
                    for (let j = 0; j < item.children.length; j++) {
                        let itemChildren = item.children[j]
                        state.params[itemChildren.field] = {
                            id: itemChildren.id, type: itemChildren.type, field: itemChildren.field
                        } as any

                        // 赋值初始值
                        if (itemChildren.type === inputTypeObject.checkbox || itemChildren.type === inputTypeObject.json || itemChildren.type === inputTypeObject.images || itemChildren.type === inputTypeObject.children) {
                            state.params[itemChildren.field].value = JSON.parse(itemChildren.value)
                        } else {
                            state.params[itemChildren.field].value = itemChildren.value
                        }
                    }
                }
            });
        })

        // 新增children
        const addChildrenFunc = (children: any) => {
            let oldValue = JSON.parse(children.value)
            let oldData = JSON.parse(children.data)
            let tmpValue = {} as any

            for (let i = 0; i < oldData.length; i++) {
                tmpValue[oldData[i].field] = ''
            }
            oldValue.push(tmpValue)
            state.params[children.field].value.push(tmpValue)
            children.value = JSON.stringify(oldValue)
        }

        // 删除children
        const delChildrenFunc = (children: any) => {
            let oldValue = JSON.parse(children.value)
            oldValue.pop()
            state.params[children.field].value.pop()
            children.value = JSON.stringify(oldValue)
        }

        // 修改配置
        const submitFunc = () => {
            api.post('/setting/update', state.params).then(() => {
                positiveNotify('保存配置成功')
            })
        }

        return {
            inputTypeObject,
            submitFunc,
            addChildrenFunc,
            delChildrenFunc,
            ...toRefs(state)
        }
    }
}
</script>

<style scoped>

</style>
