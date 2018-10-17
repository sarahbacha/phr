'use strict';

const router = require('express').Router();
module.exports = function (app) {
  var controller = require('./patient.route.controller.js')(app);
  router.get('/', controller.getPatients);
  router.get('/:id', controller.getPatientbyID);
  router.get('/history/:id', controller.getPatientHistorybyID);
  router.patch('/:id', controller.updatePatient);
  router.post('/', controller.createPatient);
  return router;
};