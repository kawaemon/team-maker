from golang:1.19-alpine as build

workdir /app
copy . .
run go build -v -tags netgo .

# ---

from gcr.io/distroless/static-debian11
copy --from=build /app/team-maker /
entrypoint ["/team-maker"]
