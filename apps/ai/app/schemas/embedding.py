from typing import List
from pydantic import BaseModel

class EmbeddingInput(BaseModel):
    sentence: str

class EmbeddingOutput(BaseModel):
    embeddings: List[float]
