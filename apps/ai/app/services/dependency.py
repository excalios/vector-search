from typing import Annotated

from app.services.embedding import EmbeddingService
from fastapi import Depends


def get_embedding_service() -> EmbeddingService:
    embedding_service = EmbeddingService()
    return embedding_service


DepEmbeddingService = Annotated[EmbeddingService, Depends(get_embedding_service)]
