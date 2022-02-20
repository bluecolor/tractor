<template lang="pug">
a-form.connection-options(:label-col='{ span: 8 }', :wrapper-col='{ span: 16 }', labelAlign='left')
  a-form-item(label='Delimiter', name='delimiter')
    a-input(v-model:value='delimiter')
  a-form-item(label='Skip rows', name='skipRows')
    a-input-number(v-model:value='skipRows', :min='0', :max='100', style='width: 100%')
</template>

<script setup>
import { defineProps } from 'vue'
const props = defineProps({
  delimiter: {
    type: String,
    default: ';'
  },
  skipRows: {
    type: Number,
    default: 0
  }
})
</script>
