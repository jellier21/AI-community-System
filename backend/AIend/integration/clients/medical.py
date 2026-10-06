import os
import requests

base_url = os.getenv("GIN_BASEURL")
headers = {
    "Content-Type": "application/json",
}


def InsertMedicalLogSQL(account, symptom, diagnosis: str):
    data = {
        "account": account,
        "symptom": symptom,
        "diagnosis": diagnosis,
    }
    response = requests.post(url=f"{base_url}/agent/medical/add", headers=headers, json=data)
    if response.status_code == 200:
        print(response.json())
        return response.json()
    else:
        print(f"Error: {response.status_code}")
        return None


def readMedicalLog(account: str):
    response = requests.get(url=f"{base_url}/agent/medical/logs?account={account}")
    if response.status_code == 200:
        data = response.json()
        print(data)
        logs = data.get("medicalLogs", []) if isinstance(data, dict) else data
        result = []
        for log in logs:
            result.append({
                "symptom": log["symptom"],
                "diagnosis": log["diagnosis"],
                "create_at": log["create_at"]
            })
        return result
    else:
        print(f"Error: {response.status_code}")
        return None


