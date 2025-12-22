from typing import List, Literal
from pydantic import BaseModel

class EmbeddingInput(BaseModel):
    sentence: str
    type: Literal['generalist'] | Literal['specialist']

class EmbeddingOutput(BaseModel):
    embeddings: List[float]
