import {ref} from 'vue'
import {request,uploadFile} from '@/util/request.js'
import {BASEURL_NAME} from "@/util/configuration.js"


const SubmitMyQuestion = async (ques, file)=>{
	// 判空：文字和图片至少有一个
	const question = (ques || '').trim()
	if (question === '' && !file) {
		uni.showToast({ title: '请先输入内容或选择图片', icon: 'none' })
		return false
	}

	// 情况1：只有文字（没有图片）→ 走普通 POST，不调 uploadFile
	if (!file) {
		const res = await request({
				url: '/api/agent/medical/diagnosis',
				method: 'POST',
				data: { question: question, account: uni.getStorageSync('user_account') },
				timeout: 120000          // 诊断要等大模型，别用默认 10 秒
			})
		if (res && res.response) {
			console.log("---------", res.response)
		}
		return res
	}

	// 情况2：有图片（文字可有可无）→ 走 uploadFile
	const res = await uploadFile({
		url: '/api/agent/medical/diagnosis',
		filePath: file,
		name: 'file',
		formData: { question: question, account: uni.getStorageSync('user_account') },
		timeout: 120000          // 同上，30 秒也不够
	})
	if (res && res.response) {
		console.log("---------", res.response)
	}
	return res
}


export{
	SubmitMyQuestion,
}
