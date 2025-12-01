## city := r.URL.Query().Get("city")
/weather?city=London

# How Go evaluates this code:
r.URL = the full URL object
.Query() = extract everything after ?
.Get("city") = get the value of city

# json.NewEncoder(w).Encode(resp)
# How Go evaluates this code:
json.NewEncoder(w) = create a new JSON encoder that writes to w (the HTTP response writer)
.Encode(resp) = encode the resp object as JSON and write it to w

# func wheaterHandler (w http.ResponseWriter, r *http.Request
- w http.ResponseWriter: the response writer used to send data back to the client
like : {"temperature":20,"condition":"Sunny"}

- r *http.Request: the incoming HTTP request from the client
like : /wheater?city=London
r = the full HTTP request object : 

# can this function takes other parameters?
No, in Go's net/http package, handler functions must have this specific signature to be compatible with the http.Handler interface.
 
# can weatherHandler take a different parameter type instead of http.ResponseWriter and *http.Request?
No, the parameters must be of type http.ResponseWriter and *http.Request to conform to the http.Handler interface.

os.getenv("API_KEY")
get the value of the environment variable named "API_KEY"