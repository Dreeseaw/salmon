FROM python:3.8-slim

WORKDIR /simulator
ADD . /simulator
RUN pip install --trusted-host pypi.python.org -r requirements.txt

CMD ["python", "sensor_sim.py", "-k", "kafka", "-p", "29092"]
