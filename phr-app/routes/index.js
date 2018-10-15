const router = require('express').Router();
var LedgerRouter = require('./ledger/index');
module.exports = function (app) {
  router.use('/', LedgerRouter(app));
  return router;
};