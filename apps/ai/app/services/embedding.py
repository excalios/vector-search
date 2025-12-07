from typing import List
from app.core.instrumentation import get_tracer
from app.core.logging import get_logger
from app.core.response import ErrorResponse, SuccessResponse, success_response
from sentence_transformers import SentenceTransformer

tracer = get_tracer("service.embedding")


class EmbeddingService:
    __log = get_logger()

    def __init__(self):
        self.general_model = SentenceTransformer(
            "jinaai/jina-embeddings-v2-base-en", # switch to en/zh for English or Chinese
            trust_remote_code=True
        )

    async def general_embed(self, sentence: str) -> SuccessResponse[List[float]] | ErrorResponse:
        with tracer.start_as_current_span("service.embedding") as span:
            self.__log.info("Service layer log", extra={"layer": "service"})
            embeddings = self.general_model.encode(sentence)
            return success_response(embeddings.tolist())
