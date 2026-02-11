from fastapi import APIRouter, UploadFile, File
from app.models import Log
from dotenv import load_dotenv
import boto3
import uuid
import os

load_dotenv()
router = APIRouter(prefix="/logs", tags=["Logs"])

s3 = boto3.client(
    "s3",
    region_name=os.getenv("AWS_REGION", "ap-south-1")
)

BUCKET_NAME = os.getenv("AWS_BUCKET_NAME")

@router.post("/")
async def create_log(file: UploadFile = File(...)):

    unique_id = str(uuid.uuid4())
    s3_key = f"logs/{unique_id}_{file.filename}"

    s3.upload_fileobj(file.file, BUCKET_NAME, s3_key)

    log = await Log.create(
        log_name=file.filename,
        log_size=file.size,
        log_s3_key=s3_key,
        log_s3_id=unique_id,
    )

    return log
