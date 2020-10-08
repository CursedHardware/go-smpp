# SMPP REPL

## Command

### Client feature

- [x] `connect` to target smpp server
- [x] `disconnect`
- [x] `send-message`, send a `submit_sm` to remote
- [x] `send-ussd`, send a `submit_sm` with `ussd_service_op` to remote
- [x] `query`, send a `query_sm` or `query_broadcast_sm` to remote
- [ ] `load-test`, Load testing to server

### Server feature

- [ ] `start-service`, Start SMPP Server
- [ ] `stop-service`, Stop SMPP Server
- [ ] `send-message`, Send a `deliver_sm` to client-side
- [ ] `load-test`, Load testing to client
