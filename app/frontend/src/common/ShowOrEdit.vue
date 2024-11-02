<template>
  <div @dblclick="handleDblClick" style="min-height: 22px">
    <n-input
        v-if="isEdit"
        ref="inputRef"
        v-model:value="inputValue"
        @keydown.enter="handleSubmit"
        @blur="handleBlur"
    />
    <template v-else>
      {{ value }}
    </template>
  </div>
</template>

<script setup>
import { ref, nextTick } from 'vue'

const props = defineProps({
  value: {
    type: [String, Number, null],
    required: true
  },
  onUpdateValue: {
    type: [Function, Array],
    required: true
  }
})

const isEdit = ref(false)
const inputRef = ref(null)
const inputValue = ref(props.value)

// 双击触发编辑
function handleDblClick() {
  isEdit.value = true
  nextTick(() => {
    inputRef.value?.focus()
  })
}

// 按回车提交
function handleSubmit() {
  props.onUpdateValue(inputValue.value)
  isEdit.value = false
}

// 失去焦点时取消编辑,恢复原值
function handleBlur() {
  isEdit.value = false
  inputValue.value = props.value
}
</script>