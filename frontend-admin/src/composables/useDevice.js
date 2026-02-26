import { ref, onMounted, onUnmounted } from 'vue'

/**
 * 设备检测 composable
 * 通过 UA 和屏幕宽度判断设备类型
 */
export function useDevice() {
  const isMobile = ref(false)
  const isTablet = ref(false)
  const isDesktop = ref(true)
  const screenWidth = ref(0)

  // 移动端 UA 关键词
  const mobileKeywords = [
    'Android',
    'webOS',
    'iPhone',
    'iPad',
    'iPod',
    'BlackBerry',
    'IEMobile',
    'Opera Mini',
    'Mobile',
    'mobile',
    'CriOS',
    'FxiOS'
  ]

  // 平板 UA 关键词
  const tabletKeywords = [
    'iPad',
    'Android.*Tablet',
    'Tablet',
    'PlayBook',
    'Silk'
  ]

  /**
   * 检测是否为移动设备 (基于 UA)
   */
  const checkMobileByUA = () => {
    const ua = navigator.userAgent
    return mobileKeywords.some(keyword => ua.includes(keyword))
  }

  /**
   * 检测是否为平板设备 (基于 UA)
   */
  const checkTabletByUA = () => {
    const ua = navigator.userAgent
    return tabletKeywords.some(keyword => {
      const regex = new RegExp(keyword, 'i')
      return regex.test(ua)
    })
  }

  /**
   * 更新设备状态
   */
  const updateDeviceStatus = () => {
    screenWidth.value = window.innerWidth

    // 基于 UA 判断
    const isMobileUA = checkMobileByUA()
    const isTabletUA = checkTabletByUA()

    // 基于屏幕宽度判断
    const isMobileWidth = screenWidth.value < 768
    const isTabletWidth = screenWidth.value >= 768 && screenWidth.value < 1024

    // 综合判断：UA 或 屏幕宽度满足即视为移动端/平板
    isMobile.value = isMobileUA || isMobileWidth
    isTablet.value = isTabletUA || isTabletWidth
    isDesktop.value = !isMobile.value && !isTablet.value

    // 如果是平板但使用移动端 UA，也算作移动端
    if (isTablet.value && isMobileUA) {
      isMobile.value = true
    }
  }

  // 防抖处理
  let resizeTimer = null
  const handleResize = () => {
    if (resizeTimer) clearTimeout(resizeTimer)
    resizeTimer = setTimeout(updateDeviceStatus, 100)
  }

  onMounted(() => {
    updateDeviceStatus()
    window.addEventListener('resize', handleResize)
  })

  onUnmounted(() => {
    window.removeEventListener('resize', handleResize)
    if (resizeTimer) clearTimeout(resizeTimer)
  })

  return {
    isMobile,
    isTablet,
    isDesktop,
    screenWidth,
    updateDeviceStatus
  }
}

/**
 * 获取设备类型字符串
 */
export function getDeviceType() {
  const ua = navigator.userAgent
  const width = window.innerWidth

  if (width < 768 || /Mobile|Android|iPhone/i.test(ua)) {
    return 'mobile'
  } else if (width < 1024 || /iPad|Tablet/i.test(ua)) {
    return 'tablet'
  }
  return 'desktop'
}
