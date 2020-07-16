# listener

- run

        go build && ./listener

- deploy

        GOOS=linux GOARCH=amd64 go build
        nohup ./listener > out.log 2>&1 &



