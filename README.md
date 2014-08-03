i3status-netusage is an i3status companion that shows the network traffic.

How to add it to the i3wm config:

    bar {
      status_command sh -c 'i3status | i3status-netusage --interface=em1'
    }

Sample output:

    19.1 KiB/s↓   21.2 KiB/s↑|0.22|2014-08-02 22:51:40
