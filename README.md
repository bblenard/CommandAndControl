# CommandAndControl

## Testing

I have provided a docker-compose file under 'DockerStuff' that builds
an environment to test the client, server, and cli components. Testing
the C2 is done by running `build build.bash`. This will compile the
three executables and move them to the server / client directories in
DockerStuff. After that is done you can start the client and server via
`docker-compose build && docker-compose up`. The server also has the cli
binary so after bringing up the environment with `docker-compose up` you
can use the cli by using `docker run -it $SERVER_CONTAINER bash` to
execute `/app/cli`
