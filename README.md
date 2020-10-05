# go-smpp

A complete implementation of SMPP v5 protocol, written in golang.

## Key features

- Message encoding auto-detection
- Multipart SMS automatic splitting and concatenating

Supported encodings:

- GSM 7Bit
- ASCII
- Latin-1
- Cyrillic
- Hebrew
- Shift-JIS
- ISO-2022-JP
- EUC-JP
- EUC-KR
- UCS-2

## Caveats

- Please read [the SMPP Specification Version 5](docs/SMPP_v5.pdf) first. [pdu](pdu) is not limited to any value range.
- If you do not like the default [conn.go](conn.go) implementation, you can easily replace it.

## Example

1. Connect to the Remote (SMPP server)

    ```go
    parent, err := net.Dial("tcp", "m2m-device:2775")
    if err != nil {
        panic(err)
    }
    conn = smpp.NewConn(context.Background(), parent)
    conn.WriteTimeout = n * time.Second // set write timeout (optional, default 15 minutes)
    conn.ReadTimeout =  n * time.Second // set read timeout  (optional, default 15 minutes)
    go conn.Watch()                     // start watchdog
    ```

2. Handshake

    ```go
    resp, err := conn.Submit(ctx, &pdu.BindTransmitter{
        SystemID:   "your system id",
        Password:   "your password",
        SystemType: "your system type",
        Version:    pdu.SMPPVersion50,
    })
    if err != nil {
        panic(err)
    }
    r := resp.(*pdu.BindTransmitterResp)
    if r.Header.CommandStatus == 0 {
        // start keep-alive
        go conn.EnquireLink(time.Minute, time.Minute)
    }
    ```

3. Send Message

    ```go
    packet := &pdu.SubmitSM{
        SourceAddr: pdu.Address{TON: 1, NPI: 1, No: "00919821"},
        DestAddr:   pdu.Address{TON: 1, NPI: 1, No: "99919821"},
    }
    err := packet.Message.Compose("Hello World!")
    if err != nil {
        panic(err)
    }
    resp, err := conn.Submit(context.Background(), packet)
    if err != nil {
        panic(err)
    }
    resp // submit_sm_resp returns
    ```

4. Event loop for receiving messages

    ```go
    for {
        packet := <-conn.PDU()
        // reply a responsable packet
        if p, ok := packet.(pdu.Responsable); ok {
            err := conn.Send(context.Background(), p.Resp())
            if err != nil {
                fmt.Println(err)
            }
        }
    }
    ```

## LICENSE

This piece of software is released under [the MIT license](LICENSE).
