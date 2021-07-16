# Device-specific Caveats

## Synway SMG4000 Series

- Only **SMPP v3.4** is implemented.

- `enquire_link` need to be invoked every 0.5 seconds.

- Firmware versions before 09/25/2020, `deliver_sm` - `dest_addr` field returns garbage.

- The use of `bind_receiver` and `bind_transmitter` is not supported.

## DBLTek GoIP Series

- Only **SMPP v3.4** is implemented.

- `enquire_link` need to be invoked every 0.5 seconds.

- System ID pattern\
  e.q: if set `goip` then every slot is `goip01` ... `goip48`,\
  but if use `goip` login, then use **all slots** sent sms.

- Multipart SMS,
  Only support `TLV 0424, 4.8.4.36 message_payload`

- Only supports some Command IDs:

  ```plaintext
  bind_receiver
  bind_transmitter
  bind_transceiver
  submit_sm
  deliver_sm
  enquire_link
  unbind
  ```
