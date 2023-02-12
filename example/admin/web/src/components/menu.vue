<template>
    <q-list v-for='(item, index) in data' :key='index' :class='listClass'>
        <q-expansion-item :icon='item.data.icon' :label='item.name' :default-opened='expansionOpenedFunc(item.children)'
                          header-class='text-body2' expand-icon-class='text-grey'
                          v-if='item.hasOwnProperty("children") && item.children !== null && item.children.length > 0'>
            <menu-list :data='item.children' :inset-level='insetLevel+0.5' :list-class='listClass'
                        :active-class='activeClass'></menu-list>
        </q-expansion-item>
        <q-item clickable v-ripple v-else :inset-level='insetLevel' @click='toRouterFunc(item)'
                :active='item.route === $route.fullPath' :active-class='activeClass'>
            <q-item-section avatar v-if='item.data.icon !== ""'>
                <q-icon :name='item.data.icon'></q-icon>
            </q-item-section>
            <q-item-section>
                <div class='text-body2'>{{ item.name }}</div>
            </q-item-section>
        </q-item>
    </q-list>
</template>

<script lang='ts'>
import { defineComponent } from 'vue';

export default defineComponent({
    name: 'MenuList',
    props: {
        data: { type: Array, required: true },
        insetLevel: { type: Number, default: 0 },
        listClass: { type: String, default: '' },
        activeClass: { type: String, default: 'bg-blue-8 text-white' }
    },
    methods: {
        //  判断是否默认打开
        expansionOpenedFunc(children: any) {
            if (children === null) {
                return false;
            }
            let findMenu = children.find((item: any) => {
                return item.route === this.$route.path;
            });
            return findMenu !== undefined;
        },

        // 跳转路由方案
        toRouterFunc(item: any) {
            if (this.$route.path === item.route.split('?')[0]) {
                void this.$router.replace({ name: 'Replace', query: { path: item.route } });
            } else {
                void this.$router.push(item.route);
            }
        }
    }
});
</script>

<style scoped>

</style>
