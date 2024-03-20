### Stats Ninja 🥷

This program will run on every server as docker container having docker unix socket mounted. 
It will collect system stats and usage stats of all running swarm services and will send to swiftwave management node.

> This program will not expose any API for polling, it will send its stats directly to swiftwave management server

### CLI

1. `enable` - Enable swiftwave stats ninja service
   ```bash
   sudo swiftwave_stats_ninja enable (http/https)://<swiftwave-management-node-ip>:<port> <analytics-token>
   ```
2. `disable` - Disable swiftwave stats ninja service
   ```bash
    sudo swiftwave_stats_ninja disable
    ```

### License
Apache License 2.0