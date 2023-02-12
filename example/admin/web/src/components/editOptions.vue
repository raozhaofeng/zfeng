<template>
    <div v-for='(option, optionIndex) in options' :key='optionIndex'>
        <div class='row items-center q-gutter-md q-mb-md'>
            <div class='col'>
                <q-input dense outlined v-model='option.label' label='option名称'
                         @update:model-value="$emit('update:data', options)"></q-input>
            </div>
            <div class='col'>
                <q-input dense outlined v-model='option.value' label='option数据'
                         @update:model-value="$emit('update:data', options)"></q-input>
            </div>
            <div class='col-1'>
                <div class='row'>
                    <q-btn flat icon='add' class='no-padding' color='primary' size='md'
                           @click="createOptionsFunc(optionIndex)"></q-btn>
                    <q-btn flat icon='close' class='no-padding' color='red' size='md'
                           @click="deleteOptionFunc"></q-btn>
                </div>
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import {reactive, toRefs, watch} from 'vue'
import {inputTypeObject, inputDefaultOptionsConf} from 'src/utils/define';

export default {
    name: 'editOptions',
    props: {
        type: {
            type: String
        },
        data: {
            type: Array
        }
    },
    setup(props: any, context: any) {
        const state = reactive({
            options: props.data
        })

        watch(props, (newProps) => {
            state.options = newProps.data;
        });

        // 新增option
        const createOptionsFunc = (optionIndex: number) => {
            if (props.type === inputTypeObject.slider || props.type === inputTypeObject.range) {
                if (state.options?.length > 2) {
                    return
                }
            }
            state.options.splice(optionIndex + 1, 0, JSON.parse(JSON.stringify(inputDefaultOptionsConf)));
            context.emit('update:data', state.options)
        }

        // 删除option
        const deleteOptionFunc = (optionIndex: number) => {
            state.options.splice(optionIndex, 1);
            context.emit('update:data', state.options)
        }

        return {
            inputTypeObject,
            inputDefaultOptionsConf,
            createOptionsFunc,
            deleteOptionFunc,
            ...toRefs(state)
        }
    }
}
</script>

<style scoped>

</style>