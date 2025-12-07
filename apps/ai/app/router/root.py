from app.core.logging import DepLogger
from fastapi import APIRouter

router = APIRouter()


@router.get("/")
async def root(logger: DepLogger):
    return {"message": "Welcome to the Machine Learning API"}


@router.get("/health-check")
async def health_check(logger: DepLogger):
    """
    Checks the health of the application and related connection.
    """
    logger.info("Performing health check...")
    logger.info("Health check successful")
    return {"status": "ok"}
