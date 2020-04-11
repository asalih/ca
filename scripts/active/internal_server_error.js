
function analyze(response) {
    
    if (response.StatusCode == 500){
        return Found(LOW, "Internal Server Error", {"Parameter Type": response.Method })
    }
  }