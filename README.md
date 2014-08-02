i3status-netusage is an i3status companion that shows the network traffic.

How to add it to the i3wm config:

    bar {
      status_command sh -c 'i3status | i3status-netusage --interface=em1'
    }

Sample output:

    19.1 KiB/s↓   21.2 KiB/s↑

Know bugs:

  - the time is not yet correctly computed
  - the first data line from i3status will not have the network traffic
  - the traffic is added after the regular i3status.
