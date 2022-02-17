<template lang="pug">
a-form.connection-options(:label-col='{ span: 8 }', :wrapper-col='{ span: 16 }', labelAlign='left')
  a-form-item(
    label='Database',
    name='database',
    :rules='[{ required: true, message: "Please enter database name" }]'
  )
    a-select(ref='select', v-model:value='database')
      a-select-option(v-for='m in databases', :value='m') {{ m }}
  a-form-item(
    label='Table',
    name='table',
    :rules='[{ required: true, message: "Please enter table name" }]'
  )
    a-input(v-model:value='table')
</template>

<script setup>
import { ref, defineProps, onBeforeMount } from 'vue'
import { useStore } from 'vuex'

const props = defineProps({
  connection: {
    type: Object,
    default: () => ({})
  },
  database: {
    type: String,
    default: ''
  },
  table: {
    type: String,
    default: ''
  }
})

const store = useStore()
const database = ref('append')
const table = ref('')
const databases = ref([])

const fetchDatabases = () => {
  const payload = {
    connection: props.connection,
    request: 'databases'
  }
  store
    .dispatch('connections/resolveConnectorRequest', payload)
    .then((result) => {
      console.log(result)
      // this.databases.value = databases
    })
    .catch((error) => {
      console.error(error)
    })
}

onBeforeMount(() => {
  console.log(props.connection)
  fetchDatabases()
})
</script>

<style lang="scss">
.connection-options {
  .ant-form-item {
    margin-bottom: 12px;
  }
}
</style>
