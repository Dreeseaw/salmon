#!/usr/bin/env python
import argparse
import json
from datetime import datetime, timedelta
from random import randint, random, choice
from numpy.random import poisson
from confluent_kafka import Producer
from math import sqrt
from time import sleep
import sys
from uuid import uuid4

# sensor class
class sensor:

    def __init__(self,id,args):

        self.id = id
        self.traveler_speeds = args.traveler_speeds
        self.traveler_types = list(self.traveler_speeds.keys())
        self.travelers = []
        self.traveler_total = 0
        self.last_arrival = datetime.now()
        self.arrival_rate = args.arrival_rate

    # poisson arrivals
    def process_arrivals(self):
        t = datetime.now()
        interval = timedelta.total_seconds(t-elf.last_arrival)
        num_arrivals = poisson(self.arrival_rate*interval)
        arrival_counts = {t: 0 for t in self.traveler_types}
        for i in range(num_arrivals):

            # make a new traveler with random type
            traveler_type = choice(self.traveler_types)
            self.traveler_total += 1
            arrival_counts[traveler_type] += 1
            t = traveler(str(uuid4()),traveler_type,self.traveler_speeds[traveler_type])
            self.travelers.append(t)
            self.last_arrival = t.timestamp

        return arrival_counts

    # move all existing travelers forward
    def process_movements(self):

        for t in self.travelers:
            t.move()

        # remove travelers that exited the field of view
        self.travelers = [x for x in self.travelers if (x.position[0]<=100 and x.position[1]<=100)]

    # format movement data and produce to Kafka
    def send(self,producer):

        for t in self.travelers:
            data = t.to_dict()
            data['sensor_id'] = self.id
            producer.produce('movements', json.dumps(data))

# traveler class
class traveler:

    def __init__(self,id,type,average_speed):

        self.id = id
        self.type = type
        self.timestamp = datetime.now()
        self.average_speed = average_speed
        self.current_speed = 0

        # keep it simple for now.
        # Travelers always start on the left side of the coordinate grid
        # and move towards the right, never going backwards, but
        # possibly going up or down.
        self.position = [random(),100*random()]

    def move(self):

        #how much time has passed?
        datetime_new = datetime.now()
        dt = timedelta.total_seconds(datetime_new-self.timestamp)

        # go a random positive distance at given average speed
        d = 2*self.average_speed*dt*random()
        dx = d*random()
        dy = sqrt(d**2 - dx**2)

        #update position, time and speed
        self.position[0] += dx
        self.position[1] += dy
        self.timestamp = datetime_new
        self.current_speed = d/dt

    def to_dict(self):

        return {'traveler_id': self.id, 'traveler_type': self.type, 'position': self.position, 'timestamp': self.timestamp.isoformat()}


# command line arguments
parser = argparse.ArgumentParser(formatter_class=argparse.ArgumentDefaultsHelpFormatter)
parser.add_argument('-k', '--kafka_host', default='localhost', help='The host running the kafka server.')
parser.add_argument('-p', '--kafka_port', default='9092', help='The port accepting connections on the kafka server.')
parser.add_argument('-s', '--sensors', default=100, type=int, help='The number of sensors to simulate.')
parser.add_argument('-n', '--arrival_rate', default=10, type=float, help='The average number of travelers arriving at each sensor per second.')
args = parser.parse_args()


args.traveler_speeds = {'pedestrian': 5, 'bicyclist': 20, 'vehicle': 40}
traveler_types = args.traveler_speeds

conf = {'bootstrap.servers': '{}:{}'.format(args.kafka_host, args.kafka_port),
             'queue.buffering.max.messages': 100000,
             'queue.buffering.max.ms' : 10,
             'batch.num.messages': 1000}

producer = Producer(conf)

print('Starting the simulator. Press CTRL-C at any time to exit.')

# create sensors
print('Creating {0} sensors...'.format(args.sensors))
sensors = [sensor(id,args) for id in range(args.sensors)]

# main loop
print('Simulating {0} travelers per second (average) on each sensor...'.format(args.arrival_rate))
loops = 0
try:
    last_time = datetime.now()
    while True:
        arrivals = {t: 0 for t in traveler_types}

        for s in sensors:

            # move everyone
            s.process_movements()

            # new arrivals
            new_arrivals = s.process_arrivals()
            arrivals = {t: arrivals[t] + new_arrivals[t] for t in traveler_types}

            try:
                s.send(producer)
                producer.poll(0)
            except BufferError as e:
                print('\nLocal producer queue is full ({0} messages awaiting delivery): try again\n'.format(len(producer)))
                sys.stdout.flush()
                producer.poll(10) # wait for more space on the queue



        #total the counts on this loop:
        loops += 1
        elapsed = timedelta.total_seconds(datetime.now()-last_time)
        print('iteration {0} took {1} seconds. New travelers detected: {2}     \r'.format(loops,elapsed,arrivals))
        sys.stdout.flush()
        last_time = datetime.now()
        sleep(0.1)

except KeyboardInterrupt:
    producer.flush()
    print('done!')s
