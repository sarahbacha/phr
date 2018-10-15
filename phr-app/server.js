// call the packages we need
var express = require('express');
var app = express();
var cors = require('cors')
const bodyParser = require('body-parser');

app.use(bodyParser.json());
app.use(cors());
app.use(bodyParser.urlencoded({
  extended: false
}));
// Enable CORS
app.use(function (req, res, next) {
  // Website you wish to allow to connect
  res.header('Access-Control-Allow-Origin', '*');

  // Request methods you wish to allow
  res.header('Access-Control-Allow-Methods', 'GET, POST, OPTIONS, PUT, PATCH, DELETE');

  // Request headers you wish to allow
  res.header('Access-Control-Allow-Headers', 'Content-Type');

  // Set to true if you need the website to include cookies in the requests sent
  // to the API (e.g. in case you use sessions)
  res.header('Access-Control-Allow-Credentials', 'true');
  next();
});
app.use('/api', require('./routes')(app));
app.use('/status', function (req, res, next) {
  res.json({
    'status': 'ok'
  });
});
app.use(function (err, req, res, next) {
  res.status(400).send('Error, this url does not exist.' + req.originalUrl);
});
var port = process.env.PORT || 8000;

// Start the server and listen on port 
app.listen(port, function () {
  console.log("Live on port: " + port);
});