# Device-specific Caveats

## Synway SMG4000 Series

1. `enquire_link` need to be invoked every 5 seconds.
2. Firmware versions before 09/25/2020, `deliver_sm` - `dest_addr` field returns garbage.
3. The use of `bind_receiver` and `bind_transmitter` is not supported.
4. Only **SMPP v3.4** is implemented.

## DBLTek GoIP Series

1. System ID pattern
<br>e.q: if set `goip` then every slot is `goip01` ... `goip48`
