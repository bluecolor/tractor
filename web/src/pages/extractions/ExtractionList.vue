<template lang="pug">
PageTitle(title='Extractions')
.no-extractions(v-if='!status.loading && extractions.length === 0')
  .message
    | Yo don't have any extractions yet.
  a-button(type='primary', size='large', @click='$router.push("/extractions/new")')
    | Add Connection
.loading(v-if='status.loading')
  .message
    | Loading...

.connection-list(v-if='!status.loading && extractions.length > 0')
  a-row(type='flex', justify='end', style='margin-bottom: 20px')
    a-col(:span='4', style='text-align: right')
      a-button(type='primary', @click='$router.push("/extractions/new")')
        | Add Extraction
  a-table(:dataSource='extractions', :columns='columns')
</template>

<script setup>
import { onBeforeMount, ref, markRaw } from 'vue'
import { useStore } from 'vuex'
import { message } from 'ant-design-vue'
import { DeleteOutlined, ExceptionOutlined, ExperimentOutlined } from '@ant-design/icons-vue'
import PageTitle from '@/components/PageTitle.vue'
const store = useStore()
const columns = markRaw([
  {
    title: 'Name',
    dataIndex: 'name',
    key: 'name',
    customRender: ({ text, record }) => {
      return <a href={`/extractions/${record.id}`}>{text}</a>
    }
  }
])
const extractions = ref([])
const status = ref({
  loading: false,
  type: undefined,
  message: undefined
})

onBeforeMount(() => {
  fetchExtractions()
})
const fetchExtractions = () => {
  status.value.loading = true
  store
    .dispatch('extractions/getExtractions')
    .then((result) => {
      extractions.value = result
    })
    .catch((err) => {
      const m = err?.response?.data?.message || 'Failed to get extractions'
      message.error(m)
    })
    .finally(() => {
      status.value.loading = false
    })
}
const removeConnection = (id) => {
  const index = extractions.value.findIndex((c) => c.id === id)
  if (index === -1) {
    return
  }
  extractions.value.splice(index, 1)
}
const deleteConnection = (id) => {
  status.value.loading = true
  store
    .dispatch('extractions/deleteExtraction', id)
    .then(() => {
      removeConnection(id)
      message.success('Extraction deleted')
    })
    .catch((err) => {
      const m = err?.response?.data?.message || 'Failed to delete connection'
      message.error(m)
    })
    .finally(() => {
      status.value.loading = false
    })
}
const onDelete = ({ id }) => {
  deleteConnection(id)
}
</script>

<style lang="scss" scoped>
.no-extractions,
.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  .message {
    margin-bottom: 1rem;
    font-size: 1.2rem;
    color: #999;
  }
}
</style>
