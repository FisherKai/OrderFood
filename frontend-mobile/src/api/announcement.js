import request from './request'

export const getActiveAnnouncements = () => {
  return request.get('/announcements/active')
}

export const getAnnouncementDetail = (id) => {
  return request.get(`/announcements/${id}`)
}
