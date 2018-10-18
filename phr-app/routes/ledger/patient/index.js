'use strict';

const router = require('express').Router();
module.exports = function (app) {
  var controller = require('./patient.route.controller.js')(app);
  router.get('/:id', controller.getPatientbyID);

  router.get('/general/:id', controller.getPatientGeneralInformation);
  router.post('/general', controller.createPatientGeneralInformation);
  router.patch('/general/:id', controller.updatePatientGeneralInformation);
  router.get('/generalhistory/:id', controller.getPatientGeneralInformationHistory);
  
  router.get('/diabetes/:id', controller.getPatientDiabetesIndications);
  router.post('/diabetes', controller.createPatientDiabetesIndications);
  router.patch('/diabetes/:id', controller.updatePatientDiabetesIndications);
  router.get('/diabeteshistory/:id', controller.getPatientDiabetesIndicationsHistory);
  return router;
};