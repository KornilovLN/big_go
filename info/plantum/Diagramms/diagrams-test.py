#!/usr/bin/env python3
# -*- coding: utf-8 -*-

#--- Пришлось установить библиотеку plantuml, Diagrams и Graphviz
#--- чтобы использовать в этом коде
#--- pip install plantuml
#--- pip install diagrams
#--- sudo apt-get install graphviz

from diagrams import Diagram, Cluster
from diagrams.programming.language import Go
from diagrams.onprem.queue import RabbitMQ
from diagrams.onprem.database import PostgreSQL
from diagrams.onprem.inmemory import Redis
from diagrams.onprem.client import User

with Diagram("Big Go Architecture", show=True):
    # Добавляем пользователя
    user = User("User")

    with Cluster("Services"):
        generator = Go("Generator")
        collector = Go("Collector")
        user1 = Go("User1")
        user2 = Go("User2")
    
    with Cluster("Infrastructure"):
        rabbitmq = RabbitMQ("RabbitMQ")
        postgres = PostgreSQL("PostgreSQL")
        redis = Redis("Redis")
    
    generator >> rabbitmq >> collector
    collector >> user1
    collector >> user2
    collector >> postgres
    collector >> redis

    # Добавленные отношения
    user >> user1
    user >> user2
    user1 >> postgres
    user2 >> postgres
    rabbitmq >> redis
    redis >> postgres    
