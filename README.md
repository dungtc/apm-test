# ELASTIC APM AGENT GOLANG TEST

#### Compose

Run docker-compose up. Compose will download the official docker containers and start Elasticsearch, Kibana, and APM Server.

#### SERVER

Add apmhttp wrapper for middleware handler
```
http.ListenAndServe(":8080", apmhttp.Wrap(mux))
```

#### CLIENT

Add apmhttp client wrapper for http request
```
client := apmhttp.WrapClient(http.DefaultClient)
```

#### Visualize

Use the APM app at http://localhost:5601/app/apm to visualize your application performance data!

