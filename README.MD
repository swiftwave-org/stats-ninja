### Stats Ninja 🥷

This program will run on every server. 
It will collect system stats and usage stats of all running swarm services and will send to swiftwave management node.

> This program will not expose any API for polling, it will send its stats directly to swiftwave management server

### CLI

1. `enable` - Enable swiftwave stats ninja service
   ```bash
   sudo swiftwave-stats-ninja enable (http/https)://<swiftwave-management-node-ip>:<port> <analytics-token>
   ```
2. `disable` - Disable swiftwave stats ninja service
   ```bash
    sudo swiftwave-stats-ninja disable
    ```

### License
Apache License 2.0