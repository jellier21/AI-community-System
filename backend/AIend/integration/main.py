import os
from fastapi import FastAPI
from dotenv import load_dotenv

load_dotenv()

from service.handle import handleApiServer

app = FastAPI()

handleApiServer(app)
