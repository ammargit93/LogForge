from fastapi import FastAPI
from tortoise.contrib.fastapi import register_tortoise
from app.routes import logs
import os
from dotenv import load_dotenv

load_dotenv()

app = FastAPI(title="LogForge API")

app.include_router(logs.router)

register_tortoise(
    app,
    db_url=f"postgres://{os.getenv('DB_USER')}:{os.getenv('DB_PASSWORD')}@localhost:5432/{os.getenv('DB_NAME')}",
    modules={"models": ["app.models"]},
    generate_schemas=True,
    add_exception_handlers=True,
)

@app.get("/")
async def root():
    return {"status": "LogForge running"}
