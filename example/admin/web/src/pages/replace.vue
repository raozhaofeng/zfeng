<template>
    <div></div>
</template>

<script lang='ts'>
import { onBeforeMount } from 'vue';
import Router from 'src/router';

export default {
    name: 'ReplacePage',
    setup() {
        onBeforeMount(() => {
            const queryPath = Router.currentRoute.value.query.path as any;
            void Router.replace(queryPath);
        });
    }
};
</script>

<style scoped>

</style>
