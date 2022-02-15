/* eslint-disable camelcase */
import request from './request'

export default {
  getConnectionTypes() {
    return request.get('/connections/types')
  },
  getProviderTypes() {
    return request.get('/connections/providers/types')
  }
}
