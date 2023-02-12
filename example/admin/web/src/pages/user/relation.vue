<template>
    <div class="q-ma-md">
        <q-card flat bordered>
            <q-card-section>
                <q-input ref='filterRef' dense outlined v-model='filter' label='搜索网站用户名...'>
                    <template v-slot:append>
                        <q-icon v-if="filter !== ''" name='clear' class='cursor-pointer' @click='resetFilterFunc'/>
                    </template>
                </q-input>
            </q-card-section>
            <q-card-section>
                <q-tree :nodes='userTree' node-key='id' :filter='filter' :filter-method='filterMethodFunc' ref="treeRef">
                    <template v-slot:header-root='prop'>
                        <q-item class="no-padding">
                            <q-item-section avatar>
                                <q-img no-spinner :src="imageSrc(prop.node.avatar)" width="50px" height="50px"></q-img>
                            </q-item-section>
                            <q-item-section>
                                <div class="text-weight-bold text-body1">{{ prop.node.username }} <span
                                    class="text-caption q-ml-sm text-grey">成员({{ prop.node.sum_people }})</span></div>
                                <div class="text-caption text-grey">总购买金额:
                                    <span>{{ prop.node.sum_amount.toFixed(2) }}</span></div>
                                <div class="text-caption text-grey">总收益金额:
                                    <span>{{ prop.node.sum_earnings.toFixed(2) }}</span></div>
                            </q-item-section>
                        </q-item>
                    </template>
                    <template v-slot:header-generic='prop'>
                        <q-item class="no-padding">
                            <q-item-section avatar>
                                <q-img no-spinner :src="imageSrc(prop.node.avatar)" width="50px" height="50px"></q-img>
                            </q-item-section>
                            <q-item-section>
                                <div class="text-weight-bold text-body1">{{ prop.node.username }} <span
                                    class="text-caption q-ml-sm text-grey">成员({{ prop.node.sum_people }})</span></div>
                                <div class="text-caption text-grey">总购买金额:
                                    <span>{{ prop.node.sum_amount.toFixed(2) }}</span></div>
                                <div class="text-caption text-grey">总收益金额:
                                    <span>{{ prop.node.sum_earnings.toFixed(2) }}</span></div>
                            </q-item-section>
                        </q-item>
                    </template>
                </q-tree>
            </q-card-section>
        </q-card>
    </div>
</template>

<script lang="ts">
import {onMounted, ref, reactive, toRefs} from 'vue';
import {api} from 'boot/axios';
import {imageSrc} from 'src/utils';

export default {
    name: 'UserRelation',
    setup() {
        const filterRef = ref(null) as any;
        const treeRef = ref(null) as any;
        const filter = ref('') as any;
        const state = reactive({
            userTree: [] as any,
            queryURL: '/user/relation'
        });

        // 初始化请求数据
        onMounted(() => {
            queryTreeDataFunc();
        });


        // 请求数据
        const queryTreeDataFunc = () => {
            void api.post(state.queryURL).then((res: any) => {
                state.userTree = res
                setTimeout(() => {
                    treeRef.value.expandAll()
                }, 100)
            });
        };

        // 重置过滤器
        const resetFilterFunc = () => {
            filter.value = '';
            filterRef.value.focus();
        };

        // 自定义过滤器
        const filterMethodFunc = (node: any, filter: any) => {
            const filterLower = filter.toLowerCase();
            return node.username && node.username.toLowerCase().indexOf(filterLower) > -1;
        };

        return {
            filter,
            filterRef,
            treeRef,
            imageSrc,
            resetFilterFunc,
            filterMethodFunc,
            ...toRefs(state)
        };
    }
}
</script>

<style scoped>

</style>