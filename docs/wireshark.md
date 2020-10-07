# Wireshark

1. Discard `enquire_link` and `enquire_link_resp` packets

    ```plain
    smpp and !(smpp.command_id in {0x00000015 0x80000015})
    ```

2. Capturing

    ```shell
    tcpdump -w smpp.pcap port 2775
    # or
    tshark -w smpp.pcap port 2775
    ```

## References

- <https://smpp.org/smpp-testing-development.html>
- <https://wiki.wireshark.org/SMPP>
- <https://www.wireshark.org/docs/dfref/s/smpp.html>
