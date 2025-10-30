import fastapi
import random
from typing import Optional


app = fastapi.FastAPI(
    title="Temperature API",
    description="API для получения данных о температуре от датчиков",
    version="1.0.0"
)


def get_location_by_sensor_id(sensor_id: str) -> str:
    switch = {
        "1": "Living Room",
        "2": "Bedroom", 
        "3": "Kitchen"
    }
    return switch.get(sensor_id, "Unknown")


def get_sensor_id_by_location(location: str) -> str:
    switch = {
        "Living Room": "1",
        "Bedroom": "2",
        "Kitchen": "3"
    }
    return switch.get(location, "0")


@app.get("/temperature")
def get_temperature(location: Optional[str] = None, sensor_id: Optional[str] = None):
    if not location and not sensor_id:
        raise fastapi.HTTPException(
            status_code=400, 
            detail="Необходимо указать либо location, либо sensor_id"
        )
    
    if not location:
        location = get_location_by_sensor_id(sensor_id)
    
    if not sensor_id:
        sensor_id = get_sensor_id_by_location(location)
    
    temperature = round(random.uniform(10.0, 30.0), 1)
    
    return {
        "temperature": temperature,
        "location": location,
        "sensor_id": sensor_id
    }


@app.get("/temperature/{sensor_id}")
def get_temperature_by_sensor_id(sensor_id: str, location: Optional[str] = None):
    if not location:
        location = get_location_by_sensor_id(sensor_id)
    
    temperature = round(random.uniform(10.0, 30.0), 1)
    
    return {
        "temperature": temperature,
        "location": location,
        "sensor_id": sensor_id
    }


@app.get("/health")
def health_check():
    """Проверка здоровья сервиса"""
    return {"status": "ok"}
