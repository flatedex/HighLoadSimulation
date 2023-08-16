# Simulation for high-ish load and long-time operations

## Idea
The idea is to simulate long-time operation, like sms-notification, and high-load operation, like authorization.

You can not be sure that outer service will operate in certain time gap, service may work for several seconds.
For interval between request from client to perfoming operation in outer service and giving back response,
client is blocked because he can not get response due to lont-time operation in outer service.

Soluiton to this problem might be way to use web-sockets to inform client that service done his work.

## Tecnologies

This project uses **microsevice architecture** and **RabbitMQ** as message broker, **Rabbit** provides communication between services with it's internal queues and distribution.

Server provides communication between client and services and this supposed to be it's **only** role. 
