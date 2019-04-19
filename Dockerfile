from alpine:3.9

RUN apk update ; \
    apk upgrade ; \
    apk add bash

COPY build/rump /usr/local/bin/

RUN chmod +x /usr/local/bin/rump

ENTRYPOINT ["/bin/bash"]
