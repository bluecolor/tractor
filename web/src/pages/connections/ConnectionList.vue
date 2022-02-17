<template lang="pug">
PageTitle(title='Connections')
.no-connections(v-if='!status.loading && connections.length === 0')
  .message
    | Yo don't have any connections yet.
  a-button(type='primary', size='large', @click='$router.push("/connections/new")')
    | Add Connection
.loading(v-if='status.loading')
  .message
    | Loading...

.connection-list(v-if='!status.loading && connections.length > 0')
  a-row(type='flex', justify='end', style='margin-bottom: 20px')
    a-col(:span='4', style='text-align: right')
      a-button(type='primary', @click='$router.push("/connections/new")')
        | Add Connection
  a-table(:dataSource='connections', :columns='columns')
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
      return <a href={`/connections/${record.id}`}>{text}</a>
    }
  },
  {
    title: 'Type',
    dataIndex: ['connectionType', 'name'],
    key: 'connectionType'
  },
  {
    title: 'Source',
    dataIndex: 'asSource',
    key: 'asSource',
    customRender: ({ text, record }) => {
      return (
        <span>
          <a-tag color={record.asSource ? 'green' : 'red'} key={record.id}>
            {record.asSource ? 'Yes' : 'No'}
          </a-tag>
        </span>
      )
    }
  },
  {
    title: 'Target',
    dataIndex: 'asTarget',
    key: 'asTarget',
    customRender: ({ text, record }) => {
      return (
        <span>
          <a-tag color={record.asSource ? 'green' : 'red'} key={record.id}>
            {record.asSource ? 'Yes' : 'No'}
          </a-tag>
        </span>
      )
    }
  },
  {
    title: '',
    key: 'delete',
    align: 'center',
    dataIndex: 'id',
    customRender: ({ text, record }) => {
      return (
        <a-button-group>
          <a-tooltip title="Test">
            <a-button
              onClick={() => onDelete(record)}
              size="small"
              icon={<ExperimentOutlined />}
            ></a-button>
          </a-tooltip>
          <a-tooltip title="Delete" placement="right">
            <a-popconfirm
              title="Are you sure delete this connection?"
              ok-text="Yes"
              cancel-text="No"
              onConfirm={() => onDelete(record)}
            >
              <a-button size="small" icon={<DeleteOutlined style="color:red" />}></a-button>
            </a-popconfirm>
          </a-tooltip>
        </a-button-group>
      )
    }
  }
])
const connections = ref([])
const status = ref({
  loading: false,
  type: undefined,
  message: undefined
})

onBeforeMount(() => {
  fetchConnections()
})
const fetchConnections = () => {
  status.value.loading = true
  store
    .dispatch('connections/getConnections')
    .then((result) => {
      connections.value = result
    })
    .catch((err) => {
      const m = err?.response?.data?.message || 'Failed to get connections'
      message.error(m)
    })
    .finally(() => {
      status.value.loading = false
    })
}
const removeConnection = (id) => {
  const index = connections.value.findIndex((c) => c.id === id)
  if (index === -1) {
    return
  }
  connections.value.splice(index, 1)
}
const deleteConnection = (id) => {
  status.value.loading = true
  store
    .dispatch('connections/deleteConnection', id)
    .then(() => {
      removeConnection(id)
      message.success('Connection deleted')
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
.no-connections,
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
