from typing import List
from app.core.instrumentation import get_tracer
from app.core.response import ErrorResponse, SuccessResponse
from app.services.dependency import DepEmbeddingService
from app.schemas.embedding import EmbeddingInput
from fastapi import APIRouter

router = APIRouter(
    prefix="/embedding",
    tags=["embedding"],
)

tracer = get_tracer("api.embedding")


@router.post(
    "/general",
    responses={
        200: {"model": SuccessResponse[List[float]]},
        400: {"model": ErrorResponse},
        502: {"model": ErrorResponse},
        503: {"model": ErrorResponse},
    },
)
async def general(service: DepEmbeddingService, input: EmbeddingInput):
    with tracer.start_as_current_span("route.embedding.general") as span:
        span.set_attribute("http.method", "GET")
        span.set_attribute("http.route", "/general")
        embeddings = await service.general_embed(input.sentence)
        return embeddings
