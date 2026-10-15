import os

from langchain.agents import create_agent
from langchain_openai import ChatOpenAI
from langgraph.checkpoint.memory import InMemorySaver

from service.llm.tools.tools import ReadMedicalLogClient


def ImageModel_Init():
    return ChatOpenAI(
        model=os.getenv("IMAGE_MODEL_NAME"),
        base_url=os.getenv("IMAGE_BASEURL"),
        api_key=os.getenv("IMAGE_KEY"),
        max_tokens=1200,
        timeout=180,
        max_retries=1,
        temperature=0,
    )


def TextModel_Init():
    return ChatOpenAI(
        model=os.getenv("TEXTMODEL_NAME"),
        base_url=os.getenv("TEXT_BASEURL"),
        api_key=os.getenv("TEXT_KEY"),
        max_tokens=2048,
        timeout=180,
        max_retries=1,
        temperature=0,
    )


def UseToolAnswer(tools: list):
    print(f"注册工具: {tools}")
    model = ImageModel_Init()
    return create_agent(
        model=model,
        tools=tools,
        system_prompt="你是一个乐于助人的助手，一直使用工具build_articledata来回答问题。",
    )


def AnalyseImageMessage():
    model = ImageModel_Init()
    return create_agent(
        model=model,
        system_prompt="【角色】：你是一个观察分析大师，并且外科医学知识丰富，识别图片内容并给出伤口的详细描述，包括病状、面积、大小、颜色等。"
                      "【硬性规则】：仅描述伤口内容即可，无需过度分析，不要添加任何标记符号。",
    )


# 模块级共享检查点：所有请求复用同一实例，记忆才能跨请求生效
_memory_checkpointer = InMemorySaver()

_physician_agent = None


def PhysicianAgent():
    global _physician_agent
    if _physician_agent is None:  # 只编译一次，避免每次请求重新编译
        model = TextModel_Init()
        _physician_agent = create_agent(
            model=model,
            tools=[ReadMedicalLogClient],
            system_prompt=(
                "【角色】：你是一个诊断大师，根据传入的患者的症状和体征和需求，给出诊断结果和治疗方案。"
                "【硬性规则】："
                "1. 用户的账号已在对话消息中提供（格式为『用户账号：xxx』）。诊断前可调用 ReadMedicalLogClient 工具查询该账号的历史诊断记录；若记录为空则证明为首次诊断。"
                "2. 若工具查询失败、拿不到账号或没有历史记录，直接根据当前症状诊断，【绝不】向用户索要账号。"
                "3. 根据诊断结果和患者需求，给出诊断方案和建议。"
                "4. 不要过多生成！"
            ),
            checkpointer=_memory_checkpointer  # ← 关键：配上检查点，thread_id 才会生效
        )
    return _physician_agent


def JudgeWheaterMedical():
    model = TextModel_Init()
    return create_agent(
        model=model,
        system_prompt=(
            "【角色】：你是一个判断分析者，判断一段医疗咨询回复是否包含实质性的诊断结果或治疗建议。"
            "【判定标准】："
            "1. 回复中给出了具体的病因分析、可能原因、诊断结论、治疗方案、用药或就医建议等实质性内容 → 返回『是』"
            "2. 回复只是在反问用户（如『请问您需要什么帮助』『请告诉我具体症状』）、追问更多信息、或没有给出任何诊断分析 → 返回『否』"
            "【硬性规则】：只输出一个汉字『是』或『否』，禁止输出标点、空格、解释或任何其他内容。"
        ),

    )


def SeeImageMessage():
    """图片内容分析 agent：识别图片内容并给出文字描述"""
    model = ImageModel_Init()
    return create_agent(
        model=model,
        system_prompt="【角色】：你是一个图像理解专家，能够识别图片内容并给出详细文字描述，包括主体、场景、文字、物体等。"
                      "【硬性规则】：仅描述图片内容，不要添加分析、建议或评价。",
    )


_solution_teacher_agent = None


def SolutionTeacher():
    """百科问答 agent（带记忆：按账号 thread_id 保存会话历史）"""
    global _solution_teacher_agent
    if _solution_teacher_agent is None:  # 只编译一次，共享检查点记忆才生效
        model = TextModel_Init()
        _solution_teacher_agent = create_agent(
            model=model,
            system_prompt="【角色】：你是一个知识渊博的百科问答专家，能够解答科学、历史、文化、技术、生活常识等各种知识性问题。"
                          "【硬性规则】：1. 答案准确、清晰、有条理。2. 结合提供的图片描述和用户问题作答。3. 不确定的内容如实说明，不编造。",
            checkpointer=_memory_checkpointer  # ← 关键：配上检查点，thread_id 生效
        )
    return _solution_teacher_agent


def MapReader():
    """图片识别 agent：识别图片内容并给出详细解读"""
    model = ImageModel_Init()
    return create_agent(
        model=model,
        system_prompt="【角色】：你是一个图片识别专家，能够全面识别并解读图片内容。"
                      "【输出要求】：1. 识别图片中的主体、场景、人物/物体、文字信息等，用清晰条理的中文描述。"
                      "2. 指出关键细节；如果是地图、路牌、票据、植物、动物等特殊图片，给出对应解读。"
                      "【硬性规则】：仅基于图片内容回答，不确定的内容如实说明，不编造。",
    )


def WeeklyReportAgent():
    """社区周报生成 agent"""
    model = TextModel_Init()
    return create_agent(
        model=model,
        system_prompt="【角色】：你是一个社区运营分析专家，根据社区一周的公告记录和成员医疗咨询记录，生成社区周报。"
                      "【任务】：1. 总结本周社区动态（公告发布的主题和重点）。"
                      "2. 总结成员健康关注点（常见的症状、诊断倾向、咨询次数）。"
                      "3. 给出下周社区管理建议（健康宣教方向、活动组织、需要关注的成员健康问题）。"
                      "【格式】：用 markdown 分三部分输出："
                      "## 📊 本周概览\n## 🏥 健康关注\n## 💡 下周建议"
                      "【硬性规则】：所有数据必须来自输入内容，不要编造；数据为空时如实说明。",
    )
