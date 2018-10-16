'use strict';

const router = require('express').Router();
module.exports = function (app) {
  var controller = require('./account.route.controller.js')(app);
  // router.get('/', controller.getAccounts);
  // router.get('/:id', controller.getAccountbyID);
  // router.patch('/:id', controller.updateAccount);
  router.post('/', controller.createAccount);
  router.post('/login', controller.logintoAccount);
  return router;
};