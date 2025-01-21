<!--
  - Copyright 2025 Bronya0 <tangssst@163.com>.
  - Author Github: https://github.com/Bronya0
  -
  - Licensed under the Apache License, Version 2.0 (the "License");
  - you may not use this file except in compliance with the License.
  - You may obtain a copy of the License at
  -
  -     https://www.apache.org/licenses/LICENSE-2.0
  -
  - Unless required by applicable law or agreed to in writing, software
  - distributed under the License is distributed on an "AS IS" BASIS,
  - WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  - See the License for the specific language governing permissions and
  - limitations under the License.
  -->

<template>
<!--  https://www.naiveui.com/zh-CN/os-theme/components/menu-->
    <n-menu
        :mode="'vertical'"
        :value="props.value"
        @update:value="handleMenuSelect"
        :options="props.options"
        style="--wails-draggable:no-drag"
    />

</template>


<script setup>

import emitter from "../utils/eventBus";

const props = defineProps(['options', 'value']);

const handleMenuSelect = (key, item) => {
  emitter.emit('menu_select', key)
}

</script>

<style>
.aside-btn .n-button .n-button__icon {
  font-size: 22px;
}
</style>