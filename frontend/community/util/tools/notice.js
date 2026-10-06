import {request} from '@/util/request.js'

// 发布公告
const PublishNotice = async (community_uid, title, content)=>{
	const res = await request({
		url: '/api/notice/publish',
		method: 'POST',
		data: { community_uid, title, content }
	})
	return res
}

// 查公告列表
const GetNoticeList = async (community_uid)=>{
	const res = await request({
		url: '/api/notice/list',
		data: { community_uid }
	})
	return res.info || []
}

// 查最新一条
const GetLatestNotice = async (community_uid)=>{
	const res = await request({
		url: '/api/notice/latest',
		data: { community_uid }
	})
	return res.info || null
}

export { PublishNotice, GetNoticeList, GetLatestNotice }
