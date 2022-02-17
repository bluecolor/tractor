import request from './request'

export default {
  getConnectionTypes() {
    return request.get('/connections/types')
  },
  getProviderTypes() {
    return request.get('/connections/providers/types')
  },
  createConnection(payload) {
    return request.post('/connections', payload)
  },
  getConnections() {
    return request.get('/connections')
  },
  getConnection(id) {
    return request.get(`/connections/${id}`)
  },
  deleteConnection(id) {
    return request.delete(`/connections/${id}`)
  },
  updateConnection(id, payload) {
    return request.put(`/connections/${id}`, payload)
  },
  resolveConnectorRequest(payload) {
    return request.post('/connections/connectors/resolve', payload)
  }
}
