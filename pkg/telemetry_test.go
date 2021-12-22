package cheek

// func TestPhoneHome(t *testing.T) {

// 	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		defer r.Body.Close()
// 		body, err := io.ReadAll(r.Body)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 		}
// 		// mirror this
// 		w.Header().Set("Content-Type", "application/json")
// 		fmt.Fprintln(w, string(body))
// 	}))

// 	defer testServer.Close()

// 	viper.Set("phoneHomeURL", testServer.URL)
// 	viper.Set("noTelemetry", false)

// 	et := ET{}
// 	_, err := et.PhoneHome()
// 	if err != nil {
// 		// no error means success
// 		t.Fatal(err)
// 	}

// 	// respect to not phone home
// 	viper.Set("noTelemetry", true)
// 	resp_body, err := et.PhoneHome()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	assert.Equal(t, resp_body, []byte{})
// }
