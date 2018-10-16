const router = require('express').Router();
var PatientRouter = require('./patient/index');
var AccountRouter = require('./account/index');
module.exports = function (app) {
  router.use('/patient', PatientRouter(app));
  router.use('/account', AccountRouter(app));
  return router;
};