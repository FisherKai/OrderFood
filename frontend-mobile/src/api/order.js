import request from './request'

export const createOrder = (data) => {
  return request.post('/orders', data)
}

export const createReservation = (data) => {
  return request.post('/orders/reserve', data)
}

export const getUserOrders = (params) => {
  return request.get('/orders', { params })
}
