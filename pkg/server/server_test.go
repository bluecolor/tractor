package server

// func TestServer(t *testing.T) {
// 	config := conf.Config{}
// 	repository, err := repo.NewRepository(config.DB)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	ts := httptest.NewServer(routes.BuildRoutes(repository))
// 	defer ts.Close()
// }

// func TestHelloHandler(t *testing.T) {

// 	ts := httptest.NewServer(
// 		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			fmt.Fprintln(w, "Hello, client")
// 		}))
// 	defer ts.Close()

// 	client := ts.Client()
// 	//res, err := client.Get(ts.URL)

// 	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
// 	// pass 'nil' as the third parameter.
// 	req, err := http.NewRequest("GET", "http://localhsot", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(helloHandler(client, "http://localhost"))

// 	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
// 	// directly and pass in our Request and ResponseRecorder.
// 	handler.ServeHTTP(rr, req)

// 	// Check the status code is what we expect.
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusOK)
// 	}

// 	// Check the response body is what we expect.
// 	expected := `{"alive": true}`
// 	if rr.Body.String() != expected {
// 		t.Errorf("handler returned unexpected body: got %v want %v",
// 			rr.Body.String(), expected)
// 	}
// }
