FROM python:3.8-slim

WORKDIR /consumer
ADD . /consumer
RUN pip install --trusted-host pypi.python.org -r requirements.txt

CMD ["python", "sensor_consumer.py", "-k", "kafka", "-p", "29092"]
