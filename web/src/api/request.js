import axios from 'axios'
import Cookies from 'js-cookie'

const request = axios.create({
  baseURL: 'http://localhost:3000/api/v1',
  timeout: 6000
})

const errorHandler = (error) => {
  if (error.response) {
    const data = error.response.data
    if (error.response.status === 403) {
      console.error(data.message)
    }
    if (error.response.status === 401 && !(data.result && data.result.isLogin)) {
      console.error('Authorization verification failed')
    }
  }
  return Promise.reject(error)
}

// request interceptor
request.interceptors.request.use((config) => {
  config.headers['Content-Type'] = 'application/json'
  config.crossDomain = true
  const token = Cookies.get('access_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
}, errorHandler)

// response interceptor
request.interceptors.response.use((response) => {
  return Promise.resolve(response.data)
}, errorHandler)

export default request
