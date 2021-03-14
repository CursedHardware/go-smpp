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

- Please read [the SMPP Specification Version 5](docs/SMPP_v5.pdf) first.
  [pdu](pdu) is not limited to any value range.
- If you do not like the default [session.go](session.go) implementation, you can easily replace it.
- [Device-specific Caveats](docs/device-specific-caveats.md)

## Command line tools

1. [smpp-receiver](cmd/smpp-receiver)
   <br>SMPP Simple Receiver tool

2. [smpp-repl](cmd/smpp-repl)
   <br>SMPP Simple Test tool

## Example

1. Connect to the Remote (SMPP server)

   ```go
   session, err = smpp.Dial(context.Background(), "m2m-device:2775")
   if err != nil {
       panic(err)
   }
   session.WriteTimeout = n * time.Second // set write timeout (optional, default 15 minutes)
   session.ReadTimeout =  n * time.Second // set read timeout  (optional, default 15 minutes)
   ```

2. Handshake

   ```go
   resp, err := session.Submit(context.Background(), &pdu.BindTransceiver{
       SystemID:   "your system id",
       Password:   "your password",
       SystemType: "your system type",
       Version:    pdu.SMPPVersion50,
   })
   if err != nil {
       panic(err)
   }
   r := resp.(*pdu.BindTransceiverResp)
   if r.Header.CommandStatus == 0 {
       // start keep-alive
       go session.EnquireLink(time.Minute, time.Minute)
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
   resp, err := session.Submit(context.Background(), packet)
   if err != nil {
       panic(err)
   }
   resp // submit_sm_resp returns
   ```

4. Event loop for receiving messages

   ```go
   for packet := range session.PDU() {
       // reply a responsable packet
       if p, ok := packet.(pdu.Responsable); ok {
           err := session.Send(p.Resp())
           if err != nil {
               fmt.Println(err)
           }
       }
   }
   ```

## LICENSE

This piece of software is released under [the MIT license](LICENSE).

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
