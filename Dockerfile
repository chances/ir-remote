FROM balenalib/raspberry-pi-alpine

COPY /app/ir-remote .

CMD ./ir-remote
