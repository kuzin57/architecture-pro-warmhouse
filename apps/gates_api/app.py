import fastapi
from fastapi import HTTPException
from pydantic import BaseModel
from enum import Enum
import random

app = fastapi.FastAPI(title="Gates API", description="API для управления воротами")

class GateState(str, Enum):
    OPEN = "open"
    CLOSED = "closed"

class GateStatus(BaseModel):
    is_active: bool
    state: GateState | None = None

class ChangeStateRequest(BaseModel):
    state: GateState

# Имитация состояния ворот в памяти
gate_state = {
    "is_active": True,
    "state": GateState.CLOSED
}

@app.get("/status", response_model=GateStatus)
def get_gate_status():
    """
    Получить статус ворот
    """
    if not gate_state["is_active"]:
        return GateStatus(is_active=False, state=None)
    
    return GateStatus(
        is_active=gate_state["is_active"],
        state=gate_state["state"]
    )

@app.post("/changestate")
def change_gate_state(request: ChangeStateRequest):
    """
    Изменить состояние ворот
    """
    if not gate_state["is_active"]:
        raise HTTPException(
            status_code=400, 
            detail="Датчик ворот неактивен. Невозможно изменить состояние."
        )
    
    # Имитация случайной ошибки при изменении состояния (5% вероятность)
    if random.random() < 0.05:
        raise HTTPException(
            status_code=500,
            detail="Ошибка при изменении состояния ворот. Попробуйте еще раз."
        )
    
    gate_state["state"] = request.state
    
    return {
        "message": f"Состояние ворот успешно изменено на {request.state.value}",
        "new_state": request.state.value
    }

@app.post("/activate")
def activate_gate():
    """
    Активировать датчик ворот
    """
    gate_state["is_active"] = True
    gate_state["state"] = GateState.CLOSED  # По умолчанию закрыты
    
    return {
        "message": "Датчик ворот активирован",
        "is_active": True
    }

@app.post("/deactivate")
def deactivate_gate():
    """
    Деактивировать датчик ворот
    """
    gate_state["is_active"] = False
    gate_state["state"] = None
    
    return {
        "message": "Датчик ворот деактивирован",
        "is_active": False
    }

@app.get("/health")
def health_check():
    """
    Проверка здоровья сервиса
    """
    return {"status": "healthy", "service": "gates_api"}
