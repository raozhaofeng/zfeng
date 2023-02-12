<template>
    <div v-if='scope.col.type === inputTypeObject.image'>
        <q-img :src='imageSrc(scope.value)'
               loading='lazy' spinner-color='white' style='max-height: 50px; max-width: 50px'>
            <q-popup-proxy transition-show='flip-up' transition-hide='flip-down'>
                <q-card style="width: 360px;">
                    <q-card-section>
                        <q-img no-spinner :src='imageSrc(scope.value)'></q-img>
                    </q-card-section>
                </q-card>
            </q-popup-proxy>
        </q-img>
    </div>
    <div v-else-if="scope.col.type === inputTypeObject.images">
        <q-img :src='imageSrc(JSON.parse(scope.value)[0].value)'
               v-if="scope.value !== ''"
               loading='lazy' spinner-color='white' style='max-height: 50px; max-width: 50px'>
            <q-popup-proxy transition-show='flip-up' transition-hide='flip-down'>
                <q-card style="width: 360px;">
                    <q-card-section v-for="(image, imageIndex) in JSON.parse(scope.value)" :key="imageIndex">
                        <q-img no-spinner :src='imageSrc(image.value)'></q-img>
                    </q-card-section>
                </q-card>
            </q-popup-proxy>
        </q-img>
    </div>
    <div v-else-if='scope.col.type === inputTypeObject.file'>
        [文件]
    </div>
    <div v-else-if='scope.col.type === inputTypeObject.datePicker'>
        {{ date.formatDate(scope.value * 1000, 'YYYY-MM-DD HH:mm:ss') }}
    </div>
    <div v-else-if='scope.col.type === inputTypeObject.select || scope.col.type === inputTypeObject.selectNumber'>
        <div v-for='(selected, selectedIndex) in scope.col.data' :key='selectedIndex'>
            <q-badge outline :color='quasarColorsList[selectedIndex]'
                     v-if='parseInt(selected.value) === scope.value' :label='selected.label'/>
        </div>
    </div>
    <div
        v-else-if='scope.col.type === inputTypeObject.editText || scope.col.type === inputTypeObject.editNumber || scope.col.type === inputTypeObject.editTextarea'
        class='ellipsis'>
        {{ scope.value === null || scope.value === '' ? '...' : scope.value }}
        <q-popup-edit v-model='scope.row[scope.col.field]'
                      :model-value='scope.row[scope.col.field]'>
            <template v-slot:default='popupScope'>
                <q-input v-model='popupScope.value' dense autofocus counter
                         v-if="scope.col.type === inputTypeObject.editTextarea"
                         :model-value='popupScope.value' type='textarea' @keyup.enter.stop/>
                <q-input v-model='popupScope.value' dense autofocus counter
                         v-else-if="scope.col.type === inputTypeObject.editText"
                         :model-value='popupScope.value' @keyup.enter.stop/>
                <q-input v-model.number='popupScope.value' dense autofocus counter
                         v-else-if="scope.col.type === inputTypeObject.editNumber"
                         :model-value='popupScope.value' type='number' @keyup.enter.stop/>
                <div class='row justify-end'>
                    <div>
                        <q-btn label='取消' color='primary' v-close-popup flat></q-btn>
                    </div>
                    <div>
                        <q-btn label='确定' color='primary' v-close-popup flat
                               @click='updatePopupEditFunc(scope.row, scope.col.field, popupScope.value)'></q-btn>
                    </div>
                </div>
            </template>
        </q-popup-edit>
    </div>
    <div v-else-if='scope.col.type === inputTypeObject.editToggle'>
        <q-toggle v-model='scope.row[scope.col.field]'
                  :model-value='scope.row[scope.col.field]'
                  :true-value='parseInt(scope.col?.data[0].value)'
                  :false-value='parseInt(scope.col?.data[1].value)'
                  @update:model-value='updatePopupEditFunc(scope.row, scope.col.field, scope.row[scope.col.field])'/>
    </div>
    <div v-else-if="scope.col.type === 'options'">
        <div class="row q-gutter-xs" :style="{width: '120px'}">
            <div v-for='(option, optionIndex) in scope.col.items' :key='optionIndex'>
                <q-btn :label='option.label' :color='option.color' size='xs'
                       v-if='editBtnEvalIsShowFunc(option.isShow, {scope: scope}) && isAuth(option.url)'
                       @click='showDialogFormFunc(option, scope.row)'></q-btn>
            </div>
        </div>
    </div>
    <div v-else class='ellipsis'>
        {{ scope.value }}
        <q-popup-proxy>
            <q-card>
                <q-card-section>
                    {{ scope.value }}
                </q-card-section>
            </q-card>
        </q-popup-proxy>
    </div>

    <dialog-form ref='dialogFormRef' :title='dialogForm.label' :items='dialogForm.items'
                 :extend-values="dialogForm.extendValues" :type="dialogForm.type"
                 :dynamic-data='dialogForm.dynamicData' :url="dialogForm.url" :after-func='dialogForm.afterFunc'
                 :values='dialogForm.values'></dialog-form>
</template>

<script lang="ts">
import {reactive, ref, toRefs, watch} from 'vue'
import {api} from 'boot/axios';
import {date} from 'quasar';
import {editBtnTypeObject, inputTypeObject, quasarColorsObject} from 'src/utils/define';
import {confirmBoxDialog, imageSrc} from 'src/utils';
import {positiveNotify} from 'src/utils';
import DialogForm from 'src/components/dialogForm.vue';
import {editBtnEvalIsShowFunc} from 'src/utils/define';
import {isAuth} from 'src/utils/auth';

export default {
    name: 'TableBodyCell',
    props: {
        updateKey: {
            type: String
        },
        updateUrl: {
            type: String
        },
        data: {
            type: Object
        },
        extendValues: {
            type: Object
        },
        dynamicData: {
            type: Object
        }
    },
    components: {DialogForm},
    emits: ['refresh', 'update:data'],
    setup(props: any, context: any) {
        const dialogFormRef = ref(null) as any;
        const state = reactive({
            dialogForm: {} as any,
            quasarColorsList: Object.values(quasarColorsObject),
            scope: props.data
        })

        watch(props, (newProps) => {
            state.scope = newProps.data;
        });

        // 编辑提交
        const updatePopupEditFunc = (row: any, field: string, value: any) => {
            const params = {} as any;
            params[props.updateKey] = row[props.updateKey];
            params[field] = value;
            api.post(props.updateUrl, params).then(() => {
                row[field] = value
                context.emit('update:data', state.scope)
                positiveNotify('保存成功')
            });
        }

        // 显示模态框方法
        const showDialogFormFunc = (conf: any, row: any) => {
            switch (conf.type) {
                case editBtnTypeObject.delete:
                    const params = {} as any
                    params[props.updateKey] = []
                    params[props.updateKey].push(row[props.updateKey])
                    confirmDeleteFunc(conf.url, '是否确认删除【' + row[props.updateKey] + '】', params)
                    break;
                case editBtnTypeObject.update:
                case editBtnTypeObject.view:
                    state.dialogForm = conf
                    state.dialogForm.afterFunc = () => {
                        context.emit('refresh')
                    }
                    state.dialogForm.extendValues = props.extendValues
                    state.dialogForm.values = {}
                    state.dialogForm.extendValues[props.updateKey] = row[props.updateKey]
                    state.dialogForm.items.forEach((item: any) => {
                        state.dialogForm.values[item.field] = row[item.field]
                    })
                    state.dialogForm.dynamicData = props.dynamicData
                    dialogFormRef.value.isShow = true
            }
        }

        // 确认删除数据
        const confirmDeleteFunc = async (url: string, msg: string, params: any) => {
            await confirmBoxDialog('删除数据', msg, () => {
                api.post(url, params).then(() => {
                    context.emit('refresh')
                })
            })
        }

        return {
            isAuth,
            date,
            imageSrc,
            inputTypeObject,
            dialogFormRef,
            editBtnEvalIsShowFunc,
            updatePopupEditFunc,
            showDialogFormFunc,
            ...toRefs(state)
        }
    }
}
</script>

<style scoped>

</style>