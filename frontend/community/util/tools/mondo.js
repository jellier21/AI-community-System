import {request, uploadFile} from '@/util/request.js'

const SubmitTeacherQuestion = async (ques, file)=>{
	// 判空：文字和图片至少有一个
	const question = (ques || '').trim()
	if (question === '' && !file) {
		uni.showToast({ title: '请先输入内容或选择图片', icon: 'none' })
		return false
	}

	// 情况1：只有文字 → 普通 POST
	if (!file) {
		const res = await request({
			url: '/api/agent/teacher/ask',
			method: 'POST',
			data: { question: question, account: uni.getStorageSync('user_account') },
			timeout: 120000
		})
		return res
	}

	// 情况2：有图片（文字可有可无）→ uploadFile 传文件
	const res = await uploadFile({
		url: '/api/agent/teacher/ask',
		filePath: file,
		name: 'file',
		formData: { question: question, account: uni.getStorageSync('user_account') },
		timeout: 120000
	})
	return res
}

export{
	SubmitTeacherQuestion,
}
