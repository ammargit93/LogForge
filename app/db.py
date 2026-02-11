from tortoise import Tortoise
from dotenv import load_dotenv
import os

load_dotenv()

async def init_db():
    await Tortoise.init(
        db_url=f"postgres://{os.getenv('DB_USER')}:{os.getenv('DB_PASSWORD')}@localhost:5432/{os.getenv('DB_NAME')}",
        modules={"models": ["app.models"]},
    )
    await Tortoise.generate_schemas()
