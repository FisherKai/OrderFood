import request from './request'

export const getDishes = (params) => {
  return request.get('/dishes', { params })
}

export const getDishDetail = (id) => {
  return request.get(`/dishes/${id}`)
}

export const getCategories = () => {
  return request.get('/categories')
}
