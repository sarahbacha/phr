'use strict';

const router = require('express').Router();
module.exports = function (app) {
  var controller = require('./patient.route.controller.js')(app);
  router.get('/', controller.getPatients);
  router.get('/:id', controller.getPatientbyID);
  router.patch('/:id', controller.updateIndication);
  router.post('/', controller.createPatient);
  return router;
};