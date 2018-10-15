const router = require('express').Router();
var PatientRouter = require('./patient/index');
module.exports = function (app) {
  router.use('/patient', PatientRouter(app));
  return router;
};