  function analyze(response) {
    
    if (response.Headers["Server"] != "undefined") {
        console.log("Server Identified: " + response.Headers["Server"] + " in " + response.URL)
      return true
    }

    return false
  }