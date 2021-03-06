'use strict';

var Fabric_Client = require('fabric-client');
var path = require('path');
var util = require('util');
var os = require('os');
const uuidv1 = require('uuid/v1');

module.exports = function (app) {
  return {
    getPatientbyID: getPatientbyID,
    getPatientGeneralInformation: getPatientGeneralInformation,
    getPatientDiabetesIndications: getPatientDiabetesIndications,
    createPatientGeneralInformation: createPatientGeneralInformation,
    createPatientDiabetesIndications: createPatientDiabetesIndications,
    updatePatientGeneralInformation: updatePatientGeneralInformation,
    updatePatientDiabetesIndications: updatePatientDiabetesIndications,
    getPatientGeneralInformationHistory: getPatientGeneralInformationHistory,
    getPatientDiabetesIndicationsHistory: getPatientDiabetesIndicationsHistory,
    getPatientAccess: getPatientAccess,
    givePatientAccess: givePatientAccess
  };

  function getPatientbyID(req, res) {
    var fabric_client = new Fabric_Client();
    var key = req.params.id

    console.log('getPatientbyID ', key);
    // setup the fabric network
    var channel = fabric_client.newChannel('mychannel');
    var peer = fabric_client.newPeer('grpc://localhost:7051');
    channel.addPeer(peer);

    //
    var member_user = null;
    var store_path = path.join(os.homedir(), '.hfc-key-store');
    console.log('Store path:' + store_path);
    var tx_id = null;

    // create the key value store as defined in the fabric-client/config/default.json 'key-value-store' setting
    Fabric_Client.newDefaultKeyValueStore({
      path: store_path
    }).then((state_store) => {
      // assign the store to the fabric client
      fabric_client.setStateStore(state_store);
      var crypto_suite = Fabric_Client.newCryptoSuite();
      // use the same location for the state store (where the users' certificate are kept)
      // and the crypto store (where the users' keys are kept)
      var crypto_store = Fabric_Client.newCryptoKeyStore({
        path: store_path
      });
      crypto_suite.setCryptoKeyStore(crypto_store);
      fabric_client.setCryptoSuite(crypto_suite);

      // get the enrolled user from persistence, this user will sign all requests
      return fabric_client.getUserContext('user1', true);
    }).then((user_from_store) => {
      if (user_from_store && user_from_store.isEnrolled()) {
        console.log('Successfully loaded user from persistence');
        member_user = user_from_store;
      } else {
        throw new Error('Failed to get user.... run registerUser.js');
      }

      // queryPatient - requires 1 argument, ex: args: ['4'],
      const request = {
        chaincodeId: 'phr-app',
        txId: tx_id,
        fcn: 'queryPatient',
        args: [key]
      };

      // send the query proposal to the peer
      return channel.queryByChaincode(request);
    }).then((query_responses) => {
      console.log("Query has completed, checking results");
      // query_responses could have more than one  results if there multiple peers were used as targets
      if (query_responses && query_responses.length == 1) {
        if (query_responses[0] instanceof Error) {
          console.error("error from query = ", query_responses[0]);
          res.send("Could not locate patient")

        } else {
          console.log("query_responses is ", query_responses.toString());
          console.log("Response is ", query_responses[0].toString());
          res.send(query_responses[0].toString())
        }
      } else {
        console.log("No payloads were returned from query");
        res.send("Could not locate patient")
      }
    }).catch((err) => {
      console.error('Failed to query successfully :: ' + err);
      res.send("Could not locate patient")
    });
  }

  function getPatientGeneralInformation(req, res) {
    var fabric_client = new Fabric_Client();
    var key = req.params.id

    console.log('PatientGeneralInformationID', key);
    // setup the fabric network
    var channel = fabric_client.newChannel('mychannel');
    var peer = fabric_client.newPeer('grpc://localhost:7051');
    channel.addPeer(peer);

    //
    var member_user = null;
    var store_path = path.join(os.homedir(), '.hfc-key-store');
    console.log('Store path:' + store_path);
    var tx_id = null;

    // create the key value store as defined in the fabric-client/config/default.json 'key-value-store' setting
    Fabric_Client.newDefaultKeyValueStore({
      path: store_path
    }).then((state_store) => {
      // assign the store to the fabric client
      fabric_client.setStateStore(state_store);
      var crypto_suite = Fabric_Client.newCryptoSuite();
      // use the same location for the state store (where the users' certificate are kept)
      // and the crypto store (where the users' keys are kept)
      var crypto_store = Fabric_Client.newCryptoKeyStore({
        path: store_path
      });
      crypto_suite.setCryptoKeyStore(crypto_store);
      fabric_client.setCryptoSuite(crypto_suite);

      // get the enrolled user from persistence, this user will sign all requests
      return fabric_client.getUserContext('user1', true);
    }).then((user_from_store) => {
      if (user_from_store && user_from_store.isEnrolled()) {
        console.log('Successfully loaded user from persistence');
        member_user = user_from_store;
      } else {
        throw new Error('Failed to get user.... run registerUser.js');
      }

      // queryPatient - requires 1 argument, ex: args: ['4'],
      const request = {
        chaincodeId: 'phr-app',
        txId: tx_id,
        fcn: 'queryPatientGeneralInformation',
        args: [key]
      };

      // send the query proposal to the peer
      return channel.queryByChaincode(request);
    }).then((query_responses) => {
      console.log("Query has completed, checking results");
      // query_responses could have more than one  results if there multiple peers were used as targets
      if (query_responses && query_responses.length == 1) {
        if (query_responses[0] instanceof Error) {
          console.error("error from query = ", query_responses[0]);
          res.send("Could not locate patient")

        } else {
          console.log("Response is ", query_responses[0].toString());
          res.send(query_responses[0].toString())
        }
      } else {
        console.log("No payloads were returned from query");
        res.send("Could not locate patient")
      }
    }).catch((err) => {
      console.error('Failed to query successfully :: ' + err);
      res.send("Could not locate patient")
    });
  }

  function getPatientDiabetesIndications(req, res) {
    var fabric_client = new Fabric_Client();
    var key = req.params.id

    console.log('PatientDiabetesIndicationsID', key);
    // setup the fabric network
    var channel = fabric_client.newChannel('mychannel');
    var peer = fabric_client.newPeer('grpc://localhost:7051');
    channel.addPeer(peer);

    //
    var member_user = null;
    var store_path = path.join(os.homedir(), '.hfc-key-store');
    console.log('Store path:' + store_path);
    var tx_id = null;

    // create the key value store as defined in the fabric-client/config/default.json 'key-value-store' setting
    Fabric_Client.newDefaultKeyValueStore({
      path: store_path
    }).then((state_store) => {
      // assign the store to the fabric client
      fabric_client.setStateStore(state_store);
      var crypto_suite = Fabric_Client.newCryptoSuite();
      // use the same location for the state store (where the users' certificate are kept)
      // and the crypto store (where the users' keys are kept)
      var crypto_store = Fabric_Client.newCryptoKeyStore({
        path: store_path
      });
      crypto_suite.setCryptoKeyStore(crypto_store);
      fabric_client.setCryptoSuite(crypto_suite);

      // get the enrolled user from persistence, this user will sign all requests
      return fabric_client.getUserContext('user1', true);
    }).then((user_from_store) => {
      if (user_from_store && user_from_store.isEnrolled()) {
        console.log('Successfully loaded user from persistence');
        member_user = user_from_store;
      } else {
        throw new Error('Failed to get user.... run registerUser.js');
      }

      // queryPatient - requires 1 argument, ex: args: ['4'],
      const request = {
        chaincodeId: 'phr-app',
        txId: tx_id,
        fcn: 'queryPatientDiabetesIndications',
        args: [key]
      };

      // send the query proposal to the peer
      return channel.queryByChaincode(request);
    }).then((query_responses) => {
      console.log("Query has completed, checking results");
      // query_responses could have more than one  results if there multiple peers were used as targets
      if (query_responses && query_responses.length == 1) {
        if (query_responses[0] instanceof Error) {
          console.error("error from query = ", query_responses[0]);
          res.send("Could not locate patient")

        } else {
          console.log("Response is ", query_responses[0].toString());
          res.send(query_responses[0].toString())
        }
      } else {
        console.log("No payloads were returned from query");
        res.send("Could not locate patient")
      }
    }).catch((err) => {
      console.error('Failed to query successfully :: ' + err);
      res.send("Could not locate patient")
    });
  }

  function createPatientGeneralInformation(req, res) {
    console.log("submit recording of a patient: ");
    var patient = req.body;
    console.log(patient);

    var id = uuidv1();
    var accountid = patient.accountid;
    var reporterid = patient.reporterid;
    var dob = patient.dob;
    var sex = patient.sex;
    var race = patient.race;
    console.log(id, accountid, dob, sex, race);

    var fabric_client = new Fabric_Client();

    // setup the fabric network
    var channel = fabric_client.newChannel('mychannel');
    var peer = fabric_client.newPeer('grpc://localhost:7051');
    channel.addPeer(peer);
    var order = fabric_client.newOrderer('grpc://localhost:7050')
    channel.addOrderer(order);

    var member_user = null;
    var store_path = path.join(os.homedir(), '.hfc-key-store');
    console.log('Store path:' + store_path);
    var tx_id = null;

    // create the key value store as defined in the fabric-client/config/default.json 'key-value-store' setting
    Fabric_Client.newDefaultKeyValueStore({
      path: store_path
    }).then((state_store) => {
      // assign the store to the fabric client
      fabric_client.setStateStore(state_store);
      var crypto_suite = Fabric_Client.newCryptoSuite();
      // use the same location for the state store (where the users' certificate are kept)
      // and the crypto store (where the users' keys are kept)
      var crypto_store = Fabric_Client.newCryptoKeyStore({
        path: store_path
      });
      crypto_suite.setCryptoKeyStore(crypto_store);
      fabric_client.setCryptoSuite(crypto_suite);

      // get the enrolled user from persistence, this user will sign all requests
      return fabric_client.getUserContext('user1', true);
    }).then((user_from_store) => {
      if (user_from_store && user_from_store.isEnrolled()) {
        console.log('Successfully loaded user from persistence');
        member_user = user_from_store;
      } else {
        throw new Error('Failed to get user.... run registerUser.js');
      }

      // get a transaction id object based on the current user assigned to fabric client
      tx_id = fabric_client.newTransactionID();
      console.log("Assigning transaction_id: ", tx_id._transaction_id);

      // recordPatient - requires 5 args, ID, gender, age, perc, date , ex: args: ['1', 'female', '24', '9', '10-10-2018'], 
      // send proposal to endorser
      const request = {
        //targets : --- letting this default to the peers assigned to the channel
        chaincodeId: 'phr-app',
        fcn: 'recordPatientGeneralInformation',
        args: [id, accountid, reporterid, dob, sex, race],
        chainId: 'mychannel',
        txId: tx_id
      };
      // send the transaction proposal to the peers
      return channel.sendTransactionProposal(request);
    }).then((results) => {
      var proposalResponses = results[0];
      console.log('proposalResponses ', proposalResponses);
      console.log('proposal ', proposal);
      var proposal = results[1];
      let isProposalGood = false;
      if (proposalResponses && proposalResponses[0].response &&
        proposalResponses[0].response.status === 200) {
        isProposalGood = true;
        console.log('Transaction proposal was good');
      } else {
        console.error('Transaction proposal was bad');
      }
      if (isProposalGood) {
        console.log(util.format(
          'Successfully sent Proposal and received ProposalResponse: Status - %s, message - "%s"',
          proposalResponses[0].response.status, proposalResponses[0].response.message));

        // build up the request for the orderer to have the transaction committed
        var request = {
          proposalResponses: proposalResponses,
          proposal: proposal
        };

        // set the transaction listener and set a timeout of 30 sec
        // if the transaction did not get committed within the timeout period,
        // report a TIMEOUT status
        var transaction_id_string = tx_id.getTransactionID(); //Get the transaction ID string to be used by the event processing
        var promises = [];

        var sendPromise = channel.sendTransaction(request);
        promises.push(sendPromise); //we want the send transaction first, so that we know where to check status

        // get an eventhub once the fabric client has a user assigned. The user
        // is required because the event registration must be signed
        let event_hub = channel.newChannelEventHub('localhost:7051');

        // using resolve the promise so that result status may be processed
        // under the then clause rather than having the catch clause process
        // the status
        let txPromise = new Promise((resolve, reject) => {
          let handle = setTimeout(() => {
            event_hub.disconnect();
            resolve({
              event_status: 'TIMEOUT'
            }); //we could use reject(new Error('Trnasaction did not complete within 30 seconds'));
          }, 3000);
          event_hub.connect();
          event_hub.registerTxEvent(transaction_id_string, (tx, code) => {
            // this is the callback for transaction event status
            // first some clean up of event listener
            clearTimeout(handle);
            event_hub.unregisterTxEvent(transaction_id_string);
            event_hub.disconnect();

            // now let the application know what happened
            var return_status = {
              event_status: code,
              tx_id: transaction_id_string
            };
            if (code !== 'VALID') {
              console.error('The transaction was invalid, code = ' + code);
              resolve(return_status); // we could use reject(new Error('Problem with the tranaction, event status ::'+code));
            } else {
              console.log('The transaction has been committed on peer ' + event_hub.getPeerAddr());
              resolve(return_status);
            }
          }, (err) => {
            //this is the callback if something goes wrong with the event registration or processing
            reject(new Error('There was a problem with the eventhub ::' + err));
          });
        });
        promises.push(txPromise);

        return Promise.all(promises);
      } else {
        console.error('Failed to send Proposal or receive valid response. Response null or status is not 200. exiting...');
        throw new Error('Failed to send Proposal or receive valid response. Response null or status is not 200. exiting...');
      }
    }).then((results) => {
      console.log('Send transaction promise and event listener promise have completed');
      let res_msg = "";
      // check the results in the order the promises were added to the promise all list
      if (results && results[0] && results[0].status === 'SUCCESS') {
        console.log('Successfully sent transaction to the orderer.');
        // res.send(tx_id.getTransactionID());
      } else {
        console.error('Failed to order the transaction. Error code: ' + response.status);
        res_msg += 'Failed to order the transaction. Error code: ' + response.status;
      }

      if (results && results[1] && results[1].event_status === 'VALID') {
        console.log('Successfully committed the change to the ledger by the peer');
        // res.send(tx_id.getTransactionID());
        res_msg += (tx_id.getTransactionID());
      } else {
        console.log('Transaction failed to be committed to the ledger due to ::' + results[1].event_status);
        res_msg += 'Transaction failed to be committed to the ledger due to ::' + results[1].event_status;
      }
      res.send(res_msg);
    }).catch((err) => {
      console.error('Failed to invoke successfully :: ' + err);
      if (res.headerSent === false) {
        res.send("Error: adding new patient, check server logs.");
      }
    });
  }

  function createPatientDiabetesIndications(req, res) {
    console.log("submit recording of a patient: ");
    var patient = req.body;
    console.log(patient);

    var id = uuidv1();
    var accountid = patient.accountid;
    var reporterid = patient.reporterid;
    var HbA1c = patient.HbA1c;
    var HbA1cDate = patient.HbA1cDate;

    console.log(id, accountid, HbA1c, HbA1cDate);

    var fabric_client = new Fabric_Client();

    // setup the fabric network
    var channel = fabric_client.newChannel('mychannel');
    var peer = fabric_client.newPeer('grpc://localhost:7051');
    channel.addPeer(peer);
    var order = fabric_client.newOrderer('grpc://localhost:7050')
    channel.addOrderer(order);

    var member_user = null;
    var store_path = path.join(os.homedir(), '.hfc-key-store');
    console.log('Store path:' + store_path);
    var tx_id = null;

    // create the key value store as defined in the fabric-client/config/default.json 'key-value-store' setting
    Fabric_Client.newDefaultKeyValueStore({
      path: store_path
    }).then((state_store) => {
      // assign the store to the fabric client
      fabric_client.setStateStore(state_store);
      var crypto_suite = Fabric_Client.newCryptoSuite();
      // use the same location for the state store (where the users' certificate are kept)
      // and the crypto store (where the users' keys are kept)
      var crypto_store = Fabric_Client.newCryptoKeyStore({
        path: store_path
      });
      crypto_suite.setCryptoKeyStore(crypto_store);
      fabric_client.setCryptoSuite(crypto_suite);

      // get the enrolled user from persistence, this user will sign all requests
      return fabric_client.getUserContext('user1', true);
    }).then((user_from_store) => {
      if (user_from_store && user_from_store.isEnrolled()) {
        console.log('Successfully loaded user from persistence');
        member_user = user_from_store;
      } else {
        throw new Error('Failed to get user.... run registerUser.js');
      }

      // get a transaction id object based on the current user assigned to fabric client
      tx_id = fabric_client.newTransactionID();
      console.log("Assigning transaction_id: ", tx_id._transaction_id);

      // recordPatient - requires 5 args, ID, gender, age, perc, date , ex: args: ['1', 'female', '24', '9', '10-10-2018'], 
      // send proposal to endorser
      const request = {
        //targets : --- letting this default to the peers assigned to the channel
        chaincodeId: 'phr-app',
        fcn: 'recordPatientDiabetesIndications',
        args: [id, accountid, reporterid, HbA1c, HbA1cDate],
        chainId: 'mychannel',
        txId: tx_id
      };
      // send the transaction proposal to the peers
      return channel.sendTransactionProposal(request);
    }).then((results) => {
      var proposalResponses = results[0];
      console.log('proposalResponses ', proposalResponses);
      console.log('proposal ', proposal);
      var proposal = results[1];
      let isProposalGood = false;
      if (proposalResponses && proposalResponses[0].response &&
        proposalResponses[0].response.status === 200) {
        isProposalGood = true;
        console.log('Transaction proposal was good');
      } else {
        console.error('Transaction proposal was bad');
      }
      if (isProposalGood) {
        console.log(util.format(
          'Successfully sent Proposal and received ProposalResponse: Status - %s, message - "%s"',
          proposalResponses[0].response.status, proposalResponses[0].response.message));

        // build up the request for the orderer to have the transaction committed
        var request = {
          proposalResponses: proposalResponses,
          proposal: proposal
        };

        // set the transaction listener and set a timeout of 30 sec
        // if the transaction did not get committed within the timeout period,
        // report a TIMEOUT status
        var transaction_id_string = tx_id.getTransactionID(); //Get the transaction ID string to be used by the event processing
        var promises = [];

        var sendPromise = channel.sendTransaction(request);
        promises.push(sendPromise); //we want the send transaction first, so that we know where to check status

        // get an eventhub once the fabric client has a user assigned. The user
        // is required because the event registration must be signed
        let event_hub = channel.newChannelEventHub('localhost:7051');

        // using resolve the promise so that result status may be processed
        // under the then clause rather than having the catch clause process
        // the status
        let txPromise = new Promise((resolve, reject) => {
          let handle = setTimeout(() => {
            event_hub.disconnect();
            resolve({
              event_status: 'TIMEOUT'
            }); //we could use reject(new Error('Trnasaction did not complete within 30 seconds'));
          }, 3000);
          event_hub.connect();
          event_hub.registerTxEvent(transaction_id_string, (tx, code) => {
            // this is the callback for transaction event status
            // first some clean up of event listener
            clearTimeout(handle);
            event_hub.unregisterTxEvent(transaction_id_string);
            event_hub.disconnect();

            // now let the application know what happened
            var return_status = {
              event_status: code,
              tx_id: transaction_id_string
            };
            if (code !== 'VALID') {
              console.error('The transaction was invalid, code = ' + code);
              resolve(return_status); // we could use reject(new Error('Problem with the tranaction, event status ::'+code));
            } else {
              console.log('The transaction has been committed on peer ' + event_hub.getPeerAddr());
              resolve(return_status);
            }
          }, (err) => {
            //this is the callback if something goes wrong with the event registration or processing
            reject(new Error('There was a problem with the eventhub ::' + err));
          });
        });
        promises.push(txPromise);

        return Promise.all(promises);
      } else {
        console.error('Failed to send Proposal or receive valid response. Response null or status is not 200. exiting...');
        throw new Error('Failed to send Proposal or receive valid response. Response null or status is not 200. exiting...');
      }
    }).then((results) => {
      console.log('Send transaction promise and event listener promise have completed');
      let res_msg = "";
      // check the results in the order the promises were added to the promise all list
      if (results && results[0] && results[0].status === 'SUCCESS') {
        console.log('Successfully sent transaction to the orderer.');
        // res.send(tx_id.getTransactionID());
      } else {
        console.error('Failed to order the transaction. Error code: ' + response.status);
        res_msg += 'Failed to order the transaction. Error code: ' + response.status;
      }

      if (results && results[1] && results[1].event_status === 'VALID') {
        console.log('Successfully committed the change to the ledger by the peer');
        // res.send(tx_id.getTransactionID());
        res_msg += (tx_id.getTransactionID());
      } else {
        console.log('Transaction failed to be committed to the ledger due to ::' + results[1].event_status);
        res_msg += 'Transaction failed to be committed to the ledger due to ::' + results[1].event_status;
      }
      res.send(res_msg);
    }).catch((err) => {
      console.error('Failed to invoke successfully :: ' + err);
      if (res.headerSent === false) {
        res.send("Error: adding new patient, check server logs.");
      }
    });
  }

  function updatePatientGeneralInformation(req, res) {
    console.log("changing Patient general information: ");

    var id = req.params.id;
    var patient = req.body;

    var accountid = patient.accountid;
    var reporterid = patient.reporterid;
    var dob = patient.dob;
    var sex = patient.sex;
    var race = patient.race;

    var fabric_client = new Fabric_Client();

    // setup the fabric network
    var channel = fabric_client.newChannel('mychannel');
    var peer = fabric_client.newPeer('grpc://localhost:7051');
    channel.addPeer(peer);
    var order = fabric_client.newOrderer('grpc://localhost:7050')
    channel.addOrderer(order);

    var member_user = null;
    var store_path = path.join(os.homedir(), '.hfc-key-store');
    console.log('Store path:' + store_path);
    var tx_id = null;

    // create the key value store as defined in the fabric-client/config/default.json 'key-value-store' setting
    Fabric_Client.newDefaultKeyValueStore({
      path: store_path
    }).then((state_store) => {
      // assign the store to the fabric client
      fabric_client.setStateStore(state_store);
      var crypto_suite = Fabric_Client.newCryptoSuite();
      // use the same location for the state store (where the users' certificate are kept)
      // and the crypto store (where the users' keys are kept)
      var crypto_store = Fabric_Client.newCryptoKeyStore({
        path: store_path
      });
      crypto_suite.setCryptoKeyStore(crypto_store);
      fabric_client.setCryptoSuite(crypto_suite);

      // get the enrolled user from persistence, this user will sign all requests
      return fabric_client.getUserContext('user1', true);
    }).then((user_from_store) => {
      if (user_from_store && user_from_store.isEnrolled()) {
        console.log('Successfully loaded user from persistence');
        member_user = user_from_store;
      } else {
        throw new Error('Failed to get user.... run registerUser.js');
      }

      // get a transaction id object based on the current user assigned to fabric client
      tx_id = fabric_client.newTransactionID();
      console.log("Assigning transaction_id: ", tx_id._transaction_id);

      // changePatientDiabetesIndication - 3 arguments, patient id, new HbA1c Percentage and Date.  , ex: args: ['1', '9','10-10-2018'],
      // send proposal to endorser
      var request = {
        //targets : --- letting this default to the peers assigned to the channel
        chaincodeId: 'phr-app',
        fcn: 'editPatientGeneralInformation',
        args: [id, accountid, reporterid, dob, sex, race],
        chainId: 'mychannel',
        txId: tx_id
      };

      // send the transaction proposal to the peers
      return channel.sendTransactionProposal(request);
    }).then((results) => {
      var proposalResponses = results[0];
      var proposal = results[1];
      let isProposalGood = false;
      if (proposalResponses && proposalResponses[0].response &&
        proposalResponses[0].response.status === 200) {
        isProposalGood = true;
        console.log('Transaction proposal was good');
      } else {
        console.error('Transaction proposal was bad');
      }
      if (isProposalGood) {
        console.log(util.format(
          'Successfully sent Proposal and received ProposalResponse: Status - %s, message - "%s"',
          proposalResponses[0].response.status, proposalResponses[0].response.message));

        // build up the request for the orderer to have the transaction committed
        var request = {
          proposalResponses: proposalResponses,
          proposal: proposal
        };

        // set the transaction listener and set a timeout of 30 sec
        // if the transaction did not get committed within the timeout period,
        // report a TIMEOUT status
        var transaction_id_string = tx_id.getTransactionID(); //Get the transaction ID string to be used by the event processing
        var promises = [];

        var sendPromise = channel.sendTransaction(request);
        promises.push(sendPromise); //we want the send transaction first, so that we know where to check status

        // get an eventhub once the fabric client has a user assigned. The user
        // is required bacause the event registration must be signed
        let event_hub = channel.newChannelEventHub('localhost:7051');

        // using resolve the promise so that result status may be processed
        // under the then clause rather than having the catch clause process
        // the status
        let txPromise = new Promise((resolve, reject) => {
          let handle = setTimeout(() => {
            event_hub.disconnect();
            resolve({
              event_status: 'TIMEOUT'
            }); //we could use reject(new Error('Trnasaction did not complete within 30 seconds'));
          }, 3000);
          event_hub.connect();
          event_hub.registerTxEvent(transaction_id_string, (tx, code) => {
            // this is the callback for transaction event status
            // first some clean up of event listener
            clearTimeout(handle);
            event_hub.unregisterTxEvent(transaction_id_string);
            event_hub.disconnect();

            // now let the application know what happened
            var return_status = {
              event_status: code,
              tx_id: transaction_id_string
            };
            if (code !== 'VALID') {
              console.error('The transaction was invalid, code = ' + code);
              resolve(return_status); // we could use reject(new Error('Problem with the tranaction, event status ::'+code));
            } else {
              console.log('The transaction has been committed on peer ' + event_hub.getPeerAddr());
              resolve(return_status);
            }
          }, (err) => {
            //this is the callback if something goes wrong with the event registration or processing
            reject(new Error('There was a problem with the eventhub ::' + err));
          });
        });
        promises.push(txPromise);

        return Promise.all(promises);
      } else {
        console.error('Failed to send Proposal or receive valid response. Response null or status is not 200. exiting...');
        res.send("Error: no patient found");
        // throw new Error('Failed to send Proposal or receive valid response. Response null or status is not 200. exiting...');
      }
    }).then((results) => {
      console.log('Send transaction promise and event listener promise have completed');
      let res_msg = "";
      // check the results in the order the promises were added to the promise all list
      if (results && results[0] && results[0].status === 'SUCCESS') {
        console.log('Successfully sent transaction to the orderer.');
        // res.json(tx_id.getTransactionID())
      } else {
        console.error('Failed to order the transaction. Error code: ' + response.status);
        res_msg += 'Failed to order the transaction. Error code: ' + response.status;
      }

      if (results && results[1] && results[1].event_status === 'VALID') {
        console.log('Successfully committed the change to the ledger by the peer');
        // res.send(tx_id.getTransactionID());
        res_msg += (tx_id.getTransactionID());
      } else {
        console.log('Transaction failed to be committed to the ledger due to ::' + results[1].event_status);
        res_msg += 'Transaction failed to be committed to the ledger due to ::' + results[1].event_status;
      }
      res.send(res_msg);
    }).catch((err) => {
      console.error('Failed to invoke successfully :: ' + err);
      if (res.headerSent === false) {
        res.send("Error: no patient found");
      }
    });

  }

  function updatePatientDiabetesIndications(req, res) {
    console.log("changing Patient diabetes indication: ");

    var id = req.params.id;
    var patient = req.body;

    var accountid = patient.accountid;
    var reporterid = patient.reporterid;
    var HbA1c = patient.HbA1c;
    var HbA1cDate = patient.HbA1cDate;

    var fabric_client = new Fabric_Client();

    // setup the fabric network
    var channel = fabric_client.newChannel('mychannel');
    var peer = fabric_client.newPeer('grpc://localhost:7051');
    channel.addPeer(peer);
    var order = fabric_client.newOrderer('grpc://localhost:7050')
    channel.addOrderer(order);

    var member_user = null;
    var store_path = path.join(os.homedir(), '.hfc-key-store');
    console.log('Store path:' + store_path);
    var tx_id = null;

    // create the key value store as defined in the fabric-client/config/default.json 'key-value-store' setting
    Fabric_Client.newDefaultKeyValueStore({
      path: store_path
    }).then((state_store) => {
      // assign the store to the fabric client
      fabric_client.setStateStore(state_store);
      var crypto_suite = Fabric_Client.newCryptoSuite();
      // use the same location for the state store (where the users' certificate are kept)
      // and the crypto store (where the users' keys are kept)
      var crypto_store = Fabric_Client.newCryptoKeyStore({
        path: store_path
      });
      crypto_suite.setCryptoKeyStore(crypto_store);
      fabric_client.setCryptoSuite(crypto_suite);

      // get the enrolled user from persistence, this user will sign all requests
      return fabric_client.getUserContext('user1', true);
    }).then((user_from_store) => {
      if (user_from_store && user_from_store.isEnrolled()) {
        console.log('Successfully loaded user from persistence');
        member_user = user_from_store;
      } else {
        throw new Error('Failed to get user.... run registerUser.js');
      }

      // get a transaction id object based on the current user assigned to fabric client
      tx_id = fabric_client.newTransactionID();
      console.log("Assigning transaction_id: ", tx_id._transaction_id);

      // changePatientDiabetesIndication - 3 arguments, patient id, new HbA1c Percentage and Date.  , ex: args: ['1', '9','10-10-2018'],
      // send proposal to endorser
      var request = {
        //targets : --- letting this default to the peers assigned to the channel
        chaincodeId: 'phr-app',
        fcn: 'editPatientDiabetesIndications',
        args: [id, accountid, reporterid, HbA1c, HbA1cDate],
        chainId: 'mychannel',
        txId: tx_id
      };

      // send the transaction proposal to the peers
      return channel.sendTransactionProposal(request);
    }).then((results) => {
      var proposalResponses = results[0];
      var proposal = results[1];
      let isProposalGood = false;
      if (proposalResponses && proposalResponses[0].response &&
        proposalResponses[0].response.status === 200) {
        isProposalGood = true;
        console.log('Transaction proposal was good');
      } else {
        console.error('Transaction proposal was bad');
      }
      if (isProposalGood) {
        console.log(util.format(
          'Successfully sent Proposal and received ProposalResponse: Status - %s, message - "%s"',
          proposalResponses[0].response.status, proposalResponses[0].response.message));

        // build up the request for the orderer to have the transaction committed
        var request = {
          proposalResponses: proposalResponses,
          proposal: proposal
        };

        // set the transaction listener and set a timeout of 30 sec
        // if the transaction did not get committed within the timeout period,
        // report a TIMEOUT status
        var transaction_id_string = tx_id.getTransactionID(); //Get the transaction ID string to be used by the event processing
        var promises = [];

        var sendPromise = channel.sendTransaction(request);
        promises.push(sendPromise); //we want the send transaction first, so that we know where to check status

        // get an eventhub once the fabric client has a user assigned. The user
        // is required bacause the event registration must be signed
        let event_hub = channel.newChannelEventHub('localhost:7051');

        // using resolve the promise so that result status may be processed
        // under the then clause rather than having the catch clause process
        // the status
        let txPromise = new Promise((resolve, reject) => {
          let handle = setTimeout(() => {
            event_hub.disconnect();
            resolve({
              event_status: 'TIMEOUT'
            }); //we could use reject(new Error('Trnasaction did not complete within 30 seconds'));
          }, 3000);
          event_hub.connect();
          event_hub.registerTxEvent(transaction_id_string, (tx, code) => {
            // this is the callback for transaction event status
            // first some clean up of event listener
            clearTimeout(handle);
            event_hub.unregisterTxEvent(transaction_id_string);
            event_hub.disconnect();

            // now let the application know what happened
            var return_status = {
              event_status: code,
              tx_id: transaction_id_string
            };
            if (code !== 'VALID') {
              console.error('The transaction was invalid, code = ' + code);
              resolve(return_status); // we could use reject(new Error('Problem with the tranaction, event status ::'+code));
            } else {
              console.log('The transaction has been committed on peer ' + event_hub.getPeerAddr());
              resolve(return_status);
            }
          }, (err) => {
            //this is the callback if something goes wrong with the event registration or processing
            reject(new Error('There was a problem with the eventhub ::' + err));
          });
        });
        promises.push(txPromise);

        return Promise.all(promises);
      } else {
        console.error('Failed to send Proposal or receive valid response. Response null or status is not 200. exiting...');
        res.send("Error: no patient found");
        // throw new Error('Failed to send Proposal or receive valid response. Response null or status is not 200. exiting...');
      }
    }).then((results) => {
      console.log('Send transaction promise and event listener promise have completed');
      let res_msg = "";
      // check the results in the order the promises were added to the promise all list
      if (results && results[0] && results[0].status === 'SUCCESS') {
        console.log('Successfully sent transaction to the orderer.');
        // res.json(tx_id.getTransactionID())
      } else {
        console.error('Failed to order the transaction. Error code: ' + response.status);
        res_msg += 'Failed to order the transaction. Error code: ' + response.status;
      }

      if (results && results[1] && results[1].event_status === 'VALID') {
        console.log('Successfully committed the change to the ledger by the peer');
        // res.send(tx_id.getTransactionID());
        res_msg += (tx_id.getTransactionID());
      } else {
        console.log('Transaction failed to be committed to the ledger due to ::' + results[1].event_status);
        res_msg += 'Transaction failed to be committed to the ledger due to ::' + results[1].event_status;
      }
      res.send(res_msg);
    }).catch((err) => {
      console.error('Failed to invoke successfully :: ' + err);
      if (res.headerSent === false) {
        res.send("Error: no patient found");
      }
    });

  }

  function getPatientGeneralInformationHistory(req, res) {
    var fabric_client = new Fabric_Client();
    var key = req.params.id

    console.log('getPatientGeneralInformationHistory ', key);
    // setup the fabric network
    var channel = fabric_client.newChannel('mychannel');
    var peer = fabric_client.newPeer('grpc://localhost:7051');
    channel.addPeer(peer);

    //
    var member_user = null;
    var store_path = path.join(os.homedir(), '.hfc-key-store');
    console.log('Store path:' + store_path);
    var tx_id = null;

    // create the key value store as defined in the fabric-client/config/default.json 'key-value-store' setting
    Fabric_Client.newDefaultKeyValueStore({
      path: store_path
    }).then((state_store) => {
      // assign the store to the fabric client
      fabric_client.setStateStore(state_store);
      var crypto_suite = Fabric_Client.newCryptoSuite();
      // use the same location for the state store (where the users' certificate are kept)
      // and the crypto store (where the users' keys are kept)
      var crypto_store = Fabric_Client.newCryptoKeyStore({
        path: store_path
      });
      crypto_suite.setCryptoKeyStore(crypto_store);
      fabric_client.setCryptoSuite(crypto_suite);

      // get the enrolled user from persistence, this user will sign all requests
      return fabric_client.getUserContext('user1', true);
    }).then((user_from_store) => {
      if (user_from_store && user_from_store.isEnrolled()) {
        console.log('Successfully loaded  from persistence');
        member_user = user_from_store;
      } else {
        throw new Error('Failed to get user.... run registerUser.js');
      }

      // queryPatient - requires 1 argument, ex: args: ['4'],
      const request = {
        chaincodeId: 'phr-app',
        txId: tx_id,
        fcn: 'queryPatientGeneralInformationHistory',
        args: [key]
      };

      // send the query proposal to the peer
      return channel.queryByChaincode(request);
    }).then((query_responses) => {
      console.log("Query has completed, checking results");
      // query_responses could have more than one  results if there multiple peers were used as targets
      if (query_responses && query_responses.length == 1) {
        if (query_responses[0] instanceof Error) {
          console.error("error from query = ", query_responses[0]);
          res.send("Could not locate patient")

        } else {
          console.log("Response is ", query_responses[0].toString());
          res.send(query_responses[0].toString())
        }
      } else {
        console.log("No payloads were returned from query");
        res.send("Could not locate patient")
      }
    }).catch((err) => {
      console.error('Failed to query successfully :: ' + err);
      res.send("Could not locate patient")
    });
  }

  function getPatientDiabetesIndicationsHistory(req, res) {
    var fabric_client = new Fabric_Client();
    var key = req.params.id

    console.log('getPatientGeneralInformationHistory ', key);
    // setup the fabric network
    var channel = fabric_client.newChannel('mychannel');
    var peer = fabric_client.newPeer('grpc://localhost:7051');
    channel.addPeer(peer);

    //
    var member_user = null;
    var store_path = path.join(os.homedir(), '.hfc-key-store');
    console.log('Store path:' + store_path);
    var tx_id = null;

    // create the key value store as defined in the fabric-client/config/default.json 'key-value-store' setting
    Fabric_Client.newDefaultKeyValueStore({
      path: store_path
    }).then((state_store) => {
      // assign the store to the fabric client
      fabric_client.setStateStore(state_store);
      var crypto_suite = Fabric_Client.newCryptoSuite();
      // use the same location for the state store (where the users' certificate are kept)
      // and the crypto store (where the users' keys are kept)
      var crypto_store = Fabric_Client.newCryptoKeyStore({
        path: store_path
      });
      crypto_suite.setCryptoKeyStore(crypto_store);
      fabric_client.setCryptoSuite(crypto_suite);

      // get the enrolled user from persistence, this user will sign all requests
      return fabric_client.getUserContext('user1', true);
    }).then((user_from_store) => {
      if (user_from_store && user_from_store.isEnrolled()) {
        console.log('Successfully loaded  from persistence');
        member_user = user_from_store;
      } else {
        throw new Error('Failed to get user.... run registerUser.js');
      }

      // queryPatient - requires 1 argument, ex: args: ['4'],
      const request = {
        chaincodeId: 'phr-app',
        txId: tx_id,
        fcn: 'queryPatientDiabetesIndicationsHistory',
        args: [key]
      };

      // send the query proposal to the peer
      return channel.queryByChaincode(request);
    }).then((query_responses) => {
      console.log("Query has completed, checking results");
      // query_responses could have more than one  results if there multiple peers were used as targets
      if (query_responses && query_responses.length == 1) {
        if (query_responses[0] instanceof Error) {
          console.error("error from query = ", query_responses[0]);
          res.send("Could not locate patient")

        } else {
          console.log("Response is ", query_responses[0].toString());
          res.send(query_responses[0].toString())
        }
      } else {
        console.log("No payloads were returned from query");
        res.send("Could not locate patient")
      }
    }).catch((err) => {
      console.error('Failed to query successfully :: ' + err);
      res.send("Could not locate patient")
    });
  }

  function getPatientAccess(req, res) {
    var fabric_client = new Fabric_Client();
    var key = req.params.id

    console.log('getPatientAccess ', key);
    // setup the fabric network
    var channel = fabric_client.newChannel('mychannel');
    var peer = fabric_client.newPeer('grpc://localhost:7051');
    channel.addPeer(peer);

    //
    var member_user = null;
    var store_path = path.join(os.homedir(), '.hfc-key-store');
    console.log('Store path:' + store_path);
    var tx_id = null;

    // create the key value store as defined in the fabric-client/config/default.json 'key-value-store' setting
    Fabric_Client.newDefaultKeyValueStore({
      path: store_path
    }).then((state_store) => {
      // assign the store to the fabric client
      fabric_client.setStateStore(state_store);
      var crypto_suite = Fabric_Client.newCryptoSuite();
      // use the same location for the state store (where the users' certificate are kept)
      // and the crypto store (where the users' keys are kept)
      var crypto_store = Fabric_Client.newCryptoKeyStore({
        path: store_path
      });
      crypto_suite.setCryptoKeyStore(crypto_store);
      fabric_client.setCryptoSuite(crypto_suite);

      // get the enrolled user from persistence, this user will sign all requests
      return fabric_client.getUserContext('user1', true);
    }).then((user_from_store) => {
      if (user_from_store && user_from_store.isEnrolled()) {
        console.log('Successfully loaded  from persistence');
        member_user = user_from_store;
      } else {
        throw new Error('Failed to get user.... run registerUser.js');
      }

      // queryPatient - requires 1 argument, ex: args: ['4'],
      const request = {
        chaincodeId: 'phr-app',
        txId: tx_id,
        fcn: 'getPatientReportingAccess',
        args: [key]
      };

      // send the query proposal to the peer
      return channel.queryByChaincode(request);
    }).then((query_responses) => {
      console.log("Query has completed, checking results");
      // query_responses could have more than one  results if there multiple peers were used as targets
      if (query_responses && query_responses.length == 1) {
        if (query_responses[0] instanceof Error) {
          console.error("error from query = ", query_responses[0]);
          res.send("Could not locate patient")

        } else {
          console.log("Response is ", query_responses[0].toString());
          res.send(query_responses[0].toString())
        }
      } else {
        console.log("No payloads were returned from query");
        res.send("Could not locate patient")
      }
    }).catch((err) => {
      console.error('Failed to query successfully :: ' + err);
      res.send("Could not locate patient")
    });
  }

  function givePatientAccess(req, res) {
    console.log("give patient access ");

    var accessinfo = req.body;
    var accountid = accessinfo.accountid;
    var reporterid = accessinfo.reporterid;
    var status = accessinfo.status;

    var fabric_client = new Fabric_Client();

    // setup the fabric network
    var channel = fabric_client.newChannel('mychannel');
    var peer = fabric_client.newPeer('grpc://localhost:7051');
    channel.addPeer(peer);
    var order = fabric_client.newOrderer('grpc://localhost:7050')
    channel.addOrderer(order);

    var member_user = null;
    var store_path = path.join(os.homedir(), '.hfc-key-store');
    console.log('Store path:' + store_path);
    var tx_id = null;

    // create the key value store as defined in the fabric-client/config/default.json 'key-value-store' setting
    Fabric_Client.newDefaultKeyValueStore({
      path: store_path
    }).then((state_store) => {
      // assign the store to the fabric client
      fabric_client.setStateStore(state_store);
      var crypto_suite = Fabric_Client.newCryptoSuite();
      // use the same location for the state store (where the users' certificate are kept)
      // and the crypto store (where the users' keys are kept)
      var crypto_store = Fabric_Client.newCryptoKeyStore({
        path: store_path
      });
      crypto_suite.setCryptoKeyStore(crypto_store);
      fabric_client.setCryptoSuite(crypto_suite);

      // get the enrolled user from persistence, this user will sign all requests
      return fabric_client.getUserContext('user1', true);
    }).then((user_from_store) => {
      if (user_from_store && user_from_store.isEnrolled()) {
        console.log('Successfully loaded user from persistence');
        member_user = user_from_store;
      } else {
        throw new Error('Failed to get user.... run registerUser.js');
      }

      // get a transaction id object based on the current user assigned to fabric client
      tx_id = fabric_client.newTransactionID();
      console.log("Assigning transaction_id: ", tx_id._transaction_id);

      // recordPatient - requires 5 args, ID, gender, age, perc, date , ex: args: ['1', 'female', '24', '9', '10-10-2018'], 
      // send proposal to endorser
      const request = {
        //targets : --- letting this default to the peers assigned to the channel
        chaincodeId: 'phr-app',
        fcn: 'recordPatientReporting',
        args: [accountid, reporterid, status],
        chainId: 'mychannel',
        txId: tx_id
      };
      // send the transaction proposal to the peers
      return channel.sendTransactionProposal(request);
    }).then((results) => {
      var proposalResponses = results[0];
      console.log('proposalResponses ', proposalResponses);
      console.log('proposal ', proposal);
      var proposal = results[1];
      let isProposalGood = false;
      if (proposalResponses && proposalResponses[0].response &&
        proposalResponses[0].response.status === 200) {
        isProposalGood = true;
        console.log('Transaction proposal was good');
      } else {
        console.error('Transaction proposal was bad');
      }
      if (isProposalGood) {
        console.log(util.format(
          'Successfully sent Proposal and received ProposalResponse: Status - %s, message - "%s"',
          proposalResponses[0].response.status, proposalResponses[0].response.message));

        // build up the request for the orderer to have the transaction committed
        var request = {
          proposalResponses: proposalResponses,
          proposal: proposal
        };

        // set the transaction listener and set a timeout of 30 sec
        // if the transaction did not get committed within the timeout period,
        // report a TIMEOUT status
        var transaction_id_string = tx_id.getTransactionID(); //Get the transaction ID string to be used by the event processing
        var promises = [];

        var sendPromise = channel.sendTransaction(request);
        promises.push(sendPromise); //we want the send transaction first, so that we know where to check status

        // get an eventhub once the fabric client has a user assigned. The user
        // is required because the event registration must be signed
        let event_hub = channel.newChannelEventHub('localhost:7051');

        // using resolve the promise so that result status may be processed
        // under the then clause rather than having the catch clause process
        // the status
        let txPromise = new Promise((resolve, reject) => {
          let handle = setTimeout(() => {
            event_hub.disconnect();
            resolve({
              event_status: 'TIMEOUT'
            }); //we could use reject(new Error('Trnasaction did not complete within 30 seconds'));
          }, 3000);
          event_hub.connect();
          event_hub.registerTxEvent(transaction_id_string, (tx, code) => {
            // this is the callback for transaction event status
            // first some clean up of event listener
            clearTimeout(handle);
            event_hub.unregisterTxEvent(transaction_id_string);
            event_hub.disconnect();

            // now let the application know what happened
            var return_status = {
              event_status: code,
              tx_id: transaction_id_string
            };
            if (code !== 'VALID') {
              console.error('The transaction was invalid, code = ' + code);
              resolve(return_status); // we could use reject(new Error('Problem with the tranaction, event status ::'+code));
            } else {
              console.log('The transaction has been committed on peer ' + event_hub.getPeerAddr());
              resolve(return_status);
            }
          }, (err) => {
            //this is the callback if something goes wrong with the event registration or processing
            reject(new Error('There was a problem with the eventhub ::' + err));
          });
        });
        promises.push(txPromise);

        return Promise.all(promises);
      } else {
        console.error('Failed to send Proposal or receive valid response. Response null or status is not 200. exiting...');
        throw new Error('Failed to send Proposal or receive valid response. Response null or status is not 200. exiting...');
      }
    }).then((results) => {
      console.log('Send transaction promise and event listener promise have completed');
      let res_msg = "";
      // check the results in the order the promises were added to the promise all list
      if (results && results[0] && results[0].status === 'SUCCESS') {
        console.log('Successfully sent transaction to the orderer.');
        // res.send(tx_id.getTransactionID());
      } else {
        console.error('Failed to order the transaction. Error code: ' + response.status);
        res_msg += 'Failed to order the transaction. Error code: ' + response.status;
      }

      if (results && results[1] && results[1].event_status === 'VALID') {
        console.log('Successfully committed the change to the ledger by the peer');
        // res.send(tx_id.getTransactionID());
        res_msg += (tx_id.getTransactionID());
      } else {
        console.log('Transaction failed to be committed to the ledger due to ::' + results[1].event_status);
        res_msg += 'Transaction failed to be committed to the ledger due to ::' + results[1].event_status;
      }
      res.send(res_msg);
    }).catch((err) => {
      console.error('Failed to invoke successfully :: ' + err);
      if (res.headerSent === false) {
        res.send("Error: adding new patient, check server logs.");
      }
    });
  }
};