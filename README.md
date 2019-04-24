# HowlerMonkey

HowlerMonkey allows users to create a lightweight Push-Service that powered by Server-sent Events (SSE).

## Server-sent Events
Mozilla MDN Web Docs have an easy-to-understand explanation about SSE:

> Traditionally, a web page has to send a request to the server to receive new data; that is, the page requests data from the server. With server-sent events (SSE), it's possible for a server to send new data to a web page at any time, by pushing messages to the web page. These incoming messages can be treated as Events + data inside the web page.

## Why SSE OVER WebSockets
Many use WebSockets due to its capability of bi-directional and full-duplex communication. However, you do not always need to establish bi-directional communication. A typical example is an email or feed push service. Since my initial idea is to develop a lightweight HTTP Push-Service, SSE becomes my natural choice.

## Acknowledgement
This project modifies Kyle L. Jensen's codes in his 
[golang-html5-sse-example](https://github.com/kljensen/golang-html5-sse-example), which is based on 
[Leroy Campbell's SSE example in Go](https://gist.github.com/artisonian/3836281) and 
[the HTML5Rocks SSE tutorial](http://www.html5rocks.com/en/tutorials/eventsource/basics/). Many credits should go to Kyle and Leroy.

## How-To

### Install
```$xslt
go get github.com/iamharvey/HowlerMonkey@latest
```

### Start A Server
An SSE service can be considered as 'broker'. You can simply create a broker using `NewBroker()` 
and wrap it up as a HTTP server.
```$xslt
    ...
	// Create a new broker
	b := HowlerMonkey.NewBroker()

	// Start the broker
	b.Start()

	// Set up router
	r := chi.NewRouter()

	r.Get("/", home)
	r.Get("/events", b.GetEvents)
	r.Get("/send/{event}", b.SendEvent)

	s := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 60,
		ReadTimeout:  time.Second * 60,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	// Starts server gracefully
	log.Printf("[INFO] Starting server at http://%s\n", addr)
	Octopus.GracefulServe(s, false)
}
```
A full server example can be found [here](https://github.com/iamharvey/HowlerMonkey/blob/master/server/main.go).

### How to get events from server
You need create an **EventSource** object with the event source address (e.g. http://localhost:5678/events) 
first. Then, you can handle the events sent from the server by calling an EventHandler 
(e.g. **onmessage**, **onopen**, **onerror**). Here is an example: 

```$xslt
<div id="events"></div>
<script>
    if(typeof(EventSource) !== "undefined") {
        var source = new EventSource("http://localhost:5678/events");
        source.onmessage = function(event) {
            document.getElementById("events").innerHTML +=  (Date() + " " + event.data + "<br>");
        };
    } else {
        document.getElementById("events").innerHTML = "Your browser does not support server-sent events.";
    }
</script>
```

### Third-Party Modules Used In The Project
- [Octopus](https://github.com/NBCFB/Octopus) is a project that allows user to gracefully start, upgrade and shutdown an 
http server.
- [Chi](https://github.com/go-chi/chi) is an awesome lightweight router for GO HTTP services.


## License (Unlicense)
Since Kyle makes his codes UNLICENSE. I retain his choice of 
[UNLICENSED](https://github.com/iamharvey/HowlerMonkey/blob/master/LICENSE) for this modified version too.
