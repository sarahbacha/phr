'use strict';

const router = require('express').Router();
module.exports = function (app) {
  var controller = require('./account.route.controller.js')(app);
  router.post('/', controller.createAccount);
  router.post('/login', controller.logintoAccount);
  return router;
};