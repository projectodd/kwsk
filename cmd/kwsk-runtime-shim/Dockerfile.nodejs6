FROM openwhisk/nodejs6action

# Move the node server to port 8081
RUN sed -ie "s/8080/8081/" app.js

# Add our new Golang server shim, which runs on port 8080
COPY kwsk-runtime-shim /usr/local/bin/

# Add a little wrapper script that starts the shim and any other
# command passed in its arguments
COPY kwsk-wrapper.sh /usr/local/bin/

CMD kwsk-wrapper.sh node --expose-gc app.js
