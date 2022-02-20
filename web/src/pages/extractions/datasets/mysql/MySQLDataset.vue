<template lang="pug">
//- a-form.connection-options(:label-col='{ span: 8 }', :wrapper-col='{ span: 16 }', labelAlign='left')
//-   a-form-item(label='Database', name='database')
//-     a-select(ref='select', v-model:value='database')
//-       a-select-option(v-for='m in databases', :value='m') {{ m }}
//-   a-form-item(label='Table', name='table')
//-     a-select(ref='select', v-model:value='table')
//-       a-select-option(v-for='m in tables', :value='m') {{ m }}
</template>

<script>
export default {
  props: {
    connection: {
      type: Object,
      required: true
    },
    dataset: {
      type: Object,
      default: () => ({
        name: '',
        fields: '',
        config: {
          database: '',
          table: ''
        }
      })
    }
  },
  data() {}
}
</script>

<script setup>
// import { ref, defineProps, onBeforeMount, watch, defineEmits } from 'vue'
// import { useStore } from 'vuex'

// const props = defineProps({
//   connection: {
//     type: Object,
//     default: () => ({})
//   },
//   database: {
//     type: String,
//     default: ''
//   },
//   table: {
//     type: String,
//     default: ''
//   }
// })

// const store = useStore()
// const database = ref('')
// const table = ref('')
// const databases = ref([])
// const tables = ref([])
// const emit = defineEmits(['set-fields'])

// const fetchDatabases = () => {
//   const payload = {
//     connection: props.connection,
//     request: 'databases'
//   }
//   store
//     .dispatch('connections/resolveConnectorRequest', payload)
//     .then((result) => {
//       databases.value = result
//       if (databases.value.length > 0) {
//         let db = props.connection.config.databases
//         if (result.length > 0) {
//           db = db ?? result[0]
//         }
//         database.value = db
//       }
//     })
//     .catch((error) => {
//       console.error(error)
//     })
// }
// const fetchTables = () => {
//   const payload = {
//     connection: props.connection,
//     request: 'tables',
//     options: { database: database.value }
//   }
//   store
//     .dispatch('connections/resolveConnectorRequest', payload)
//     .then((result) => {
//       tables.value = result
//     })
//     .catch((error) => {
//       console.error(error)
//     })
// }
// const fetchFields = () => {
//   const connectionId = props.connection.id
//   const d = database.value
//   const t = table.value

//   if (!connectionId || !d || !t) {
//     return
//   }
//   store
//     .dispatch('connections/getFields', {
//       connectionId,
//       database: d,
//       table: t
//     })
//     .then((result) => {
//       console.log(result)
//     })
//     .catch((error) => {
//       console.error(error)
//     })
// }

// watch(database, fetchTables)
// watch(table, fetchFields)

// onBeforeMount(() => {
//   fetchDatabases()
// })
</script>

<style lang="scss">
.connection-options {
  .ant-form-item {
    margin-bottom: 12px;
  }
}
</style>
