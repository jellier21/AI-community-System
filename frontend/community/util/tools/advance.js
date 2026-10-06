import {request} from '@/util/request.js'

// 提交建议
const SubmitAdvance = async (community_uid, title, content)=>{
	const res = await request({
		url: '/api/advance/submit',
		method: 'POST',
		data: { community_uid, title, content }
	})
	return res
}

// 查看建议列表
const GetAdvanceList = async (community_uid)=>{
	const res = await request({
		url: '/api/advance/list',
		data: { community_uid }
	})
	return res.info || []
}

export { SubmitAdvance, GetAdvanceList }
