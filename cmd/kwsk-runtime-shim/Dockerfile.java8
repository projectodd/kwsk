FROM adoptopenjdk/openjdk8-openj9:jdk8u162-b12_openj9-0.8.0

ARG OPENWHISK_RUNTIME_JAVA_VERSION="8@1.1.2"

# Fetch upstream image source and put it where upstream expects it
RUN curl -L -o java.tar.gz https://github.com/apache/incubator-openwhisk-runtime-java/archive/$OPENWHISK_RUNTIME_JAVA_VERSION.tar.gz \
  && mkdir /upstream \
  && tar --strip-components=1 -xf java.tar.gz -C /upstream \
  && mkdir /javaAction \
  && cp -R /upstream/core/java8/proxy/* /javaAction

# Move the java server to port 8081
RUN sed -ie "s/8080/8081/" /javaAction/src/main/java/openwhisk/java/action/Proxy.java

############ BEGIN upstream commands ############
RUN rm -rf /var/lib/apt/lists/* && apt-get clean && apt-get update \
	&& apt-get install -y --no-install-recommends locales \
	&& rm -rf /var/lib/apt/lists/* \
	&& locale-gen en_US.UTF-8
ENV LANG="en_US.UTF-8" \
	LANGUAGE="en_US:en" \
	LC_ALL="en_US.UTF-8" \
	VERSION=8 \
	UPDATE=162 \
	BUILD=12
RUN cd /javaAction \
	&& rm -rf .classpath .gitignore .gradle .project .settings Dockerfile build \
	&& ./gradlew oneJar \
	&& rm -rf /javaAction/src \
	&& ./compileClassCache.sh
############ END upstream commands ############

# Add our new Golang server shim, which runs on port 8080
COPY kwsk-runtime-shim /usr/local/bin/

# Add a little wrapper script that starts the shim and any other
# command passed in its arguments
COPY kwsk-wrapper.sh /usr/local/bin/

CMD ["kwsk-wrapper.sh", "java", "-Dfile.encoding=UTF-8", "-Xshareclasses:cacheDir=/javaSharedCache,readonly", "-Xquickstart", "-jar", "/javaAction/build/libs/javaAction-all.jar"]
