
function analyze(response) {
    
    if (response.Headers["Server"] != "undefined") {
      return Found(INFORMATION, "Server Identified", {"Server": response.Headers["Server"]})
    }
  }