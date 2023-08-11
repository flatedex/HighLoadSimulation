'use strict'

const express = require('express');
const path = require('path');
const bodyParser = require('body-parser');
const router = express.Router();

const PORT = 8050;
const HOST = "127.0.0.1";

const app = express();

var amqp = require('amqplib/callback_api');

app.set('view engine', 'ejs');
app.set('views', './views');
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: false }));

app.get('/', (req, res) => {
    res.render("index");
});
app.post('/', (req, res) => {
    let name = req.body.Name;
    let message = req.body.Message;

    amqp.connect('amqp://guest:guest@localhost:5672/', function(error, connection) {
        if(error) {
            throw error;
        }
        connection.createChannel(function(error1, channel){
            if(error1){
                throw error;
            }
            let queue = "myQueue";
            let msg = message;

            channel.assertQueue(queue, {
                durable: false
            });

            channel.sendToQueue(queue, Buffer.from(msg));
            console.log(`Message by ${name} sent`);
        });
    });
});

app.listen(PORT, HOST, () => {
    console.log(`Server running on http://${HOST}:${PORT}`);
});