
function analyze(response) {
    
    var nsCookies = [];
    for(var i = 0; i < response.Cookies.length; i++){
        var elem = response.Cookies[i];

        if(!elem.HttpOnly){
            nsCookies.push(elem.Name)
        }

    }

    if (nsCookies.length > 0) {
      return Found(LOW, "Cookie Not Marked as HttpOnly", {"Identified Cookie(s)": nsCookies.toString(), "Cookie Source": "HTTP Header" })
    }
  }