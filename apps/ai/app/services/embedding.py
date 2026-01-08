from typing import List
from app.core.instrumentation import get_tracer
from app.core.logging import get_logger
from app.core.response import ErrorResponse, SuccessResponse, success_response
import sentence_transformers
from sentence_transformers import SentenceTransformer

from app.schemas.embedding import EmbeddingInput

tracer = get_tracer("service.embedding")


class EmbeddingService:
    __log = get_logger()

    def __init__(self):
        self.general_model = SentenceTransformer(
            "jinaai/jina-embeddings-v2-base-en", # switch to en/zh for English or Chinese
            trust_remote_code=True
        )
        self.specialist_model = SentenceTransformer("ncbi/MedCPT-Query-Encoder")

    def general_embed(self, input: EmbeddingInput) -> SuccessResponse[List[float]] | ErrorResponse:
        with tracer.start_as_current_span("service.embedding") as span:
            self.__log.info("Service layer log", extra={"layer": "service"})
            embeddings = None
            match input.type:
                case "generalist":
                    embeddings = self.general_model.encode(input.sentence)
                case "specialist":
                    embeddings = self.specialist_model.encode(input.sentence)
            return success_response(embeddings.tolist())
