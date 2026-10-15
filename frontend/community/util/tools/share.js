import {request} from '@/util/request.js'

const PublishShare = async (community_uid, content, images='')=>{
	return await request({
		url: '/api/share/publish',
		method: 'POST',
		data: { community_uid, content, images }
	})
}

const GetMyShares = async ()=>{
	const res = await request({ url: '/api/share/my' })
	return res.info || []
}

const GetCommunityShares = async (community_uid)=>{
	const res = await request({ url: '/api/share/list', data: { community_uid } })
	return res.info || []
}

const GetLatestShares = async ()=>{
	const res = await request({ url: '/api/share/latest' })
	return res.info || []
}

const DeleteShare = async (id)=>{
	return await request({
		url: '/api/share/delete',
		method: 'POST',
		data: { id }
	})
}

const UpdateShare = async (id, content, images='')=>{
	return await request({
		url: '/api/share/update',
		method: 'POST',
		data: { id, content, images }
	})
}

// ====== 七牛云上传 ======
// 选一张本地图片并上传，返回完整 URL
const uploadOneImage = async (filePath)=>{
	// 1. 拿 uptoken
	const tokenRes = await request({ url: '/api/image/upload/token' })
	const uptoken = tokenRes.uptoken
	const domain = tokenRes.domain

	// 2. 直传七牛云（华北 z2 现行上传域名 up-z2.qiniup.com）
	return new Promise((resolve, reject)=>{
		uni.uploadFile({
			url: 'https://up-z2.qiniup.com',
			filePath: filePath,
			name: 'file',
			formData: { token: uptoken, key: 'share/' + Date.now() + '.jpg' },
			success: (uploadRes)=>{
				console.log('七牛返回 status:', uploadRes.statusCode, 'body:', uploadRes.data)
				let data = {}
				try{
					data = typeof uploadRes.data === 'string' ? JSON.parse(uploadRes.data) : uploadRes.data
				}catch(e){
					reject(new Error('七牛返回非JSON: ' + uploadRes.data))
					return
				}
				if(uploadRes.statusCode !== 200){
					reject(new Error('上传失败: ' + (data.error || uploadRes.data)))
					return
				}
				const key = data.key || data.hash
				if(!key){
					reject(new Error('七牛返回无key: ' + JSON.stringify(data)))
					return
				}
				let fullDomain = domain.replace(/\/$/, '')
				if(!/^https?:\/\//.test(fullDomain)) fullDomain = 'https://' + fullDomain
				resolve(fullDomain + '/' + key)
			},
			fail: (err)=>{
				console.log('七牛uploadFile fail:', err)
				reject(err)
			}
		})
	})
}

// 选多张图，全部上传，返回 URL 数组
const chooseAndUploadImages = async (count=9)=>{
	const chooseRes = await new Promise((resolve, reject)=>{
		uni.chooseImage({
			count: count,
			success: resolve,
			fail: reject
		})
	})
	const urls = []
	for(const path of chooseRes.tempFilePaths){
		uni.showLoading({ title: '上传中...' })
		const url = await uploadOneImage(path)
		urls.push(url)
	}
	uni.hideLoading()
	return urls
}

export {
	PublishShare, GetMyShares, GetCommunityShares, GetLatestShares,
	DeleteShare, UpdateShare,
	chooseAndUploadImages, uploadOneImage
}

