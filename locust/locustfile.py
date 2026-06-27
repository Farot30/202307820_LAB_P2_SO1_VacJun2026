from locust import HttpUser, task, between
import random
import datetime

# Equipos
TEAMS = ["GTM", "MEX", "BRA", "ARG", "ESP"]

class MundialUser(HttpUser):
    # Tiempo de espera entre peticiones de un usuario simulado
    wait_time = between(1, 3)

    @task
    def enviar_prediccion(self):
        # 1. Seleccionar equipos (local y visitante no pueden ser iguales)
        home_team = random.choice(TEAMS)
        away_team = random.choice([t for t in TEAMS if t != home_team])

        # 2. Generar goles aleatorios (rango sugerido: 0 a 5)
        home_goals = random.randint(0, 5)
        away_goals = random.randint(0, 5)

        # 3. Generar usuario aleatorio ("user_N" entre 1 y 1000)
        username = f"user_{random.randint(1, 1000)}"

        # 4. Generar timestamp actual en formato ISO
        timestamp = datetime.datetime.utcnow().isoformat() + "Z"

        # Construir el JSON exactamente como lo pide el enunciado
        payload = {
            "home_team": home_team,
            "away_team": away_team,
            "home_goals": home_goals,
            "away_goals": away_goals,
            "username": username,
            "timestamp": timestamp
        }
 
        # Si estás en modo local apuntaremos a la API en Rust más adelante.
        self.client.post("/grpc-202307820", json=payload)