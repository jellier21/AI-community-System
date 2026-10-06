from langchain.tools import tool
from clients.medical import InsertMedicalLogSQL, readMedicalLog


@tool
def build_articledata(title: str, content: str) -> str:
    """生成有标题和内容的文章
    Args:
        title: 文章的标题
        content: 文章的内容
    """
    print("title:", title, "\ncontent:", content, flush=True)
    return f"标题：{title}\n内容：{content}"


@tool
def InsertMedicalClient(account: str, symptom: str, diagnosis: str) -> str:
    """诊断完成后，将症状和诊断结果插入数据库
    Args:
        account: 患者的账号
        symptom: 患者的症状描述文本。
        diagnosis: 大模型给出的最终诊断结果。
    """
    result = InsertMedicalLogSQL(account, symptom, diagnosis)
    print("调用InsertMedicalLogSQL结果:", result, flush=True)
    return f"诊断记录已保存: {result}"


@tool
def ReadMedicalLogClient(account: str) -> str:
    """为用户诊断时读取账号account的用户诊断记录
    Args:
        account: 患者的账号
    Returns:
        str: 患者的json格式诊断记录
    """
    result = readMedicalLog(account)
    print("调用readMedicalLog结果:", result, flush=True)
    return f"诊断记录: {result}"
