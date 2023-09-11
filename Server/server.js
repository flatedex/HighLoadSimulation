'use strict'


const express = require('express');
const bcrypt = require('bcrypt');
const path = require('path');
const bodyParser = require('body-parser');
const router = express.Router();


const PORT = 8050;
const HOST = "0.0.0.0";
const saltRounds = 10;

const app = express();

var amqp = require('amqplib');

app.set('view engine', 'ejs');
app.set('views', './views');
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: false }));

app.get('/', (req, res) => {
    res.render("index");
});
app.post('/', async (req, res) => {
    let name = req.body.name;
    let password = req.body.password;
    let phone = req.body.phone;
    let user = {};

    bcrypt.genSalt(saltRounds, function(error, salt){
        bcrypt.hash(password, salt, function(error, hash){
            user = {
                name: name,
                password: hash,
                phone: phone,
            };
        });
    });

    let connection = await amqp.connect('amqp://guest:guest@rabbitmq:5672/'); 
    
    let channel = await connection.createChannel();
         
    let queue = "codeGenerator";

    let msg = JSON.stringify(user);

    await channel.assertQueue(queue, {
        durable: true
    });

    channel.sendToQueue(queue, Buffer.from(msg));
    console.log(`Message by ${name} sent`);
    res.render("index");
});

app.post('/confirm', async (req, res) => {
    let name = req.body.code;

});

app.listen(PORT, HOST, () => {
    console.log(`Server running on http://${HOST}:${PORT}`);
});