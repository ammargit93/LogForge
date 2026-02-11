from fastapi import APIRouter
from app.models import Log

router = APIRouter(prefix="/logs", tags=["Logs"])


@router.post("/")
async def create_log(message: str):
    log = await Log.create(message=message)
    return {"id": log.id, "message": log.message}


@router.get("/")
async def list_logs():
    logs = await Log.all()
    return logs
