import request from './request'

export const createReview = (data) => {
  return request.post('/reviews', data)
}

export const getDishReviews = (dishId, params) => {
  return request.get(`/reviews/${dishId}`, { params })
}
