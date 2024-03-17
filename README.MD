### Stats Ninja 🥷

This program will run on every server as docker container having docker unix socket mounted. 
It will collect system stats and usage stats of all running swarm services and will send to swiftwave management node.

> This program will not expose any API for polling, it will send its stats directly to swiftwave management server

### Configuration

**Environment Variables**
- SUBMISSION_ENDPOINT: The endpoint to submit the stats to
- AUTHORIZATION_HEADER_VAL (optional): The value of the authorization header
- DOCKER_HOST: unix or tcp socket to connect to

**Volume Bind**
- bind **&lt;host docker socket path&gt;** to -> **/var/run/docker.sock**

### License
Apache License 2.0