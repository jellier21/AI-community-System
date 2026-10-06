from PIL import Image
from langchain_core.messages import HumanMessage

from .util.model import GlmModel46_Init, UseToolAnswer, AnalyseImageMessage, PhysicianAgent, JudgeWheaterMedical, \
    SolutionTeacher, SeeImageMessage, MapReader, WeeklyReportAgent
from .tools.tools import build_articledata


def TextAskAnswer(word):
    llm = GlmModel46_Init()
    return llm.invoke(word)


def TextAskAnswerV2(word):
    answer = UseToolAnswer([build_articledata])
    return {
        "answer": answer.invoke(word)
    }


def ImageAskAnswerV4(image_b64, prompt=None, content_type=None):
    llm = AnalyseImageMessage()
    if image_b64:
        mime = content_type or 'image/jpeg'
        content = []
        if prompt:
            content.append({"type": "text", "text": prompt})
        content.append({"type": "image_url",
                        "image_url": {"url": f"data:{mime};base64,{image_b64}"}})
        result = llm.invoke({"messages": [HumanMessage(content=content)]})
        answer = (result["messages"][-1].content
                  if isinstance(result, dict)
                  else getattr(result, "content", str(result)))
        return answer
    return "未传入图片，不要分析，仅回复 <未传入图片> 。"


def DiagnosisBySymptoms(symptoms, question=None, account="default"):
    config = {"configurable": {"thread_id": account}}
    llm = PhysicianAgent()
    content = []
    if account and account != "default":
        content.append({"role": "user", "content": f"用户账号：{account}（可用于查询历史诊断记录，无需展示给用户）"})
    if question:
        content.append({"role": "user", "content": question})
    content.append({"role": "user", "content": symptoms})
    result = llm.invoke({"messages": content}, config=config)
    return result["messages"][-1].content


# 医疗话题关键词（用户输入命中即视为医疗话题）
MEDICAL_KEYWORDS = [
    "疼", "痛", "伤口", "肿", "发烧", "发热", "咳嗽", "头晕", "头痛",
    "恶心", "呕吐", "腹泻", "脖子", "膝盖", "腰", "关节", "肌肉",
    "骨折", "扭伤", "拉伤", "血压", "血糖", "过敏", "感染", "炎症",
    "症状", "诊断", "治疗", "检查", "药", "医院", "医生",
    "腿", "脚", "手", "背", "胃", "心", "眼", "耳", "鼻", "皮肤",
    "失眠", "乏力", "食欲", "结痂", "发炎", "酸痛", "麻木",
]


def JudgeMedical(user_input: str, ai_answer: str) -> bool:
    """入库条件：用户输入是医疗话题，或 AI 回答包含实质性诊断"""
    # 1. 用户输入命中医疗关键词 → 直接入库（问诊过程中的反问轮也保留）
    if user_input and any(k in user_input for k in MEDICAL_KEYWORDS):
        print("用户输入为医疗话题，入库", flush=True)
        return True
    # 2. 用户输入非医疗时，用智能体判断 AI 回答是否实质诊断
    if not ai_answer or not ai_answer.strip():
        return False
    try:
        agent = JudgeWheaterMedical()
        result = agent.invoke(ai_answer)
        if isinstance(result, dict):
            messages = result.get("messages") or []
            answer = messages[-1].content if messages else ""
        else:
            answer = getattr(result, "content", str(result))
        answer = str(answer).strip()
        print("医疗入库判断原始返回:", repr(answer), flush=True)
        if answer.startswith("不是"):
            return False
        if answer.startswith("是"):
            return True
        return False
    except Exception as e:
        print("医疗入库判断失败，保守不入库:", e, flush=True)
        return False


def SeeImageDesc(image_b64):
    """分析图片内容，返回文字描述"""
    llm = SeeImageMessage()
    content = [
        {"type": "image_url", "image_url": {"url": f"data:image/jpeg;base64,{image_b64}"}},
    ]
    result = llm.invoke({"messages": [HumanMessage(content=content)]})
    if isinstance(result, dict):
        messages = result.get("messages") or []
        return messages[-1].content if messages else str(result)
    return getattr(result, "content", str(result))


def AnswerTeacher(question='', image_desc='', account=''):
    """百科问答：图片描述 + 文字问题 + 账号 → 答案（按账号记忆）"""
    llm = SolutionTeacher()
    parts = []
    if account:
        parts.append(f"用户账号：{account}")
    if image_desc:
        parts.append(f"图片内容描述：{image_desc}")
    if question:
        parts.append(f"用户问题：{question}")
    config = {"configurable": {"thread_id": account or "default"}}   # ← 用账号做会话线程
    result = llm.invoke({"messages": [HumanMessage(content="\n".join(parts))]}, config=config)
    if isinstance(result, dict):
        messages = result.get("messages") or []
        return messages[-1].content if messages else str(result)
    return getattr(result, "content", str(result))


def MapReading(image_b64):
    """识图：识别图片内容并给出解读"""
    llm = MapReader()
    content = [
        {"type": "image_url", "image_url": {"url": f"data:image/jpeg;base64,{image_b64}"}},
    ]
    result = llm.invoke({"messages": [HumanMessage(content=content)]})   # ← 必须 dict 格式
    if isinstance(result, dict):
        messages = result.get("messages") or []
        return messages[-1].content if messages else str(result)
    return getattr(result, "content", str(result))


def GenerateWeeklyReport(community_name, notices, medicals):
    """生成社区 AI 周报：公告列表 + 医疗记录列表 → 周报"""
    llm = WeeklyReportAgent()
    if not notices:
        notices_txt = "本周无公告"
    else:
        notices_txt = "\n".join(
            [f"- [{n.get('create_at', '')}] {n.get('title', '')}: {n.get('content', '')}" for n in notices])
    if not medicals:
        medicals_txt = "本周无医疗咨询"
    else:
        medicals_txt = "\n".join(
            [f"- [{m.get('create_at', '')}] 症状: {m.get('symptom', '')} / 诊断: {m.get('diagnosis', '')}" for m in medicals])

    content = (
        f"社区名称：{community_name}\n"
        f"本周公告（共{len(notices)}条）：\n{notices_txt}\n\n"
        f"本周医疗咨询（共{len(medicals)}次）：\n{medicals_txt}"
    )
    result = llm.invoke({"messages": [HumanMessage(content=content)]})   # ← dict 格式
    if isinstance(result, dict):
        messages = result.get("messages") or []
        return messages[-1].content if messages else str(result)
    return getattr(result, "content", str(result))

