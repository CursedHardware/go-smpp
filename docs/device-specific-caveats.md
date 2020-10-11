# Device-specific Caveats

## Synway SMG4000 Series

- `enquire_link` need to be invoked every 5 seconds.
- Firmware versions before 09/25/2020, `deliver_sm` - `dest_addr` field returns garbage.
- The use of `bind_receiver` and `bind_transmitter` is not supported.
- Only **SMPP v3.4** is implemented.

## DBLTek GoIP Series

- System ID pattern<br/>e.q: if set `goip` then every slot is `goip01` ... `goip48`
- Only supports some Command IDs:<br>`bind_receiver`, `bind_transmitter`, `bind_transceiver`<br>`submit_sm`, `deliver_sm`<br>`enquire_link`, `unbind`
- Only **SMPP v3.4** is implemented.
