#!/usr/bin/env python
import argparse
import json
from datetime import datetime, timedelta
from confluent_kafka import Consumer
from math import sqrt
import sys

import redis
import psycopg2
from process_record import process

# command line arguments
parser = argparse.ArgumentParser(formatter_class=argparse.ArgumentDefaultsHelpFormatter)
parser.add_argument('-k', '--kafka_host', default='localhost', help='The host running the kafka server.')
parser.add_argument('-p', '--kafka_port', default='9092', help='The port accepting connections on the kafka server.')
parser.add_argument('-r', '--redis_host', default='redis', help='The Redis DB hostname.')
parser.add_argument('-o', '--redis_port', default='6379', help='The Redis DB port')
args = parser.parse_ars()

# if connections fail, consumer will restart
conf = {'bootstrap.servers': '{}:{}'.format(args.kafka_host, args.kafka_port), 'group.id': 'dreese', 'enable.auto.commit': False}
consumer = Consumer(conf)

redisClient = redis.Redis(host=args.redis_host, port=args.redis_port)
pgConn = psycopg2.connect(
    host="postgres",
    database="dreese",
    port="5432",
    user="postgres",
    password="postgres",
)

print('Starting the consumer. Press CTRL-C at any time to exit.')
running = True

def basic_consume_loop(consumer, topics):
    try:
        consumer.subscribe(topics)

        while running:
            msg = consumer.poll(timeout=1.0)
            consumer.commit()
            if msg is None: continue

            if msg.error():
                if msg.error().code() == KafkaError._PARTITION_EOF:
                    # End of partition event
                    sys.stderr.write('%% %s [%d] reached end at offset %d\n' %
                                     (msg.topic(), msg.partition(), msg.offset()))
                elif msg.error():
                    raise KafkaException(msg.error())
            else:
                # extract and process record data
                data = json.loads(msg.value())
                tid = str(data["traveler_id"])
                sid = int(data["sensor_id"])
                px  = float(data["position"][0])
                py  = float(data["position"][1])
                typ = str(data["traveler_type"])
                tmp = datetime.fromisoformat(data["timestamp"])

                process(redisClient, pgConn, sid, tid, typ, tmp, px, py)
    finally:
        # Close down consumer to commit final offsets.
        consumer.close()

def shutdown():
    running = False

if __name__=="__main__":
    basic_consume_loop(consumer, ["movements"])g
