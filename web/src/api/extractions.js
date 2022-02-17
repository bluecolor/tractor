import request from './request'

export default {
  getExtractions() {
    return request.get('/extractions')
  }
}
