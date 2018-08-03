FROM openwhisk/python2action

# Move the python server to port 8081
ENV FLASK_PROXY_PORT 8081

# Add our new Golang server shim, which runs on port 8080
COPY kwsk-runtime-shim.cgo /usr/local/bin/kwsk-runtime-shim

# Add a little wrapper script that starts the shim and any other
# command passed in its arguments
COPY kwsk-wrapper.sh /usr/local/bin/

CMD ["/bin/bash", "-c", "cd pythonAction && kwsk-wrapper.sh python -u pythonrunner.py"]
