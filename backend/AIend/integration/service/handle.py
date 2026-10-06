import base64
from fastapi import APIRouter, Request
from fastapi import File, UploadFile, Form

import logging

from clients.medical import InsertMedicalLogSQL
from service.llm.invoke import TextAskAnswer, ImageAskAnswerV4, \
    DiagnosisBySymptoms, JudgeMedical, SeeImageDesc, AnswerTeacher, MapReading, GenerateWeeklyReport


def handleApiServer(app):
    @app.get("/")
    async def read_root(message: str = "Hello World"):
        return {"message": message}

    api = APIRouter(
        prefix="/api",
        tags=["agent api 实例"]
    )
    app.include_router(api)

    @api.get('/text')
    async def test(message: str = 'message为空'):
        answer = TextAskAnswer(message).content.strip()
        return {
            "message": answer
        }

    medical = APIRouter(
        prefix="/medical",
        tags=["医疗相关接口"]
    )
    api.include_router(medical)

    @medical.post('/cure')
    async def test_image(file: UploadFile = File(None), ask: str = Form('根据你的理解诊断传入症状'),
                         account: str = Form(None)):
        try:
            if file is None:
                # 纯文字诊断：直接根据用户描述的症状给出诊断
                symptom = ask
                diagnosis = DiagnosisBySymptoms(ask, account=account)
                # 判断 AI 给出的回答是否为医疗诊治内容
                if JudgeMedical(ask, diagnosis):
                    InsertMedicalLogSQL(account, symptom, diagnosis)
                else:
                    print(f"非医疗回答，跳过写库: {diagnosis}", flush=True)

                return {"symptom": symptom, "diagnosis": diagnosis}

            img_b64 = base64.b64encode(await file.read()).decode('utf-8')
            answer = ImageAskAnswerV4(img_b64, '请描述这张图片')
            symptom = answer.get("content") if isinstance(answer, dict) else str(answer)
            diagnosis = DiagnosisBySymptoms(symptom, ask, account)
            # 图片本身就是医疗素材，判断用户问题；ask 为默认值时判断图片描述出的症状
            if JudgeMedical(ask if ask != '根据你的理解诊断传入症状' else symptom, diagnosis):
                InsertMedicalLogSQL(account, symptom, diagnosis)
            else:
                print(f"非医疗回答，跳过写库: {diagnosis}", flush=True)
            return {"symptom": symptom, "diagnosis": diagnosis}

        except Exception as e:
            logging.error(f"图片分析错误: {e}")
            return {"image": "图片分析失败，请检查输入"}

    teacher = APIRouter(
        prefix="/teacher",
        tags=["百科问答"]
    )
    api.include_router(teacher)

    @teacher.post('/api')
    async def teacher_api(file: UploadFile = File(None), ask: str = Form(''), account: str = Form('')):
        try:
            if not account:
                return {"code": 400, "error": "账号为必填项"}
            question = (ask or '').strip()
            if file is None and not question:
                return {"code": 400, "error": "图片和文字问题至少填一项"}

            image_desc = ''
            if file is not None:
                img_b64 = base64.b64encode(await file.read()).decode('utf-8')
                image_desc = SeeImageDesc(img_b64)

            answer = AnswerTeacher(question, image_desc, account)
            return {"code": 200, "answer": answer, "image_desc": image_desc}
        except Exception as e:
            logging.error(f"百科问答错误: {e}")
            return {"code": 500, "error": str(e)}

    map_reader = APIRouter(
        prefix="/map",
        tags=["识图"]
    )
    api.include_router(map_reader)

    @map_reader.post('/recognize')
    async def map_recognize(file: UploadFile = File(None), account: str = Form('')):
        try:
            if file is None:
                return {"code": 400, "error": "请上传图片"}
            img_b64 = base64.b64encode(await file.read()).decode('utf-8')
            result = MapReading(img_b64)
            return {"code": 200, "result": result}
        except Exception as e:
            logging.error(f"识图错误: {e}")
            return {"code": 500, "error": str(e)}


    report = APIRouter(
        prefix="/report",
        tags=["社区周报"]
    )
    api.include_router(report)

    @report.post('/generate')
    async def report_generate(request: Request):
        try:
            data = await request.json()
            community_name = data.get('community_name', '')
            notices = data.get('notices', []) or []
            medicals = data.get('medicals', []) or []
            report_text = GenerateWeeklyReport(community_name, notices, medicals)
            return {"code": 200, "report": report_text}
        except Exception as e:
            logging.error(f"周报生成错误: {e}")
            return {"code": 500, "error": str(e)}
