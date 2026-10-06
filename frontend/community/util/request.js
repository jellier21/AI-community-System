import { BASEURL_NAME } from '@/util/configuration.js'

const request = ({ url, method = 'GET', data = {}, header = {}, timeout = 10000 } = {}) => {
  return new Promise((resolve, reject) => {
    const token = uni.getStorageSync('jwtToken')
    uni.request({
      url: BASEURL_NAME + url,
      method: method.toUpperCase(),
      header: {
        'Content-Type': 'application/x-www-form-urlencoded',
        ...(token ? { Authorization: token } : {}),
        ...header
      },
      data,
      timeout,
      success: (res) => {
        if (res.statusCode >= 200 && res.statusCode < 300) {
          resolve(res.data)                  // 成功：把数据交出去
        } else if (res.statusCode === 401) {
			console.log("========Token已失效========")
          uni.removeStorageSync('jwtToken')
          uni.navigateTo({ url: '/pages/index/index' })
          reject(new Error('登录已过期'))
        } else {
          reject(new Error(`请求失败：${res.statusCode}`))
        }
      },
      fail: (err) => {
        uni.showToast({ title: '网络异常', icon: 'none' })
        reject(err)                          // 失败：把错误交出去
      }
    })
  })
}

const uploadFile = ({ url, filePath, name = 'file', formData = {}, header = {}, timeout = 30000 } = {}) => {
  return new Promise((resolve, reject) => {
    const token = uni.getStorageSync('jwtToken')

    uni.uploadFile({
      url: BASEURL_NAME + url,
      filePath,
      name,
      formData,               // ← 文字就走这里
      header: {
        ...(token ? { Authorization: token } : {}),
        ...header
      },
      timeout,
      success: (res) => {
        if (res.statusCode >= 200 && res.statusCode < 300) {
          try { resolve(JSON.parse(res.data)) }
          catch { resolve(res.data) }
        } else if (res.statusCode === 401) {
          uni.removeStorageSync('jwtToken')
          uni.reLaunch({ url: '/pages/index/index' })
          reject(new Error('登录已过期'))
        } else {
          reject(new Error(`上传失败：${res.statusCode}`))
        }
      },
      fail: (err) => {
        uni.showToast({ title: '上传失败', icon: 'none' })
        reject(err)
      }
    })
  })
}


const checkJwtTimer = async () => {
	try {
		const res = await request({ url: '/api/jwt/monitor' })
		console.log(res)
		return false          // 正常
	} catch (err) {
		console.log('err:', err)
		return true           // 异常/Token 失效
	}
}
// 调用：const bad = await checkJwtTimer()

export {
	request,
	uploadFile,
	checkJwtTimer
}
