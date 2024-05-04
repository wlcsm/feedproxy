# RSS/Atom Feed proxy

I want to be able to filter RSS feeds based on keywords. Since NetNewsWire doesn't have this capability, I have decided to setup my own proxy for it.

This is a proxy that runs locally.

```
feedproxy [PORT]
```

You then direct your RSS/Atom feed requests at it by encoding the website address in the path and adding a "keyword" query parameter to perform the filter

Example:
```
http://127.0.0.1:8080/news.ycombinator.com/rss?keyword=quantum
```

This will then send a request to `https://news.ycombinator.com/rss` with all the HTTP headers from the original request. Then filters out all items that don't contain the keyword "quantum" (case insensitive).
