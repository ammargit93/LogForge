from tortoise import fields
from tortoise.models import Model


class Log(Model):
    id = fields.IntField(pk=True)
    log_name = fields.TextField()
    log_size = fields.IntField()
    log_s3_key = fields.TextField()
    log_s3_id = fields.TextField()
    created_at = fields.DatetimeField(auto_now_add=True)
    
